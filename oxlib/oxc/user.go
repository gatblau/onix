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
	"strings"
)

type UserList struct {
	Values []User
}

func (list *UserList) reader() (*bytes.Reader, error) {
	jsonBytes, err := ToJson(list)
	return bytes.NewReader(jsonBytes), err
}

// User the user resource
type User struct {
	Key       string `json:"key"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Pwd       string `json:"pwd"`
	Expires   string `json:"expires"`
	Service   bool   `json:"service"`
	ACL       string `json:"acl"`
	Version   int64  `json:"version"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
	ChangedBy string `json:"changedBy"`
}

// Get the Role in the http Response
func (user *User) decode(response *http.Response) (*User, error) {
	result := new(User)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}

// Get the FQN for the item resource
func (user *User) uri(baseUrl string) (string, error) {
	if len(user.Key) == 0 {
		return "", fmt.Errorf("the user does not have a key: cannot construct User resource URI")
	}
	return fmt.Sprintf("%s/user/%s", baseUrl, user.Key), nil
}

// Get a JSON bytes reader for the Serializable
func (user *User) reader() (*bytes.Reader, error) {
	jsonBytes, err := user.bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*jsonBytes), err
}

// Get a []byte representing the Serializable
func (user *User) bytes() (*[]byte, error) {
	b, err := ToJson(user)
	return &b, err
}

func (user *User) valid() error {
	if len(user.Key) == 0 {
		return fmt.Errorf("user key is missing")
	}
	if len(user.Name) == 0 {
		return fmt.Errorf("user name is missing")
	}
	if user.Service {
		if len(user.Email) > 0 {
			return fmt.Errorf("user email is not allowed on service accounts")
		}
	} else {
		if len(user.Email) == 0 {
			return fmt.Errorf("user email is required on user accounts")
		}
	}
	return nil
}

func (user *User) Controls() []AccessControl {
	controls := make([]AccessControl, 0)
	acList := strings.Split(user.ACL, ",")
	for _, c := range acList {
		parts := strings.Split(c, ":")
		if len(parts) != 3 {
			// cannot process then skip
			continue
		}
		controls = append(controls, AccessControl{
			Realm:  parts[0],
			URI:    parts[1],
			Method: parts[2],
		})
	}
	return controls
}

func (user *User) Allowed(acRealm, acURI, acMethod string) bool {
	controls := user.Controls()
	for _, control := range controls {
		if (control.Realm == acRealm || control.Realm == "*") &&
			(control.URI == acURI || control.URI == "*") &&
			(control.Method == acMethod || control.Method == "*") {
			return true
		}
	}
	return false
}

type AccessControl struct {
	Realm  string
	URI    string
	Method string
}
