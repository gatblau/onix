package build

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"archive/zip"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/merge"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// merge two or more maps
// the latter map overrides the former if duplicate keys exist across the two maps
func mergeMaps(ms ...map[string]string) map[string]string {
	res := map[string]string{}
	for _, m := range ms {
		for k, v := range m {
			res[k] = v
		}
	}
	return res
}

// zip a file or a folder
func zipSource(source, target string, excludeSource []string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer func() {
		err := zipfile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	archive := zip.NewWriter(zipfile)
	defer func() {
		err := archive.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	info, err := os.Stat(source)
	if err != nil {
		return nil
	}
	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}
	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// do not add to the zip file excluded sources
		if contains(source, excludeSource) {
			return nil
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}
		if info.IsDir() {
			header.Name += string(os.PathSeparator)
		} else {
			header.Method = zip.Deflate
		}
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() {
			err := file.Close()
			if err != nil {
				log.Print(err)
				runtime.Goexit()
			}
		}()
		_, err = io.Copy(writer, file)
		return err
	})
	return err
}

// detect the content type of given file
func findContentType(f *os.File) (string, error) {
	// get the first 512 bytes to sniff the content type
	buffer := make([]byte, 512)
	_, err := f.Read(buffer)
	if err != nil {
		return "", err
	}
	return http.DetectContentType(buffer), nil
}

// gets the error message for a shell exit status
func exitMsg(exitCode int) string {
	switch exitCode {
	case 1:
		return "error 1 - general error"
	case 2:
		return "error 2 - misuse of shell built-ins (check for permission or access problem)"
	case 126:
		return "error 126 - command invoked cannot execute (check for permission problem)"
	case 127:
		return "error 127 - command not found (check for typos or missing commands)"
	case 128:
		return "error 128 - invalid argument to exit (check when you are not returning something that is not integer args in the range 0 – 255)"
	case 130:
		return "error 130 - script terminated by CTRL-C"
	default:
		return fmt.Sprintf("exit code %d", exitCode)
	}
}

// wait a time duration for a file or folder to be created on the path
func waitForTargetToBeCreated(path string) {
	elapsed := 0
	found := false
	for {
		_, err := os.Stat(path)
		if !os.IsNotExist(err) {
			found = true
			break
		}
		if elapsed > 30 {
			break
		}
		elapsed++
		time.Sleep(500 * time.Millisecond)
	}
	if !found {
		core.RaiseErr("target '%s' not found after command execution", path)
	}
}

// copy a single file
func copyFile(src, dst string) error {
	var err error
	var srcFd *os.File
	var dstFd *os.File
	var srcInfo os.FileInfo

	if srcFd, err = os.Open(src); err != nil {
		return err
	}
	defer func() {
		err := srcFd.Close()
		if err != nil {
			log.Print(err)
			runtime.Goexit()
		}
	}()
	if dstFd, err = os.Create(dst); err != nil {
		return err
	}
	defer func() {
		err := dstFd.Close()
		if err != nil {
			log.Print(err)
			runtime.Goexit()
		}
	}()
	if _, err = io.Copy(dstFd, srcFd); err != nil {
		return err
	}
	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

// copy the files in a folder recursively
func copyFolder(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcInfo os.FileInfo
	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}
	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}
	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcFp := path.Join(src, fd.Name())
		dstFp := path.Join(dst, fd.Name())
		if fd.IsDir() {
			if err = copyFolder(srcFp, dstFp); err != nil {
				core.ErrorLogger.Printf(err.Error())
			}
		} else {
			if err = copyFile(srcFp, dstFp); err != nil {
				core.ErrorLogger.Printf(err.Error())
			}
		}
	}
	return nil
}

func renameFile(src string, dst string) (err error) {
	err = copyFile(src, dst)
	if err != nil {
		return fmt.Errorf("failed to copy source file %s to %s: %s", src, dst, err)
	}
	err = os.RemoveAll(src)
	if err != nil {
		return fmt.Errorf("failed to cleanup source file %s: %s", src, err)
	}
	return nil
}

// converts a byte count into a pretty label
func bytesToLabel(size int64) string {
	var suffixes [5]string
	suffixes[0] = "B"
	suffixes[1] = "KB"
	suffixes[2] = "MB"
	suffixes[3] = "GB"
	suffixes[4] = "TB"
	base := math.Log(float64(size)) / math.Log(1024)
	getSize := round(math.Pow(1024, base-math.Floor(base)), .5, 2)
	getSuffix := suffixes[int(math.Floor(base))]
	return strconv.FormatFloat(getSize, 'f', -1, 64) + getSuffix
}

func round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

// executes a command and sends output and error streams to stdout and stderr
func execute(cmd string, dir string, env *merge.Envar, interactive bool) (err error) {
	// executes the command
	_, err = ExeAsync(cmd, dir, env, interactive)
	// if there is an error return it
	if err != nil {
		return err
	}
	// return without error
	return nil
}

func contains(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func findGitPath(path string) (string, error) {
	for {
		_, err := os.Stat(filepath.Join(path, ".git"))
		if os.IsNotExist(err) {
			path = filepath.Dir(path)
			if strings.HasSuffix(path, string(os.PathSeparator)) {
				return "", fmt.Errorf("cannot find .git path")
			}
		} else {
			return path, nil
		}
	}
}

// check the the specified function is in the manifest
func isExported(m *data.Manifest, fx string) bool {
	for _, function := range m.Functions {
		if function.Name == fx {
			return true
		}
	}
	return false
}
