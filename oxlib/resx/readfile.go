/*
  Onix Config Manager - Onix Library
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package resx

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ReadFile use this function to generically fetch a file from a URI regardless of the type of location the file is in
// supported URI schemes are:
// - if prefix http:// or https:// then fetches from http endpoint
// - if prefix is none then reads from file system
// - if prefix is s3:// or s3s:// then fetches from s3 bucket endpoint
// - if prefix is ftp then returns a not supported error
// credentials are valid for http and s3 URIs and follow the syntax "user:pwd"
func ReadFile(uri, creds string) ([]byte, error) {
	if strings.HasPrefix(uri, "http") {
		return getHttpFile(uri, creds)
	} else if strings.HasPrefix(uri, "s3") {
		return getS3File(uri, creds)
	} else if strings.HasPrefix(uri, "ftp") {
		return getFtpFile(uri, creds)
	} else {
		return getFsFile(uri)
	}
	return nil, nil
}

// getFtpFile reads a file from an ftp endpoint
func getFtpFile(uri string, creds string) ([]byte, error) {
	return nil, fmt.Errorf("ftp scheme is not currently supported")
}

// getFsFile reads a file from the file system
func getFsFile(uri string) ([]byte, error) {
	path, err := filepath.Abs(uri)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(path)
}

// getHttpFile reads a file from a http endpoint
func getHttpFile(uri, creds string) ([]byte, error) {
	// if credentials are provided
	if len(creds) > 0 {
		// add them to the uri scheme
		u, err := addCredsToHttpURI(uri, creds)
		if err != nil {
			return nil, err
		}
		uri = u
	}
	// create an http client with defined timeout
	client := http.Client{
		Timeout: 60 * time.Second,
	}
	// create a new http request
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	// execute the request
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("cannot fetch '%s': %s\n", uri, resp.Status)
	}
	// return the byte content in the response body
	return ioutil.ReadAll(resp.Body)
}

// addCredsToHttpURI add credentials to http(s) URI
func addCredsToHttpURI(uri string, creds string) (string, error) {
	// if there are no credentials or the uri is a file path
	if len(creds) == 0 || strings.HasPrefix(uri, "http") {
		// skip and return as is
		return uri, nil
	}
	parts := strings.Split(uri, "/")
	if !strings.HasPrefix(parts[0], "http") {
		return uri, fmt.Errorf("invalid URI scheme, http(s) expected when specifying credentials\n")
	}
	parts[2] = fmt.Sprintf("%s@%s", creds, parts[2])
	return strings.Join(parts, "/"), nil
}

// getS3File reads a file from an S3 bucket
func getS3File(uri, creds string) ([]byte, error) {
	s3Client, bucketName, objectName, err := newS3Client(uri, creds)
	if err != nil {
		return nil, err
	}
	reader, err := s3Client.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	// create a new buffer to convert the reader to bytes
	buf := new(bytes.Buffer)
	// read into the buffer
	_, err = buf.ReadFrom(reader)
	// return the byte slice
	return buf.Bytes(), err
}

func newS3Client(uri, creds string) (client *minio.Client, bucketName, objectName string, err error) {
	var (
		endpoint string
		useSSL   bool
	)
	// if scheme is s3s use SSL
	if strings.HasPrefix(uri, "s3s://") {
		useSSL = true
		endpoint = uri[len("s3s://"):]
	} else if strings.HasPrefix(uri, "s3://") {
		useSSL = false
		endpoint = uri[len("s3://"):]
	} else {
		return nil, "", "", fmt.Errorf("invalid URI scheme: it should be s3:// or s3s://, uri was '%s'", uri)
	}
	p := strings.Split(endpoint, "/")
	if len(p) < 3 {
		return nil, "", "", fmt.Errorf("invalid URI, format should be [s3|s3s]://endpoint/bucket-name/object-name; it was: '%s'", uri)
	}
	endpoint = p[0]
	bucketName = p[1]
	objectName = p[2]
	for i := 3; i < len(p); i++ {
		objectName += "/" + p[i]
	}
	// store minio credentials
	var c *credentials.Credentials
	// if credentials provided
	if len(creds) > 0 {
		parts := strings.Split(creds, ":")
		if len(parts) != 2 {
			return nil, "", "", fmt.Errorf("invalid credentials, format should be ID:SECRET, provided '%s'", creds)
		}
		c = credentials.NewStaticV4(parts[0], parts[1], "")
	}
	// Initialize minio client object.
	client, err = minio.New(endpoint, &minio.Options{
		Creds:  c,
		Secure: useSSL,
	})
	return
}
