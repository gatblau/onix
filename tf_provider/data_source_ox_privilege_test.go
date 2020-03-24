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
	"strconv"
	"testing"
)

const (
	PrivilegeDataSourceName = "data.ox_privilege.ox_privilege_1_data"
	PrivilegeDsKey          = "test_acc_ox_privilege_1_data"
	PrivilegeDsRole         = "ADMIN"
	PrivilegeDsPartition    = "INS"
	PrivilegeDsCanCreate    = true
	PrivilegeDsCanRead      = true
	PrivilegeDsCanDelete    = false
)

func init() {
	// defines a sweeper to clean up dangling test resources
	resource.AddTestSweepers("PrivilegeDataSource", &resource.Sweeper{
		Name: "PrivilegeDataSource",
		F: func(region string) error {
			_, err := cfg.Client.DeletePrivilege(&oxc.Privilege{Key: PrivilegeDsKey})
			return err
		},
	})
}

func TestPrivilegeDataSource(t *testing.T) {
	resourceName := PrivilegeDataSourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { preparePrivilegeDataSourceTest(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// read
			{
				Config: oxPrivilegeDataSource(),
				Check: resource.ComposeTestCheckFunc(
					// check item resource attributes in Terraform state
					resource.TestCheckResourceAttr(resourceName, "key", PrivilegeDsKey),
					resource.TestCheckResourceAttr(resourceName, "role", PrivilegeDsRole),
					resource.TestCheckResourceAttr(resourceName, "partition", PrivilegeDsPartition),
					resource.TestCheckResourceAttr(resourceName, "can_read", strconv.FormatBool(PrivilegeDsCanRead)),
					resource.TestCheckResourceAttr(resourceName, "can_create", strconv.FormatBool(PrivilegeDsCanCreate)),
					resource.TestCheckResourceAttr(resourceName, "can_delete", strconv.FormatBool(PrivilegeDsCanDelete)),
				),
			},
		},
	})
}

func oxPrivilegeDataSource() string {
	return fmt.Sprintf(
		`data "ox_privilege" "ox_privilege_1_data" {
					key = "%s"	
				}`, PrivilegeDsKey)
}

// create supporting database entities
func preparePrivilegeDataSourceTest(t *testing.T) {
	// use the default partition INS and default role ADMIN
	_, err := cfg.Client.PutPrivilege(
		&oxc.Privilege{
			Key:       PrivilegeDsKey,
			Role:      PrivilegeDsRole,
			Partition: PrivilegeDsPartition,
			Create:    PrivilegeDsCanCreate,
			Read:      PrivilegeDsCanRead,
			Delete:    PrivilegeDsCanDelete,
		})
	if err != nil {
		t.Error(err)
	}
}
