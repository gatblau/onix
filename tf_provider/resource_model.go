/*
   Onix Config Manager - Terraform Provider
   Copyright (c) 2018-2019 by www.gatblau.org

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

/*
   MODEL RESOURCE
*/

func ModelResource() *schema.Resource {
	return &schema.Resource{
		Create: createOrUpdateModel,
		Read:   readModel,
		Update: createOrUpdateModel,
		Delete: deleteModel,
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
				Required: true,
			},
			"partition": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func createOrUpdateModel(data *schema.ResourceData, m interface{}) error {
	return put(data, m, modelPayload(data), "%s/model/%s", "key", "")
}

func deleteModel(data *schema.ResourceData, m interface{}) error {
	return delete(m, modelPayload(data), "%s/model/%s", "key", "")
}

func modelPayload(data *schema.ResourceData) Payload {
	key := data.Get("key").(string)
	name := data.Get("name").(string)
	description := data.Get("description").(string)
	partition := data.Get("partition").(string)
	return &Model{
		Key:         key,
		Name:        name,
		Description: description,
		Partition:   partition,
	}
}

func ModelDataSource() *schema.Resource {
	return &schema.Resource{
		Read: readModel,

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
				Required: true,
			},
		},
	}
}

func readModel(d *schema.ResourceData, m interface{}) error {
	return nil
}
