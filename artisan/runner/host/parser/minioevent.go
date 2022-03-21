/*
  Onix Config Manager - Artisan Host Runner
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package parser

import (
	"fmt"
	"net/url"
	"strings"
)

type MinioEvent struct {
	EventName string `json:"EventName"`
	Key       string `json:"Key"`
	Records   []struct {
		AwsRegion         string `json:"awsRegion"`
		EventName         string `json:"eventName"`
		EventSource       string `json:"eventSource"`
		EventTime         string `json:"eventTime"`
		EventVersion      string `json:"eventVersion"`
		RequestParameters struct {
			AccessKey       string `json:"accessKey"`
			Region          string `json:"region"`
			SourceIPAddress string `json:"sourceIPAddress"`
		} `json:"requestParameters"`
		ResponseElements struct {
			Content_length       string `json:"content-length"`
			XAmzRequestId        string `json:"x-amz-request-id"`
			XMinioDeploymentId   string `json:"x-minio-deployment-id"`
			XMinioOriginEndpoint string `json:"x-minio-origin-endpoint"`
		} `json:"responseElements"`
		S3 struct {
			Bucket struct {
				Arn           string `json:"arn"`
				Name          string `json:"name"`
				OwnerIdentity struct {
					PrincipalID string `json:"principalId"`
				} `json:"ownerIdentity"`
			} `json:"bucket"`
			ConfigurationID string `json:"configurationId"`
			Object          struct {
				ContentType  string `json:"contentType"`
				ETag         string `json:"eTag"`
				Key          string `json:"key"`
				Sequencer    string `json:"sequencer"`
				Size         int64  `json:"size"`
				UserMetadata struct {
					Content_type string `json:"content-type"`
				} `json:"userMetadata"`
				VersionID string `json:"versionId"`
			} `json:"object"`
			S3SchemaVersion string `json:"s3SchemaVersion"`
		} `json:"s3"`
		Source struct {
			Host      string `json:"host"`
			Port      string `json:"port"`
			UserAgent string `json:"userAgent"`
		} `json:"source"`
		UserIdentity struct {
			PrincipalID string `json:"principalId"`
		} `json:"userIdentity"`
	} `json:"Records"`
}

func (e MinioEvent) GetObjectDownloadURL() (string, error) {

	object := e.Records[0].S3.Object
	bucket := e.Records[0].S3.Bucket.Name
	//endpoint := os.Getenv("OBJECT_STORE_DOMAIN")
	endpoint := strings.Replace(e.Records[0].ResponseElements.XMinioOriginEndpoint, "http", "s3", 1)
	key, unescapeErr := url.PathUnescape(object.Key)
	if unescapeErr != nil {
		return "", unescapeErr
	}
	url := fmt.Sprintf("%s/%s/%s", endpoint, bucket, key)
	return url, nil
}
