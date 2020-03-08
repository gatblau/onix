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

// terraform resource for an Onix Model
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
			"partition": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"managed": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

// terraform data source for an Onix Model
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

func createModel(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into a Model
	model := newModel(data)

	// put the Model to the Web API
	err := model.put(meta)

	// set Model Id key
	data.SetId(model.Key)

	return err
}

func readModel(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into a Model
	model := newModel(data)

	// get the resource
	model, err := model.get(meta)

	return err
}

func updateModel(data *schema.ResourceData, meta interface{}) error {
	return createModel(data, meta)
}

func deleteModel(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into an Model
	model := newModel(data)

	// delete the model
	return model.delete(meta)
}
