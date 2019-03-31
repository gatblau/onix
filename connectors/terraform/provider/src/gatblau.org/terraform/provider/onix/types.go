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
package main

import (
	"bytes"
	"encoding/json"
)

type ItemType struct {
	Key         string                 `json:"key"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	AttrValid   map[string]interface{} `json:"attrValid"`
	Filter      map[string]interface{} `json:"filter"`
	MetaSchema  map[string]interface{} `json:"metaSchema"`
	Model       string                 `json:"modelKey"`
}

func (itemType *ItemType) ToJSON() (*bytes.Reader, error) {
	return getJSONBytesReader(itemType)
}

type Item struct {
	Key         string                 `json:"key"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Status      int                    `json:"status"`
	Type        string                 `json:"type"`
	Tag         []interface{}          `json:"tag"`
	Meta        map[string]interface{} `json:"meta"`
	Attribute   map[string]interface{} `json:"attribute"`
}

func (item *Item) ToJSON() (*bytes.Reader, error) {
	return getJSONBytesReader(item)
}

type Link struct {
	Key          string                 `json:"key"`
	Description  string                 `json:"description"`
	Type         string                 `json:"type"`
	Tag          []interface{}          `json:"tag"`
	Meta         map[string]interface{} `json:"meta"`
	Attribute    map[string]interface{} `json:"attribute"`
	StartItemKey string                 `json:"startItemKey"`
	EndItemKey   string                 `json:"endItemKey"`
}

func (link *Link) ToJSON() (*bytes.Reader, error) {
	return getJSONBytesReader(link)
}

type LinkType struct {
	Key         string                 `json:"key"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Attribute   map[string]interface{} `json:"attribute"`
	MetaSchema  map[string]interface{} `json:"metaSchema"`
	Model       string                 `json:"modelKey"`
}

func (linkType *LinkType) ToJSON() (*bytes.Reader, error) {
	return getJSONBytesReader(linkType)
}

type LinkRule struct {
	Key              string `json:"key"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	LinkTypeKey      string `json:"linkTypeKey"`
	StartItemTypeKey string `json:"startItemTypeKey"`
	EndItemTypeKey   string `json:"endItemTypeKey"`
}

func (linkRule *LinkRule) ToJSON() (*bytes.Reader, error) {
	return getJSONBytesReader(linkRule)
}

type Model struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (model *Model) ToJSON() (*bytes.Reader, error) {
	return getJSONBytesReader(model)
}

func getJSONBytesReader(data interface{}) (*bytes.Reader, error) {
	jsonBytes, err := json.Marshal(data)
	return bytes.NewReader(jsonBytes), err
}
