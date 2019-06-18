/*
   Onix CMDB - Copyright (c) 2018-2019 by www.gatblau.org

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
package provider

import (
	"bytes"
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
)

// Check for errors in the result and the passed in error
func check(result *Result, err error) error {
	if err != nil {
		return err
	} else if result.Error {
		return errors.New(result.Message)
	} else {
		return nil
	}
}

// Configuration information for the Terraform provider
type Config struct {
	URI    string
	User   string
	Pwd    string
	Client Client
}

// Interface implemented by all payload objects to enable
// generic key extraction and conversion to byte Reader
type Payload interface {
	KeyValue() string
	ToJSON() (*bytes.Reader, error)
}

// Executes an HTTP PUT request to the Onix WAPI passing the following parameters:
// - data: a Terraform *schema.ResourceData
// - m: the Terraform provider metadata
// - payload: the payload object
// - resourceName: the WAPI resource name (e.g. item, itemtype, link, etc.)
func put(data *schema.ResourceData, m interface{}, payload Payload, resourceName string) error {
	// get the Config instance from the meta object passed-in
	cfg := m.(Config)

	// converts the passed-in payload to a bytes Reader
	bytes, err := payload.ToJSON()

	// any errors are returned immediately
	if err != nil {
		return err
	}

	// make an http put request to the service
	result, err := cfg.Client.Put(resourceName, payload.KeyValue(), bytes)

	// any errors are returned
	if e := check(result, err); e != nil {
		return e
	}

	// sets the id in the resource data
	data.SetId(payload.KeyValue())

	// return no error
	return nil
}

func delete(data *schema.ResourceData, m interface{}, payload Payload, resourceName string) error {
	// get the Config instance from the meta object passed-in
	cfg := m.(Config)

	// make an http put request to the service
	result, err := cfg.Client.Delete(resourceName, payload.KeyValue())

	// any errors are returned
	if e := check(result, err); e != nil {
		return e
	}

	return nil
}
