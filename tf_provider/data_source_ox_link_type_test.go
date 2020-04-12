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
	LinkTypeDatasourceName = "data.ox_link_type.ox_link_type_1"
	LinkTypeDsModelKey     = "test_acc_ox_link_type_model_1_data"
	LinkTypeDsModelName    = "ox_link_type_model_1_data"
	LinkTypeDsModelDesc    = "ox_link_type_model_1_data Description"
	LinkTypeDsKey          = "test_acc_ox_link_type_1_data"
	LinkTypeDsName         = "ox_link_type_1_data_name"
	LinkTypeDsDesc         = "ox_link_type_1_data_description"
	LinkTypeDsEncryptMeta  = false
	LinkTypeDsEncryptTxt   = false
)

var LinkTypeDsStyle = make(map[string]interface{})

func init() {
	// defines a sweeper to clean up dangling test resources
	resource.AddTestSweepers("LinkTypeDataSource", &resource.Sweeper{
		Name: "LinkTypeDataSource",
		F: func(region string) error {
			_, err := cfg.Client.DeleteLinkType(&oxc.LinkType{Key: LinkTypeDsKey})
			_, err = cfg.Client.DeleteModel(&oxc.Model{Key: LinkTypeDsModelKey})
			return err
		},
	})
}

func TestLinkTypeDataSource(t *testing.T) {
	resourceName := LinkTypeDatasourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { prepareLinkTypeDataSourceTest(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// read
			{
				Config: oxLinkTypeDataSource(),
				Check: resource.ComposeTestCheckFunc(
					// check item type attributes in Terraform state
					resource.TestCheckResourceAttr(resourceName, "key", LinkTypeDsKey),
					resource.TestCheckResourceAttr(resourceName, "name", LinkTypeDsName),
					resource.TestCheckResourceAttr(resourceName, "description", LinkTypeDsDesc),
					resource.TestCheckResourceAttr(resourceName, "model_key", LinkTypeDsModelKey),
					resource.TestCheckResourceAttr(resourceName, "encrypt_meta", strconv.FormatBool(LinkTypeDsEncryptMeta)),
					resource.TestCheckResourceAttr(resourceName, "encrypt_txt", strconv.FormatBool(LinkTypeDsEncryptTxt)),
				),
			},
		},
	})
}

func oxLinkTypeDataSource() string {
	return fmt.Sprintf(
		`data "ox_link_type" "ox_link_type_1" {
					key = "%s"	
				}`, LinkTypeDsKey)
}

// create supporting database entities
func prepareLinkTypeDataSourceTest(t *testing.T) {
	// need a model first
	result, err := cfg.Client.PutModel(
		&oxc.Model{
			Key:         LinkTypeDsModelKey,
			Name:        LinkTypeDsModelName,
			Description: LinkTypeDsModelDesc,
		})
	if err != nil {
		t.Error(err)
	}
	if result.Error {
		t.Error(result.Message)
	}

	// an item type in the model
	result, err = cfg.Client.PutLinkType(
		&oxc.LinkType{
			Key:         LinkTypeDsKey,
			Name:        LinkTypeDsName,
			Description: LinkTypeDsDesc,
			Model:       LinkTypeDsModelKey,
			Style:       LinkTypeDsStyle,
		})
	if err != nil {
		t.Error(err)
	}
	if result.Error {
		t.Error(result.Message)
	}
}
