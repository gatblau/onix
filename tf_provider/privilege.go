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
		},
	}
}

func createPrivilege(data *schema.ResourceData, meta interface{}) error {
	// get the Ox client
	c := meta.(Config).Client

	// read the tf data into an Privilege
	privilege := newPrivilege(data)

	// put the Privilege to the Web API
	err := err(c.PutPrivilege(privilege))

	// set Privilege Id key
	data.SetId(privilege.Key)

	return err
}

func readPrivilege(data *schema.ResourceData, meta interface{}) error {
	// get the Ox client
	c := meta.(Config).Client

	// read the tf data into an Privilege
	privilege := &Privilege{Key: data.Get("key").(string)}

	// get the restful resource
	privilege, err := c.GetPrivilege(privilege)

	// populate the tf resource data
	if err == nil {
		populatePrivilege(data, privilege)
	}

	return err
}

func updatePrivilege(data *schema.ResourceData, meta interface{}) error {
	// same as create - Web PI is idempotent
	return createPrivilege(data, meta)
}

func deletePrivilege(data *schema.ResourceData, meta interface{}) error {
	// get the Ox client
	c := meta.(Config).Client

	// read the resource data into a privilege
	privilege := newPrivilege(data)

	// delete the privilege
	return err(c.DeletePrivilege(privilege))
}

// create a new Privilege from a terraform resource
func newPrivilege(data *schema.ResourceData) *Privilege {
	return &Privilege{
		Key:       data.Get("key").(string),
		Role:      data.Get("role").(string),
		Partition: data.Get("partition").(string),
		Create:    data.Get("can_create").(bool),
		Read:      data.Get("can_read").(bool),
		Delete:    data.Get("can_delete").(bool),
	}
}

// populate the Privilege with the data in the terraform resource
func populatePrivilege(data *schema.ResourceData, privilege *Privilege) {
	data.SetId(privilege.Id)
	data.Set("key", privilege.Key)
	data.Set("role", privilege.Role)
	data.Set("partition", privilege.Partition)
	data.Set("can_create", privilege.Create)
	data.Set("can_read", privilege.Read)
	data.Set("can_delete", privilege.Delete)
}
