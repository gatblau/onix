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

// the Item Type resource
type ItemType struct {
	Id           string                 `json:"id"`
	Key          string                 `json:"key"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Filter       map[string]interface{} `json:"filter"`
	MetaSchema   map[string]interface{} `json:"metaSchema"`
	Model        string                 `json:"modelKey"`
	NotifyChange bool                   `json:"notifyChange"`
	Tag          []interface{}          `json:"tag"`
	EncryptMeta  bool                   `json:"encryptMeta"`
	EncryptTxt   bool                   `json:"encryptTxt"`
	Managed      bool                   `json:"managed"`
	Version      int64                  `json:"version"`
	Created      string                 `json:"created"`
	Updated      string                 `json:"updated"`
}

// create a new Item from a terraform resource
func newItemType(data *schema.ResourceData) *ItemType {
	return &ItemType{
		Key:          data.Get("key").(string),
		Name:         data.Get("name").(string),
		Description:  data.Get("description").(string),
		Model:        data.Get("model_key").(string),
		Filter:       data.Get("filter").(map[string]interface{}),
		MetaSchema:   data.Get("meta_schema").(map[string]interface{}),
		NotifyChange: data.Get("notify_change").(bool),
		EncryptMeta:  data.Get("encrypt_meta").(bool),
		EncryptTxt:   data.Get("encrypt_txt").(bool),
		Managed:      data.Get("managed").(bool),
		Tag:          data.Get("tag").([]interface{}),
	}
}

// create a new Item from a terraform resource
func (itemType *ItemType) toJSON() (*bytes.Reader, error) {
	return GetJSONBytesReader(itemType)
}

// get the Item Type in the http Response
func (itemType *ItemType) decode(response *http.Response) (*ItemType, error) {
	result := new(ItemType)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}

// populate the Item with the data in the terraform resource
func (itemType *ItemType) populate(data *schema.ResourceData) {
	data.SetId(itemType.Id)
	data.Set("key", itemType.Key)
	data.Set("name", itemType.Name)
	data.Set("description", itemType.Description)
	data.Set("filter", itemType.Filter)
	data.Set("meta_schema", itemType.MetaSchema)
	data.Set("notify_change", itemType.NotifyChange)
	data.Set("tag", itemType.Tag)
	data.Set("encrypt_meta", itemType.EncryptMeta)
	data.Set("encrypt_txt", itemType.EncryptTxt)
	data.Set("managed", itemType.Managed)
	data.Set("model", itemType.Model)
	data.Set("version", itemType.Version)
	data.Set("created", itemType.Created)
	data.Set("updated", itemType.Updated)
}

// get the FQN for the item resource
func (itemType *ItemType) uri(baseUrl string) string {
	return fmt.Sprintf("%s/itemtype/%s", baseUrl, itemType.Key)
}

// issue a put http request with the Item Type data as payload to the resource URI
func (itemType *ItemType) put(meta interface{}) error {
	cfg := meta.(Config)

	// converts the passed-in payload to a bytes Reader
	bytes, err := itemType.toJSON()

	// any errors are returned immediately
	if err != nil {
		return err
	}

	// make an http put request to the service
	result, err := cfg.Client.Put(itemType.uri(cfg.Client.BaseURL), bytes)

	// any errors are returned
	if e := check(result, err); e != nil {
		return e
	}

	return nil
}

// issue a delete http request to the resource URI
func (itemType *ItemType) delete(meta interface{}) error {
	// get the Config instance from the meta object passed-in
	cfg := meta.(Config)

	// make an http delete request to the service
	result, err := cfg.Client.Delete(itemType.uri(cfg.Client.BaseURL))

	// any errors are returned
	if e := check(result, err); e != nil {
		return e
	}

	return nil
}

// issue a get http request to the resource URI
func (itemType *ItemType) get(meta interface{}) (*ItemType, error) {
	cfg := meta.(Config)

	// make an http put request to the service
	result, err := cfg.Client.Get(itemType.uri(cfg.Client.BaseURL))

	if err != nil {
		return nil, err
	}

	it, err := itemType.decode(result)

	defer func() {
		if ferr := result.Body.Close(); ferr != nil {
			err = ferr
		}
	}()

	return it, err
}
