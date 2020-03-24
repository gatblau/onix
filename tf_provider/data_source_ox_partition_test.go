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
	"fmt"
	"github.com/gatblau/oxc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const (
	PartitionDataSourceName = "data.ox_partition.ox_partition_1_data"
	PartitionDsKey          = "test_acc_ox_partition_1_data"
	PartitionDsName         = "ox_partition_1_data name"
	PartitionDsDesc         = "ox_partition_1_data description"
)

func init() {
	// defines a sweeper to clean up dangling test resources
	resource.AddTestSweepers("PartitionDataSource", &resource.Sweeper{
		Name: "PartitionDataSource",
		F: func(region string) error {
			_, err := cfg.Client.DeletePartition(&oxc.Partition{Key: PartitionDsKey})
			return err
		},
	})
}

func TestPartitionDataSource(t *testing.T) {
	resourceName := PartitionDataSourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { preparePartitionDataSourceTest(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// read
			{
				Config: oxPartitionDataSource(),
				Check: resource.ComposeTestCheckFunc(
					// check item resource attributes in Terraform state
					resource.TestCheckResourceAttr(resourceName, "key", PartitionDsKey),
					resource.TestCheckResourceAttr(resourceName, "name", PartitionDsName),
					resource.TestCheckResourceAttr(resourceName, "description", PartitionDsDesc),
				),
			},
		},
	})
}

func oxPartitionDataSource() string {
	return fmt.Sprintf(
		`data "ox_partition" "ox_partition_1_data" {
					key = "%s"	
				}`, PartitionDsKey)
}

// create supporting database entities
func preparePartitionDataSourceTest(t *testing.T) {
	_, err := cfg.Client.PutPartition(
		&oxc.Partition{
			Key:         PartitionDsKey,
			Name:        PartitionDsName,
			Description: PartitionDsDesc,
		})
	if err != nil {
		t.Error(err)
	}
}
