/*
    Onix CMDB - Copyright (c) 2018-2019 by www.gatblau.org

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
package data

import (
	"fmt"
	"io"
	"net/http"
)

type OnixClient struct {
	BaseURL string
	Token   string
}

type Result struct {
	Changed   bool   `json:"changed"`
	Error     bool   `json:"error"`
	Message   string `json:"message"`
	Operation string `json:"operation"`
	Ref       string `json:"ref"`
}

func (o *OnixClient) Put(resourceName, key string, payload io.Reader) (result *Result) {

	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/%s/%s", o.BaseURL, resourceName, key), payload)

	req.Header.Set("Content-Type", "application/json")

	resp, _ := http.DefaultClient.Do(req)

	defer resp.Body.Close()

	onixResponse := new(Result)

	return onixResponse
}
