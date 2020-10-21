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
	UserDataSourceName = "data.ox_user.ox_user_1_data"
	UserDsKey          = "test_acc_ox_user_1_data"
	UserDsName         = "ox_user_1_data name"
	UserDsEmail        = "ox_user_1_data@email.com"
)

func init() {
	// defines a sweeper to clean up dangling test resources
	resource.AddTestSweepers("UserDataSource", &resource.Sweeper{
		Name: "UserDataSource",
		F: func(region string) error {
			_, err := cfg.Client.DeleteUser(&oxc.User{Key: UserDsKey})
			return err
		},
	})
}

func TestUserDataSource(t *testing.T) {
	resourceName := UserDataSourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { prepareUserDataSourceTest(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// read
			{
				Config: oxUserDataSource(),
				Check: resource.ComposeTestCheckFunc(
					// check item resource attributes in Terraform state
					resource.TestCheckResourceAttr(resourceName, "key", UserDsKey),
					resource.TestCheckResourceAttr(resourceName, "name", UserDsName),
					resource.TestCheckResourceAttr(resourceName, "email", UserDsEmail),
				),
			},
		},
	})
}

func oxUserDataSource() string {
	return fmt.Sprintf(
		`data "ox_user" "ox_user_1_data" {
					key = "%s"	
				}`, UserDsKey)
}

// create supporting database entities
func prepareUserDataSourceTest(t *testing.T) {
	_, err := cfg.Client.PutUser(
		&oxc.User{
			Key:   UserDsKey,
			Name:  UserDsName,
			Email: UserDsEmail,
		}, false)
	if err != nil {
		t.Error(err)
	}
}
