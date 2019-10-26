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
		Create: createOrUpdateItem,
		Read:   readItem,
		Update: createOrUpdateItem,
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

func createOrUpdateItem(data *schema.ResourceData, m interface{}) error {
	return put(data, m, itemPayload(data), "item")
}

func deleteItem(data *schema.ResourceData, m interface{}) error {
	return delete(m, itemPayload(data), "item")
}

func itemPayload(data *schema.ResourceData) Payload {
	return &Item{
		Key:         data.Get("key").(string),
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
		Type:        data.Get("type").(string),
		Meta:        data.Get("meta").(map[string]interface{}),
		Attribute:   data.Get("attribute").(map[string]interface{}),
		Tag:         data.Get("tag").([]interface{}),
		Partition:   data.Get("partition").(string),
	}
}
