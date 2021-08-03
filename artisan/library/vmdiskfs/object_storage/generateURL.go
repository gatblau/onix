package object_storage

import (
	"go/types"
	"net/url"
	"time"
)

func GenerateUrl(useSSL bool, bucket string, filename string) (*url.URL, bool, error) {
	// Set request parameters for content-disposition
	reqParams := make(url.Values)
	// URL will expire in 24 hours - 1 day
	expiry := time.Second * 24 * 60 * 60

	// initializing object storage client
	s3Client := ObjectStorageProvider(useSSL)

	// Generate a pre-signed url with expiration 1 day
	preSignedURL, err := s3Client.PresignedGetObject(bucket, filename, expiry, reqParams)
	if err != nil {
		return &url.URL{}, false, err
	} else {
		return preSignedURL, true, types.Error{}
	}
}
