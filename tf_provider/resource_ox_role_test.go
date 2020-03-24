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
	"strconv"
	"testing"
)

// all the constants for the item test
const (
	RoleResourceName = "ox_role.ox_role_1"
	RoleRsKey        = "test_acc_ox_role_1"
	RoleRsName       = "ox_role name"
	RoleRsDesc       = "ox_role description"
	RoleRsLevel      = 1
)

func TestRoleResource(t *testing.T) {
	resourceName := RoleResourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { prepareRoleResourceTest(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			//create
			{
				Config: oxRoleResource(),
				Check: resource.ComposeTestCheckFunc(
					// check item resource attributes in Terraform state
					resource.TestCheckResourceAttr(resourceName, "key", RoleRsKey),
					resource.TestCheckResourceAttr(resourceName, "name", RoleRsName),
					resource.TestCheckResourceAttr(resourceName, "description", RoleRsDesc),
					resource.TestCheckResourceAttr(resourceName, "level", strconv.Itoa(RoleRsLevel)),

					// check for side effects in Onix database
					checkRoleResourceCreated(RoleResourceName),
				),
			},
		},
		CheckDestroy: checkRoleResourceDestroyed,
	})
}

// create supporting database entities
func prepareRoleResourceTest(t *testing.T) {
}

// check the role has been created in the database
func checkRoleResourceCreated(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		// retrieve the partition from the database
		key := rs.Primary.Attributes["key"]
		partition, err := cfg.Client.GetRole(&oxc.Role{Key: key})
		if err != nil {
			return fmt.Errorf("can't read role %s", key)
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
		if err := checkEntityAttr(rs, "level", strconv.Itoa(partition.Level)); err != nil {
			return err
		}

		return nil
	}
}

// check the role has been destroyed destroyed in the database
func checkRoleResourceDestroyed(s *terraform.State) error {
	role, _ := cfg.Client.GetRole(&oxc.Role{Key: RoleRsKey})
	// should get an error as role should not exist anymore
	if role != nil {
		return fmt.Errorf("role %s still exists", RoleRsKey)
	}
	// clean up the database after the test
	return cleanUpRoleRs()
}

// remove supporting database entities
func cleanUpRoleRs() error {
	return nil
}

func oxRoleResource() string {
	return fmt.Sprintf(
		`resource "ox_role" "ox_role_1" {
  key         = "%s"
  name        = "%s"
  description = "%s"
  level       = %s
}`, RoleRsKey, RoleRsName, RoleRsDesc, strconv.Itoa(RoleRsLevel))
}
