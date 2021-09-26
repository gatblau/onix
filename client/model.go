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

type ModelList struct {
	Values []Model
}

func (list *ModelList) reader() (*bytes.Reader, error) {
	jsonBytes, err := ToJson(list)
	return bytes.NewReader(jsonBytes), err
}

// the Model resource
type Model struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Partition   string `json:"partition"`
	Managed     bool   `json:"managed"`
	Version     int64  `json:"version"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
	ChangedBy   string `json:"changedBy"`
}

// Get the Model in the http Response
func (model *Model) decode(response *http.Response) (*Model, error) {
	result := new(Model)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}

// Get the FQN for the model resource
func (model *Model) uri(baseUrl string) (string, error) {
	if len(model.Key) == 0 {
		return "", fmt.Errorf("the model does not have a key: cannot construct Model resource URI")
	}
	return fmt.Sprintf("%s/model/%s", baseUrl, model.Key), nil
}

// Get a JSON bytes reader for the Serializable
func (model *Model) reader() (*bytes.Reader, error) {
	jsonBytes, err := model.bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*jsonBytes), err
}

// Get a []byte representing the Serializable
func (model *Model) bytes() (*[]byte, error) {
	b, err := ToJson(model)
	return &b, err
}

func (model *Model) valid() error {
	if len(model.Key) == 0 {
		return fmt.Errorf("model key is missing")
	}
	if len(model.Name) == 0 {
		return fmt.Errorf("model name is missing")
	}
	return nil
}
