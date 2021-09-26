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

type ItemTypeAttributeList struct {
	Values []ItemTypeAttribute
}

func (list *ItemTypeAttributeList) reader() (*bytes.Reader, error) {
	jsonBytes, err := ToJson(list)
	return bytes.NewReader(jsonBytes), err
}

type ItemTypeAttribute struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	DefValue    string `json:"defValue"`
	Required    bool   `json:"required"`
	Regex       string `json:"regex"`
	ItemTypeKey string `json:"itemTypeKey"`
	Version     int64  `json:"version"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
	ChangedBy   string `json:"changedBy"`
}

// Get the Item Type Attribute in the http Response
func (typeAttr *ItemTypeAttribute) decode(response *http.Response) (*ItemTypeAttribute, error) {
	result := new(ItemTypeAttribute)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}

// Get the FQN for the item resource
func (typeAttr *ItemTypeAttribute) uri(baseUrl string) (string, error) {
	if len(typeAttr.ItemTypeKey) == 0 {
		return "", fmt.Errorf("the item type attribute does not have an item type key: cannot construct itemtype attr resource URI")
	}
	if len(typeAttr.Key) == 0 {
		return "", fmt.Errorf("the item type attribute does not have a key: cannot construct itemtype attr resource URI")
	}
	return fmt.Sprintf("%s/itemtype/%s/attribute/%s", baseUrl, typeAttr.ItemTypeKey, typeAttr.Key), nil
}

// Get a JSON bytes reader for the Serializable
func (typeAttr *ItemTypeAttribute) reader() (*bytes.Reader, error) {
	jsonBytes, err := typeAttr.bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*jsonBytes), err
}

// Get a []byte representing the Serializable
func (typeAttr *ItemTypeAttribute) bytes() (*[]byte, error) {
	b, err := ToJson(typeAttr)
	return &b, err
}

func (typeAttr *ItemTypeAttribute) valid() error {
	if len(typeAttr.Key) == 0 {
		return fmt.Errorf("item type attribute key is missing")
	}
	if len(typeAttr.ItemTypeKey) == 0 {
		return fmt.Errorf("item type key is missing")
	}
	return nil
}
