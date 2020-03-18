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
	LinkResourceName    = "ox_link.ox_link_1"
	LinkRsItem1Key      = "test_acc_ox_link_1_item_1"
	LinkRsItem1Name     = "test_acc_ox_link_1_item_1 name"
	LinkRsItem1Desc     = "test_acc_ox_link_1_item_1 description"
	LinkRsItem2Key      = "test_acc_ox_link_1_item_2"
	LinkRsItem2Name     = "test_acc_ox_link_1_item_2 name"
	LinkRsItem2Desc     = "test_acc_ox_link_1_item_2 description"
	LinkRsKey           = "test_acc_ox_link_1"
	LinkRsDesc          = "ox_link Description"
	LinkRsMeta          = `{ "OS" = "RHEL7.3" }`
	LinkRsAttr          = `{ "RAM" : "3", "CPU" : "1" }`
	LinkRsLinkRuleKey   = "test_acc_item_type_ox_link-test_acc_item_type_ox_link"
	LinkRsLinkRuleName  = "test_acc_item_type_ox_link->test_acc_item_type_ox_link"
	LinkRsLinkRuleDesc  = "test_acc_item_type_ox_link->test_acc_item_type_ox_link description"
	LinkRsModelKey      = "test_acc_model_ox_link"
	LinkRsModelName     = "Model - ox_link"
	LinkRsModelDesc     = "Model for testing ox_link."
	LinkRsItemTypeKey   = "test_acc_item_type_ox_link"
	LinkRsItemTypeName  = "test_acc_item_type_ox_link name"
	LinkRsItemTypeDesc  = "test_acc_item_type_ox_link description"
	LinkRsLinkTypeKey   = "test_acc_link_type_ox_link"
	LinkRsLinkTypeName  = "test_acc_link_type_ox_link name"
	LinkRsLinkTypeDesc  = "test_acc_link_type_ox_link description"
	LinkRsTypeAttr1Key  = "test_acc_link_type_attr_cpu_ox_link"
	LinkRsTypeAttr1Name = "CPU"
	LinkRsTypeAttr1Desc = "CPU attr for testing ox_link."
	LinkRsTypeAttr2Key  = "test_acc_link_type_attr_ram_ox_link"
	LinkRsTypeAttr2Name = "RAM"
	LinkRsTypeAttr2Desc = "RAM attr for testing ox_link."
)

func TestLinkResource(t *testing.T) {
	resourceName := LinkResourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { prepareLinkResourceTest(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			//create
			{
				Config: oxLinkResource(),
				Check: resource.ComposeTestCheckFunc(
					// check item resource attributes in Terraform state
					resource.TestCheckResourceAttr(resourceName, "key", LinkRsKey),
					resource.TestCheckResourceAttr(resourceName, "description", LinkRsDesc),
					resource.TestCheckResourceAttr(resourceName, "type", LinkRsLinkTypeKey),
					resource.TestCheckResourceAttr(resourceName, "start_item_key", LinkRsItem1Key),
					resource.TestCheckResourceAttr(resourceName, "end_item_key", LinkRsItem2Key),

					// check for side effects in Onix database
					checkLinkResourceCreated(LinkResourceName),
				),
			},
		},
		CheckDestroy: checkLinkResourceDestroyed,
	})
}

