//   Onix Config DatabaseProvider - Dbman
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

type Plan struct {
	Releases []Info `json:"releases"`
}

type Info struct {
	DbVersion  string `json:"dbVersion"`
	AppVersion string `json:"appVersion"`
	Path       string `json:"path"`
}

// get a JSON bytes reader for the Plan
func (plan *Plan) json() (*bytes.Reader, error) {
	jsonBytes, err := plan.bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*jsonBytes), err
}

// get a []byte representing the Plan
func (plan *Plan) bytes() (*[]byte, error) {
	b, err := oxc.ToJson(plan)
	return &b, err
}

// get the Plan in the http Response
func (plan *Plan) decode(response *http.Response) (*Plan, error) {
	result := new(Plan)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}

// return the info for the release of the specified app version
// Info: release information
func (plan *Plan) info(appVersion string) *Info {
	for _, release := range plan.Releases {
		if release.AppVersion == appVersion {
			return &release
		}
	}
	return nil
}
