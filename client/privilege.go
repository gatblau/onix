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

type PrivilegeList struct {
	Values []Privilege
}

func (list *PrivilegeList) reader() (*bytes.Reader, error) {
	jsonBytes, err := ToJson(list)
	return bytes.NewReader(jsonBytes), err
}

// the Privilege resource
type Privilege struct {
	Key       string `json:"key"`
	Role      string `json:"roleKey"`
	Partition string `json:"partitionKey"`
	Create    bool   `json:"canCreate"`
	Read      bool   `json:"canRead"`
	Delete    bool   `json:"canDelete"`
	Version   int64  `json:"version"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
	ChangedBy string `json:"changedBy"`
}

// Get the Privilege in the http Response
func (privilege *Privilege) decode(response *http.Response) (*Privilege, error) {
	result := new(Privilege)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}

// Get the FQN for the privilege resource
func (privilege *Privilege) uri(baseUrl string) (string, error) {
	if len(privilege.Key) == 0 {
		return "", fmt.Errorf("the privilege does not have a key: cannot construct privilege resource URI")
	}
	return fmt.Sprintf("%s/privilege/%s", baseUrl, privilege.Key), nil
}

// Get a JSON bytes reader for the Serializable
func (privilege *Privilege) reader() (*bytes.Reader, error) {
	jsonBytes, err := privilege.bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*jsonBytes), err
}

// Get a []byte representing the Serializable
func (privilege *Privilege) bytes() (*[]byte, error) {
	b, err := ToJson(privilege)
	return &b, err
}

func (privilege *Privilege) valid() error {
	if len(privilege.Key) == 0 {
		return fmt.Errorf("privilege key is missing")
	}
	if len(privilege.Role) == 0 {
		return fmt.Errorf("privilege role is missing")
	}
	if len(privilege.Partition) == 0 {
		return fmt.Errorf("privilege partition is missing")
	}
	return nil
}
