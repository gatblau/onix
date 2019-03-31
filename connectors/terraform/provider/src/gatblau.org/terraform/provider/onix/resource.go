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

import (
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
)

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
				Required: false,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeInt,
				Required: false,
			},
			"meta": &schema.Schema{
				Type:     schema.TypeMap,
				Required: false,
			},
			"tag": &schema.Schema{
				Type:     schema.TypeList,
				Required: false,
			},
			"attribute": &schema.Schema{
				Type:     schema.TypeMap,
				Required: false,
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

/*
	ITEM TYPE RESOURCE
 */
func ItemTypeResource() *schema.Resource {
	return &schema.Resource{
		Create: createItemType,
		Read:   readItemType,
		Update: updateItemType,
		Delete: deleteItemType,
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

func createItemType(d *schema.ResourceData, m interface{}) error {
	return nil
}

func updateItemType(d *schema.ResourceData, m interface{}) error {
	return nil
}

func deleteItemType(d *schema.ResourceData, m interface{}) error {
	return nil
}

/*
	LINK RESOURCE
 */

func LinkResource() *schema.Resource {
	return &schema.Resource{
		Create: createLink,
		Read:   readLink,
		Update: updateLink,
		Delete: deleteLink,
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

func createLink(d *schema.ResourceData, m interface{}) error {
	return nil
}

func updateLink(d *schema.ResourceData, m interface{}) error {
	return nil
}

func deleteLink(d *schema.ResourceData, m interface{}) error {
	return nil
}

/*
   LINK TYPE RESOURCE
*/

func LinkTypeResource() *schema.Resource {
	return &schema.Resource{
		Create: createLinkType,
		Read:   readLinkType,
		Update: updateLinkType,
		Delete: deleteLinkType,
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

func createLinkType(d *schema.ResourceData, m interface{}) error {
	return nil
}

func updateLinkType(d *schema.ResourceData, m interface{}) error {
	return nil
}

func deleteLinkType(d *schema.ResourceData, m interface{}) error {
	return nil
}

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
		},
	}
}

func createLinkRule(d *schema.ResourceData, m interface{}) error {
	return nil
}

func updateLinkRule(d *schema.ResourceData, m interface{}) error {
	return nil
}

func deleteLinkRule(d *schema.ResourceData, m interface{}) error {
	return nil
}

/*
   MODEL RESOURCE
*/

func ModelResource() *schema.Resource {
	return &schema.Resource{
		Create: createModel,
		Read:   readModel,
		Update: updateModel,
		Delete: deleteModel,
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

func createModel(d *schema.ResourceData, m interface{}) error {
	return nil
}

func updateModel(d *schema.ResourceData, m interface{}) error {
	return nil
}

func deleteModel(d *schema.ResourceData, m interface{}) error {
	return nil
}

func check(result *Result, err error) (error) {
	if err != nil {
		return err
	} else if result.Error {
		return errors.New(result.Message)
	} else {
		return nil
	}
}
