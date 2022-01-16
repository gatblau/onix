/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
	"path/filepath"
)

func Tar(files []TarFile, buf io.Writer, preserveDirStruct bool) error {
	// create a new Writer for tar; writing to the tar writer will write to the "buf" writer
	// note: to additionally apply gzip can create gw := gzip.NewWriter(buf) and pass to "tw" instead of "buf"
	tw := tar.NewWriter(buf)
	defer tw.Close()

	// iterate over files and add them to the tar archive
	for _, file := range files {
		if len(file.Path) > 0 {
			err := addFileToTar(tw, file.Path, preserveDirStruct)
			if err != nil {
				return err
			}
		} else if len(file.Bytes) > 0 {
			err := addBytesToTar(tw, file.Name, file.Bytes, false)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func addFileToTar(tw *tar.Writer, filename string, preserveDirStruct bool) error {
	// open the file which will be written into the archive
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// get FileInfo about our file providing file size, mode, etc.
	info, err := file.Stat()
	if err != nil {
		return err
	}

	// create a tar Header from the FileInfo data
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	if preserveDirStruct {
		// use the full path as name (FileInfoHeader only takes the basename)
		// If we don't do this the directory structure would not be preserved
		// https://golang.org/src/archive/tar/common.go?#L626
		header.Name = filename
	}

	// Write file header to the tar archive
	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}

	// Copy file content to tar archive
	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}

	return nil
}

func addBytesToTar(tw *tar.Writer, filename string, file []byte, preserveDirStruct bool) error {
	tmp, err := NewTempDir()
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(tmp, filename), file, 0755)
	if err != nil {
		return err
	}
	info, err := os.Stat(filepath.Join(tmp, filename))
	if err != nil {
		return err
	}
	err = os.RemoveAll(tmp)
	if err != nil {
		return err
	}
	// create a tar Header from the FileInfo data
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	if preserveDirStruct {
		// use the full path as name (FileInfoHeader only takes the basename)
		// If we don't do this the directory structure would not be preserved
		// https://golang.org/src/archive/tar/common.go?#L626
		header.Name = filename
	}

	// Write file header to the tar archive
	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}

	// Copy file content to tar archive
	_, err = io.Copy(tw, bytes.NewReader(file))
	if err != nil {
		return err
	}

	return nil
}

func Untar(tarballReader io.Reader, outputPath string) error {
	tarReader := tar.NewReader(tarballReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(outputPath, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}
		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}
	return nil
}

type TarFile struct {
	Path  string
	Bytes []byte
	Name  string
}
