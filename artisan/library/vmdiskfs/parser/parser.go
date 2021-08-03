package parser

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// ReadPutNotification - reading put notification received from object storage notification service
func ReadPutNotification() []byte {
	pwd := os.Getenv("PIPELINE_HOME")
	readNotification, err := ioutil.ReadFile(pwd + "/.artisan/files/context")
	if err != nil {
		log.Printf("Error Reading file %s", err)
	}
	return readNotification
}

// ParseNotification - filters notification based on bucket name and filename,
// supports currently AWS S3 and MinIO
func ParseNotification(notification []byte, provider string) (string, string) {
	minio := MinIO{}
	awss3 := AwsS3{}
	if provider == "minio" {
		err := json.Unmarshal(notification, &minio)
		if err != nil {
			log.Println("----------Put Notification message is not MinIO Style-----------")
		}
		return minio.Records[0].S3.Bucket.Name, minio.Records[0].S3.Object.Key
	} else if provider == "aws" {
		awsStr1 := strings.ReplaceAll(string(notification), "\\", "")
		awsStr2 := strings.ReplaceAll(awsStr1, "\"{", "{")
		awsStr3 := strings.ReplaceAll(awsStr2, "}\"", "}")
		if err := json.Unmarshal([]byte(awsStr3), &awss3); err != nil {
			log.Printf("Not valid content: %s", err)
		}
		return awss3.Message.Records[0].S3.Bucket.Name, awss3.Message.Records[0].S3.Object.Key
	} else {
		return "na", "na"
	}
}

// ParsedInformationWriter - Bucket and file names saves to file
func ParsedInformationWriter(provider string) bool {
	notification := ReadPutNotification()
	bucket, fname := ParseNotification(notification, provider)

	if len(bucket) == 0 && len(fname) == 0 {
		log.Printf("Bucket and File name are: %i", 0)
		return false
	} else {
		file, err := os.Create("scripts/parsed.txt")
		if err != nil {
			log.Printf("File create Error: %s", err)
		}
		if _, err := file.WriteString(bucket + "\n"); err != nil {
			log.Println(err)
		}
		if _, err := file.WriteString(fname + "\n"); err != nil {
			log.Println(err)
		}
		return true
	}
}
