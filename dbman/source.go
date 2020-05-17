//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// all entities interface for payload serialisation
type entity interface {
	json() (*bytes.Reader, error)
	bytes() (*[]byte, error)
}

// Scripts HTTP client
type Source struct {
	cfg *Config
}

const (
	GET = "GET"
)

func NewSource(cfg *Config) *Source {
	return &Source{
		cfg: cfg,
	}
}

// Make a GET HTTP request to the WAPI
func (s *Source) get(url string) (*http.Response, error) {
	// create request
	req, err := http.NewRequest(GET, url, nil)
	if err != nil {
		return nil, err
	}
	// add http headers
	err = s.addHttpHeaders(req, nil)
	if err != nil {
		return nil, err
	}
	// issue http request
	resp, err := http.DefaultClient.Do(req)
	// do we have a nil response?
	if resp == nil {
		return resp, errors.New(fmt.Sprintf("error: response was empty for resource: %s", url))
	}
	// check error status codes
	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("error: response returned status: %s. resource: %s", resp.Status, url))
	}
	return resp, err
}

// add http headers to the request object
func (s *Source) addHttpHeaders(req *http.Request, payload entity) error {
	// add headers to disable caching
	req.Header.Add("Cache-Control", `no-cache"`)
	req.Header.Add("Pragma", "no-cache")
	// if there is an access token defined
	if len(s.cfg.SchemaUsername) > 0 && len(s.cfg.SchemaToken) > 0 {
		credentials := base64.StdEncoding.EncodeToString([]byte(
			fmt.Sprintf("%s:%s", s.cfg.SchemaUsername, s.cfg.SchemaToken)))
		req.Header.Add("Authorization", credentials)
	}
	return nil
}

func (s *Source) getRequestBody(payload entity) (*bytes.Reader, error) {
	// if no payload exists
	if payload == nil {
		// returns an empty reader
		return bytes.NewReader([]byte{}), nil
	}
	// gets a byte reader to pass to the request body
	return payload.json()
}

// convert the passed-in object to a JSON byte slice
// NOTE: json.Marshal is purposely not used as it will escape any < > characters
func jsonBytes(object interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	// switch off the escaping!
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(object)
	return buffer.Bytes(), err
}
