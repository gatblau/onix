package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"net/http"
)

type ItemTypeAttribute struct {
	Id          string `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	DefValue    string `json:"defValue"`
	Managed     bool   `json:"managed"`
	Required    bool   `json:"required"`
	Regex       string `json:"regex"`
	ItemTypeKey string `json:"itemTypeKey"`
	Version     int64  `json:"version"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
}

func (typeAttr *ItemTypeAttribute) toJSON() (*bytes.Reader, error) {
	return GetJSONBytesReader(typeAttr)
}

func newItemTypeAttr(data *schema.ResourceData) *ItemTypeAttribute {
	return &ItemTypeAttribute{
		Key:         data.Get("key").(string),
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
		Type:        data.Get("type").(string),
		DefValue:    data.Get("def_value").(string),
		Managed:     data.Get("managed").(bool),
		Required:    data.Get("managed").(bool),
		Regex:       data.Get("regex").(string),
		ItemTypeKey: data.Get("item_type_key").(string),
	}
}

// get the Item Type Attribute in the http Response
func (typeAttr *ItemTypeAttribute) decode(response *http.Response) (*ItemTypeAttribute, error) {
	result := new(ItemTypeAttribute)
	err := json.NewDecoder(response.Body).Decode(result)
	return result, err
}

// populate the Item Type Attribute with the data in the terraform resource
func (typeAttr *ItemTypeAttribute) populate(data *schema.ResourceData) {
	data.SetId(typeAttr.Id)
	data.Set("key", typeAttr.Key)
	data.Set("description", typeAttr.Description)
	data.Set("type", typeAttr.Type)
	data.Set("def_value", typeAttr.DefValue)
	data.Set("managed", typeAttr.Managed)
	data.Set("required", typeAttr.Required)
	data.Set("regex", typeAttr.Regex)
	data.Set("item_type_key", typeAttr.ItemTypeKey)
}

// get the FQN for the item type attribute resource
func (typeAttr *ItemTypeAttribute) uri(baseUrl string) string {
	return fmt.Sprintf("%s/itemtype/%s/attribute/%s", baseUrl, typeAttr.ItemTypeKey, typeAttr.Key)
}

// issue a put http request with the Item Type Attribute data as payload to the resource URI
func (typeAttr *ItemTypeAttribute) put(meta interface{}) error {
	cfg := meta.(Config)

	// converts the passed-in payload to a bytes Reader
	bytes, err := typeAttr.toJSON()

	// any errors are returned immediately
	if err != nil {
		return err
	}

	// make an http put request to the service
	result, err := cfg.Client.Put(typeAttr.uri(cfg.Client.BaseURL), bytes)

	// any errors are returned
	if e := check(result, err); e != nil {
		return e
	}

	return nil
}

// issue a delete http request to the resource URI
func (typeAttr *ItemTypeAttribute) delete(meta interface{}) error {
	// get the Config instance from the meta object passed-in
	cfg := meta.(Config)

	// make an http delete request to the service
	result, err := cfg.Client.Delete(typeAttr.uri(cfg.Client.BaseURL))

	// any errors are returned
	if e := check(result, err); e != nil {
		return e
	}

	return nil
}

// issue a get http request to the resource URI
func (typeAttr *ItemTypeAttribute) get(meta interface{}) (*ItemTypeAttribute, error) {
	cfg := meta.(Config)

	// make an http put request to the service
	result, err := cfg.Client.Get(typeAttr.uri(cfg.Client.BaseURL))

	if err != nil {
		return nil, err
	}

	return typeAttr.decode(result)
}
