/*
   Onix Config Manager - Dbman- Onix Database Manager
   Copyright (c) 2018-2020 by www.gatblau.org

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
	"net/http"
)

type Index struct {
	Releases []struct {
		DbVersion  string `json:"dbVersion"`
		AppVersion string `json:"appVersion"`
		Path       string `json:"path"`
	} `json:"releases"`
}

// get a JSON bytes reader for the Index
func (ix *Index) json() (*bytes.Reader, error) {
	jsonBytes, err := ix.bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*jsonBytes), err
}

// get a []byte representing the Index
func (ix *Index) bytes() (*[]byte, error) {
	b, err := jsonBytes(ix)
	return &b, err
}

// get the Index in the http Response
func (ix *Index) decode(response *http.Response) (*Index, error) {
	result := new(Index)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}
