//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package util

import (
	"bytes"
	"encoding/json"
	"github.com/gatblau/oxc"
	"net/http"
)

type Index struct {
	Releases []ReleaseInfo `json:"releases"`
}

type ReleaseInfo struct {
	DbVersion  string `json:"dbVersion"`
	AppVersion string `json:"appVersion"`
	Path       string `json:"path"`
}

// get a JSON bytes reader for the Index
func (ix *Index) json() (*bytes.Reader, error) {
	jsonBytes, err := ix.bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*jsonBytes), err
}

// get a []byte representing the Index
func (ix *Index) bytes() (*[]byte, error) {
	b, err := oxc.ToJson(ix)
	return &b, err
}

// get the Index in the http Response
func (ix *Index) decode(response *http.Response) (*Index, error) {
	result := new(Index)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}
