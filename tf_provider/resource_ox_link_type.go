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
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"meta_schema": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"model_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"tag": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"encrypt_meta": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"encrypt_txt": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"managed": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func createLinkType(data *schema.ResourceData, meta interface{}) error {
	// get the Ox client
	c := meta.(Config).Client

	// read the tf data into an Item
	linkType := newLinkType(data)

	// put the Link Type to the Web API
	err := err(c.PutLinkType(linkType))
	if err != nil {
		return err
	}

	// set Link Type Id key
	data.SetId(linkType.Key)

	return readLinkType(data, meta)
}

func updateLinkType(data *schema.ResourceData, meta interface{}) error {
	return createLinkType(data, meta)
}

func deleteLinkType(data *schema.ResourceData, meta interface{}) error {
	// get the Ox client
	c := meta.(Config).Client

	// read the resource data into an Item
	linkType := newLinkType(data)

	// delete the linkType
	return err(c.DeleteLinkType(linkType))
}

func newLinkType(data *schema.ResourceData) *LinkType {
	return &LinkType{
		Key:         data.Get("key").(string),
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
		Model:       data.Get("model_key").(string),
		MetaSchema:  data.Get("meta_schema").(map[string]interface{}),
		EncryptMeta: data.Get("encrypt_meta").(bool),
		EncryptTxt:  data.Get("encrypt_txt").(bool),
		Managed:     data.Get("managed").(bool),
		Tag:         data.Get("tag").([]interface{}),
		Version:     getVersion(data),
	}
}

// populate the LinkType with the data in the terraform resource
func populateLinkType(data *schema.ResourceData, linkType *LinkType) {
	data.SetId(linkType.Key)
	data.Set("key", linkType.Key)
	data.Set("name", linkType.Name)
	data.Set("description", linkType.Description)
	data.Set("meta_schema", linkType.MetaSchema)
	data.Set("model_key", linkType.Model)
	data.Set("encrypt_txt", linkType.EncryptTxt)
	data.Set("encrypt_meta", linkType.EncryptMeta)
	data.Set("tag", linkType.Tag)
	data.Set("managed", linkType.Managed)
	data.Set("created", linkType.Created)
	data.Set("updated", linkType.Updated)
	// TODO: enable below after upgrading client
	// data.Set("changed_by", linkType.ChangedBy)
}
