/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
)

// Make a GET HTTP request to the specified URL
func Get(url, user, pwd string) (*http.Response, error) {
	// create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// add http request headers
	if len(user) > 0 && len(pwd) > 0 {
		req.Header.Add("Authorization", basicToken(user, pwd))
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

func basicToken(user string, pwd string) string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, pwd))))
}
