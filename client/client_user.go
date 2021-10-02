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

// PutUser issue a Put http request with the User data as payload to the resource URI
// notify: if true, emails new users to make them aware of the new account
//   requires the service to have email integration enabled
func (c *Client) PutUser(user *User, notify bool) (*Result, error) {
	// validates user
	if err := user.valid(); err != nil {
		return nil, err
	}
	uri, err := user.uri(c.conf.BaseURI)
	if err != nil {
		return nil, err
	}
	if notify {
		uri += "?notify=true"
	} else {
		uri += "?notify=false"
	}
	resp, err := c.Put(uri, user, c.addHttpHeaders)
	return result(resp, err)
}

// DeleteUser issue a Delete http request to the resource URI
func (c *Client) DeleteUser(user *User) (*Result, error) {
	uri, err := user.uri(c.conf.BaseURI)
	if err != nil {
		return nil, err
	}
	resp, err := c.Delete(uri, c.addHttpHeaders)
	return result(resp, err)
}

// GetUser issue a Get http request to the resource URI
func (c *Client) GetUser(user *User) (*User, error) {
	uri, err := user.uri(c.conf.BaseURI)
	if err != nil {
		return nil, err
	}
	result, err := c.Get(uri, c.addHttpHeaders)
	if err != nil {
		return nil, err
	}
	i, err := user.decode(result)
	defer func() {
		if ferr := result.Body.Close(); ferr != nil {
			err = ferr
		}
	}()
	return i, err
}
