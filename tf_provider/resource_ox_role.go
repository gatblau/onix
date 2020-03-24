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

func RoleResource() *schema.Resource {
	return &schema.Resource{
		Create: createRole,
		Read:   readRole,
		Update: updateRole,
		Delete: deleteRole,
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
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"level": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"created": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"updated": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"changed_by": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func createRole(data *schema.ResourceData, meta interface{}) error {
	// get the Ox client
	c := meta.(Config).Client

	// read the tf data into an Role
	role := newRole(data)

	// put the Role to the Web API
	err := err(c.PutRole(role))
	if err != nil {
		return err
	}

	// set Role Id key
	data.SetId(role.Key)

	return readRole(data, meta)
}

func updateRole(data *schema.ResourceData, meta interface{}) error {
	// same as create - Web PI is idempotent
	return createRole(data, meta)
}

func deleteRole(data *schema.ResourceData, meta interface{}) error {
	// get the Ox client
	c := meta.(Config).Client

	// read the resource data into a role
	role := newRole(data)

	// delete the role
	return err(c.DeleteRole(role))
}

// create a new Role from a terraform resource
func newRole(data *schema.ResourceData) *Role {
	return &Role{
		Key:         data.Get("key").(string),
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
		Level:       data.Get("level").(int),
		Version:     getVersion(data),
	}
}

// populate the Role with the data in the terraform resource
func populateRole(data *schema.ResourceData, role *Role) {
	data.SetId(role.Key)
	data.Set("key", role.Key)
	data.Set("name", role.Name)
	data.Set("description", role.Description)
	data.Set("level", role.Level)
	data.Set("created", role.Created)
	data.Set("updated", role.Updated)
	data.Set("changed_by", role.ChangedBy)
}
