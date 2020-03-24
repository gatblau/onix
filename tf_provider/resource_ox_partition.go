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

func PartitionResource() *schema.Resource {
	return &schema.Resource{
		Create: createPartition,
		Read:   readPartition,
		Update: updatePartition,
		Delete: deletePartition,
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
				Optional: true,
			},
			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"created": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"updated": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"changed_by": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func createPartition(data *schema.ResourceData, meta interface{}) error {
	// get the Ox client
	c := meta.(Config).Client

	// read the tf data into a Partition
	partition := newPartition(data)

	// put the Partition to the Web API
	err := err(c.PutPartition(partition))
	if err != nil {
		return err
	}

	// set Partition Id key
	data.SetId(partition.Key)

	return readPartition(data, meta)
}

func updatePartition(data *schema.ResourceData, meta interface{}) error {
	// same as create - Web PI is idempotent
	return createPartition(data, meta)
}

func deletePartition(data *schema.ResourceData, meta interface{}) error {
	// get the Ox client
	c := meta.(Config).Client

	// read the resource data into a Partition
	partition := newPartition(data)

	// delete the partition
	return err(c.DeletePartition(partition))
}

// create a new Partition from a terraform resource
func newPartition(data *schema.ResourceData) *Partition {
	return &Partition{
		Key:         data.Get("key").(string),
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
		Version:     getVersion(data),
	}
}

// populate the Partition with the data in the terraform resource
func populatePartition(data *schema.ResourceData, partition *Partition) {
	data.SetId(partition.Key)
	data.Set("key", partition.Key)
	data.Set("name", partition.Name)
	data.Set("description", partition.Description)
	data.Set("owner", partition.Owner)
	data.Set("created", partition.Created)
	data.Set("updated", partition.Updated)
	data.Set("changed_by", partition.ChangedBy)
}
