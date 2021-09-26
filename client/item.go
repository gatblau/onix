package client

/*
   Onix Configuration Manager - HTTP Client
   Copyright (c) 2018-2021 by www.gatblau.org

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

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ItemList struct {
	Values []Item
}

func (list *ItemList) json() (*bytes.Reader, error) {
	jsonBytes, err := ToJson(list)
	return bytes.NewReader(jsonBytes), err
}

// the Item resource
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
	Version     int64                  `json:"version"`
	Created     string                 `json:"created"`
	Updated     string                 `json:"updated"`
	EncKeyIx    int64                  `json:"encKeyIx"`
	ChangedBy   string                 `json:"changedBy"`
}

// Get the Item in the http Response
func (item *Item) decode(response *http.Response) (*Item, error) {
	result := new(Item)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}

// Get the ItemList in the http Response
func decodeItemList(response *http.Response) (*ItemList, error) {
	result := new(ItemList)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}

// Get the FQN for the item resource
func (item *Item) uri(baseUrl string) (string, error) {
	if len(item.Key) == 0 {
		return "", fmt.Errorf("the item does not have a key: cannot construct Item resource URI")
	}
	return fmt.Sprintf("%s/item/%s", baseUrl, item.Key), nil
}

// Get a JSON bytes reader for the Serializable
func (item *Item) reader() (*bytes.Reader, error) {
	jsonBytes, err := item.bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*jsonBytes), err
}

// Get a []byte representing the Serializable
func (item *Item) bytes() (*[]byte, error) {
	b, err := ToJson(item)
	return &b, err
}

func (item *Item) valid() error {
	if len(item.Key) == 0 {
		return fmt.Errorf("item key is missing")
	}
	if len(item.Name) == 0 {
		return fmt.Errorf("item name is missing")
	}
	return nil
}

// Get the FQN for the item / children resource
func (item *Item) uriItemChildren(baseUrl string) (string, error) {
	if len(item.Key) == 0 {
		return "", fmt.Errorf("the item does not have a key: cannot construct Item/Children resource URI")
	}
	return fmt.Sprintf("%s/item/%s/children", baseUrl, item.Key), nil
}

func (item *Item) uriItemFirstLevelChildren(baseUrl, childType string) (string, error) {
	if len(item.Key) == 0 {
		return "", fmt.Errorf("the item does not have a key: cannot construct Item/Children resource URI")
	}
	return fmt.Sprintf("%s/item/%s/list/%s", baseUrl, item.Key, childType), nil
}

func uriItemsOfType(baseUrl, itemType string) (string, error) {
	return fmt.Sprintf("%s/item?type=%s", baseUrl, itemType), nil
}

// GetBoolAttr return the value of the boolean attribute for the specified name or false if the attribute does not exist
func (item *Item) GetBoolAttr(name string) bool {
	attr := item.Attribute[name]
	if attr != nil {
		value, ok := attr.(bool)
		if ok {
			return value
		}
		return false
	}
	return false
}

// GetStringAttr return the value of the string attribute for the specified name or empty if the attribute does not exist
func (item *Item) GetStringAttr(name string) string {
	attr := item.Attribute[name]
	if attr != nil {
		value, ok := attr.(string)
		if ok {
			return value
		}
		return ""
	}
	return ""
}
