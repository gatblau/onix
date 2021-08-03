package parser

import "time"

// put notification structure for MinIO object storage bucket
type MinIO struct {
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

// put notification structure for AWS S3 object storage bucket
type AwsS3 struct {
	Message struct {
		Records []struct {
			EventVersion string    `json:"eventVersion"`
			EventSource  string    `json:"eventSource"`
			AwsRegion    string    `json:"awsRegion"`
			EventTime    time.Time `json:"eventTime"`
			EventName    string    `json:"eventName"`
			UserIdentity struct {
				PrincipalID string `json:"principalId"`
			} `json:"userIdentity"`
			RequestParameters struct {
				SourceIPAddress string `json:"sourceIPAddress"`
			} `json:"requestParameters"`
			ResponseElements struct {
				XAmzRequestID string `json:"x-amz-request-id"`
				XAmzID2       string `json:"x-amz-id-2"`
			} `json:"responseElements"`
			S3 struct {
				S3SchemaVersion string `json:"s3SchemaVersion"`
				ConfigurationID string `json:"configurationId"`
				Bucket          struct {
					Name          string `json:"name"`
					OwnerIdentity struct {
						PrincipalID string `json:"principalId"`
					} `json:"ownerIdentity"`
					Arn string `json:"arn"`
				} `json:"bucket"`
				Object struct {
					Key       string `json:"key"`
					Size      int    `json:"size"`
					ETag      string `json:"eTag"`
					Sequencer string `json:"sequencer"`
				} `json:"object"`
			} `json:"s3"`
		} `json:"Records"`
	} `json:"Message"`
	MessageID        string    `json:"MessageId"`
	Signature        string    `json:"Signature"`
	SignatureVersion string    `json:"SignatureVersion"`
	SigningCertURL   string    `json:"SigningCertURL"`
	Subject          string    `json:"Subject"`
	Timestamp        time.Time `json:"Timestamp"`
	TopicArn         string    `json:"TopicArn"`
	Type             string    `json:"Type"`
	UnsubscribeURL   string    `json:"UnsubscribeURL"`
}
