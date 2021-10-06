package oxc

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

type LinkTypeList struct {
	Values []Item
}

func (list *LinkTypeList) reader() (*bytes.Reader, error) {
	jsonBytes, err := ToJson(list)
	return bytes.NewReader(jsonBytes), err
}

type LinkType struct {
	Key         string                 `json:"key"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	MetaSchema  map[string]interface{} `json:"metaSchema"`
	Model       string                 `json:"modelKey"`
	Tag         []interface{}          `json:"tag"`
	EncryptMeta bool                   `json:"encryptMeta"`
	EncryptTxt  bool                   `json:"encryptTxt"`
	Style       map[string]interface{} `json:"style"`
	Version     int64                  `json:"version"`
	Created     string                 `json:"created"`
	Updated     string                 `json:"updated"`
	ChangedBy   string                 `json:"changedBy"`
}

// Get the Link Type in the http Response
func (linkType *LinkType) decode(response *http.Response) (*LinkType, error) {
	result := new(LinkType)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}

// Get the FQN for the link type resource
func (linkType *LinkType) uri(baseUrl string) (string, error) {
	if len(linkType.Key) == 0 {
		return "", fmt.Errorf("the link type does not have a key: cannot construct linktype resource URI")
	}
	return fmt.Sprintf("%s/linktype/%s", baseUrl, linkType.Key), nil
}

// Get a JSON bytes reader for the Serializable
func (linkType *LinkType) reader() (*bytes.Reader, error) {
	jsonBytes, err := linkType.bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*jsonBytes), err
}

// Get a []byte representing the Serializable
func (linkType *LinkType) bytes() (*[]byte, error) {
	b, err := ToJson(linkType)
	return &b, err
}

func (linkType *LinkType) valid() error {
	if len(linkType.Key) == 0 {
		return fmt.Errorf("link key is missing")
	}
	if len(linkType.Model) == 0 {
		return fmt.Errorf("model is missing")
	}
	return nil
}
