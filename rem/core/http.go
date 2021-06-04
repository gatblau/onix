package core

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

/*
  Onix Config Manager - Host Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

func makeRequest(uri, method, user, pwd string, body io.Reader) ([]byte, error) {
	// create an http client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		// set the client timeout period
		Timeout: 1 * time.Minute,
	}
	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", basicAuthToken(user, pwd))
	// submits the request
	resp, err := client.Do(req)
	// do we have a nil response?
	if resp == nil {
		return nil, errors.New(fmt.Sprintf("error: response was empty for resource: %s, check the service is up and running", uri))
	}
	// check for response status
	if resp.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("error: response returned status: %s", resp.Status))
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
