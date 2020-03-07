/*
   Onix Config Manager - Copyright (c) 2018-2020 by www.gatblau.org

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
	"strings"
)

type ItemType struct {
	Key          string                 `json:"key"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Filter       map[string]interface{} `json:"filter"`
	MetaSchema   map[string]interface{} `json:"metaSchema"`
	Model        string                 `json:"modelKey"`
	NotifyChange bool                   `json:"notifyChange"`
	Tag          []interface{}          `json:"tag"`
	EncryptMeta  bool                   `json:"encryptMeta"`
	EncryptTxt   bool                   `json:"encryptTxt"`
	Managed      bool                   `json:"managed"`
}

func (itemType *ItemType) ToJSON() (*bytes.Reader, error) {
	return getJSONBytesReader(itemType)
}

func (itemType *ItemType) Get(key string) string {
	switch strings.ToLower(key) {
	case "key":
		{
			return itemType.Key
		}
	default:
		panic(fmt.Sprintf("key %s not supported", key))
	}
}

type ItemTypeAttribute struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	DefValue    string `json:"defValue"`
	Managed     bool   `json:"managed"`
	Required    bool   `json:"required"`
	Regex       string `json:"regex"`
	ItemTypeKey string `json:"itemTypeKey"`
}

func (typeAttr *ItemTypeAttribute) ToJSON() (*bytes.Reader, error) {
	return getJSONBytesReader(typeAttr)
}

func (typeAttr *ItemTypeAttribute) Get(key string) string {
	switch strings.ToLower(key) {
	case "key":
		{
			return typeAttr.Key
		}
	case "item_type_key":
		{
			return typeAttr.ItemTypeKey
		}
	default:
		panic(fmt.Sprintf("key %s not supported", key))
	}
}

type Item struct {
	Key         string                 `json:"key"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Status      int                    `json:"status"`
	Type        string                 `json:"type"`
	Tag         []interface{}          `json:"tag"`
	Meta        map[string]interface{} `json:"meta"`
	Txt         string                 `json:"txt"`
	Attribute   map[string]interface{} `json:"attribute"`
	Partition   string                 `json:"partition"`
}

func (item *Item) ToJSON() (*bytes.Reader, error) {
	return getJSONBytesReader(item)
}

func (item *Item) Get(key string) string {
	switch strings.ToLower(key) {
	case "key":
		{
			return item.Key
		}
	default:
		panic(fmt.Sprintf("key %s not supported", key))
	}
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

func (link *Link) Get(key string) string {
	switch strings.ToLower(key) {
	case "key":
		{
			return link.Key
		}
	default:
		panic(fmt.Sprintf("key %s not supported", key))
	}
}

type LinkType struct {
	Key         string                 `json:"key"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Attribute   map[string]interface{} `json:"attribute"`
	MetaSchema  map[string]interface{} `json:"metaSchema"`
	Model       string                 `json:"modelKey"`
	Tag         []interface{}          `json:"tag"`
	EncryptMeta bool                   `json:"encryptMeta"`
	EncryptTxt  bool                   `json:"encryptTxt"`
	Managed     bool                   `json:"managed"`
}

func (linkType *LinkType) ToJSON() (*bytes.Reader, error) {
	return getJSONBytesReader(linkType)
}

func (linkType *LinkType) Get(key string) string {
	switch strings.ToLower(key) {
	case "key":
		{
			return linkType.Key
		}
	default:
		panic(fmt.Sprintf("key %s not supported", key))
	}
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

func (rule *LinkRule) Get(key string) string {
	switch strings.ToLower(key) {
	case "key":
		{
			return rule.Key
		}
	default:
		panic(fmt.Sprintf("key %s not supported", key))
	}
}

type Model struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Partition   string `json:"partition"`
	Managed     bool   `json:"managed"`
}

func (model *Model) ToJSON() (*bytes.Reader, error) {
	return getJSONBytesReader(model)
}

func (model *Model) Get(key string) string {
	switch strings.ToLower(key) {
	case "key":
		{
			return model.Key
		}
	default:
		panic(fmt.Sprintf("key %s not supported", key))
	}
}

func getJSONBytesReader(data interface{}) (*bytes.Reader, error) {
	jsonBytes, err := json.Marshal(data)
	return bytes.NewReader(jsonBytes), err
}
