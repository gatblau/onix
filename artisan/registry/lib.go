package registry

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/oxlib/httpserver"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// MoveFile use instead of os.Rename() to avoid issues moving a file whose source and destination paths are
// on different file systems or drive
// e.g. when running in Kubernetes by Tekton
func MoveFile(src string, dst string) (err error) {
	err = CopyFile(src, dst)
	if err != nil {
		return fmt.Errorf("failed to copy source file %s to %s: %s", src, dst, err)
	}
	err = os.RemoveAll(src)
	if err != nil {
		return fmt.Errorf("failed to cleanup source file %s: %s", src, err)
	}
	return nil
}

// credit https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}
	_, err = os.Stat(dst)
	// if the destination does not exist
	if os.IsNotExist(err) {
		// create the destination folder
		err = os.MkdirAll(dst, si.Mode())
		if err != nil {
			return
		}
	}
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

func MoveFolderContent(srcFolder, dstFolder string) error {
	srcFolder = core.ToAbs(srcFolder)
	dstFolder = core.ToAbs(dstFolder)
	file, err := os.Open(srcFolder)
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()

	files, err := file.Readdir(-1)
	if err != nil {
		return err
	}
	for _, info := range files {
		if info.IsDir() {
			err = CopyDir(path.Join(srcFolder, info.Name()), path.Join(dstFolder, info.Name()))
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(path.Join(srcFolder, info.Name()), path.Join(dstFolder, info.Name()))
			if err != nil {
				return err
			}
		}
	}
	return os.RemoveAll(srcFolder)
}

// Upload content to an http endpoint
func upload(client *http.Client, url string, values map[string]io.Reader, user string, pwd string) (err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for fieldname, reader := range values {
		var fw io.Writer
		// add a file
		if x, ok := reader.(*os.File); ok {
			fmt.Printf("adding file %s", x.Name())
			if fw, err = w.CreateFormFile(fieldname, x.Name()); err != nil {
				return
			}
		} else {
			// add other fields
			fmt.Printf("adding field %s", fieldname)
			if fw, err = w.CreateFormField(fieldname); err != nil {
				return
			}
		}
		_, err = io.Copy(fw, reader)
		// // if the reader can be closed
		// if x, ok := reader.(io.Closer); ok {
		// 	x.Close()
		// }
		// if there is an error return it
		if err != nil {
			return err
		}
	}
	// don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("accept", "application/json")
	if len(user) > 0 && len(pwd) > 0 {
		req.Header.Add("authorization", httpserver.BasicToken(user, pwd))
	}
	// Submit the request
	res, err := client.Do(req)
	if err != nil {
		return
	}

	// Check the response
	if res.StatusCode > 299 {
		return fmt.Errorf("failed to push, the remote server responded with status code %d: %s", res.StatusCode, res.Status)
	}
	return nil
}

func openFile(path string) *os.File {
	r, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return r
}

// unzip a package
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()
	err = os.MkdirAll(dest, 0755)
	if err != nil {
		return err
	}
	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()
		path := filepath.Join(dest, f.Name)
		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()
			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}
	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}
	return nil
}

// remove an item from a slice
func removeItem(slice []string, item string) []string {
	var ix int = -1
	for i := 0; i < len(slice); i++ {
		if slice[i] == item {
			ix = i
			break
		}
	}
	if ix > -1 {
		return remove(slice, ix)
	}
	return slice
}

func remove(slice []string, ix int) []string {
	slice[ix] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}
