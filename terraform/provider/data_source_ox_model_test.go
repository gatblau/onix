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
	ModelDataSourceName = "data.ox_model.ox_model_1_data"
	ModelDsKey          = "test_acc_ox_model_1_data"
	ModelDsName         = "ox model 1 data"
	ModelDsDesc         = "ox_model_1_data Description"
	ModelDsPartition    = "REF"
)

func init() {
	// defines a sweeper to clean up dangling test resources
	resource.AddTestSweepers("ModelDataSource", &resource.Sweeper{
		Name: "ModelDataSource",
		F: func(region string) error {
			_, err := cfg.Client.DeleteModel(&oxc.Model{Key: ModelDsKey})
			return err
		},
	})
}

func TestModelDataSource(t *testing.T) {
	resourceName := ModelDataSourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { prepareModelDataSourceTest(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// read
			{
				Config: oxModelDataSource(),
				Check: resource.ComposeTestCheckFunc(
					// check item resource attributes in Terraform state
					resource.TestCheckResourceAttr(resourceName, "key", ModelDsKey),
					resource.TestCheckResourceAttr(resourceName, "name", ModelDsName),
					resource.TestCheckResourceAttr(resourceName, "description", ModelDsDesc),
					resource.TestCheckResourceAttr(resourceName, "partition", ModelDsPartition),
				),
			},
		},
	})
}

func oxModelDataSource() string {
	return fmt.Sprintf(
		`data "ox_model" "ox_model_1_data" {
					key = "%s"	
				}`, ModelDsKey)
}

// create supporting database entities
func prepareModelDataSourceTest(t *testing.T) {
	// need a model first
	_, err := cfg.Client.PutModel(
		&oxc.Model{
			Key:         ModelDsKey,
			Name:        ModelDsName,
			Description: ModelDsDesc,
		})
	if err != nil {
		t.Error(err)
	}
}
