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

func createLinkTypeAttribute(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into an Link Type Attribute
	linkTypeAttr := newLinkTypeAttr(data)

	// put the Link Type Attribute to the Web API
	err := linkTypeAttr.put(meta)

	// set Link Type Attribute Id key
	data.SetId(linkTypeAttr.Key)

	return err
}

func readLinkTypeAttr(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into an Link Type Attribute
	linkTypeAttr := newLinkTypeAttr(data)

	// get the resource
	linkTypeAttr, err := linkTypeAttr.get(meta)

	return err
}

func updateLinkTypeAttribute(data *schema.ResourceData, meta interface{}) error {
	return createLinkTypeAttribute(data, meta)
}

func deleteLinkTypeAttribute(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into an Link Type Attribute
	linkTypeAttr := newLinkTypeAttr(data)

	// delete the Link Type Attr
	return linkTypeAttr.delete(meta)
}
