/*
   Onix Web API Client Library - Copyright (c) 2020 by www.gatblau.org

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
package wapic

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	Log    *log.Entry
	Token  string
	Config *Onix
	self   *http.Client
}

// creates a new Onix REST web client
func New(log *log.Entry, cfg *Onix) (*Client, error) {
	client := new(Client)
	client.Log = log
	client.Config = cfg
	err := client.setAuthenticationToken()
	if err != nil {
		return client, err
	}
	client.self = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	return client, err
}

// sets up the authentication Token used by the client
func (c *Client) setAuthenticationToken() error {
	var err error = nil
	switch c.Config.AuthMode {
	case "basic":
		c.Log.Tracef("Setting basic authentication token.")
		c.Token = NewBasicToken(c.Config.Username, c.Config.Password)
	case "oidc":
		c.Log.Tracef("Requesting bearer authentication token.")
		c.Token, err = NewBearerToken(c.Config.TokeURI, c.Config.ClientId, c.Config.ClientSecret, c.Config.Username, c.Config.Password)
		if err != nil {
			c.Log.Errorf("Failed to authenticate with OpenId server.", err)
		} else {
			c.Log.Tracef("Bearer token acquired.")
		}
	case "none":
		c.Log.Tracef("No authentication is used to connect to the Onix Config Manager.")
		c.Token = ""
	default:
		c.Log.Errorf("Cannot understand authentication mode selected: %s.", c.Config.AuthMode)
	}
	return err
}

// makes a generic HTTP request
func (c *Client) makeRequest(method string, resourceName string, key string, payload io.Reader) (*Result, error) {
	var (
		req *http.Request
		err error
	)
	// creates the request
	if len(key) > 0 {
		// with key
		req, err = http.NewRequest(method, fmt.Sprintf("%s/%s/%s", c.Config.URL, resourceName, key), payload)
	} else {
		// without key
		req, err = http.NewRequest(method, fmt.Sprintf("%s/%s", c.Config.URL, resourceName), payload)
	}
	// any errors are returned
	if err != nil {
		return &Result{Message: err.Error(), Error: true}, err
	}

	if method != "DELETE" {
		// requires a response in json format
		req.Header.Set("Content-Type", "application/json")
	}

	// if an authentication Token has been specified then add it to the request header
	if c.Token != "" && len(c.Token) > 0 {
		req.Header.Set("Authorization", c.Token)
	}

	// submits the request
	response, err := c.self.Do(req)

	// checks there are no errors in the response
	if response != nil && response.StatusCode > 300 {
		// if the response contains an error then returns
		return &Result{Message: response.Status, Error: true}, errors.New(response.Status)
	} else if err != nil {
		return &Result{Message: err.Error(), Error: true}, err
	}

	defer func() {
		if ferr := response.Body.Close(); ferr != nil {
			err = ferr
		}
	}()

	// decodes the response
	result := new(Result)
	err = json.NewDecoder(response.Body).Decode(result)

	// returns the result
	return result, err
}

// issues an GET HTTP request to the WAPI
func (c *Client) Get(resourceName string, key string, filter map[string]string) (interface{}, error) {
	var (
		req *http.Request
		err error
	)
	if len(key) > 0 {
		// if a resource key is passed, then query such resource
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", c.Config.URL, resourceName, key), nil)
	} else {
		// otherwise issue a find query with params (filters)
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/%s", c.Config.URL, resourceName), nil)
		// if there are query string params
		if filter != nil {
			// adds them to the request
			qParams := url.Values{}
			for k, v := range filter {
				qParams.Add(k, v)
			}
			req.URL.RawQuery = qParams.Encode()
		}
	}
	req.Header.Set("Content-Type", "application/json")
	// only add authorisation header if there is a token
	if len(c.Token) > 0 {
		req.Header.Set("Authorization", c.Token)
	}
	resp, err := c.self.Do(req)
	if resp != nil {
		defer func() {
			if ferr := resp.Body.Close(); ferr != nil {
				err = ferr
			}
		}()
	}
	if err != nil {
		return nil, err
	}
	// if the response status is OK (200)
	if resp.StatusCode == 200 {
		// if no key was passed-in, then assumes a query
		if len(key) == 0 {
			result := new(ItemList)
			err = json.NewDecoder(resp.Body).Decode(result)
			return result, err
		}

		// we have a key
		switch {
		case resourceName == "item":
			result := new(Item)
			err = json.NewDecoder(resp.Body).Decode(result)
			return *result, err
		case resourceName == "itemtype":
			result := new(ItemType)
			err = json.NewDecoder(resp.Body).Decode(result)
			return *result, err
		case resourceName == "link":
			result := new(Link)
			err = json.NewDecoder(resp.Body).Decode(result)
			return *result, err
		case resourceName == "linktype":
			result := new(LinkType)
			err = json.NewDecoder(resp.Body).Decode(result)
			return *result, err
		case resourceName == "model":
			result := new(Model)
			err = json.NewDecoder(resp.Body).Decode(result)
			return *result, err
		}
		// if the response status is something other than not found
	} else if resp.StatusCode != 404 {
		// return an error with the status message
		return nil, errors.New(resp.Status)
	}
	// the model was not found
	return nil, nil
}

// issues an http put request to the Onix Config Manager passing the specified item
// - payload: the payload object
// - resourceName: the WAPI resource name (e.g. item, itemtype, link, etc.)
// returns the payload key and a success flag
func (c *Client) Put(payload Payload, resourceName string) (string, *Result, error) {
	var (
		err    error
		result *Result
	)
	// converts the passed-in payload to a JSON bytes reader
	bytes, err := payload.ToJSON()

	if err != nil {
		c.Log.Errorf("Failed to marshall %s data: %s.", resourceName, err)
		return "", nil, err
	}
	// makes the http PUT request
	result, err = c.makeRequest(PUT, resourceName, payload.KeyValue(), bytes)
	if err != nil {
		c.Log.Errorf("Failed to PUT %s: %s.", resourceName, err)
		return "", nil, err
	}
	if result.Error {
		c.Log.Errorf("Failed to PUT %s: %s.", resourceName, result.Message)
		return "", result, err
	}
	if result.Changed {
		c.Log.Tracef("%s: %s update successful.", resourceName, payload.KeyValue())
		return payload.KeyValue(), result, err
	}
	c.Log.Tracef("%s: %s, Nothing to update.", resourceName, payload.KeyValue())
	return payload.KeyValue(), result, err
}
