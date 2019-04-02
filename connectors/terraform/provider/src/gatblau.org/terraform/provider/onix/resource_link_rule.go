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
	LINK RULE RESOURCE
*/

func LinkRuleResource() *schema.Resource {
	return &schema.Resource{
		Create: createOrUpdateLinkRule,
		Read:   readLinkRule,
		Update: createOrUpdateLinkRule,
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

func createOrUpdateLinkRule(data *schema.ResourceData, m interface{}) error {
	return put(data, m, linkRulePayload(data), "linkrule")
}

func deleteLinkRule(data *schema.ResourceData, m interface{}) error {
	return delete(data, m, linkRulePayload(data), "linkrule")
}

func linkRulePayload(data *schema.ResourceData) Payload {
	return &LinkRule{
		Key:              data.Get("key").(string),
		Name:             data.Get("name").(string),
		Description:      data.Get("description").(string),
		LinkTypeKey:      data.Get("link_type_key").(string),
		StartItemTypeKey: data.Get("start_item_type_key").(string),
		EndItemTypeKey:   data.Get("end_item_type_key").(string),
	}
}
