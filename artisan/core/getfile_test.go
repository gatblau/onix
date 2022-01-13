/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"fmt"
	"testing"
)

// TestGetFileS3 test the retrieval of a file from an S£ bucket
// launch minio service: docker run -p 9000:9000 -p 9001:9001 quay.io/minio/minio server /data --console-address ":9001"
// create user called "abcdefgh" with password "12345678"
// create a bucket called "test"
// upload file.txt to the bucket
func TestGetFileS3(t *testing.T) {
	content, err := GetFile("s3://127.0.0.1:9000/test/file.txt", "abcdefgh:12345678")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("downloaded %d bytes", len(content))
}
