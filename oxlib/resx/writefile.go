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
	"github.com/minio/minio-go/v7/pkg/notification"
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
	// get the directory part of the path
	dir := filepath.Dir(path)
	// if it does not exist
	_, err = os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		// creates target directory
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return os.WriteFile(path, content, 0755)
}

// writeS3File write a file to an S3 bucket
func writeS3File(content []byte, uri string, creds string) error {
	ctx := context.Background()
	s3Client, bucketName, objectName, err := NewS3Client(uri, creds)
	if err != nil {
		return err
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

// EnsureBucketNotification check that a bucket exists and if not it creates it
// If a bucket is created, and notification information is provided then creates a bucket notification
func EnsureBucketNotification(uri, creds, filterSuffix string, arn *notification.Arn) (*minio.Client, error) {
	ctx := context.Background()
	s3Client, bucketName, _, err := NewS3Client(uri, creds)
	if err != nil {
		return nil, err
	}
	// Check to see if we already own this bucket (which happens if you run this twice)
	var exists bool
	exists, err = s3Client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("Failed to check if bucket exists: %s\n", err)
	}
	// if the bucket does not exist, attempts to create it
	if !exists {
		err = s3Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: ""})
		if err != nil {
			return nil, fmt.Errorf("Failed to create bucket: %s\n", err)
		}
		if arn != nil {
			// creates the notification configuration
			cfg := notification.NewConfig(*arn)
			cfg.AddEvents(notification.ObjectCreatedPut)
			cfg.AddFilterSuffix(filterSuffix)
			config := notification.Configuration{}
			config.AddQueue(cfg)
			// set the bucket notification
			err = s3Client.SetBucketNotification(ctx, bucketName, config)
			if err != nil {
				return nil, fmt.Errorf("Failed to set bucket notification for ARN '%s': %s\n", arn.String(), err)
			}
		}
	}
	return s3Client, nil
}
