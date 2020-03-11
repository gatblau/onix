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

// the Privilege resource
type Privilege struct {
	Id        string `json:"id"`
	Key       string `json:"key"`
	Role      string `json:"roleKey"`
	Partition string `json:"partitionKey"`
	Create    bool   `json:"canCreate"`
	Read      bool   `json:"canRead"`
	Delete    bool   `json:"canDelete"`
	Version   int64  `json:"version"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
}

// create a new Privilege from a terraform resource
func newPrivilege(data *schema.ResourceData) *Privilege {
	return &Privilege{
		Key:       data.Get("key").(string),
		Role:      data.Get("role").(string),
		Partition: data.Get("partition").(string),
		Create:    data.Get("can_create").(bool),
		Read:      data.Get("can_read").(bool),
		Delete:    data.Get("can_delete").(bool),
	}
}

// get a JSON bytes reader for the item
func (privilege *Privilege) toJSON() (*bytes.Reader, error) {
	return GetJSONBytesReader(privilege)
}

// get the Privilege in the http Response
func (privilege *Privilege) decode(response *http.Response) (*Privilege, error) {
	result := new(Privilege)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}

// populate the Privilege with the data in the terraform resource
func (privilege *Privilege) populate(data *schema.ResourceData) {
	data.SetId(privilege.Id)
	data.Set("key", privilege.Key)
	data.Set("role", privilege.Role)
	data.Set("partition", privilege.Partition)
	data.Set("can_create", privilege.Create)
	data.Set("can_read", privilege.Read)
	data.Set("can_delete", privilege.Delete)
}

// get the FQN for the item resource
func (privilege *Privilege) uri(baseUrl string) string {
	return fmt.Sprintf("%s/privilege/%s", baseUrl, privilege.Key)
}

// issue a put http request with the Privilege data as payload to the resource URI
func (privilege *Privilege) put(meta interface{}) error {
	cfg := meta.(Config)
	// converts the passed-in payload to a bytes Reader
	bytes, err := privilege.toJSON()

	// any errors are returned immediately
	if err != nil {
		return err
	}

	// make an http put request to the service
	result, err := cfg.Client.Put(privilege.uri(cfg.Client.BaseURL), bytes)

	// any errors are returned
	if e := check(result, err); e != nil {
		return e
	}

	return nil
}

// issue a delete http request to the resource URI
func (privilege *Privilege) delete(meta interface{}) error {
	// get the Config instance from the meta object passed-in
	cfg := meta.(Config)

	// make an http delete request to the service
	result, err := cfg.Client.Delete(privilege.uri(cfg.Client.BaseURL))

	// any errors are returned
	if e := check(result, err); e != nil {
		return e
	}

	return nil
}

// issue a get http request to the resource URI
func (privilege *Privilege) get(meta interface{}) (*Privilege, error) {
	cfg := meta.(Config)

	// make an http put request to the service
	result, err := cfg.Client.Get(privilege.uri(cfg.Client.BaseURL))

	if err != nil {
		return nil, err
	}

	i, err := privilege.decode(result)

	defer func() {
		if ferr := result.Body.Close(); ferr != nil {
			err = ferr
		}
	}()

	return i, err
}
