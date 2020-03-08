/*
   Onix Config Manager - Terraform Provider
   Copyright (c) 2018-2020 by www.gatblau.org

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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"net/http"
)

type Link struct {
	Id           string                 `json:"id"`
	Key          string                 `json:"key"`
	Description  string                 `json:"description"`
	Type         string                 `json:"type"`
	Tag          []interface{}          `json:"tag"`
	Meta         map[string]interface{} `json:"meta"`
	Attribute    map[string]interface{} `json:"attribute"`
	StartItemKey string                 `json:"startItemKey"`
	EndItemKey   string                 `json:"endItemKey"`
	Version      int64                  `json:"version"`
	Created      string                 `json:"created"`
	Updated      string                 `json:"updated"`
}

func (link *Link) toJSON() (*bytes.Reader, error) {
	return GetJSONBytesReader(link)
}

func newLink(data *schema.ResourceData) *Link {
	return &Link{
		Key:          data.Get("key").(string),
		Description:  data.Get("description").(string),
		Type:         data.Get("type").(string),
		Meta:         data.Get("meta").(map[string]interface{}),
		Attribute:    data.Get("attribute").(map[string]interface{}),
		Tag:          data.Get("tag").([]interface{}),
		StartItemKey: data.Get("start_item_key").(string),
		EndItemKey:   data.Get("end_item_key").(string),
	}
}

// get the Link in the http Response
func (link *Link) decode(response *http.Response) (*Link, error) {
	result := new(Link)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}

// populate the Link with the data in the terraform resource
func (link *Link) populate(data *schema.ResourceData) {
	data.SetId(link.Id)
	data.Set("key", link.Key)
	data.Set("description", link.Description)
	data.Set("type", link.Type)
	data.Set("meta", link.Meta)
	data.Set("tag", link.Tag)
	data.Set("attribute", link.Attribute)
	data.Set("start_item_key", link.StartItemKey)
	data.Set("end_item_key", link.EndItemKey)
}

// get the FQN for the item resource
func (link *Link) uri(baseUrl string) string {
	return fmt.Sprintf("%s/link/%s", baseUrl, link.Key)
}

// issue a put http request with the Link data as payload to the resource URI
func (link *Link) put(meta interface{}) error {
	cfg := meta.(Config)

	// converts the passed-in payload to a bytes Reader
	bytes, err := link.toJSON()

	// any errors are returned immediately
	if err != nil {
		return err
	}

	// make an http put request to the service
	result, err := cfg.Client.Put(link.uri(cfg.Client.BaseURL), bytes)

	// any errors are returned
	if e := check(result, err); e != nil {
		return e
	}

	return nil
}

// issue a delete http request to the resource URI
func (link *Link) delete(meta interface{}) error {
	// get the Config instance from the meta object passed-in
	cfg := meta.(Config)

	// make an http delete request to the service
	result, err := cfg.Client.Delete(link.uri(cfg.Client.BaseURL))

	// any errors are returned
	if e := check(result, err); e != nil {
		return e
	}

	return nil
}

// issue a get http request to the resource URI
func (link *Link) get(meta interface{}) (*Link, error) {
	cfg := meta.(Config)

	// make an http put request to the service
	result, err := cfg.Client.Get(link.uri(cfg.Client.BaseURL))

	if err != nil {
		return nil, err
	}

	return link.decode(result)
}
