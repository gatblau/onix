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
	"os"
	"path/filepath"
	"strings"
)

// WriteFile writes a file to a URI
// supported URI schemes are:
// - if prefix is none then writes to file system
// - if prefix is s3:// or s3s:// then put the file to s3 bucket endpoint
// - if prefix http:// or https:// then returns a not supported error
// - if prefix is ftp then returns a not supported error
// credentials are valid for s3 URIs and follow the syntax "user:pwd"
func WriteFile(content []byte, uri, creds string) error {
	if strings.HasPrefix(uri, "http") {
		return writeHttpFile(content, uri, creds)
	} else if strings.HasPrefix(uri, "s3") {
		return writeS3File(content, uri, creds)
	} else if strings.HasPrefix(uri, "ftp") {
		return writeFtpFile(content, uri, creds)
	} else {
		return writeFsFile(content, uri)
	}
	return nil
}

// write a file to the file system
func writeFsFile(content []byte, uri string) error {
	path, err := filepath.Abs(uri)
	if err != nil {
		return err
	}
	return os.WriteFile(path, content, 0755)
}

// writeS3File write a file to an S3 bucket
func writeS3File(content []byte, uri string, creds string) error {
	s3Client, bucketName, objectName, err := newS3Client(uri, creds)
	if err != nil {
		return err
	}
	_, err = s3Client.PutObject(
		context.Background(),
		bucketName,
		objectName,
		bytes.NewReader(content),
		int64(len(content)),
		minio.PutObjectOptions{
			ContentType: "application/octet-stream",
		})
	return err
}

func writeFtpFile(content []byte, uri string, creds string) error {
	return fmt.Errorf("ftp scheme is not currently supported")
}

func writeHttpFile(content []byte, uri string, creds string) error {
	return fmt.Errorf("http scheme is not currently supported")
}
