/*
  Onix Config Manager - Artisan's Doorman Proxy
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package main

import (
	"fmt"
	"strings"
	"time"
)

// Notification a notification to be sent by the service
type Notification struct {
	// Recipient of the notification if type is email
	Recipient string `yaml:"recipient" json:"recipient" example:"info@email.com"`
	// Type of the notification (e.g. email, snow, etc.)
	Type string `yaml:"type" json:"type" example:"email"`
	// Subject of the notification
	Subject string `yaml:"subject" json:"subject" example:"New Notification"`
	// Content of the template
	Content string `yaml:"content" json:"content" example:"A new event has been received."`
}

func (n Notification) Valid() error {
	if len(n.Type) == 0 {
		return fmt.Errorf("notification type is required and has not been provided")
	}
	if strings.ToUpper(n.Type) != "EMAIL" {
		return fmt.Errorf("notification type %s is not suppported", n.Type)
	}
	if len(n.Recipient) == 0 {
		return fmt.Errorf("notification recipient is required and has not been provided")
	}
	if len(n.Subject) == 0 {
		return fmt.Errorf("notification subject is required and has not been provided")
	}
	if len(n.Content) == 0 {
		return fmt.Errorf("notification content is required and has not been provided")
	}
	return nil
}

// MinioS3Event a notification event sent by MinIO when a file is uploaded to an S3 bucket
type MinioS3Event struct {
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
