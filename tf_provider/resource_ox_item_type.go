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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Type:     schema.TypeString,
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
			"style": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
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
	if err != nil {
		return err
	}

	// set Item Type Id key
	data.SetId(itemType.Key)

	return readItemType(data, meta)
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
		NotifyChange: FromString(data.Get("notify_change").(string)),
		EncryptMeta:  data.Get("encrypt_meta").(bool),
		EncryptTxt:   data.Get("encrypt_txt").(bool),
		Tag:          data.Get("tag").([]interface{}),
		Style:        data.Get("style").(map[string]interface{}),
		Version:      getVersion(data),
	}
}

// populate the Item with the data in the terraform resource
func populateItemType(data *schema.ResourceData, itemType *ItemType) {
	data.SetId(itemType.Key)
	data.Set("key", itemType.Key)
	data.Set("name", itemType.Name)
	data.Set("description", itemType.Description)
	data.Set("filter", itemType.Filter)
	data.Set("meta_schema", itemType.MetaSchema)
	data.Set("notify_change", itemType.NotifyChange)
	data.Set("tag", itemType.Tag)
	data.Set("encrypt_meta", itemType.EncryptMeta)
	data.Set("encrypt_txt", itemType.EncryptTxt)
	data.Set("style", itemType.Style)
	data.Set("model_key", itemType.Model)
	data.Set("created", itemType.Created)
	data.Set("updated", itemType.Updated)
	data.Set("changed_by", itemType.ChangedBy)
}
