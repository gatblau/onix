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

// the Role resource
type Role struct {
	Id          string `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Owner       string `json:"owner"`
	Level       int    `json:"level"`
	Version     int64  `json:"version"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
}

// create a new Role from a terraform resource
func newRole(data *schema.ResourceData) *Role {
	return &Role{
		Key:         data.Get("key").(string),
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
		Level:       data.Get("level").(int),
	}
}

// get a JSON bytes reader for the Role
func (role *Role) toJSON() (*bytes.Reader, error) {
	return GetJSONBytesReader(role)
}

// get the Role in the http Response
func (role *Role) decode(response *http.Response) (*Role, error) {
	result := new(Role)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}

// populate the Role with the data in the terraform resource
func (role *Role) populate(data *schema.ResourceData) {
	data.SetId(role.Id)
	data.Set("key", role.Key)
	data.Set("name", role.Name)
	data.Set("description", role.Description)
	data.Set("level", role.Level)
}

// get the FQN for the item resource
func (role *Role) uri(baseUrl string) string {
	return fmt.Sprintf("%s/role/%s", baseUrl, role.Key)
}

// issue a put http request with the Role data as payload to the resource URI
func (role *Role) put(meta interface{}) error {
	cfg := meta.(Config)
	// converts the passed-in payload to a bytes Reader
	bytes, err := role.toJSON()

	// any errors are returned immediately
	if err != nil {
		return err
	}

	// make an http put request to the service
	result, err := cfg.Client.Put(role.uri(cfg.Client.BaseURL), bytes)

	// any errors are returned
	if e := check(result, err); e != nil {
		return e
	}

	return nil
}

// issue a delete http request to the resource URI
func (role *Role) delete(meta interface{}) error {
	// get the Config instance from the meta object passed-in
	cfg := meta.(Config)

	// make an http delete request to the service
	result, err := cfg.Client.Delete(role.uri(cfg.Client.BaseURL))

	// any errors are returned
	if e := check(result, err); e != nil {
		return e
	}

	return nil
}

// issue a get http request to the resource URI
func (role *Role) get(meta interface{}) (*Role, error) {
	cfg := meta.(Config)

	// make an http put request to the service
	result, err := cfg.Client.Get(role.uri(cfg.Client.BaseURL))

	if err != nil {
		return nil, err
	}

	i, err := role.decode(result)

	defer func() {
		if ferr := result.Body.Close(); ferr != nil {
			err = ferr
		}
	}()

	return i, err
}
