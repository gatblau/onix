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
	ITEM DATA SOURCE
 */
func ItemDataSource() *schema.Resource {
	return &schema.Resource{
		Read: readItem,

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
			"itemtype": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
			},
			"meta": &schema.Schema{
				Type:     schema.TypeMap,
				Required: false,
				Optional: true,
			},
			"tag": &schema.Schema{
				Type:     schema.TypeList,
				Required: false,
				Optional: true,
			},
			"attribute": &schema.Schema{
				Type:     schema.TypeMap,
				Required: false,
				Optional: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
			},
			"created": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"updated": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"changedBy": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
		},
	}
}

func readItem(d *schema.ResourceData, m interface{}) error {
	return nil
}

/*
	ITEM TYPE DATA SOURCE
 */

func ItemTypeDataSource() *schema.Resource {
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

func readItemType(d *schema.ResourceData, m interface{}) error {
	return nil
}

/*
	LINK DATA SOURCE
 */

func LinkDataSource() *schema.Resource {
	return &schema.Resource{
		Read: readLink,

		Schema: map[string]*schema.Schema{
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func readLink(d *schema.ResourceData, m interface{}) error {
	return nil
}

/*
	LINK TYPE DATA SOURCE
 */

func LinkTypeDataSource() *schema.Resource {
	return &schema.Resource{
		Read: readLinkType,

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

func readLinkType(d *schema.ResourceData, m interface{}) error {
	return nil
}

/*
	LINK RULE DATA SOURCE
 */

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

func readLinkRule(d *schema.ResourceData, m interface{}) error {
	return nil
}

/*
	MODEL DATA SOURCE
 */

func ModelDataSource() *schema.Resource {
	return &schema.Resource{
		Read: readModel,

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

func readModel(d *schema.ResourceData, m interface{}) error {
	return nil
}
