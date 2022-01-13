/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

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

// GetFile fetch a file from a URI
func GetFile(uri, creds string) ([]byte, error) {
	if strings.HasPrefix(uri, "http") {
		return getHttpFile(uri, creds)
	} else if strings.HasPrefix(uri, "s3") {
		return getS3File(uri, creds)
	} else if strings.HasPrefix(uri, "ftp") {
		return nil, fmt.Errorf("ftp scheme is not currently supported")
	} else {
		return getFsFile(uri, creds)
	}
	return nil, nil
}

// getFsFile reads a file from the file system
func getFsFile(uri, creds string) ([]byte, error) {
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
	defer resp.Body.Close()
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
	var (
		endpoint, bucketName, objectName string
		useSSL                           bool
	)
	// if scheme is s3s use SSL
	if strings.HasPrefix(uri, "s3s://") {
		useSSL = true
		endpoint = uri[len("s3s://"):]
	} else if strings.HasPrefix(uri, "s3://") {
		useSSL = false
		endpoint = uri[len("s3://"):]
	} else {
		return nil, fmt.Errorf("invalid URI scheme: it should be s3:// or s3s://, uri was '%s'", uri)
	}
	p := strings.Split(endpoint, "/")
	if len(p) < 3 {
		return nil, fmt.Errorf("invalid URI, format should be [s3|s3s]://endpoint/bucket-name/object-name; it was: '%s'", uri)
	}
	endpoint = p[0]
	bucketName = p[1]
	objectName = p[2]
	// store minio credentials
	var c *credentials.Credentials
	// if credentials provided
	if len(creds) > 0 {
		parts := strings.Split(creds, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid credentials, format should be ID:SECRET, provided '%s'", creds)
		}
		c = credentials.NewStaticV4(parts[0], parts[1], "")
	}
	// Initialize minio client object.
	s3Client, err := minio.New(endpoint, &minio.Options{
		Creds:  c,
		Secure: useSSL,
	})
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
