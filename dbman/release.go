//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// describes a database release
type Release struct {
	Release string `json:"release"`
	Schemas []struct {
		File     string `json:"file"`
		Checksum string `json:"checksum"`
	} `json:"schemas"`
	Functions []struct {
		File     string `json:"file"`
		Checksum string `json:"checksum"`
	} `json:"functions"`
	Upgrade []struct {
		File     string `json:"file"`
		Checksum string `json:"checksum"`
	} `json:"upgrade"`
}

// get a JSON bytes reader for the Index
func (r *Release) json() (*bytes.Reader, error) {
	jsonBytes, err := r.bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*jsonBytes), err
}

// get a []byte representing the Index
func (r *Release) bytes() (*[]byte, error) {
	b, err := jsonBytes(r)
	return &b, err
}

// get the Index in the http Response
func (r *Release) decode(response *http.Response) (*Release, error) {
	result := new(Release)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}