// create supporting database entities
func prepareLinkResourceTest(t *testing.T) {
	// need a model first
	result, err := cfg.Client.PutModel(
		&oxc.Model{
			Key:         LinkRsModelKey,
			Name:        LinkRsModelName,
			Description: LinkRsModelDesc,
		})
	if err != nil {
		t.Error(err)
	}
	if result.Error {
		t.Error(result.Message)
	}

	// an item type in the model
	result, err = cfg.Client.PutItemType(
		&oxc.ItemType{
			Key:         LinkRsItemTypeKey,
			Name:        LinkRsItemTypeName,
			Description: LinkRsItemTypeDesc,
			Model:       LinkRsModelKey,
		})
	if err != nil {
		t.Error(err)
	}
	if result.Error {
		t.Error(result.Message)
	}

	// an link type in the model
	result, err = cfg.Client.PutLinkType(
		&oxc.LinkType{
			Key:         LinkRsLinkTypeKey,
			Name:        LinkRsLinkTypeName,
			Description: LinkRsLinkTypeDesc,
			Model:       LinkRsModelKey,
		})
	if err != nil {
		t.Error(err)
	}
	if result.Error {
		t.Error(result.Message)
	}

	// and two attributes for the link type
	result, err = cfg.Client.PutLinkTypeAttr(
		&oxc.LinkTypeAttribute{
			Key:         LinkRsTypeAttr1Key,
			Name:        LinkRsTypeAttr1Name,
			Description: LinkRsTypeAttr1Desc,
			LinkTypeKey: LinkRsLinkTypeKey,
		})
	if err != nil {
		t.Error(err)
	}
	if result.Error {
		t.Error(result.Message)
	}

	result, err = cfg.Client.PutLinkTypeAttr(
		&oxc.LinkTypeAttribute{
			Key:         LinkRsTypeAttr2Key,
			Name:        LinkRsTypeAttr2Name,
			Description: LinkRsTypeAttr2Desc,
			LinkTypeKey: LinkRsLinkTypeKey,
		})
	if err != nil {
		t.Error(err)
	}
	if result.Error {
		t.Error(result.Message)
	}

	// a link rule to allow to connect two item types
	result, err = cfg.Client.PutLinkRule(
		&oxc.LinkRule{
			Key:              LinkRsLinkRuleKey,
			Name:             LinkRsLinkRuleName,
			Description:      LinkRsLinkRuleDesc,
			StartItemTypeKey: LinkRsItemTypeKey,
			EndItemTypeKey:   LinkRsItemTypeKey,
			LinkTypeKey:      LinkRsLinkTypeKey,
		})
	if err != nil {
		t.Error(err)
	}
	if result.Error {
		t.Error(result.Message)
	}

	// item 1
	result, err = cfg.Client.PutItem(
		&oxc.Item{
			Key:         LinkRsItem1Key,
			Name:        LinkRsItem1Name,
			Description: LinkRsItem1Desc,
			Type:        LinkRsItemTypeKey,
		})
	if err != nil {
		t.Error(err)
	}
	if result.Error {
		t.Error(result.Message)
	}

	// item 2
	result, err = cfg.Client.PutItem(
		&oxc.Item{
			Key:         LinkRsItem2Key,
			Name:        LinkRsItem2Name,
			Description: LinkRsItem2Desc,
			Type:        LinkRsItemTypeKey,
		})
	if err != nil {
		t.Error(err)
	}
	if result.Error {
		t.Error(result.Message)
	}
}

// check the link has been created in the database
func checkLinkResourceCreated(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		// retrieve the link from the database
		key := rs.Primary.Attributes["key"]
		link, err := cfg.Client.GetLink(&oxc.Link{Key: key})
		if err != nil {
			return fmt.Errorf("can't read link %s", key)
		}

		// compares the state with the database values
		if err := checkEntityAttr(rs, "description", link.Description); err != nil {
			return err
		}
		if err := checkEntityAttr(rs, "type", link.Type); err != nil {
			return err
		}
		if err := checkEntityAttr(rs, "start_item_key", link.StartItemKey); err != nil {
			return err
		}
		if err := checkEntityAttr(rs, "end_item_key", link.EndItemKey); err != nil {
			return err
		}
		return nil
	}
}

// check the item has been destroyed destroyed in the database
func checkLinkResourceDestroyed(s *terraform.State) error {
	link, _ := cfg.Client.GetLink(&oxc.Link{Key: LinkRsKey})
	// should get an error as item should not exist anymore
	if link != nil {
		return fmt.Errorf("link %s still exists", LinkRsKey)
	}
	// clean up the database after the test
	return cleanUpLinkRs()
}

// remove supporting database entities
func cleanUpLinkRs() error {
	_, err := cfg.Client.DeleteItem(&oxc.Item{Key: LinkRsItem1Key})
	_, err = cfg.Client.DeleteItem(&oxc.Item{Key: LinkRsItem2Key})
	_, err = cfg.Client.DeleteLinkType(&oxc.LinkType{Key: LinkRsLinkTypeKey})
	_, err = cfg.Client.DeleteItemType(&oxc.ItemType{Key: LinkRsItemTypeKey})
	_, err = cfg.Client.DeleteModel(&oxc.Model{Key: LinkRsLinkTypeKey})
	return err
}

func oxLinkResource() string {
	return fmt.Sprintf(`resource "ox_link" "ox_link_1" {
  key         = "%s"
  description = "%s"
  type        = "%s"
  start_item_key = "%s"
  end_item_key = "%s"
  meta        = %s
  attribute   = %s
}`, LinkRsKey,
		LinkRsDesc,
		LinkRsLinkTypeKey,
		LinkRsItem1Key,
		LinkRsItem2Key,
		LinkRsMeta,
		LinkRsAttr)
}
