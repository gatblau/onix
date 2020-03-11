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
	"github.com/hashicorp/terraform/helper/schema"
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
		},
	}
}

func RoleDataSource() *schema.Resource {
	return &schema.Resource{
		Read: readRole,

		Schema: map[string]*schema.Schema{
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"level": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"created": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"changedby": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func createRole(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into a Role
	role := newRole(data)

	// put the Role to the Web API
	err := role.put(meta)

	// set Role Id key
	data.SetId(role.Key)

	return err
}

func readRole(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into a Role
	role := newRole(data)

	// get the resource
	role, err := role.get(meta)

	return err
}

func updateRole(data *schema.ResourceData, meta interface{}) error {
	// same as create - Web PI is idempotent
	return createRole(data, meta)
}

func deleteRole(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into a Role
	role := newRole(data)

	// delete the role
	return role.delete(meta)
}
