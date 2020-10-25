/*
  Onix Config Manager - Art
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	h "net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"
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

// convert the passed in parameter to a Json Byte Array
func toJsonBytes(s interface{}) []byte {
	// serialise the seal to json
	source, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	// indent the json to make it readable
	dest := new(bytes.Buffer)
	err = json.Indent(dest, source, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return dest.Bytes()
}

// takes the combined checksum of the seal information and the compressed file
func checksum(path string, sealData *manifest) []byte {
	// read the compressed file
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	// serialise the seal info to json
	info := toJsonBytes(sealData)
	hash := sha256.New()
	// copy the seal manifest into the hash
	if _, err := io.Copy(hash, bytes.NewReader(info)); err != nil {
		log.Fatal(err)
	}
	// copy the compressed file into the hash
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}
	return hash.Sum(nil)
}

// zip a file or a folder
func zipSource(source, target string) error {
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
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}
		if info.IsDir() {
			header.Name += "/"
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
				log.Fatal(err)
			}
		}()
		_, err = io.Copy(writer, file)
		return err
	})
	return err
}

// detect the content type of a given file
func findContentType(f *os.File) (string, error) {
	// get the first 512 bytes to sniff the content type
	buffer := make([]byte, 512)
	_, err := f.Read(buffer)
	if err != nil {
		return "", err
	}
	return h.DetectContentType(buffer), nil
}

// executes a single command with arguments
func execute(cmd string, dir string, env []string) {
	strArr := strings.Split(cmd, " ")
	var c *exec.Cmd
	if len(strArr) == 1 {
		c = exec.Command(strArr[0])
	} else {
		c = exec.Command(strArr[0], strArr[1:]...)
	}
	c.Dir = dir
	c.Env = env
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	log.Printf("executing: %s\n", strings.Join(c.Args, " "))
	if err := c.Start(); err != nil {
		log.Fatal(err)
	}
	err := c.Wait()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0
			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				log.Fatal(exitMsg(status.ExitStatus()))
			}
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
		log.Fatal(err)
	}
}

// gets the error message for a shell exit status
func exitMsg(exitCode int) string {
	switch exitCode {
	case 1:
		return "exit code 1 - general error"
	case 2:
		return "exit code 2 - misuse of shell built-ins"
	case 126:
		return "exit code 126 - command invoked cannot execute"
	case 127:
		return "exit code 127 - command not found"
	case 128:
		return "exit code 128 - invalid argument to exit"
	case 130:
		return "exit code 130 - script terminated by CTRL-C"
	default:
		return fmt.Sprintf("exist code %d", exitCode)
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
		log.Fatal("error: target not found after command execution")
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
			log.Fatal(err)
		}
	}()

	if dstFd, err = os.Create(dst); err != nil {
		return err
	}
	defer func() {
		err := dstFd.Close()
		if err != nil {
			log.Fatal(err)
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
func (b *Builder) copyFiles(src string, dst string) error {
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
			if err = b.copyFiles(srcFp, dstFp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = copyFile(srcFp, dstFp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}
