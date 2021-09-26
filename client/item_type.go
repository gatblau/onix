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
	"strings"
)

type ItemTypeList struct {
	Values []ItemType
}

func (list *ItemTypeList) reader() (*bytes.Reader, error) {
	jsonBytes, err := ToJson(list)
	return bytes.NewReader(jsonBytes), err
}

type ChangeNotifyType string

const (
	NotifyTypeNone ChangeNotifyType = "N"
	NotifyTypeType ChangeNotifyType = "T"
	NotifyTypeItem ChangeNotifyType = "I"
)

func FromString(changeNotifyType string) ChangeNotifyType {
	switch strings.ToUpper(changeNotifyType) {
	case "N":
		return NotifyTypeNone
	case "T":
		return NotifyTypeType
	case "I":
		return NotifyTypeItem
	default:
		return NotifyTypeNone
	}
}

func (n *ChangeNotifyType) ToString() string {
	switch v := interface{}(n).(type) {
	case string:
		return fmt.Sprint(v)
	default:
		return fmt.Sprint(NotifyTypeNone)
	}
}

// the Item Type resource
type ItemType struct {
	Key          string                 `json:"key"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Filter       map[string]interface{} `json:"filter"`
	MetaSchema   map[string]interface{} `json:"metaSchema"`
	Model        string                 `json:"modelKey"`
	NotifyChange ChangeNotifyType       `json:"notifyChange"`
	Tag          []interface{}          `json:"tag"`
	EncryptMeta  bool                   `json:"encryptMeta"`
	EncryptTxt   bool                   `json:"encryptTxt"`
	Style        map[string]interface{} `json:"style"`
	Version      int64                  `json:"version"`
	Created      string                 `json:"created"`
	Updated      string                 `json:"updated"`
	ChangedBy    string                 `json:"changedBy"`
}

// Get the Item Type in the http Response
func (itemType *ItemType) decode(response *http.Response) (*ItemType, error) {
	result := new(ItemType)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}

// Get the FQN for the item type resource
func (itemType *ItemType) uri(baseUrl string) (string, error) {
	if len(itemType.Key) == 0 {
		return "", fmt.Errorf("the item type does not have a key: cannot construct itemtype resource URI")
	}
	return fmt.Sprintf("%s/itemtype/%s", baseUrl, itemType.Key), nil
}

// Get a JSON bytes reader for the Serializable
func (itemType *ItemType) reader() (*bytes.Reader, error) {
	jsonBytes, err := itemType.bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*jsonBytes), err
}

// Get a []byte representing the Serializable
func (itemType *ItemType) bytes() (*[]byte, error) {
	b, err := ToJson(itemType)
	return &b, err
}

func (itemType *ItemType) valid() error {
	if len(itemType.Key) == 0 {
		return fmt.Errorf("item type key is missing")
	}
	if len(itemType.Name) == 0 {
		return fmt.Errorf("item type name is missing")
	}
	return nil
}
