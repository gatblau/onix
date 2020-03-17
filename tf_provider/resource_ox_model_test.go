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

// all the constants for the model test
const (
	ModelResourceName = "ox_model.ox_model_1"
	ModelRsKey        = "test_acc_ox_model_1"
	ModelRsName       = "ox model"
	ModelRsDesc       = "ox_model Description"
	ModelRsPartition  = "REF"
	ModelRsManaged    = "true"
)

func TestModelResource(t *testing.T) {
	resourceName := ModelResourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { prepareModelResourceTest(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			//create
			{
				Config: oxModelResource(),
				Check: resource.ComposeTestCheckFunc(
					// check model resource attributes in Terraform state
					resource.TestCheckResourceAttr(resourceName, "key", ModelRsKey),
					resource.TestCheckResourceAttr(resourceName, "name", ModelRsName),
					resource.TestCheckResourceAttr(resourceName, "description", ModelRsDesc),
				),
			},
		},
		CheckDestroy: checkModelResourceDestroyed,
	})
}

// create supporting database entities
func prepareModelResourceTest(t *testing.T) {
	// nothing to do
}

// check the model has been destroyed destroyed in the database
func checkModelResourceDestroyed(s *terraform.State) error {
	model, _ := cfg.Client.GetModel(&oxc.Model{Key: ModelRsKey})
	// should get an error as model should not exist anymore
	if model != nil {
		return fmt.Errorf("Model %s still exists", ModelRsKey)
	}
	// clean up the database after the test
	return cleanUpModelRs()
}

// remove supporting database entities
func cleanUpModelRs() error {
	_, err := cfg.Client.DeleteModel(&oxc.Model{Key: ModelRsKey})
	return err
}

func oxModelResource() string {
	return fmt.Sprintf(
		`resource "ox_model" "ox_model_1" {
  key         = "%s"
  name        = "%s"
  description = "%s"
  partition   = "%s"
  managed     = %s
}`, ModelRsKey, ModelRsName, ModelRsDesc, ModelRsPartition, ModelRsManaged)
}
