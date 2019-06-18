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

package provider

import "github.com/hashicorp/terraform/helper/schema"

/*
	LINK RESOURCE
*/

func LinkResource() *schema.Resource {
	return &schema.Resource{
		Create: createOrUpdateLink,
		Read:   readLink,
		Update: createOrUpdateLink,
		Delete: deleteLink,
		Schema: map[string]*schema.Schema{
			"key": &schema.Schema{
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
			"start_item_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"end_item_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func createOrUpdateLink(data *schema.ResourceData, m interface{}) error {
	return put(data, m, linkPayload(data), "link")
}

func deleteLink(data *schema.ResourceData, m interface{}) error {
	return delete(data, m, linkPayload(data), "link")
}

func linkPayload(data *schema.ResourceData) Payload {
	return &Link{
		Key:          data.Get("key").(string),
		Description:  data.Get("description").(string),
		Type:         data.Get("type").(string),
		Meta:         data.Get("meta").(map[string]interface{}),
		Attribute:    data.Get("attribute").(map[string]interface{}),
		Tag:          data.Get("tag").([]interface{}),
		StartItemKey: data.Get("start_item_key").(string),
		EndItemKey:   data.Get("end_item_key").(string),
	}
}
