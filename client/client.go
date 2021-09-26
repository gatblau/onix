package client

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
import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	DELETE = "DELETE"
	PUT    = "PUT"
	GET    = "GET"
	POST   = "POST"
)

// all entities interface for payload serialisation
type Serializable interface {
	// the reader for the payload format
	reader() (*bytes.Reader, error)
	bytes() (*[]byte, error)
}

// a payload for an http operation in
type StringPayload string

func (s StringPayload) reader() (*bytes.Reader, error) {
	return bytes.NewReader([]byte(s)), nil
}

func (s StringPayload) bytes() (*[]byte, error) {
	b := []byte(s)
	return &b, nil
}

// modify the http request for example by adding any relevant http headers
// payload is provided for example, in case a Content-MD5 header has to be added to the request
type HttpRequestProcessor func(req *http.Request, payload Serializable) error

// Onix HTTP client
type Client struct {
	conf  *ClientConf
	self  *http.Client
	token string
}

// Result data retrieved by PUT and DELETE WAPI resources
type Result struct {
	Changed   bool   `json:"changed"`
	Error     bool   `json:"error"`
	Message   string `json:"message"`
	Operation string `json:"operation"`
	Ref       string `json:"ref"`
}

// creates a new result from an http response
func newResult(response *http.Response) (*Result, error) {
	result := new(Result)
	// attempt to de-serialise the response
	err := json.NewDecoder(response.Body).Decode(result)
	// close the response
	defer func() {
		if ferr := response.Body.Close(); ferr != nil {
			err = ferr
		}
	}()
	// if the de-serialisation fails, it could be that the server is returning an error code
	if err != nil {
		// if there is an error code
		if response.StatusCode > 299 {
			// added to the result
			result.Error = true
			result.Message = response.Status
			err = nil
		}
		return result, err
	}
	// else returns the result sent by the server
	return result, nil
}

// Response to an OAUth 2.0 token request
type OAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	IdToken     string `json:"id_token"`
}

// creates a new Onix Web API client
func NewClient(conf *ClientConf) (*Client, error) {
	// checks the passed-in configuration is correct
	err := checkConf(conf)
	if err != nil {
		return nil, err
	}

	// obtains an authentication token for the client
	token, err := conf.getAuthToken()
	if err != nil {
		return nil, err
	}

	// gets an instance of the client
	client := &Client{
		// the configuration information
		conf: conf,
		// the authentication token
		token: token,
		// the http client instance
		self: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: conf.InsecureSkipVerify,
				},
			},
			// set the client timeout period
			Timeout: conf.Timeout,
		},
	}
	return client, err
}

// Make a generic HTTP request
func (c *Client) MakeRequest(method string, url string, payload Serializable, processor HttpRequestProcessor) (*http.Response, error) {
	// prepares the request body, if no body exists, a nil reader is retrieved
	reader, err := c.getRequestBody(payload)
	if err != nil {
		return nil, err
	}

	// creates the request
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}

	// add the http headers to the request
	if processor != nil {
		err = processor(req, payload)
		if err != nil {
			return nil, err
		}
	}

	// submits the request
	resp, err := http.DefaultClient.Do(req)

	// do we have a nil response?
	if resp == nil {
		return resp, errors.New(fmt.Sprintf("error: response was empty for resource: %s, check the service is up and running", url))
	}
	// check for response status
	if resp.StatusCode >= 300 {
		err = errors.New(fmt.Sprintf("error: response returned status: %s", resp.Status))
	}
	return resp, err
}

// Make a PUT HTTP request to the specified URL
func (c *Client) Put(url string, payload Serializable, processor HttpRequestProcessor) (*http.Response, error) {
	return c.MakeRequest(PUT, url, payload, processor)
}

// Make a POST HTTP request to the specified URL
func (c *Client) Post(url string, payload Serializable, processor HttpRequestProcessor) (*http.Response, error) {
	return c.MakeRequest(POST, url, payload, processor)
}

// Make a DELETE HTTP request to the specified URL
func (c *Client) Delete(url string, processor HttpRequestProcessor) (*http.Response, error) {
	return c.MakeRequest(DELETE, url, nil, processor)
}

// Make a GET HTTP request to the specified URL
func (c *Client) Get(url string, processor HttpRequestProcessor) (*http.Response, error) {
	// create request
	req, err := http.NewRequest(GET, url, nil)
	if err != nil {
		return nil, err
	}
	// add http request headers
	if processor != nil {
		err = processor(req, nil)
		if err != nil {
			return nil, err
		}
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
func (c *Client) addHttpHeaders(req *http.Request, payload Serializable) error {
	// add authorization header if there is a token defined
	if len(c.token) > 0 {
		req.Header.Set("Authorization", c.token)
	}
	// all content type should be in JSON format
	req.Header.Set("Content-Type", "application/json")
	// if there is a payload
	if payload != nil {
		// Get the bytes in the Serializable
		data, err := payload.bytes()
		if err != nil {
			return err
		}
		// set the length of the payload
		req.ContentLength = int64(len(*data))
		// generate checksum of the payload data using the MD5 hashing algorithm
		checksum := md5.Sum(*data)
		// base 64 encode the checksum
		b64checksum := base64.StdEncoding.EncodeToString(checksum[:])
		// add Content-MD5 header (see https://tools.ietf.org/html/rfc1864)
		req.Header.Set("Content-MD5", b64checksum)
	}
	return nil
}

func (c *Client) getRequestBody(payload Serializable) (*bytes.Reader, error) {
	// if no payload exists
	if payload == nil {
		// returns an empty reader
		return bytes.NewReader([]byte{}), nil
	}
	// gets a byte reader to pass to the request body
	return payload.reader()
}

// convert the passed-in object to a JSON byte slice
// NOTE: json.Marshal is purposely not used as it will escape any < > characters
func ToJson(object interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	// switch off the escaping!
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(object)
	return buffer.Bytes(), err
}
