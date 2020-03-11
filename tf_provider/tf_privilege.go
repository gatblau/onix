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

func PrivilegeResource() *schema.Resource {
	return &schema.Resource{
		Create: createPrivilege,
		Read:   readPrivilege,
		Update: updatePrivilege,
		Delete: deletePrivilege,
		Schema: map[string]*schema.Schema{
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"partition": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"role": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"can_create": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"can_read": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"can_delete": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func PrivilegeDataSource() *schema.Resource {
	return &schema.Resource{
		Read: readPrivilege,

		Schema: map[string]*schema.Schema{
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"partition": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"role": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"can_create": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"can_read": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"can_delete": &schema.Schema{
				Type:     schema.TypeBool,
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

func createPrivilege(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into a Privilege
	privilege := newPrivilege(data)

	// put the Privilege to the Web API
	err := privilege.put(meta)

	// set Privilege Id key
	data.SetId(privilege.Key)

	return err
}

func readPrivilege(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into a Privilege
	privilege := newPrivilege(data)

	// get the resource
	privilege, err := privilege.get(meta)

	return err
}

func updatePrivilege(data *schema.ResourceData, meta interface{}) error {
	// same as create - Web PI is idempotent
	return createPrivilege(data, meta)
}

func deletePrivilege(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into a Privilege
	partition := newPrivilege(data)

	// delete the partition
	return partition.delete(meta)
}
