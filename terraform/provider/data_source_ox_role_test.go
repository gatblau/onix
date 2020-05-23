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
	RoleDataSourceName = "data.ox_role.ox_role_1_data"
	RoleDsKey          = "test_acc_ox_role_1_data"
	RoleDsName         = "ox_role_1_data name"
	RoleDsDesc         = "ox_role_1_data description"
)

func init() {
	// defines a sweeper to clean up dangling test resources
	resource.AddTestSweepers("RoleDataSource", &resource.Sweeper{
		Name: "RoleDataSource",
		F: func(region string) error {
			_, err := cfg.Client.DeleteRole(&oxc.Role{Key: RoleDsKey})
			return err
		},
	})
}

func TestRoleDataSource(t *testing.T) {
	resourceName := RoleDataSourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { prepareRoleDataSourceTest(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// read
			{
				Config: oxRoleDataSource(),
				Check: resource.ComposeTestCheckFunc(
					// check item resource attributes in Terraform state
					resource.TestCheckResourceAttr(resourceName, "key", RoleDsKey),
					resource.TestCheckResourceAttr(resourceName, "name", RoleDsName),
					resource.TestCheckResourceAttr(resourceName, "description", RoleDsDesc),
				),
			},
		},
	})
}

func oxRoleDataSource() string {
	return fmt.Sprintf(
		`data "ox_role" "ox_role_1_data" {
					key = "%s"	
				}`, RoleDsKey)
}

// create supporting database entities
func prepareRoleDataSourceTest(t *testing.T) {
	_, err := cfg.Client.PutRole(
		&oxc.Role{
			Key:         RoleDsKey,
			Name:        RoleDsName,
			Description: RoleDsDesc,
		})
	if err != nil {
		t.Error(err)
	}
}
