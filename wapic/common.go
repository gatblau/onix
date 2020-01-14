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
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	DELETE = "DELETE"
	PUT    = "PUT"
	GET    = "GET"
	POST   = "POST"
)

// Response to an OAUth 2.0 token request
type OAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	IdToken     string `json:"id_token"`
}

// creates a new Basic Authentication Token
func NewBasicToken(user string, pwd string) string {
	return fmt.Sprintf("Basic %s",
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, pwd))))
}

// creates an OAuth Bearer token
func NewBearerToken(tokenURI string, clientId string, secret string, user string, pwd string) (string, error) {
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
	req.Header.Add("authorization", NewBasicToken(clientId, secret))    // authenticates with client id and secret
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

func getJSONBytesReader(data interface{}) (*bytes.Reader, error) {
	jsonBytes, err := json.Marshal(data)
	return bytes.NewReader(jsonBytes), err
}
