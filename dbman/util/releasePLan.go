//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gatblau/oxc"
	"gopkg.in/yaml.v2"
	"net/http"
	"strings"
)

type ReleasePlan struct {
	Releases []ReleaseInfo `json:"releases"`
}

type ReleaseInfo struct {
	DbVersion  string `json:"dbVersion"`
	AppVersion string `json:"appVersion"`
	Path       string `json:"path"`
}

// get a JSON bytes reader for the ReleasePlan
func (plan *ReleasePlan) json() (*bytes.Reader, error) {
	jsonBytes, err := plan.bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*jsonBytes), err
}

// get a []byte representing the ReleasePlan
func (plan *ReleasePlan) bytes() (*[]byte, error) {
	b, err := oxc.ToJson(plan)
	return &b, err
}

// get the ReleasePlan in the http Response
func (plan *ReleasePlan) decode(response *http.Response) (*ReleasePlan, error) {
	result := new(ReleasePlan)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}

func (plan *ReleasePlan) Format(format string) string {
	switch strings.ToLower(format) {
	case "yml":
		fallthrough
	case "yaml":
		o, err := yaml.Marshal(plan)
		if err != nil {
			fmt.Printf("oops! cannot convert release plan into yaml: %v", err)
		}
		return string(o)
	case "json":
		o, err := json.MarshalIndent(plan, "", " ")
		if err != nil {
			fmt.Printf("oops! cannot convert release plan into json: %v", err)
		}
		return string(o)
	default:
		fmt.Printf("oops! output format %v not supported, try yaml or json", format)
	}
	return ""
}
