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
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/gatblau/oxc"
	"strconv"
)

// the Terraform state object
type TfState struct {
	Version          int                    `json:"version"`
	TerraformVersion string                 `json:"terraform_version"`
	Serial           int                    `json:"serial"`
	Lineage          string                 `json:"lineage"`
	Outputs          map[string]interface{} `json:"outputs"`
	Resources        []TfResource           `json:"resources"`
}

type TfResource struct {
	Mode      string       `json:"mode"`
	Type      string       `json:"type"`
	Name      string       `json:"name"`
	Provider  string       `json:"provider"`
	Instances []TfInstance `json:"instances"`
}

type TfInstance struct {
	SchemaVersion int                    `json:"schema_version"`
	Attributes    map[string]interface{} `json:"attributes"`
	Private       string                 `json:"private"`
	Dependencies  []string               `json:"dependencies"`
}

type TfInstances struct {
	Instances []TfInstance `json:"instances"`
}

// persist the state in Onix
func (state *TfState) save(client *Client, key string) error {
	err := state.saveStateItem(client, key)
	if err != nil {
		return err
	}
	err = state.saveResources(client, key)
	if err != nil {
		return err
	}
	return err
}

// retrieve the terraform state
func (state *TfState) loadState(client *Client, key string) error {
	keyItem := &Item{Key: stateKey(key)}
	// load the state
	item, err := client.GetItem(keyItem)
	if err != nil || item == nil {
		return err
	}
	value := item.Attribute["version"]
	if value != nil {
		intValue, _ := strconv.Atoi(value.(string))
		state.Version = intValue
	}
	value = item.Attribute["lineage"]
	if value != nil {
		state.Lineage = value.(string)
	}
	value = item.Attribute["serial"]
	if value != nil {
		intValue, _ := strconv.Atoi(value.(string))
		state.Serial = intValue
	}
	state.Outputs = item.Meta
	// load the resources associated to the state
	list, err := client.GetItemChildren(keyItem)
	for _, item := range list.Values {
		resx := TfResource{
			Mode:      item.Attribute["mode"].(string),
			Type:      item.Attribute["type"].(string),
			Name:      item.Attribute["name"].(string),
			Provider:  item.Attribute["provider"].(string),
			Instances: state.instance(item.Meta),
		}
		state.Resources = append(state.Resources, resx)
	}
	return nil
}

// save the state item
func (state *TfState) saveStateItem(client *Client, key string) error {
	attrs := map[string]interface{}{}
	attrs["version"] = state.Version
	attrs["terraform_version"] = state.TerraformVersion
	attrs["serial"] = state.Serial
	attrs["lineage"] = state.Lineage
	result, err := client.PutItem(&Item{
		Key:         stateKey(key),
		Name:        fmt.Sprintf("STATE -> %s", key),
		Description: "",
		Type:        TfStateType,
		Attribute:   attrs,
		Meta:        state.Outputs,
		Tag:         []interface{}{"terraform", "state"},
	})
	return state.check(err, result)
}

// save terraform resources and link them to the state item
func (state *TfState) saveResources(client *Client, key string) error {
	for i := 0; i < len(state.Resources); i++ {
		attrs := map[string]interface{}{}
		attrs["name"] = state.Resources[i].Name
		attrs["mode"] = state.Resources[i].Mode
		attrs["type"] = state.Resources[i].Type
		attrs["provider"] = state.Resources[i].Provider
		meta := map[string]interface{}{}
		meta["instances"] = state.Resources[i].Instances
		itemKey := fmt.Sprintf("TF_RESOURCE_%s_%s", key, state.Resources[i].Name)
		result, err := client.PutItem(&Item{
			Key:         itemKey,
			Name:        fmt.Sprintf("RESOURCE -> %s", state.Resources[i].Name),
			Description: "",
			Type:        TfResourceType,
			Attribute:   attrs,
			Meta:        meta,
			Tag:         []interface{}{"terraform", "resource"},
		})
		if err = state.check(err, result); err != nil {
			return err
		}
		result, err = client.PutLink(&Link{
			Key:          fmt.Sprintf("FT_LINK:%s->%s", key, itemKey),
			Description:  "",
			Type:         TfLinkType,
			Tag:          []interface{}{"terraform"},
			StartItemKey: fmt.Sprintf("TF_STATE_%s", key),
			EndItemKey:   itemKey,
		})
		if err = state.check(err, result); err != nil {
			return err
		}
	}
	return nil
}

// check for errors in result and returns a single error
func (state *TfState) check(err error, result *Result) error {
	if err != nil {
		return err
	}
	if result.Error {
		return errors.New(fmt.Sprintf("%s - %s", result.Ref, result.Message))
	}
	return nil
}

// get the JSON string representation of the TfState
func (state *TfState) toJSONString() string {
	result, _ := json.Marshal(state)
	return string(result[:])
}

func (state *TfState) instance(meta map[string]interface{}) []TfInstance {
	metaBytes, err := json.Marshal(meta)
	if err != nil {
		return nil
	}
	// metaStr := string(metaBytes)
	var instances TfInstances
	err = json.Unmarshal(metaBytes, &instances)
	if err != nil {
		return nil
	}
	return instances.Instances
}

func stateKey(key string) string {
	return fmt.Sprintf("TF_STATE_%s", key)
}
