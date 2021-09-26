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
	"fmt"
	"net/http"
)

// clear all data in the database
func (c *Client) Clear() (*Result, error) {
	resp, err := c.Delete(fmt.Sprintf("%s/clear", c.conf.BaseURI), c.addHttpHeaders)
	return result(resp, err)
}

// generic function to check for errors and retrieve a result
func result(resp *http.Response, err error) (*Result, error) {
	// if there is a response
	if resp != nil {
		// extract the response result
		result, err2 := newResult(resp)
		if err2 != nil {
			return result, err2
		}
		return result, err
	}
	return nil, err
}
