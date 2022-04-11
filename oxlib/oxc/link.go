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

type LinkList struct {
	Values []Link
}

func (list *LinkList) reader() (*bytes.Reader, error) {
	jsonBytes, err := ToJson(list)
	return bytes.NewReader(jsonBytes), err
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
	Version      int64                  `json:"version"`
	Created      string                 `json:"created"`
	Updated      string                 `json:"updated"`
	ChangedBy    string                 `json:"changedBy"`
}

// Get the Link in the http Response
func (link *Link) decode(response *http.Response) (*Link, error) {
	result := new(Link)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}

// Get the FQN for the link type resource
func (link *Link) uri(baseUrl string) (string, error) {
	if len(link.Key) == 0 {
		return "", fmt.Errorf("the link does not have a key: cannot construct link resource URI")
	}
	return fmt.Sprintf("%s/link/%s", baseUrl, link.Key), nil
}

func uriLinks(baseUrl string) (string, error) {
	return fmt.Sprintf("%s/link", baseUrl), nil
}

// Get the LinkList in the http Response
func decodeLinkList(response *http.Response) (*LinkList, error) {
	result := new(LinkList)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}

// Get a JSON bytes reader for the Serializable
func (link *Link) reader() (*bytes.Reader, error) {
	jsonBytes, err := link.bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*jsonBytes), err
}

// Get a []byte representing the Serializable
func (link *Link) bytes() (*[]byte, error) {
	b, err := ToJson(link)
	return &b, err
}

func (link *Link) valid() error {
	if len(link.Key) == 0 {
		return fmt.Errorf("link key is missing")
	}
	if len(link.Type) == 0 {
		return fmt.Errorf("link type is missing")
	}
	if len(link.StartItemKey) == 0 {
		return fmt.Errorf("start item key is missing")
	}
	if len(link.EndItemKey) == 0 {
		return fmt.Errorf("end item key is missing")
	}
	return nil
}
