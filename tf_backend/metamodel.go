/*
   Terraform Http Backend - Onix - Copyright (c) 2018 by www.gatblau.org

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
	"fmt"
	. "gatblau.org/onix/wapic"
	log "github.com/sirupsen/logrus"
)

type MetaModel struct {
	client *Client
	log    *log.Entry
}

// creates a new instance of the Terraform model
func NewModel(log *log.Entry, client *Client) *MetaModel {
	model := new(MetaModel)
	model.log = log
	model.client = client
	return model
}

// checks the terraform model is defined in Onix
func (m *MetaModel) exists() (bool, error) {
	model, err := m.client.Get("model", "TERRAFORM", nil)
	if err != nil {
		return false, err
	}
	return model != nil, nil
}

func (m *MetaModel) create() *Result {
	_, result, _ := m.client.Put(m.getData(), "data")
	return result
}

// gets the terraform meta model data
func (c *MetaModel) getData() Payload {
	return &Data{
		Models: []Model{
			Model{
				Key:         "TERRAFORM",
				Name:        "Terraform Model",
				Description: "Defines the item and link types that describe Terraform resources.",
			},
		},
		ItemTypes: []ItemType{
			ItemType{
				Key:         "TF_STATE",
				Name:        "Terraform State",
				Description: "State about a group of managed infrastructure and configuration resources. This state is used by Terraform to map real world resources to your configuration, keep track of metadata, and to improve performance for large infrastructures.",
				Model:       "TERRAFORM",
			},
			ItemType{
				Key:         "TF_RESOURCE",
				Name:        "Terraform Resource",
				Description: "Each resource block describes one or more infrastructure objects, such as virtual networks, compute instances, or higher-level components such as DNS records.",
				Model:       "TERRAFORM",
			},
		},
		LinkTypes: []LinkType{
			LinkType{
				Key:         "TF_STATE_LINK",
				Name:        "Terraform State Link Type",
				Description: "Links Terraform resources that are part of a state.",
				Model:       "TERRAFORM",
			},
		},
		LinkRules: []LinkRule{
			LinkRule{
				Key:              fmt.Sprintf("%s->%s", "TF_STATE", "TF_RESOURCE"),
				Name:             "Terraform State to Resource Rule",
				Description:      "Allow the linking of a Terraform State item to one or more Terraform Resource items using Terraform State Links.",
				LinkTypeKey:      "TF_STATE_LINK",
				StartItemTypeKey: "TF_STATE",
				EndItemTypeKey:   "TF_RESOURCE",
			},
		},
	}
}
