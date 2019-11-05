/*
   Sentinel - Copyright (c) 2019 by www.gatblau.org

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
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
)

type WebhookPub struct {
	uri            []string
	authentication []string
	token          []string
	log            *logrus.Entry
	hooks          int
	client         []*http.Client
}

// gets the configuration for the publisher
func (pub *WebhookPub) Init(c *Config, log *logrus.Entry) {
	hooks := len(c.Publishers.Webhook)
	pub.log = log
	pub.uri = make([]string, hooks)
	pub.token = make([]string, hooks)
	pub.client = make([]*http.Client, hooks)

	// loads the configuration for the registered web hooks
	for i := 0; i < len(c.Publishers.Webhook); i++ {
		if contains(pub.uri, c.Publishers.Webhook[i].URI) {
			pub.log.Warnf("Webhook publisher contains two duplicate endpoint URIs: %s. \nDuplicate value will be omitted.", c.Publishers.Webhook[i].URI)
			pub.uri[i] = "" // set to empty to omit
		} else {
			pub.uri[i] = c.Publishers.Webhook[i].URI
		}
		if c.Publishers.Webhook[i].Authentication == "basic" {
			pub.token[i] = fmt.Sprintf("Basic %s",
				base64.StdEncoding.EncodeToString(
					[]byte(fmt.Sprintf("%s:%s",
						c.Publishers.Webhook[i].Username,
						c.Publishers.Webhook[i].Password))))
		} else {
			pub.token[i] = ""
		}
		pub.client[i] = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: c.Publishers.Webhook[i].InsecureSkipVerify,
				},
			},
		}
	}
}

// publishes events to the registered web hooks
func (pub *WebhookPub) Publish(event Event) {
	if event.Change.Type != "DELETE" && event.Object == nil {
		// if the CREATE OR UPDATE event metadata is nil, then is likely the resource is already gone
		// therefore it does not publish the event
		pub.log.Tracef("Event metadata not found when trying to post to '%s', skipping publication.", pub.uri)
		return
	}
	for i := 0; i < len(pub.uri); i++ {
		err := pub.post(pub.client[i], pub.uri[i], pub.token[i], event)
		if err != nil {
			pub.log.Errorf("Failed to post %s %s for %s: %s.",
				event.Change.Kind,
				event.Change.Type,
				event.Change.key,
				err)
		} else {
			pub.log.Tracef("%s %s for %s posted to webhook %s.",
				event.Change.Kind,
				event.Change.Type,
				event.Change.key,
				pub.uri[i])
		}
	}
}

// Make a POST to the webhook
func (pub *WebhookPub) post(client *http.Client, uri string, token string, object Event) error {
	// if the uri is empty omitting post
	if uri == "" {
		return errors.New("post to duplicate URI omitted, check Webhook configuration for duplicate URI values")
	}

	// gets a byte reader
	payload, err := getJSONBytesReader(object)

	if err != nil {
		pub.log.Errorf("Failed to create byte reader: %s", err)
	}

	// creates the request
	req, err := http.NewRequest("POST", uri, payload)

	// any errors are returned
	if err != nil {
		return err
	}

	// requires a response in json format
	req.Header.Set("Content-Type", "application/json")

	// if an authentication token has been specified then add it to the request header
	if token != "" && len(token) > 0 {
		req.Header.Set("Authorization", token)
	}

	// submits the request
	response, err := client.Do(req)

	// if the response contains an error then returns
	if err != nil {
		return err
	}
	if response == nil {
		return errors.New(fmt.Sprintf("No response for Request URI %s with payload %s", uri, object))
	}
	// if the response has an error then returns
	if response.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Request to URI %s failed with status: '%s'. Response body: '%s'. Request payload: '%s'", uri, response.Status, pub.toByteArray(response.Body), pub.toByteArray(req.Body)))
	} else {
		pub.log.Tracef("Payload '%s' posted to uri '%s' with status '%s'", payload, uri, response.Status)
	}
	defer func() {
		if ferr := response.Body.Close(); ferr != nil {
			err = ferr
		}
	}()

	// returns the result
	return err
}

// unmarshal the http response into a json like structure
func (pub *WebhookPub) toByteArray(r io.ReadCloser) []byte {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		pub.log.Warnf("Failed to read response or response body.")
	}
	return bytes
}
