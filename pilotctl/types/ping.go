package types

/*
Onix Config Manager - Pilot Control
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

import (
	"bytes"
	"fmt"
	"time"
)

// NewPingResponse creates a new ping response
func NewPingResponse(cmdInfo CmdInfo, pingInterval time.Duration) (*PingResponse, error) {
	// create a signature for the envelope
	envelope := PingResponseEnvelope{
		Command:  cmdInfo,
		Interval: pingInterval,
	}
	signature, err := sign(envelope)
	if err != nil {
		return nil, fmt.Errorf("cannot sign ping response: %s", err)
	}
	return &PingResponse{
		Signature: signature,
		Envelope:  envelope,
	}, nil
}

// PingResponse a command for execution with a job reference
type PingResponse struct {
	// the envelope signature
	Signature string `json:"signature"`
	// the signed content sent to pilot
	Envelope PingResponseEnvelope `json:"envelope"`
}

// PingResponseEnvelope contains the signed content sent to pilot
type PingResponseEnvelope struct {
	// the information about the command to execute
	Command CmdInfo `json:"value"`
	// the ping interval
	Interval time.Duration `json:"interval"`
}

type PingRequest struct {
	Result *JobResult `json:"result,omitempty"`
	Events *Events    `json:"events,omitempty"`
}

func (r *PingRequest) Reader() (*bytes.Reader, error) {
	jsonBytes, err := r.Bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*jsonBytes), err
}

func (r *PingRequest) Bytes() (*[]byte, error) {
	b, err := ToJson(r)
	return &b, err
}
