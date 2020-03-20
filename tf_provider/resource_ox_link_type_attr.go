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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func LinkTypeAttributeResource() *schema.Resource {
	return &schema.Resource{
		Create: createLinkTypeAttribute,
		Read:   readLinkTypeAttr,
		Update: updateLinkTypeAttribute,
		Delete: deleteLinkTypeAttribute,
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
			"link_type_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func createLinkTypeAttribute(data *schema.ResourceData, meta interface{}) error {
	// get the Ox client
	c := meta.(Config).Client

	// read the resource data into an Link Type Attribute
	linkTypeAttr := newLinkTypeAttr(data)

	// put the Link Type Attr to the Web API
	err := err(c.PutLinkTypeAttr(linkTypeAttr))
	if err != nil {
		return err
	}

	// set Link Type Attribute Id key
	data.SetId(linkTypeAttr.Key)

	return readLinkTypeAttr(data, meta)
}

func updateLinkTypeAttribute(data *schema.ResourceData, meta interface{}) error {
	return createLinkTypeAttribute(data, meta)
}

func deleteLinkTypeAttribute(data *schema.ResourceData, meta interface{}) error {
	// get the Ox client
	c := meta.(Config).Client

	// read the resource data into an Link Type Attribute
	linkTypeAttr := newLinkTypeAttr(data)

	// delete the linkTypeAttr
	return err(c.DeleteLinkTypeAttr(linkTypeAttr))
}

func newLinkTypeAttr(data *schema.ResourceData) *LinkTypeAttribute {
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
		Version:     getVersion(data),
	}
}

// populate the Link Type Attribute with the data in the terraform resource
func populateLinkTypeAttr(data *schema.ResourceData, typeAttr *LinkTypeAttribute) {
	data.SetId(typeAttr.Key)
	data.Set("key", typeAttr.Key)
	data.Set("description", typeAttr.Description)
	data.Set("type", typeAttr.Type)
	data.Set("def_value", typeAttr.DefValue)
	data.Set("managed", typeAttr.Managed)
	data.Set("required", typeAttr.Required)
	data.Set("regex", typeAttr.Regex)
	data.Set("link_type_key", typeAttr.LinkTypeKey)
	data.Set("created", typeAttr.Created)
	data.Set("updated", typeAttr.Updated)
	data.Set("changed_by", typeAttr.ChangedBy)
}
