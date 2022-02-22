package parser

import (
	"errors"
	"fmt"
	"os"
)

type S3Event struct {
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
			Content_length          string `json:"content-length"`
			X_amz_request_id        string `json:"x-amz-request-id"`
			X_minio_deployment_id   string `json:"x-minio-deployment-id"`
			X_minio_origin_endpoint string `json:"x-minio-origin-endpoint"`
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

func (e S3Event) GetObjectDownloadURL() (string, error) {
	d := os.Getenv("OBJECT_STORE_DOMAIN")
	err := e.validate()
	if err != nil {
		return "", err
	} else {
		url := fmt.Sprintf("%s/%s", d, e.Key)
		return url, nil
	}
}

func (e S3Event) validate() error {
	d := os.Getenv("OBJECT_STORE_DOMAIN")
	msg := ""
	if len(d) == 0 {
		msg = fmt.Sprintf("\n%s\n%s", "Error value for environment variable OBJECT_STORE_DOMAIN missing", msg)
		return errors.New(msg)
	}
	if len(e.Key) == 0 {
		msg = fmt.Sprintf("\n%s\n%s", "Error bucket details missing from event ", msg)
		return errors.New(msg)
	}

	if len(msg) != 0 {
		return errors.New(msg)
	}
	return nil
}
