/*
   Onix Terra - Copyright (c) 2019 by www.gatblau.org

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
	_, result, _ := m.client.Put(m.payload(), "data")
	return result
}

// gets the terraform metamodel
func (c *MetaModel) payload() Payload {
	return &Data{
		Models: []Model{
			Model{
				Key:         "TERRA",
				Name:        "Terraform Model",
				Description: "Defines the item and link types that describe Terraform resources.",
			},
		},
		ItemTypes: []ItemType{
			ItemType{
				Key:         "K8SCluster",
				Name:        "Kubernetes Cluster",
				Description: "An open-source system for automating deployment, scaling, and management of containerized applications.",
				Model:       "K8SModel",
			},
		},
		LinkTypes: []LinkType{
			LinkType{
				Key:         "K8SLink",
				Name:        "Kubernetes Resource Link Type",
				Description: "Links Kubernetes resources.",
				Model:       "K8SModel",
			},
		},
		LinkRules: []LinkRule{
			LinkRule{
				Key:              fmt.Sprintf("%s->%s", "K8SCluster", "K8SNamespace"),
				Name:             "K8S Cluster to Namespace Rule",
				Description:      "A cluster contains one or more namespaces.",
				LinkTypeKey:      "",
				StartItemTypeKey: "",
				EndItemTypeKey:   "",
			},
		},
	}
}
