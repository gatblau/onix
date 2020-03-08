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

type LinkType struct {
	Id          string                 `json:"id"`
	Key         string                 `json:"key"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	MetaSchema  map[string]interface{} `json:"metaSchema"`
	Model       string                 `json:"modelKey"`
	Tag         []interface{}          `json:"tag"`
	EncryptMeta bool                   `json:"encryptMeta"`
	EncryptTxt  bool                   `json:"encryptTxt"`
	Managed     bool                   `json:"managed"`
	Version     int64                  `json:"version"`
	Created     string                 `json:"created"`
	Updated     string                 `json:"updated"`
}

func newLinkType(data *schema.ResourceData) *LinkType {
	return &LinkType{
		Key:         data.Get("key").(string),
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
		Model:       data.Get("model_key").(string),
		MetaSchema:  data.Get("meta_schema").(map[string]interface{}),
		EncryptMeta: data.Get("encrypt_meta").(bool),
		EncryptTxt:  data.Get("encrypt_txt").(bool),
		Managed:     data.Get("managed").(bool),
		Tag:         data.Get("tag").([]interface{}),
	}
}

func (linkType *LinkType) toJSON() (*bytes.Reader, error) {
	return GetJSONBytesReader(linkType)
}

// get the Link Type in the http Response
func (linkType *LinkType) decode(response *http.Response) (*LinkType, error) {
	result := new(LinkType)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}

// populate the LinkType with the data in the terraform resource
func (linkType *LinkType) populate(data *schema.ResourceData) {
	data.SetId(linkType.Id)
	data.Set("key", linkType.Key)
	data.Set("name", linkType.Name)
	data.Set("description", linkType.Description)
	data.Set("meta_schema", linkType.MetaSchema)
	data.Set("model", linkType.Model)
	data.Set("encrypt_txt", linkType.EncryptTxt)
	data.Set("encrypt_meta", linkType.EncryptMeta)
	data.Set("tag", linkType.Tag)
	data.Set("managed", linkType.Managed)
}

// get the FQN for the link type resource
func (linkType *LinkType) uri(baseUrl string) string {
	return fmt.Sprintf("%s/linktype/%s", baseUrl, linkType.Key)
}

// issue a put http request with the link type data as payload to the resource URI
func (linkType *LinkType) put(meta interface{}) error {
	cfg := meta.(Config)
	// converts the passed-in payload to a bytes Reader
	bytes, err := linkType.toJSON()

	// any errors are returned immediately
	if err != nil {
		return err
	}

	// make an http put request to the service
	result, err := cfg.Client.Put(linkType.uri(cfg.Client.BaseURL), bytes)

	// any errors are returned
	if e := check(result, err); e != nil {
		return e
	}

	return nil
}

// issue a delete http request to the resource URI
func (linkType *LinkType) delete(meta interface{}) error {
	// get the Config instance from the meta object passed-in
	cfg := meta.(Config)

	// make an http delete request to the service
	result, err := cfg.Client.Delete(linkType.uri(cfg.Client.BaseURL))

	// any errors are returned
	if e := check(result, err); e != nil {
		return e
	}

	return nil
}

// issue a get http request to the resource URI
func (linkType *LinkType) get(meta interface{}) (*LinkType, error) {
	cfg := meta.(Config)

	// make an http put request to the service
	result, err := cfg.Client.Get(linkType.uri(cfg.Client.BaseURL))

	if err != nil {
		return nil, err
	}

	return linkType.decode(result)
}
