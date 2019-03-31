/*
   Onix CMDB - Copyright (c) 2018-2019 by www.gatblau.org

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

import "github.com/hashicorp/terraform/helper/schema"

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
			"tag": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     schema.TypeString,
				Optional: true,
			},
			"attribute": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func createItem(data *schema.ResourceData, m interface{}) error {
	client := m.(Client)

	key := data.Get("key").(string)
	name := data.Get("name").(string)
	description := data.Get("description").(string)
	itemtype := data.Get("type").(string)
	meta := data.Get("meta").(map[string]interface{})
	attribute := data.Get("attribute").(map[string]interface{})
	tag := data.Get("tag").([]interface{})

	item := Item{
		Key:         key,
		Name:        name,
		Description: description,
		Type:        itemtype,
		Meta:        meta,
		Attribute:   attribute,
		Tag:         tag,
	}
	result, err := client.Put("item", key, item.ToJSON())

	if e := check(result, err); e != nil {
		return e
	}

	data.SetId(item.Key)

	return nil
}

func updateItem(d *schema.ResourceData, m interface{}) error {
	return nil
}

func deleteItem(d *schema.ResourceData, m interface{}) error {
	return nil
}
