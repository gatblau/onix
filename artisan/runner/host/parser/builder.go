package parser

import (
	"encoding/json"
)

func NewS3Event(f []byte) (*S3Event,error) {
	s3event := S3Event{}
	if err := json.Unmarshal(f, &s3event); err != nil {
		return nil, err
	}
	return &s3event, nil
}