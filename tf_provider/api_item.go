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

// the Item resource
type Item struct {
	Id          string                 `json:"id"`
	Key         string                 `json:"key"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Status      int                    `json:"status"`
	Type        string                 `json:"type"`
	Tag         []interface{}          `json:"tag"`
	Meta        map[string]interface{} `json:"meta"`
	Txt         string                 `json:"txt"`
	Attribute   map[string]interface{} `json:"attribute"`
	Partition   string                 `json:"partition"`
	Version     int64                  `json:"version"`
	Created     string                 `json:"created"`
	Updated     string                 `json:"updated"`
	EncKeyIx    int64                  `json:"encKeyIx"`
}

// create a new Item from a terraform resource
func newItem(data *schema.ResourceData) *Item {
	return &Item{
		Key:         data.Get("key").(string),
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
		Type:        data.Get("type").(string),
		Meta:        data.Get("meta").(map[string]interface{}),
		Txt:         data.Get("txt").(string),
		Attribute:   data.Get("attribute").(map[string]interface{}),
		Tag:         data.Get("tag").([]interface{}),
		Partition:   data.Get("partition").(string),
	}
}

// get a JSON bytes reader for the item
func (item *Item) toJSON() (*bytes.Reader, error) {
	return GetJSONBytesReader(item)
}

// get the Item in the http Response
func (item *Item) decode(response *http.Response) (*Item, error) {
	result := new(Item)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}

// populate the Item with the data in the terraform resource
func (item *Item) populate(data *schema.ResourceData) {
	data.SetId(item.Id)
	data.Set("key", item.Key)
	data.Set("name", item.Name)
	data.Set("description", item.Description)
	data.Set("attribute", item.Attribute)
	data.Set("txt", item.Txt)
	data.Set("meta", item.Meta)
	data.Set("partition", item.Partition)
	data.Set("status", item.Status)
	data.Set("tag", item.Tag)
	data.Set("type", item.Type)
}

// get the FQN for the item resource
func (item *Item) uri(baseUrl string) string {
	return fmt.Sprintf("%s/item/%s", baseUrl, item.Key)
}

// issue a put http request with the Item data as payload to the resource URI
func (item *Item) put(meta interface{}) error {
	cfg := meta.(Config)
	// converts the passed-in payload to a bytes Reader
	bytes, err := item.toJSON()

	// any errors are returned immediately
	if err != nil {
		return err
	}

	// make an http put request to the service
	result, err := cfg.Client.Put(item.uri(cfg.Client.BaseURL), bytes)

	// any errors are returned
	if e := check(result, err); e != nil {
		return e
	}

	return nil
}

// issue a delete http request to the resource URI
func (item *Item) delete(meta interface{}) error {
	// get the Config instance from the meta object passed-in
	cfg := meta.(Config)

	// make an http delete request to the service
	result, err := cfg.Client.Delete(item.uri(cfg.Client.BaseURL))

	// any errors are returned
	if e := check(result, err); e != nil {
		return e
	}

	return nil
}

// issue a get http request to the resource URI
func (item *Item) get(meta interface{}) (*Item, error) {
	cfg := meta.(Config)

	// make an http put request to the service
	result, err := cfg.Client.Get(item.uri(cfg.Client.BaseURL))

	if err != nil {
		return nil, err
	}

	i, err := item.decode(result)

	defer func() {
		if ferr := result.Body.Close(); ferr != nil {
			err = ferr
		}
	}()

	return i, err
}
