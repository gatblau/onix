/*
  Onix Config Manager - Artisan's Doorman
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package types

import (
	"fmt"
)

const (
	minioEventName = "s3:ObjectCreated:Put"
)

// NewSpecEvent an event indicating that doorman should start the ingestion process
type NewSpecEvent struct {
	URI    string
	Bucket string
}

func NewSpecEventFromMinio(ev S3MinioEvent) (*NewSpecEvent, error) {
	// check the event type is correct
	if ev.EventName != minioEventName {
		return nil, fmt.Errorf("received notification is of type %s, but should be %s", ev.EventName, minioEventName)
	}
	if ev.Records == nil || len(ev.Records) == 0 {
		return nil, fmt.Errorf("invalid notification does not contains Records")
	}
	if ev.Records[0].S3.Object.Key != "spec.yaml" {
		return nil, fmt.Errorf("invalid notification refers to file %s, but it should be for spec.yaml", ev.Records[0].S3.Object.Key)
	}
	return &NewSpecEvent{
		// the location of the service originating the event
		URI: ev.Records[0].ResponseElements.XMinioOriginEndpoint,
		// the name of the bucket that triggered the event
		Bucket: ev.Records[0].S3.Bucket.Name,
	}, nil
}
