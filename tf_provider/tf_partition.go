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
		},
	}
}

func PartitionDataSource() *schema.Resource {
	return &schema.Resource{
		Read: readPartition,

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
			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"created": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"changedby": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func createPartition(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into a Partition
	partition := newPartition(data)

	// put the Partition to the Web API
	err := partition.put(meta)

	// set Item Id key
	data.SetId(partition.Key)

	return err
}

func readPartition(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into a Partition
	partition := newPartition(data)

	// get the resource
	partition, err := partition.get(meta)

	return err
}

func updatePartition(data *schema.ResourceData, meta interface{}) error {
	// same as create - Web PI is idempotent
	return createPartition(data, meta)
}

func deletePartition(data *schema.ResourceData, meta interface{}) error {
	// read the resource data into a Partition
	partition := newPartition(data)

	// delete the partition
	return partition.delete(meta)
}
