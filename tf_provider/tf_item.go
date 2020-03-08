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

/*
	ITEM RESOURCE
*/
func ItemResource() *schema.Resource {
	return &schema.Resource{
		Create: createItem,
		Read:   readItem,
		Update: updateItem,
		Delete: deleteItem,
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
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"meta": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"txt": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"attribute": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"partition": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func ItemDataSource() *schema.Resource {
	return &schema.Resource{
		Read: readItem,

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
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"meta": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"tag": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     schema.TypeString,
				Optional: true,
			},
			"attribute": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
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

func createItem(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into an Item
	item := newItem(data)

	// put the Item to the Web API
	err := item.put(meta)

	// set Item Id key
	data.SetId(item.Key)
	return err
}

func readItem(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into an Item
	item := newItem(data)

	// get the resource
	item, err := item.get(meta)

	return err
}

func updateItem(data *schema.ResourceData, meta interface{}) error {
	// same as create - Web PI is idempotent
	return createItem(data, meta)
}

func deleteItem(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into an Item
	item := newItem(data)

	// delete the item
	return item.delete(meta)
}
