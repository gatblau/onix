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

// all the constants for the link test
const (
	LinkTypeAttrResourceName   = "ox_link_type_attr.ox_link_type_attr_1"
	LinkTypeAttrRsKey          = "test_acc_ox_link_type_attr_1"
	LinkTypeAttrRsName         = "test_acc_ox_link_type_attr_1 name"
	LinkTypeAttrRsDesc         = "test_acc_ox_link_type_attr_1 description"
	LinkTypeAttrRsType         = "list"
	LinkTypeAttrRsDefValue     = "A,B,C"
	LinkTypeAttrRsLinkTypeKey  = "test_acc_link_type_ox_link_type_attr"
	LinkTypeAttrRsLinkTypeName = "test_acc_link_type_ox_link_type_attr name"
	LinkTypeAttrRsLinkTypeDesc = "test_acc_link_type_ox_link_type_attr description"
	LinkTypeAttrRsModelKey     = "test_acc_model_ox_link_type_attr"
	LinkTypeAttrRsModelName    = "test_acc_model_ox_link_type_attr name"
	LinkTypeAttrRsModelDesc    = "test_acc_model_ox_link_type_attr description"
)

func TestLinkTypeAttrResource(t *testing.T) {
	resourceName := LinkTypeAttrResourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { prepareLinkTypeAttrResourceTest(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			//create
			{
				Config: oxLinkTypeAttrResource(),
				Check: resource.ComposeTestCheckFunc(
					// check link resource attributes in Terraform state
					resource.TestCheckResourceAttr(resourceName, "key", LinkTypeAttrRsKey),
					resource.TestCheckResourceAttr(resourceName, "name", LinkTypeAttrRsName),
					resource.TestCheckResourceAttr(resourceName, "description", LinkTypeAttrRsDesc),
					resource.TestCheckResourceAttr(resourceName, "link_type_key", LinkTypeAttrRsLinkTypeKey),
					resource.TestCheckResourceAttr(resourceName, "type", LinkTypeAttrRsType),
					resource.TestCheckResourceAttr(resourceName, "def_value", LinkTypeAttrRsDefValue),

					// check for side effects in Onix database
					checkLinkTypeAttrResourceCreated(LinkTypeAttrResourceName),
				),
			},
		},
		CheckDestroy: checkLinkTypeAttrResourceDestroyed,
	})
}

// create supporting database entities
func prepareLinkTypeAttrResourceTest(t *testing.T) {
	// need a model first
	_, err := cfg.Client.PutModel(
		&oxc.Model{
			Key:         LinkTypeAttrRsModelKey,
			Name:        LinkTypeAttrRsModelName,
			Description: LinkTypeAttrRsModelDesc,
		})
	if err != nil {
		t.Error(err)
	}

	// an link type in the model
	_, err = cfg.Client.PutLinkType(
		&oxc.LinkType{
			Key:         LinkTypeAttrRsLinkTypeKey,
			Name:        LinkTypeAttrRsLinkTypeName,
			Description: LinkTypeAttrRsLinkTypeDesc,
			Model:       LinkTypeAttrRsModelKey,
		})
	if err != nil {
		t.Error(err)
	}
}

// check the link has been created in the database
func checkLinkTypeAttrResourceCreated(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		// retrieve the attr from the database
		key := rs.Primary.Attributes["key"]
		attr, err := cfg.Client.GetLinkTypeAttr(&oxc.LinkTypeAttribute{Key: key, LinkTypeKey: LinkTypeAttrRsLinkTypeKey})
		if err != nil {
			return fmt.Errorf("can't read attr type attribute %s", key)
		}

		// compares the state with the database values
		if err := checkEntityAttr(rs, "name", attr.Name); err != nil {
			return err
		}
		if err := checkEntityAttr(rs, "description", attr.Description); err != nil {
			return err
		}
		if err := checkEntityAttr(rs, "type", attr.Type); err != nil {
			return err
		}

		return nil
	}
}

// check the link has been destroyed destroyed in the database
func checkLinkTypeAttrResourceDestroyed(s *terraform.State) error {
	link, _ := cfg.Client.GetLinkTypeAttr(&oxc.LinkTypeAttribute{Key: LinkTypeAttrRsKey, LinkTypeKey: LinkTypeAttrRsLinkTypeKey})
	// should get an error as link should not exist anymore
	if link != nil {
		return fmt.Errorf("link type attribute '%s' still exists", LinkTypeAttrRsKey)
	}
	// clean up the database after the test
	return cleanUpLinkTypeAttrRs()
}

// remove supporting database entities
func cleanUpLinkTypeAttrRs() error {
	_, err := cfg.Client.DeleteLinkType(&oxc.LinkType{Key: LinkTypeAttrRsLinkTypeKey})
	_, err = cfg.Client.DeleteModel(&oxc.Model{Key: LinkTypeAttrRsModelKey})
	return err
}

func oxLinkTypeAttrResource() string {
	return fmt.Sprintf(
		`resource "ox_link_type_attr" "ox_link_type_attr_1" {
  key         	= "%s"
  name        	= "%s"
  description 	= "%s"
  link_type_key = "%s"
  type        	= "%s"
  def_value   	= "%s"
}`, LinkTypeAttrRsKey,
		LinkTypeAttrRsName,
		LinkTypeAttrRsDesc,
		LinkTypeAttrRsLinkTypeKey,
		LinkTypeAttrRsType,
		LinkTypeAttrRsDefValue)
}
