package job

/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"bytes"
	"encoding/json"
	"time"
)

// Result job result information
// note: ensure it is aligned with the same struct in pilotctl
type Result struct {
	// the unique job id
	JobId int64
	// indicates of the job was successful
	Success bool
	// the execution log for the job
	Log string
	// the error if any
	Err string
	// the completion time
	Time time.Time
}

func (r *Result) Reader() (*bytes.Reader, error) {
	jsonBytes, err := r.Bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*jsonBytes), err
}

func (r *Result) Bytes() (*[]byte, error) {
	b, err := ToJson(r)
	return &b, err
}

func ToJson(object interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	// switch off the escaping!
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(object)
	return buffer.Bytes(), err
}
