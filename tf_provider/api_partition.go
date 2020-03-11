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

// the Partition resource
type Partition struct {
	Id          string `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Owner       string `json:"owner"`
	Version     int64  `json:"version"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
}

// create a new Partition from a terraform resource
func newPartition(data *schema.ResourceData) *Partition {
	return &Partition{
		Key:         data.Get("key").(string),
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
	}
}

// get a JSON bytes reader for the item
func (partition *Partition) toJSON() (*bytes.Reader, error) {
	return GetJSONBytesReader(partition)
}

// get the Partition in the http Response
func (partition *Partition) decode(response *http.Response) (*Partition, error) {
	result := new(Partition)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}

// populate the Partition with the data in the terraform resource
func (partition *Partition) populate(data *schema.ResourceData) {
	data.SetId(partition.Id)
	data.Set("key", partition.Key)
	data.Set("name", partition.Name)
	data.Set("description", partition.Description)
	data.Set("owner", partition.Owner)
}

// get the FQN for the item resource
func (partition *Partition) uri(baseUrl string) string {
	return fmt.Sprintf("%s/partition/%s", baseUrl, partition.Key)
}

// issue a put http request with the Partition data as payload to the resource URI
func (partition *Partition) put(meta interface{}) error {
	cfg := meta.(Config)
	// converts the passed-in payload to a bytes Reader
	bytes, err := partition.toJSON()

	// any errors are returned immediately
	if err != nil {
		return err
	}

	// make an http put request to the service
	result, err := cfg.Client.Put(partition.uri(cfg.Client.BaseURL), bytes)

	// any errors are returned
	if e := check(result, err); e != nil {
		return e
	}

	return nil
}

// issue a delete http request to the resource URI
func (partition *Partition) delete(meta interface{}) error {
	// get the Config instance from the meta object passed-in
	cfg := meta.(Config)

	// make an http delete request to the service
	result, err := cfg.Client.Delete(partition.uri(cfg.Client.BaseURL))

	// any errors are returned
	if e := check(result, err); e != nil {
		return e
	}

	return nil
}

// issue a get http request to the resource URI
func (partition *Partition) get(meta interface{}) (*Partition, error) {
	cfg := meta.(Config)

	// make an http put request to the service
	result, err := cfg.Client.Get(partition.uri(cfg.Client.BaseURL))

	if err != nil {
		return nil, err
	}

	i, err := partition.decode(result)

	defer func() {
		if ferr := result.Body.Close(); ferr != nil {
			err = ferr
		}
	}()

	return i, err
}
