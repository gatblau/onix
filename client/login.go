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
	"fmt"
)

// Login information for users authenticating with client devices such as web browsers
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Get a JSON bytes reader for the Serializable
func (login *Login) reader() (*bytes.Reader, error) {
	jsonBytes, err := login.bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*jsonBytes), err
}

// Get a []byte representing the Serializable
func (login *Login) bytes() (*[]byte, error) {
	b, err := ToJson(login)
	return &b, err
}

// Get the FQN for the item resource
func (login *Login) uri(baseUrl string) (string, error) {
	return fmt.Sprintf("%s/login", baseUrl), nil
}

func (login *Login) valid() error {
	if len(login.Username) == 0 {
		return fmt.Errorf("username is missing")
	}
	if len(login.Password) == 0 {
		return fmt.Errorf("user password is missing")
	}
	return nil
}
