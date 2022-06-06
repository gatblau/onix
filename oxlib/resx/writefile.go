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
	"os"
	"path/filepath"
)

// WriteFile writes a file to a URI
// supported URI schemes are:
// - if prefix is none then writes to file system
// - if prefix is s3:// or s3s:// then put the file to s3 bucket endpoint
// - if prefix http:// or https:// then returns a not supported error
// - if prefix is ftp then returns a not supported error
// credentials are valid for s3 URIs and follow the syntax "user:pwd"
func WriteFile(content []byte, uri, creds string) error {
	switch ParseUriType(uri) {
	case File:
		return writeFsFile(content, uri)
	case Https:
		return writeHttpFile(content, uri, creds)
	case Http:
		return writeHttpFile(content, uri, creds)
	case S3:
		return writeS3File(content, uri, creds)
	case S3S:
		return writeS3File(content, uri, creds)
	case Ftps:
		return writeFtpFile(content, uri, creds)
	case Ftp:
		return writeFtpFile(content, uri, creds)
	case Unknown:
		return fmt.Errorf("unkknown URI type: %s", uri)
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
	ctx := context.Background()
	s3Client, bucketName, objectName, err := newS3Client(uri, creds)
	if err != nil {
		return err
	}
	// Check to see if we already own this bucket (which happens if you run this twice)
	var exists bool
	exists, err = s3Client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("Failed to check if bucket exists: %s\n", err)
	}
	// if the bucket does not exist, attempts to create it
	if !exists {
		err = s3Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: ""})
		if err != nil {
			return fmt.Errorf("Failed to create bucket: %s\n", err)
		}
	}
	_, err = s3Client.PutObject(
		ctx,
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
