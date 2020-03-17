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

func createLinkRule(data *schema.ResourceData, meta interface{}) error {
	// get the Ox client
	c := meta.(Config).Client

	// read the tf data into a Link Rule
	rule := newLinkRule(data)

	// put the Item to the Web API
	err := err(c.PutLinkRule(rule))
	if err != nil {
		return err
	}

	// set Link Rule Id key
	data.SetId(rule.Key)

	return readLinkRule(data, meta)
}

func updateLinkRule(data *schema.ResourceData, meta interface{}) error {
	return createLinkRule(data, meta)
}

func deleteLinkRule(data *schema.ResourceData, meta interface{}) error {
	// get the Ox client
	c := meta.(Config).Client

	// read the resource data into an Item
	rule := newLinkRule(data)

	// delete the item
	return err(c.DeleteLinkRule(rule))
}

func newLinkRule(data *schema.ResourceData) *LinkRule {
	return &LinkRule{
		Key:              data.Get("key").(string),
		Name:             data.Get("name").(string),
		Description:      data.Get("description").(string),
		LinkTypeKey:      data.Get("link_type_key").(string),
		StartItemTypeKey: data.Get("start_item_type_key").(string),
		EndItemTypeKey:   data.Get("end_item_type_key").(string),
	}
}

// populate the Link Rule with the data in the terraform resource
func populateLinkRule(data *schema.ResourceData, linkRule *LinkRule) {
	data.SetId(linkRule.Key)
	data.Set("key", linkRule.Key)
	data.Set("description", linkRule.Description)
	data.Set("link_type_key", linkRule.LinkTypeKey)
	data.Set("start_item_type_key", linkRule.StartItemTypeKey)
	data.Set("end_item_type_key", linkRule.EndItemTypeKey)
}
