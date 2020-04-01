/*
   Onix Config Manager - OxTerra - Terraform Http Backend for Onix
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
	"errors"
	"fmt"
	. "github.com/gatblau/oxc"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

type TerraModel struct {
	client *Client
}

// creates a new instance of the Terraform model
func NewTerraModel(client *Client) *TerraModel {
	model := new(TerraModel)
	model.client = client
	return model
}

// checks the Terra model is defined in Onix
func (m *TerraModel) exists() (bool, error) {
	model, err := m.client.GetModel(&Model{Key: "TERRAFORM"})
	if err != nil {
		return false, err
	}
	return model != nil, nil
}

// create the Terra model in Onix
func (m *TerraModel) create() error {
	var (
		exist    bool
		attempts int
		interval time.Duration = 30 // the interval to wait for reconnection
		err      error
	)
	// tries and connects to Onix using retry if the service is not there
	for {
		// check if the model exists
		exist, err = m.exists()
		if err == nil {
			// could connect to the Web API therefore breaks the retry loop
			break
		}
		// could not connect so increment the retry count
		attempts = attempts + 1
		// issue a warning to the console output
		log.Warn().Msgf("Can't connect to Onix: %s. Attempt %s, waiting before attempting to connect again.", err, strconv.Itoa(attempts))
		// wait of a second before retrying
		time.Sleep(interval * time.Second)
	}
	// if the model is not defined in Onix
	if !exist {
		// create the model
		log.Trace().Msg("The TERRA model is not yet defined in Onix, proceeding to create it.")
		result, err := m.client.PutData(m.getModelData())
		if err != nil {
			log.Error().Msgf("Can't create TERRA model: %s", err)
			return err
		}
		if result.Error {
			log.Error().Msgf("Can't create TERRA meta-model: %s", result.Message)
			return errors.New(result.Message)
		}
	} else {
		log.Trace().Msg("TERRA model found in Onix.")
	}
	return nil
}

// gets the Terra's meta model data
func (m *TerraModel) getModelData() *GraphData {
	return &GraphData{
		Models: []Model{
			Model{
				Key:         "TERRA",
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
