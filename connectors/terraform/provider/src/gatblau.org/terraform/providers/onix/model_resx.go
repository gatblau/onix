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
	"github.com/hashicorp/terraform/helper/schema"
)

func modelResx() *schema.Resource {
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

func readModel(d *schema.ResourceData, m interface{}) error {
	return nil
}

func updateModel(d *schema.ResourceData, m interface{}) error {
	return nil
}

func deleteModel(d *schema.ResourceData, m interface{}) error {
	return nil
}
