/*
  Onix Config Manager - Artisan's Doorman Proxy
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package main

import (
	"fmt"
	util "github.com/gatblau/onix/oxlib/httpserver"
	"net/http"
)

func newRequest(method, requestURI string) (*http.Response, error, int) {
	req, err := http.NewRequest(method, requestURI, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create http request: %s", err), http.StatusInternalServerError
	}
	user, err := getDoormanUser()
	if err != nil {
		return nil, fmt.Errorf("missing configuration"), http.StatusInternalServerError
	}
	pwd, err := getDoormanPwd()
	if err != nil {
		return nil, fmt.Errorf("missing configuration"), http.StatusInternalServerError
	}
	req.Header.Add("Authorization", util.BasicToken(user, pwd))
	resp, err := http.DefaultClient.Do(req)
	// do we have a nil response?
	if resp == nil {
		return nil, fmt.Errorf("response was empty for resource: %s", requestURI), http.StatusBadGateway
	}
	// check error status codes
	if resp.StatusCode > 201 {
		return resp, fmt.Errorf("response returned status: %s; resource: %s", resp.Status, requestURI), http.StatusBadGateway
	}
	return resp, nil, -1
}
