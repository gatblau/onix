/*
   Onix Config Manager - Copyright (c) 2018-2019 by www.gatblau.org

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
package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	DELETE = "DELETE"
	PUT    = "PUT"
	GET    = "GET"
	POST   = "POST"
)

// Onix HTTP client
type Client struct {
	BaseURL string
	Token   string
}

// Result data retrieved by PUT and DELETE WAPI resources
type Result struct {
	Changed   bool   `json:"changed"`
	Error     bool   `json:"error"`
	Message   string `json:"message"`
	Operation string `json:"operation"`
	Ref       string `json:"ref"`
}

// Response to an OAUth 2.0 token request
type OAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	IdToken     string `json:"id_token"`
}

// creates a new Basic Authentication Token
func (o *Client) newBasicToken(user string, pwd string) string {
	return fmt.Sprintf("Basic %s",
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, pwd))))
}

// Set up the authentication token used by the client
func (o *Client) setAuthToken(token string) {
	o.Token = token
}

// Gets an OAuth Bearer token
func (o *Client) getBearerToken(tokenURI string, clientId string, secret string, user string, pwd string) (string, error) {
	// constructs a payload for the form POST to the authorisation server token URI
	// passing the type of grant,the username, password and scopes
	payload := strings.NewReader(
		fmt.Sprintf("grant_type=password&username=%s&password=%s&scope=openid%%20onix", user, pwd))

	// creates the http request
	req, err := http.NewRequest(POST, tokenURI, payload)

	// if any errors then return
	if err != nil {
		return "", errors.New("Failed to create request: " + err.Error())
	}

	// adds the relevant http headers
	req.Header.Add("accept", "application/json")                        // need a response in json format
	req.Header.Add("authorization", o.newBasicToken(clientId, secret))  // authenticates with client id and secret
	req.Header.Add("cache-control", "no-cache")                         // forces caches to submit the request to the origin server for validation before releasing a cached copy
	req.Header.Add("content-type", "application/x-www-form-urlencoded") // posting an http form

	// submits the request to the authorisation server
	response, err := http.DefaultClient.Do(req)

	// if any errors then return
	if err != nil {
		return "", errors.New("Failed when submitting request: " + err.Error())
	}
	if response.StatusCode != 200 {
		return "", errors.New("Failed to obtain access token: " + response.Status + " Hint: the client might be unauthorised.")
	}

	defer func() {
		if ferr := response.Body.Close(); ferr != nil {
			err = ferr
		}
	}()

	result := new(OAuthTokenResponse)

	// decodes the response
	err = json.NewDecoder(response.Body).Decode(result)

	// if any errors then return
	if err != nil {
		return "", err
	}

	// constructs and returns a bearer token
	return fmt.Sprintf("Bearer %s", result.AccessToken), nil
}

// Make a generic HTTP request
func (o *Client) MakeRequest(method string, url string, payload io.Reader) (*Result, error) {
	// creates the request
	req, err := http.NewRequest(method, url, payload)

	// any errors are returned
	if err != nil {
		return &Result{Message: err.Error(), Error: true}, err
	}

	// requires a response in json format
	req.Header.Set("Content-Type", "application/json")

	// if an authentication token has been specified then add it to the request header
	if o.Token != "" && len(o.Token) > 0 {
		req.Header.Set("Authorization", o.Token)
	}

	// submits the request
	response, err := http.DefaultClient.Do(req)

	// if the response contains an error then returns
	if err != nil {
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

// Make a PUT HTTP request to the WAPI
func (o *Client) Put(url string, payload io.Reader) (*Result, error) {
	return o.MakeRequest(PUT, url, payload)
}

// Make a DELETE HTTP request to the WAPI
func (o *Client) Delete(url string) (*Result, error) {
	return o.MakeRequest(DELETE, url, nil)
}

// Make a GET HTTP request to the WAPI
func (o *Client) Get(resourceName string, key string) (interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", o.BaseURL, resourceName, key), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", o.Token)
	resp, err := http.DefaultClient.Do(req)
	defer func() {
		if ferr := resp.Body.Close(); ferr != nil {
			err = ferr
		}
	}()
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
	return nil, nil
}
