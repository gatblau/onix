/*
   Onix Config Manager - Copyright (c) 2018-2019 by www.gatblau.org

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
package src

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

func (itemType *ItemType) KeyValue() string {
	return itemType.Key
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
	Partition   string                 `json:"partition"`
}

func (item *Item) ToJSON() (*bytes.Reader, error) {
	return getJSONBytesReader(item)
}

func (item *Item) KeyValue() string {
	return item.Key
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

func (link *Link) KeyValue() string {
	return link.Key
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

func (linkType *LinkType) KeyValue() string {
	return linkType.Key
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

func (rule *LinkRule) KeyValue() string {
	return rule.Key
}

type Model struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Partition   string `json:"partition"`
}

func (model *Model) ToJSON() (*bytes.Reader, error) {
	return getJSONBytesReader(model)
}

func (model *Model) KeyValue() string {
	return model.Key
}

func getJSONBytesReader(data interface{}) (*bytes.Reader, error) {
	jsonBytes, err := json.Marshal(data)
	return bytes.NewReader(jsonBytes), err
}
