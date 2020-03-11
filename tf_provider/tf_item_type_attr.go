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

func ItemTypeAttributeResource() *schema.Resource {
	return &schema.Resource{
		Create: createItemTypeAttribute,
		Read:   readItemTypeAttr,
		Update: updateItemTypeAttribute,
		Delete: deleteItemTypeAttribute,
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
			"def_value": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"managed": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"required": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"regex": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"item_type_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func ItemTypeAttributeDataSource() *schema.Resource {
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

func createItemTypeAttribute(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into an Item Type Attribute
	itemTypeAttr := newItemTypeAttr(data)

	// put the Item Type Attribute to the Web API
	err := itemTypeAttr.put(meta)

	// set Item Type Attribute Id key
	data.SetId(itemTypeAttr.Key)

	return err
}

func readItemTypeAttr(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into an Item
	itemTypeAttr := newItemTypeAttr(data)

	// get the resource
	itemTypeAttr, err := itemTypeAttr.get(meta)

	return err
}

func updateItemTypeAttribute(data *schema.ResourceData, meta interface{}) error {
	return createItemTypeAttribute(data, meta)
}

func deleteItemTypeAttribute(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into an Item Type Attribute
	itemTypeAttr := newItemTypeAttr(data)

	// delete the itemTypeAttr
	return itemTypeAttr.delete(meta)
}
