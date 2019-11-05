/*
   Onix Kube - Copyright (c) 2019 by www.gatblau.org

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
)

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
	return GetJSONBytesReader(link)
}

func (link *Link) KeyValue() string {
	return link.Key
}
