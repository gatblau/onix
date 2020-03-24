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
	PrivilegeResourceName = "ox_privilege.ox_privilege_1"
	PrivilegeRsKey        = " "
	PrivilegeRsRole       = "ADMIN"
	PrivilegeRsPartition  = "INS"
	PrivilegeRsCanCreate  = true
	PrivilegeRsCanRead    = true
	PrivilegeRsCanDelete  = false
)

func TestPrivilegeResource(t *testing.T) {
	resourceName := PrivilegeResourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { preparePrivilegeResourceTest(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			//create
			{
				Config: oxPrivilegeResource(),
				Check: resource.ComposeTestCheckFunc(
					// check privilege resource attributes in Terraform state
					resource.TestCheckResourceAttr(resourceName, "key", PrivilegeRsKey),
					resource.TestCheckResourceAttr(resourceName, "role", PrivilegeRsRole),
					resource.TestCheckResourceAttr(resourceName, "partition", PrivilegeRsPartition),
					resource.TestCheckResourceAttr(resourceName, "can_read", strconv.FormatBool(PrivilegeRsCanRead)),
					resource.TestCheckResourceAttr(resourceName, "can_create", strconv.FormatBool(PrivilegeRsCanCreate)),
					resource.TestCheckResourceAttr(resourceName, "can_delete", strconv.FormatBool(PrivilegeRsCanDelete)),

					// check for side effects in Onix database
					checkPrivilegeResourceCreated(PrivilegeResourceName),
				),
			},
		},
		CheckDestroy: checkPrivilegeResourceDestroyed,
	})
}

// create supporting database entities
func preparePrivilegeResourceTest(t *testing.T) {
}

// check the item has been created in the database
func checkPrivilegeResourceCreated(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		// retrieve the privilege from the database
		key := rs.Primary.Attributes["key"]
		privilege, err := cfg.Client.GetPrivilege(&oxc.Privilege{Key: key})
		if err != nil {
			return fmt.Errorf("can't read privilege %s", key)
		}

		// compares the state with the database values
		if err := checkEntityAttr(rs, "role", privilege.Role); err != nil {
			return err
		}
		if err := checkEntityAttr(rs, "partition", privilege.Partition); err != nil {
			return err
		}
		if err := checkEntityAttr(rs, "can_create", strconv.FormatBool(privilege.Create)); err != nil {
			return err
		}
		if err := checkEntityAttr(rs, "can_read", strconv.FormatBool(privilege.Read)); err != nil {
			return err
		}
		if err := checkEntityAttr(rs, "can_delete", strconv.FormatBool(privilege.Delete)); err != nil {
			return err
		}

		return nil
	}
}

// check the item has been destroyed destroyed in the database
func checkPrivilegeResourceDestroyed(s *terraform.State) error {
	partition, _ := cfg.Client.GetPrivilege(&oxc.Privilege{Key: PrivilegeRsKey})
	// should get an error as privilege should not exist anymore
	if partition != nil {
		return fmt.Errorf("privilege %s still exists", PrivilegeRsKey)
	}
	// clean up the database after the test
	return cleanUpPrivilegeRs()
}

// remove supporting database entities
func cleanUpPrivilegeRs() error {
	return nil
}

func oxPrivilegeResource() string {
	return fmt.Sprintf(
		`resource "ox_privilege" "ox_privilege_1" {
  key         = "%s"
  role        = "%s"
  partition   = "%s"
  can_create  = %s
  can_read    = %s
  can_delete  = %s
}`, PrivilegeRsKey,
		PrivilegeRsRole,
		PrivilegeRsPartition,
		strconv.FormatBool(PrivilegeRsCanCreate),
		strconv.FormatBool(PrivilegeRsCanRead),
		strconv.FormatBool(PrivilegeRsCanDelete))
}
