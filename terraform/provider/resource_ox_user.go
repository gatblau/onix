/*
   Onix Config Manager - Terraform Provider
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
	. "github.com/gatblau/oxc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func UserResource() *schema.Resource {
	return &schema.Resource{
		Create: createUser,
		Read:   readUser,
		Update: updateUser,
		Delete: deleteUser,
		Schema: map[string]*schema.Schema{
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"pwd": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: userPwdSchemaDiff,
			},
			"service": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"expires": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func createUser(data *schema.ResourceData, meta interface{}) error {
	// get the Ox client
	c := meta.(Config).Client

	// read the tf data into an User
	user := newUser(data)

	// put the User to the Web API
	err := err(c.PutUser(user, false))
	if err != nil {
		return err
	}

	// set User Id key
	data.SetId(user.Key)

	return readUser(data, meta)
}

func updateUser(data *schema.ResourceData, meta interface{}) error {
	// same as create - Web PI is idempotent
	return createUser(data, meta)
}

func deleteUser(data *schema.ResourceData, meta interface{}) error {
	// get the Ox client
	c := meta.(Config).Client

	// read the resource data into a user
	user := newUser(data)

	// delete the user
	return err(c.DeleteUser(user))
}

// create a new User from a terraform resource
func newUser(data *schema.ResourceData) *User {
	return &User{
		Key:     data.Get("key").(string),
		Name:    data.Get("name").(string),
		Email:   data.Get("email").(string),
		Pwd:     data.Get("pwd").(string),
		Service: data.Get("service").(bool),
		Expires: data.Get("expires").(string),
		Version: getVersion(data),
	}
}

// populate the User with the data in the terraform resource
func populateUser(data *schema.ResourceData, user *User) {
	data.SetId(user.Key)
	data.Set("key", user.Key)
	data.Set("name", user.Name)
	data.Set("email", user.Email)
	data.Set("pwd", user.Pwd)
	data.Set("service", user.Service)
	data.Set("expires", user.Expires)
	data.Set("created", user.Created)
	data.Set("updated", user.Updated)
	data.Set("changed_by", user.ChangedBy)
}

func userPwdSchemaDiff(k, old, new string, d *schema.ResourceData) bool {
	// suppress comparison between submitted password and retrieved password converted to *****
	return true
}
