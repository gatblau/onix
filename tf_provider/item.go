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
			"enc_key_ix": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func createItem(data *schema.ResourceData, meta interface{}) error {
	// get the Ox client
	c := meta.(Config).Client

	// read the tf data into an Item
	item := newItem(data)

	// put the Item to the Web API
	err := err(c.PutItem(item))

	// set Item Id key
	data.SetId(item.Key)

	return err
}

func readItem(data *schema.ResourceData, meta interface{}) error {
	// get the Ox client
	c := meta.(Config).Client

	// read the tf data into an Item
	item := newItem(data)

	// get the restful resource
	item, err := c.GetItem(item)

	// populate the tf resource data
	if err == nil {
		populateItem(data, item)
	}

	return err
}

func updateItem(data *schema.ResourceData, meta interface{}) error {
	// same as create - Web PI is idempotent
	return createItem(data, meta)
}

func deleteItem(data *schema.ResourceData, meta interface{}) error {
	// get the Ox client
	c := meta.(Config).Client

	// read the resource data into an Item
	item := newItem(data)

	// delete the item
	return err(c.DeleteItem(item))
}

// create a new Item from a terraform resource
func newItem(data *schema.ResourceData) *Item {
	return &Item{
		Key:         data.Get("key").(string),
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
		Type:        data.Get("type").(string),
		Meta:        data.Get("meta").(map[string]interface{}),
		Txt:         data.Get("txt").(string),
		Attribute:   data.Get("attribute").(map[string]interface{}),
		Tag:         data.Get("tag").([]interface{}),
		Partition:   data.Get("partition").(string),
	}
}

// populate the Item with the data in the terraform resource
func populateItem(data *schema.ResourceData, item *Item) {
	data.SetId(item.Id)
	data.Set("key", item.Key)
	data.Set("name", item.Name)
	data.Set("description", item.Description)
	data.Set("attribute", item.Attribute)
	data.Set("txt", item.Txt)
	data.Set("meta", item.Meta)
	data.Set("partition", item.Partition)
	data.Set("status", item.Status)
	data.Set("tag", item.Tag)
	data.Set("type", item.Type)
}
