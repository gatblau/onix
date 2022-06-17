package core

/*
  Onix Config Manager - Pilot Control
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
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
    http "net/http"
)

const (
    DELETE = "DELETE"
    PUT    = "PUT"
    GET    = "GET"
    POST   = "POST"
)

// Serializable all entities interface for payload serialisation
type Serializable interface {
    Reader() (*bytes.Reader, error)
    Bytes() (*[]byte, error)
}

// HttpRequestProcessor modify the http request for example by adding any relevant http headers
// payload is provided for example, in case a Content-MD5 header has to be added to the request
type HttpRequestProcessor func(req *http.Request, payload Serializable) error

// Client REM HTTP client
type Client struct {
    conf  *ClientConf
    self  *http.Client
    token string
}

// Result data retrieved by PUT and DELETE WAPI resources
type Result struct {
    Error   bool   `json:"error"`
    Message string `json:"message"`
}

// creates a new result from an http response
func newResult(response *http.Response) (*Result, error) {
    result := new(Result)
    // de-serialise the response
    err := json.NewDecoder(response.Body).Decode(result)
    // if the de-serialisation fails then return the error
    if err != nil {
        return result, err
    }
    // close the response
    defer func() {
        if ferr := response.Body.Close(); ferr != nil {
            err = ferr
        }
    }()
    // returns the result
    return result, nil
}

// NewClient creates a new Onix Web API client
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
                // if proxy == nil, no proxy is set
                Proxy: conf.Proxy,
            },
            // set the client timeout period
            Timeout: conf.Timeout,
        },
    }
    return client, err
}

// MakeRequest Make a generic HTTP request
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
    resp, err := c.self.Do(req)

    // do we have a nil response?
    if resp == nil {
        return resp, errors.New(fmt.Sprintf("error: response was empty for resource: %s: %s", url, err))
    }
    // check for response status
    if resp.StatusCode >= 300 {
        err = errors.New(fmt.Sprintf("error: response returned status: %s", resp.Status))
    }
    return resp, err
}

// Put Make a PUT HTTP request to the specified URL
func (c *Client) Put(url string, payload Serializable, processor HttpRequestProcessor) (*http.Response, error) {
    return c.MakeRequest(PUT, url, payload, processor)
}

// Post Make a POST HTTP request to the specified URL
func (c *Client) Post(url string, payload Serializable, processor HttpRequestProcessor) (*http.Response, error) {
    return c.MakeRequest(POST, url, payload, processor)
}

// Delete Make a DELETE HTTP request to the specified URL
func (c *Client) Delete(url string, processor HttpRequestProcessor) (*http.Response, error) {
    return c.MakeRequest(DELETE, url, nil, processor)
}

// Get Make a GET HTTP request to the specified URL
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
    if err != nil {
        return nil, err
    }
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
        data, err := payload.Bytes()
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
    return payload.Reader()
}
