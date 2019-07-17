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
package main

import (
	"bytes"
)

type LinkType struct {
	Key         string                 `json:"key"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Attribute   map[string]interface{} `json:"attribute"`
	MetaSchema  map[string]interface{} `json:"metaSchema"`
	Model       string                 `json:"modelKey"`
}

func (linkType *LinkType) ToJSON() (*bytes.Reader, error) {
	return GetJSONBytesReader(linkType)
}

func (linkType *LinkType) KeyValue() string {
	return linkType.Key
}
