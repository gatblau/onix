package main

import (
	"fmt"
	"testing"
)

var client Client

func init() {
	client = Client{BaseURL: "http://localhost:8080"}
	client.initBasicAuthToken("admin", "0n1x")
}

func check(result *Result, err error, msg string, t *testing.T) {
	if err != nil {
		t.Error(msg)
	}
	if result.Error {
		t.Error(fmt.Sprintf("%s: %s", msg, result.Message))
	}
}

func TestOnixClient_Put(t *testing.T) {
	model := Model {
		Name: "Test Model",
		Description: "Test Model",
	}
	result, err := client.Put("model", "test_model", model.ToJSON())
	check(result, err, "create test_model failed", t)

	itemType := ItemType{
		Name:        "Test Item Type",
		Description: "Test Item Type",
		Model: "test_model",
	}
	result, err = client.Put("itemtype", "test_item_type", itemType.ToJSON())
	check(result, err, "create test_item_type failed", t)

	item_1 := Item {
		Name:        "Item 1",
		Description: "Test Item 1",
		Status:      1,
		Type:        "test_item_type",
	}
	result, err = client.Put("item", "item_1", item_1.ToJSON())
	check(result, err, "create item_1 failed", t)

	item_2 := Item{
		Name:        "Item 2",
		Description: "Test Item 2",
		Status:      2,
		Type:        "test_item_type",
	}
	result, err = client.Put("item", "item_2", item_2.ToJSON())
	check(result, err, "create item_2 failed", t)

	link_type := LinkType{
		Name:        "Test Link Type",
		Description: "Test Link Type",
		Model: "test_model",
	}
	result, err = client.Put("linktype", "test_link_type", link_type.ToJSON())
	check(result, err, "create test_link_type failed", t)

	link_rule := LinkRule{
		Name:             "Test Item Type to Test Item Type rule",
		Description:      "Allow to connect two items of type test_item_type.",
		LinkTypeKey:      "test_link_type",
		StartItemTypeKey: "test_item_type",
		EndItemTypeKey:   "test_item_type",
	}
	result, err = client.Put("linkrule", "test_item_type->test_item_type", link_rule.ToJSON())
	check(result, err, "create test_item_type->test_item_type rule failed", t)

	link := Link{
		Description:  "Test Link 1",
		Type:         "test_link_type",
		StartItemKey: "item_1",
		EndItemKey:   "item_2",
	}
	result, err = client.Put("link", "link_1", link.ToJSON())
	check(result, err, "create link_1 failed", t)
}
