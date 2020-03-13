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
	. "github.com/gatblau/oxc"
	"github.com/hashicorp/terraform/helper/schema"
)

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
		},
	}
}

func createItemType(data *schema.ResourceData, meta interface{}) error {
	// get the Ox client
	c := meta.(Config).Client

	// read the resource data into an Item
	itemType := newItemType(data)

	// put the Item Type to the Web API
	err := err(c.PutItemType(itemType))

	// set Item Type Id key
	data.SetId(itemType.Key)

	return err
}

func readItemType(data *schema.ResourceData, meta interface{}) error {
	// get the Ox client
	c := meta.(Config).Client

	// read the resource data into an Item
	itemType := &ItemType{Key: data.Get("key").(string)}

	// get the resource
	itemType, err := c.GetItemType(itemType)

	// populate the tf resource data
	if err == nil {
		populateItemType(data, itemType)
	}

	return err
}

func updateItemType(data *schema.ResourceData, meta interface{}) error {
	return createItemType(data, meta)
}

func deleteItemType(data *schema.ResourceData, meta interface{}) error {
	// get the Ox client
	c := meta.(Config).Client

	// read the resource data into an Item Type
	itemType := newItemType(data)

	// delete the item
	return err(c.DeleteItemType(itemType))
}

// create a new Item from a terraform resource
func newItemType(data *schema.ResourceData) *ItemType {
	return &ItemType{
		Key:          data.Get("key").(string),
		Name:         data.Get("name").(string),
		Description:  data.Get("description").(string),
		Model:        data.Get("model_key").(string),
		Filter:       data.Get("filter").(map[string]interface{}),
		MetaSchema:   data.Get("meta_schema").(map[string]interface{}),
		NotifyChange: data.Get("notify_change").(bool),
		EncryptMeta:  data.Get("encrypt_meta").(bool),
		EncryptTxt:   data.Get("encrypt_txt").(bool),
		Managed:      data.Get("managed").(bool),
		Tag:          data.Get("tag").([]interface{}),
	}
}

// populate the Item with the data in the terraform resource
func populateItemType(data *schema.ResourceData, itemType *ItemType) {
	data.SetId(itemType.Id)
	data.Set("key", itemType.Key)
	data.Set("name", itemType.Name)
	data.Set("description", itemType.Description)
	data.Set("filter", itemType.Filter)
	data.Set("meta_schema", itemType.MetaSchema)
	data.Set("notify_change", itemType.NotifyChange)
	data.Set("tag", itemType.Tag)
	data.Set("encrypt_meta", itemType.EncryptMeta)
	data.Set("encrypt_txt", itemType.EncryptTxt)
	data.Set("managed", itemType.Managed)
	data.Set("model", itemType.Model)
	data.Set("version", itemType.Version)
	data.Set("created", itemType.Created)
	data.Set("updated", itemType.Updated)
}
