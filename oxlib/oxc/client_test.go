package oxc

/*
   Onix Configuration Manager - HTTP Client
   Copyright (c) 2018-2021 by www.gatblau.org

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

import (
	"errors"
	"fmt"
	"testing"
)

// initialises the Web API client
var client = createClient()

// create an instance of the client
func createClient() *Client {
	client, err := NewClient(&ClientConf{
		BaseURI:            "http://localhost:8080",
		InsecureSkipVerify: true,
		AuthMode:           Basic,
		Username:           "admin",
		Password:           "0n1x",
		// uncomment below & reset configuration vars
		// to test using an OAuth bearer token
		// AuthMode:           	OIDC,
		// TokenURI:     		"https://dev-447786.okta.com/oauth2/default/v1/token",
		// ClientId:			"0oalyh...356",
		// AppSecret:			"Tsed........OP0oEf9H7",
	})
	if err != nil {
		panic(err)
	}
	return client
}

func checkResult(result *Result, err error, msg string, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	} else if result != nil {
		if result.Error {
			t.Fatal(fmt.Sprintf("%s: %s", msg, result.Message))
		} else if result.Operation == "L" {
			t.Fatal(fmt.Sprintf("Fail to update - Locked Serializable: %s", result.Ref))
		}
	}
}

func TestOnixClient_Put(t *testing.T) {
	// clear all data!
	result, err := client.Clear()
	if result != nil && !result.Changed {
		err = errors.New(result.Message)
	}
	checkResult(result, err, "failed to clear database", t)

	user := &User{
		Key:     "test_user",
		Name:    "Test User",
		Email:   "test@mail.com",
		Service: false,
		Expires: "01-01-2050 10:30:00+0100",
	}

	member := &Membership{
		Key:  "test_user_membership",
		User: "test_user",
		Role: "READER",
	}

	// delete the membership if already exists
	result, err = client.DeleteMembership(member)

	// delete the user if already exists
	result, err = client.DeleteUser(user)

	result, err = client.PutUser(user, false)
	checkResult(result, err, "create test_user failed", t)

	result, err = client.PutMembership(member)
	checkResult(result, err, "create test_user_membership failed", t)

	msg := "create test_model failed"
	model := &Model{
		Key:         "test_model",
		Name:        "Test Model",
		Description: "Test Model",
	}
	result, err = client.PutModel(model)
	checkResult(result, err, msg, t)

	itemType := &ItemType{
		Key:          "test_item_type",
		Name:         "Test Item Type",
		Description:  "Test Item Type",
		Model:        "test_model",
		EncryptMeta:  false,
		EncryptTxt:   true,
		NotifyChange: NotifyTypeType,
		Style:        newStyle(),
	}
	result, err = client.PutItemType(itemType)
	checkResult(result, err, "create test_item_type failed", t)

	itemTypeAttr := &ItemTypeAttribute{
		Key:         "test_item_type_attr_1",
		Name:        "CPU",
		Description: "Description for test_item_type_attr_1",
		Type:        "integer",
		DefValue:    "2",
		Required:    false,
		Regex:       "",
		ItemTypeKey: "test_item_type",
	}

	result, err = client.PutItemTypeAttr(itemTypeAttr)
	checkResult(result, err, "create test_item_type_attr_1 failed", t)

	item1 := &Item{
		Key:         "item_1",
		Name:        "Item 1",
		Description: "Test Item 1",
		Status:      1,
		Type:        "test_item_type",
		Txt:         "This is a test text configuration.",
		Attribute:   map[string]interface{}{"CPU": 5},
	}
	result, err = client.PutItem(item1)
	checkResult(result, err, "create item_1 failed", t)

	item2 := &Item{
		Key:         "item_2",
		Name:        "Item 2",
		Description: "Test Item 2",
		Status:      2,
		Type:        "test_item_type",
		Attribute:   map[string]interface{}{"CPU": 2},
	}
	result, err = client.PutItem(item2)
	checkResult(result, err, "create item_2 failed", t)

	linkType := &LinkType{
		Key:         "test_link_type",
		Name:        "Test Link Type",
		Description: "Test Link Type",
		Model:       "test_model",
	}
	result, err = client.PutLinkType(linkType)
	checkResult(result, err, "create test_link_type failed", t)

	linkRule := &LinkRule{
		Key:              "test_link_rule_1",
		Name:             "Test Item Type to Test Item Type rule",
		Description:      "Allow to connect two items of type test_item_type.",
		LinkTypeKey:      "test_link_type",
		StartItemTypeKey: "test_item_type",
		EndItemTypeKey:   "test_item_type",
	}
	result, err = client.PutLinkRule(linkRule)
	checkResult(result, err, "create test_item_type->test_item_type rule failed", t)

	link := &Link{
		Key:          "test_link_1",
		Description:  "Test Link 1",
		Type:         "test_link_type",
		StartItemKey: "item_1",
		EndItemKey:   "item_2",
	}
	result, err = client.PutLink(link)
	checkResult(result, err, "create link_1 failed", t)

	data := getData()
	result, err = client.PutData(data)
	if err != nil {
		t.Fatal(err)
	}
	if result.Error {
		t.Fatal(result.Message)
	}

	list, err := client.GetItemChildren(&Item{Key: "item_1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Values) == 0 {
		t.Fatal("no value in list")
	}
}

func newStyle() map[string]interface{} {
	style := make(map[string]interface{})
	style["colour"] = "red"
	style["thickness"] = 4
	return style
}

func getData() *GraphData {
	return &GraphData{
		Models: []Model{
			{
				Key:         "TERRA",
				Name:        "Terraform Model",
				Description: "Defines the item and link types that describe Terraform resources.",
			},
		},
		ItemTypes: []ItemType{
			{
				Key:         "TF_STATE",
				Name:        "Terraform State",
				Description: "State about a group of managed infrastructure and configuration resources. This state is used by Terraform to map real world resources to your configuration, keep track of metadata, and to improve performance for large infrastructures.",
				Model:       "TERRA",
			},
			{
				Key:         "TF_RESOURCE",
				Name:        "Terraform Resource",
				Description: "Each resource block describes one or more infrastructure objects, such as virtual networks, compute instances, or higher-level components such as DNS records.",
				Model:       "TERRA",
			},
		},
		ItemTypeAttributes: []ItemTypeAttribute{
			{
				Key:         "TF_ITEM_ATTR_MODE",
				Name:        "mode",
				Description: "Whether the resource is a data source or a managed resource.",
				Type:        "string",
				ItemTypeKey: "TF_RESOURCE",
				Required:    true,
			},
			{
				Key:         "TF_ITEM_ATTR_TYPE",
				Name:        "type",
				Description: "The resource type.",
				Type:        "string",
				ItemTypeKey: "TF_RESOURCE",
				Required:    true,
			},
			{
				Key:         "TF_ITEM_ATTR_PROVIDER",
				Name:        "provider",
				Description: "The provider used to manage this resource.",
				Type:        "string",
				ItemTypeKey: "TF_RESOURCE",
				Required:    true,
			},
		},
		LinkTypes: []LinkType{
			{
				Key:         "TF_STATE_LINK",
				Name:        "Terraform State Link Type",
				Description: "Links Terraform resources that are part of a state.",
				Model:       "TERRA",
			},
		},
		LinkRules: []LinkRule{
			{
				Key:              fmt.Sprintf("%s->%s", "TF_STATE", "TF_RESOURCE"),
				Name:             "Terraform State to Resource Rule",
				Description:      "Allow the linking of a Terraform State item to one or more Terraform Resource items using Terraform State Links.",
				LinkTypeKey:      "TF_STATE_LINK",
				StartItemTypeKey: "TF_STATE",
				EndItemTypeKey:   "TF_RESOURCE",
			},
		},
	}
}
