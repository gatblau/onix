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

// the Model resource
type Model struct {
	Id          string `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Partition   string `json:"partition"`
	Managed     bool   `json:"managed"`
	Version     int64  `json:"version"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
}

// create a new Model from a terraform resource
func newModel(data *schema.ResourceData) *Model {
	return &Model{
		Key:         data.Get("key").(string),
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
		Partition:   data.Get("partition").(string),
		Managed:     data.Get("managed").(bool),
	}
}

// create a new Model from a terraform resource
func (model *Model) toJSON() (*bytes.Reader, error) {
	return GetJSONBytesReader(model)
}

// get the Model in the http Response
func (model *Model) decode(response *http.Response) (*Model, error) {
	result := new(Model)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}

// populate the Model with the data in the terraform resource
func (model *Model) populate(data *schema.ResourceData) {
	data.SetId(model.Id)
	data.Set("key", model.Key)
	data.Set("name", model.Name)
	data.Set("description", model.Description)
	data.Set("partition", model.Partition)
	data.Set("managed", model.Managed)
	data.Set("version", model.Version)
	data.Set("created", model.Created)
	data.Set("updated", model.Updated)
}

// get the FQN for the model resource
func (model *Model) uri(baseUrl string) string {
	return fmt.Sprintf("%s/model/%s", baseUrl, model.Key)
}

// issue a put http request with the Model data as payload to the resource URI
func (model *Model) put(meta interface{}) error {
	cfg := meta.(Config)

	// converts the passed-in payload to a bytes Reader
	bytes, err := model.toJSON()

	// any errors are returned immediately
	if err != nil {
		return err
	}

	// make an http put request to the service
	result, err := cfg.Client.Put(model.uri(cfg.Client.BaseURL), bytes)

	// any errors are returned
	if e := check(result, err); e != nil {
		return e
	}

	return nil
}

// issue a delete http request to the resource URI
func (model *Model) delete(meta interface{}) error {
	// get the Config instance from the meta object passed-in
	cfg := meta.(Config)

	// make an http delete request to the service
	result, err := cfg.Client.Delete(model.uri(cfg.Client.BaseURL))

	// any errors are returned
	if e := check(result, err); e != nil {
		return e
	}

	return nil
}

// issue a get http request to the resource URI
func (model *Model) get(meta interface{}) (*Model, error) {
	cfg := meta.(Config)

	// make an http put request to the service
	result, err := cfg.Client.Get(model.uri(cfg.Client.BaseURL))

	if err != nil {
		return nil, err
	}

	m, err := model.decode(result)

	defer func() {
		if ferr := result.Body.Close(); ferr != nil {
			err = ferr
		}
	}()

	return m, err
}
