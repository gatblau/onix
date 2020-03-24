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
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

// all the constants for the item test
const (
	PartitionResourceName = "ox_partition.ox_partition_1"
	PartitionRsKey        = "test_acc_ox_partition_1"
	PartitionRsName       = "ox_partition name"
	PartitionRsDesc       = "ox_partition description"
)

func TestPartitionResource(t *testing.T) {
	resourceName := PartitionResourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { preparePartitionResourceTest(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			//create
			{
				Config: oxPartitionResource(),
				Check: resource.ComposeTestCheckFunc(
					// check item resource attributes in Terraform state
					resource.TestCheckResourceAttr(resourceName, "key", PartitionRsKey),
					resource.TestCheckResourceAttr(resourceName, "name", PartitionRsName),
					resource.TestCheckResourceAttr(resourceName, "description", PartitionRsDesc),

					// check for side effects in Onix database
					checkPartitionResourceCreated(PartitionResourceName),
				),
			},
		},
		CheckDestroy: checkPartitionResourceDestroyed,
	})
}

// create supporting database entities
func preparePartitionResourceTest(t *testing.T) {
}

// check the item has been created in the database
func checkPartitionResourceCreated(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		// retrieve the partition from the database
		key := rs.Primary.Attributes["key"]
		partition, err := cfg.Client.GetPartition(&oxc.Partition{Key: key})
		if err != nil {
			return fmt.Errorf("can't read partition %s", key)
		}

		// compares the state with the database values
		if err := checkEntityAttr(rs, "name", partition.Name); err != nil {
			return err
		}
		if err := checkEntityAttr(rs, "description", partition.Description); err != nil {
			return err
		}
		if err := checkEntityAttr(rs, "owner", partition.Owner); err != nil {
			return err
		}

		return nil
	}
}

// check the item has been destroyed destroyed in the database
func checkPartitionResourceDestroyed(s *terraform.State) error {
	partition, _ := cfg.Client.GetPartition(&oxc.Partition{Key: PartitionRsKey})
	// should get an error as partition should not exist anymore
	if partition != nil {
		return fmt.Errorf("partition %s still exists", PartitionRsKey)
	}
	// clean up the database after the test
	return cleanUpPartitionRs()
}

// remove supporting database entities
func cleanUpPartitionRs() error {
	return nil
}

func oxPartitionResource() string {
	return fmt.Sprintf(
		`resource "ox_partition" "ox_partition_1" {
  key         = "%s"
  name        = "%s"
  description = "%s"
}`, PartitionRsKey, PartitionRsName, PartitionRsDesc)
}
