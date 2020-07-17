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
	"errors"
)

// Result data retrieved by PUT and DELETE WAPI resources
type Result struct {
	Changed   bool   `json:"changed"`
	Error     bool   `json:"error"`
	Message   string `json:"message"`
	Operation string `json:"operation"`
	Ref       string `json:"ref"`
}

type ResultList struct {
	Values []Item
}

func (list *ResultList) ToJSON() (*bytes.Reader, error) {
	return GetJSONBytesReader(list)
}

// CheckConfigSet for errors in the result and the passed in error
func (r *Result) Check(err error) error {
	if err != nil {
		return err
	} else if r.Error {
		return errors.New(r.Message)
	} else {
		return nil
	}
}
