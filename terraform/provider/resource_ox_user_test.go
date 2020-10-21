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
	UserResourceName = "ox_user.ox_user_1"
	UserRsKey        = "test_acc_ox_user_1"
	UserRsName       = "ox_user name"
	UserRsEmail      = "ox_user@email.com"
	UserRsService    = false
	UserRsPwd        = "ckmeCbo3rybvobvob3vr3ovb!="
)

func TestUserResource(t *testing.T) {
	resourceName := UserResourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { prepareUserResourceTest(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// create
			{
				Config: oxUserResource(),
				Check: resource.ComposeTestCheckFunc(
					// check item resource attributes in Terraform state
					resource.TestCheckResourceAttr(resourceName, "key", UserRsKey),
					resource.TestCheckResourceAttr(resourceName, "name", UserRsName),
					resource.TestCheckResourceAttr(resourceName, "email", UserRsEmail),
					resource.TestCheckResourceAttr(resourceName, "service", strconv.FormatBool(UserRsService)),
					// resource.TestCheckResourceAttr(resourceName, "pwd", UserRsPwd),

					// check for side effects in Onix database
					checkUserResourceCreated(UserResourceName),
				),
			},
		},
		CheckDestroy: checkUserResourceDestroyed,
	})
}

// create supporting database entities
func prepareUserResourceTest(t *testing.T) {
}

// check the user has been created in the database
func checkUserResourceCreated(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		// retrieve the user from the database
		key := rs.Primary.Attributes["key"]
		user, err := cfg.Client.GetUser(&oxc.User{Key: key})
		if err != nil {
			return fmt.Errorf("can't read user %s", key)
		}

		// compares the state with the database values
		if err := checkEntityAttr(rs, "name", user.Name); err != nil {
			return err
		}
		if err := checkEntityAttr(rs, "email", user.Email); err != nil {
			return err
		}
		if err := checkEntityAttr(rs, "pwd", user.Pwd); err != nil {
			return err
		}
		if err := checkEntityAttr(rs, "service", strconv.FormatBool(UserRsService)); err != nil {
			return err
		}

		return nil
	}
}

// check the user has been destroyed destroyed in the database
func checkUserResourceDestroyed(s *terraform.State) error {
	user, _ := cfg.Client.GetUser(&oxc.User{Key: UserRsKey})
	// should get an error as user should not exist anymore
	if user != nil {
		return fmt.Errorf("user %s still exists", UserRsKey)
	}
	// clean up the database after the test
	return cleanUpUserRs()
}

// remove supporting database entities
func cleanUpUserRs() error {
	return nil
}

func oxUserResource() string {
	return fmt.Sprintf(
		`resource "ox_user" "ox_user_1" {
  key         = "%s"
  name        = "%s"
  email		  = "%s"
  service     = %t
  pwd         = "%s"
}`, UserRsKey, UserRsName, UserRsEmail, UserRsService, UserRsPwd)
}
