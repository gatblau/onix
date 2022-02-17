/*
  Onix Config Manager - Artisan's Doorman
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package types

import "time"

// S3MinioEvent a notification event sent by MinIO when a file is uploaded to an S3 bucket
type S3MinioEvent struct {
	EventName string    `json:"EventName"`
	Key       string    `json:"Key"`
	Records   []Records `json:"Records"`
}

type RequestParameters struct {
	AccessKey       string `json:"accessKey"`
	Region          string `json:"region"`
	SourceIPAddress string `json:"sourceIPAddress"`
}

type ResponseElements struct {
	ContentLength        string `json:"content-length"`
	XAmzRequestID        string `json:"x-amz-request-id"`
	XMinioDeploymentID   string `json:"x-minio-deployment-id"`
	XMinioOriginEndpoint string `json:"x-minio-origin-endpoint"`
}

type OwnerIdentity struct {
	PrincipalID string `json:"principalId"`
}

type Bucket struct {
	Arn           string        `json:"arn"`
	Name          string        `json:"name"`
	OwnerIdentity OwnerIdentity `json:"ownerIdentity"`
}

type UserMetadata struct {
	ContentType string `json:"content-type"`
}

type Object struct {
	ContentType  string       `json:"contentType"`
	ETag         string       `json:"eTag"`
	Key          string       `json:"key"`
	Sequencer    string       `json:"sequencer"`
	Size         int          `json:"size"`
	UserMetadata UserMetadata `json:"userMetadata"`
	VersionID    string       `json:"versionId"`
}

type S3 struct {
	Bucket          Bucket `json:"bucket"`
	ConfigurationID string `json:"configurationId"`
	Object          Object `json:"object"`
	S3SchemaVersion string `json:"s3SchemaVersion"`
}

type Source struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	UserAgent string `json:"userAgent"`
}

type UserIdentity struct {
	PrincipalID string `json:"principalId"`
}

type Records struct {
	AwsRegion         string            `json:"awsRegion"`
	EventName         string            `json:"eventName"`
	EventSource       string            `json:"eventSource"`
	EventTime         time.Time         `json:"eventTime"`
	EventVersion      string            `json:"eventVersion"`
	RequestParameters RequestParameters `json:"requestParameters"`
	ResponseElements  ResponseElements  `json:"responseElements"`
	S3                S3                `json:"s3"`
	Source            Source            `json:"source"`
	UserIdentity      UserIdentity      `json:"userIdentity"`
}
