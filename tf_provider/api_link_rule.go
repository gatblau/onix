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

type LinkRule struct {
	Id               string `json:"id"`
	Key              string `json:"key"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	LinkTypeKey      string `json:"linkTypeKey"`
	StartItemTypeKey string `json:"startItemTypeKey"`
	EndItemTypeKey   string `json:"endItemTypeKey"`
	Version          int64  `json:"version"`
	Created          string `json:"created"`
	Updated          string `json:"updated"`
}

func newLinkRule(data *schema.ResourceData) *LinkRule {
	return &LinkRule{
		Key:              data.Get("key").(string),
		Name:             data.Get("name").(string),
		Description:      data.Get("description").(string),
		LinkTypeKey:      data.Get("link_type_key").(string),
		StartItemTypeKey: data.Get("start_item_type_key").(string),
		EndItemTypeKey:   data.Get("end_item_type_key").(string),
	}
}

func (linkRule *LinkRule) toJSON() (*bytes.Reader, error) {
	return GetJSONBytesReader(linkRule)
}

// get the Link Rule in the http Response
func (linkRule *LinkRule) decode(response *http.Response) (*LinkRule, error) {
	result := new(LinkRule)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}

// populate the Link Rule with the data in the terraform resource
func (linkRule *LinkRule) populate(data *schema.ResourceData) {
	data.SetId(linkRule.Id)
	data.Set("key", linkRule.Key)
	data.Set("description", linkRule.Description)
	data.Set("link_type_key", linkRule.LinkTypeKey)
	data.Set("start_item_type_key", linkRule.StartItemTypeKey)
	data.Set("end_item_type_key", linkRule.EndItemTypeKey)
}

// get the FQN for the link rule resource
func (linkRule *LinkRule) uri(baseUrl string) string {
	return fmt.Sprintf("%s/linkrule/%s", baseUrl, linkRule.Key)
}

// issue a put http request with the Link rule data as payload to the resource URI
func (linkRule *LinkRule) put(meta interface{}) error {
	cfg := meta.(Config)

	// converts the passed-in payload to a bytes Reader
	bytes, err := linkRule.toJSON()

	// any errors are returned immediately
	if err != nil {
		return err
	}

	// make an http put request to the service
	result, err := cfg.Client.Put(linkRule.uri(cfg.Client.BaseURL), bytes)

	// any errors are returned
	if e := check(result, err); e != nil {
		return e
	}

	return nil
}

// issue a delete http request to the resource URI
func (linkRule *LinkRule) delete(meta interface{}) error {
	// get the Config instance from the meta object passed-in
	cfg := meta.(Config)

	// make an http delete request to the service
	result, err := cfg.Client.Delete(linkRule.uri(cfg.Client.BaseURL))

	// any errors are returned
	if e := check(result, err); e != nil {
		return e
	}

	return nil
}

// issue a get http request to the resource URI
func (linkRule *LinkRule) get(meta interface{}) (*LinkRule, error) {
	cfg := meta.(Config)

	// make an http put request to the service
	result, err := cfg.Client.Get(linkRule.uri(cfg.Client.BaseURL))

	if err != nil {
		return nil, err
	}

	return linkRule.decode(result)
}
