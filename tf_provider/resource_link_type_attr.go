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
	LINK TYPE ATTRIBUTE RESOURCE
*/
func LinkTypeAttributeResource() *schema.Resource {
	return &schema.Resource{
		Create: createOrUpdateLinkTypeAttribute,
		Read:   readLinkTypeAttr,
		Update: createOrUpdateLinkTypeAttribute,
		Delete: deleteLinkTypeAttribute,
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
			"link_type_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func createOrUpdateLinkTypeAttribute(data *schema.ResourceData, m interface{}) error {
	return put(data, m, linkTypeAttributePayload(data), "%s/linktype/%s/attribute/%s", "link_type_key", "key")
}

func deleteLinkTypeAttribute(data *schema.ResourceData, m interface{}) error {
	return delete(m, linkTypeAttributePayload(data), "%s/linktype/%s/attribute/%s", "link_type_key", "key")
}

func linkTypeAttributePayload(data *schema.ResourceData) Payload {
	return &LinkTypeAttribute{
		Key:         data.Get("key").(string),
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
		Type:        data.Get("type").(string),
		DefValue:    data.Get("def_value").(string),
		Managed:     data.Get("managed").(bool),
		Required:    data.Get("managed").(bool),
		Regex:       data.Get("regex").(string),
		LinkTypeKey: data.Get("link_type_key").(string),
	}
}

func LinkTypeAttributeDataSource() *schema.Resource {
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

func readLinkTypeAttr(d *schema.ResourceData, m interface{}) error {
	return nil
}
