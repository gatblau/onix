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
	"github.com/hashicorp/terraform/helper/schema"
)

/*
	LINK RULE RESOURCE
*/

func LinkRuleResource() *schema.Resource {
	return &schema.Resource{
		Create: createLinkRule,
		Read:   readLinkRule,
		Update: updateLinkRule,
		Delete: deleteLinkRule,
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
			"link_type_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"start_item_type_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"end_item_type_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func LinkRuleDataSource() *schema.Resource {
	return &schema.Resource{
		Read: readLinkRule,

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

func createLinkRule(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into a Link Rule
	linkRule := newLinkRule(data)

	// put the Link Type to the Web API
	err := linkRule.put(meta)

	// set Link Type Id key
	data.SetId(linkRule.Key)

	return err
}

func readLinkRule(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into a link rule
	linkRule := newLinkRule(data)

	// get the resource
	linkRule, err := linkRule.get(meta)

	return err
}

func updateLinkRule(data *schema.ResourceData, meta interface{}) error {
	return createLinkRule(data, meta)
}

func deleteLinkRule(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into an Link Type
	linkRule := newLinkRule(data)

	// delete the link Type
	return linkRule.delete(meta)
}
