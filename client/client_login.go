package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

/*
   Onix Configuration Manager - HTTP Client
   Copyright (c) 2018-2021 by www.gatblau.org

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
   Unless required by applicable law or agreed to in writing, software distributed under
   the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
   either express or implied.
   See the License for the specific language governing permissions and limitations under the License.

   Contributors to this project, hereby assign copyright in this code to the project,
   to be licensed under the same terms as the rest of the code.
*/

// Login check that the user is authenticated using the CMDB as user store
// and returns a list of access controls for the user
func (c *Client) Login(credentials *Login) (*UserPrincipal, error) {
	// validates user
	if err := credentials.valid(); err != nil {
		return nil, err
	}
	uri, err := credentials.uri(c.conf.BaseURI)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(uri, credentials, c.addHttpHeaders)
	// if there is a technical error
	if err != nil {
		return nil, fmt.Errorf("login failed for user '%s' due to error: '%s'\n", credentials.Username, err)
	}
	// if the response was unauthorised, login failed
	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("authentication failed for user '%s'\n", credentials.Username)
	}
	// otherwise, get the list of controls from the user information
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("user was authenticated but failed to read response body for user '%s'; cannot return list of access controls; error: '%s'\n", credentials.Username, err)
	}
	var user User
	err = json.Unmarshal(bytes, &user)
	if err != nil {
		return nil, fmt.Errorf("user was authenticated but failed to unmarhsal response body for user '%s'; cannot return list of access controls; error: '%s'\n", credentials.Username, err)
	}
	controls, err := newControls(user.ACL)
	if err != nil {
		return nil, fmt.Errorf("user was authenticated but failed to parse access controls for user '%s': '%s'\n", credentials.Username, err)
	}
	// constructs a principal and returns
	return &UserPrincipal{
		Username: credentials.Username,
		Rights:   controls,
		Created:  time.Now(),
	}, nil
}

// Control represents a user access control based on a URI that is part of a realm
type Control struct {
	// the control realm (e.g. application)
	Realm string
	// the control resource URI (e.g. typically but not exclusively, a restful endpoint URI)
	URI string
	// the method(s) used by the resource (e.g. POST for create, PUT for update, DELETE for delete)
	Method []string
}

func (c *Control) hasMethod(method string) bool {
	for _, m := range c.Method {
		if strings.ToUpper(strings.Trim(m, " ")) == strings.ToUpper(strings.Trim(method, " ")) {
			return true
		}
	}
	return false
}

type Controls []Control

// Allowed returns true if the specified control matches one of the controls granted to the user
func (controls Controls) allowed(realm, uri, method string) bool {
	for _, c := range controls {
		if (c.Realm == realm || c.Realm == "*") &&
			(c.URI == uri || c.URI == "*") &&
			(c.hasMethod(method) || c.hasMethod("*")) {
			return true
		}
	}
	return false
}

// RequestAllowed returns true if the http request matches one of the controls granted to the user for the given realm
func (controls Controls) RequestAllowed(realm string, r *http.Request) bool {
	// extract user principal from the request context
	if principal := r.Context().Value("User"); principal != nil {
		if value, ok := principal.(*UserPrincipal); ok {
			return value.Rights.allowed(realm, r.RequestURI, r.Method)
		}
	}
	return false
}

func newControls(acl string) (Controls, error) {
	var result Controls
	// if acl is empty then return an empty list of controls
	if len(strings.Trim(acl, " ")) == 0 {
		return Controls{}, nil
	}
	parts := strings.Split(acl, ",")
	for _, part := range parts {
		control, err := newControl(part)
		if err != nil {
			return nil, err
		}
		result = append(result, control)
	}
	return result, nil
}

func newControl(ac string) (Control, error) {
	parts := strings.Split(ac, ":")
	if len(parts) != 3 {
		return Control{}, fmt.Errorf("Invalid control format '%s', it should be realm:uri:method\n", ac)
	}
	return Control{
		Realm:  parts[0],
		URI:    parts[1],
		Method: strings.Split(parts[2], "|"),
	}, nil
}

// UserPrincipal represents a logged on user and the access controls granted to them
type UserPrincipal struct {
	// the user Username used as a unique identifier (typically the user email address)
	Username string `json:"username"`
	// a list of rights or access controls granted to the user
	Rights Controls `json:"acl,omitempty"`
	// the time the principal was Created
	Created time.Time `json:"created"`
}
