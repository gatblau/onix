/*
  Onix Config Manager - Artie
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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
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

// the artefact id calculated as the hex encoded SHA-256 digest of the artefact seal
func artefactId(seal *seal) string {
	// serialise the seal info to json
	info := toJsonBytes(seal)
	hash := sha256.New()
	// copy the seal manifest into the hash
	if _, err := io.Copy(hash, bytes.NewReader(info)); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(hash.Sum(nil))
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
	return http.DetectContentType(buffer), nil
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
func copyFiles(src string, dst string) error {
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
			if err = copyFiles(srcFp, dstFp); err != nil {
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

// remove an element in a slice
func removeElement(a []string, value string) []string {
	i := -1
	// find the value to remove
	for ix := 0; ix < len(a); ix++ {
		if a[ix] == value {
			i = ix
			break
		}
	}
	if i == -1 {
		return a
	}
	// Remove the element at index i from a.
	a[i] = a[len(a)-1] // Copy last element to index i.
	a[len(a)-1] = ""   // Erase last element (write zero value).
	a = a[:len(a)-1]   // Truncate slice.
	return a
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
	return strconv.FormatFloat(getSize, 'f', -1, 64) + string(getSuffix)
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

// returns the elapsed time until now in human friendly format
func toElapsedLabel(rfc850time string) string {
	created, err := time.Parse(time.RFC850, rfc850time)
	if err != nil {
		log.Fatal(err)
	}
	elapsed := time.Since(created)
	seconds := elapsed.Seconds()
	minutes := elapsed.Minutes()
	hours := elapsed.Hours()
	days := hours / 24
	weeks := days / 7
	months := weeks / 4
	years := months / 12

	if math.Trunc(years) > 0 {
		return fmt.Sprintf("%d %s ago", int64(years), plural(int64(years), "year"))
	} else if math.Trunc(months) > 0 {
		return fmt.Sprintf("%d %s ago", int64(months), plural(int64(months), "month"))
	} else if math.Trunc(weeks) > 0 {
		return fmt.Sprintf("%d %s ago", int64(weeks), plural(int64(weeks), "week"))
	} else if math.Trunc(days) > 0 {
		return fmt.Sprintf("%d %s ago", int64(days), plural(int64(days), "day"))
	} else if math.Trunc(hours) > 0 {
		return fmt.Sprintf("%d %s ago", int64(hours), plural(int64(hours), "hour"))
	} else if math.Trunc(minutes) > 0 {
		return fmt.Sprintf("%d %s ago", int64(minutes), plural(int64(minutes), "minute"))
	}
	return fmt.Sprintf("%d %s ago", int64(seconds), plural(int64(seconds), "second"))
}

// turn label into plural if value is greater than one
func plural(value int64, label string) string {
	if value > 1 {
		return fmt.Sprintf("%ss", label)
	}
	return label
}

// upload an artefact to Nexus 3
func nexus3Upload(client *http.Client, remoteUrl, remoteFolder, localFolder, fileNoExtension string) error {
	// prepare the reader instances to encode
	values := map[string]io.Reader{
		"raw.directory": strings.NewReader(remoteFolder),
		// the json filename
		"raw.asset1.filename": strings.NewReader(fmt.Sprintf("%s.json", fileNoExtension)),
		// the json file (seal)
		"raw.asset1": mustOpen(fmt.Sprintf("%s/%s.json;type=application/json", localFolder, fileNoExtension)),
		// the zip filename
		"raw.asset2.filename": strings.NewReader(fmt.Sprintf("%s.zip", fileNoExtension)),
		// the zip file (artefact)
		"raw.asset2": mustOpen(fmt.Sprintf("%s/%s.zip;type=application/zip", localFolder, fileNoExtension)),
	}
	return upload(client, remoteUrl, values)
}

// upload content to an http endpoint
func upload(client *http.Client, url string, values map[string]io.Reader) (err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// add a file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return
			}
		} else {
			// add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
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

	// Submit the request
	res, err := client.Do(req)
	if err != nil {
		return
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	return
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}
