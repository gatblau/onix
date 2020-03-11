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
	ITEM TYPE RESOURCE
*/
func ItemTypeResource() *schema.Resource {
	return &schema.Resource{
		Create: createItemType,
		Read:   readItemType,
		Update: updateItemType,
		Delete: deleteItemType,
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
			"filter": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"meta_schema": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"model_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"notify_change": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tag": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"encrypt_meta": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"encrypt_txt": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"managed": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func ItemTypeDataSource() *schema.Resource {
	return &schema.Resource{
		Read: readItemType,

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

func createItemType(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into an Item
	itemType := newItemType(data)

	// put the Item Type to the Web API
	err := itemType.put(meta)

	// set Item Type Id key
	data.SetId(itemType.Key)
	return err
}

func readItemType(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into an Item
	itemType := newItemType(data)

	// get the resource
	itemType, err := itemType.get(meta)

	return err
}

func updateItemType(data *schema.ResourceData, meta interface{}) error {
	return createItemType(data, meta)
}

func deleteItemType(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into an Item
	item := newItem(data)

	// delete the item
	return item.delete(meta)
}
