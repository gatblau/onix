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
	"errors"
	"fmt"
	"github.com/gatblau/oxc"
	"github.com/hashicorp/terraform/helper/schema"
)

// Configuration information for the Terraform provider
type Config struct {
	URI    string
	User   string
	Pwd    string
	Client oxc.Client
	Token  string
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"uri": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pwd": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"auth_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "basic",
			},
			"client_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Default:   "",
				Sensitive: true,
			},
			"token_uri": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ox_model":          ModelResource(),
			"ox_item_type":      ItemTypeResource(),
			"ox_item_type_attr": ItemTypeAttributeResource(),
			"ox_link_type":      LinkTypeResource(),
			"ox_link_type_attr": LinkTypeAttributeResource(),
			"ox_link_rule":      LinkRuleResource(),
			"ox_partition":      PartitionResource(),
			"ox_role":           RoleResource(),
			"ox_privilege":      PrivilegeResource(),
			"ox_item":           ItemResource(),
			"ox_link":           LinkResource(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ox_model_data":          ModelDataSource(),
			"ox_item_type_data":      ItemTypeDataSource(),
			"ox_item_type_attr_data": ItemTypeAttributeDataSource(),
			"ox_link_type_data":      ItemTypeDataSource(),
			"ox_link_type_attr_data": LinkTypeAttributeDataSource(),
			"ox_link_rule_data":      LinkRuleDataSource(),
			"ox_partition_data":      PartitionDataSource(),
			"ox_role_data":           RoleDataSource(),
			"ox_privilege_data":      PrivilegeDataSource(),
			"ox_item_data":           ItemDataSource(),
			"ox_link_data":           LinkDataSource(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	uri := d.Get("uri").(string)
	user := d.Get("user").(string)
	pwd := d.Get("pwd").(string)
	authMode := d.Get("auth_mode").(string)
	tokenURI := d.Get("token_uri").(string)
	clientId := d.Get("client_id").(string)
	secret := d.Get("secret").(string)

	cfg := Config{
		URI:    uri,
		User:   user,
		Pwd:    pwd,
		Client: oxc.Client{BaseURL: uri},
	}

	switch authMode {
	case "none":
		// ensure no token value is specified
		cfg.Client.SetAuthToken("")

	case "basic":
		// sets a basic authentication token
		cfg.Client.SetAuthToken(cfg.Client.NewBasicToken(user, pwd))

	case "oidc":
		// sets an OAuth Bearer token
		bearerToken, err := cfg.Client.GetBearerToken(tokenURI, clientId, secret, user, pwd)
		if err != nil {
			return "", err
		}
		cfg.Client.SetAuthToken(bearerToken)

	default:
		// can't recognise the auth_mode provided
		return "", errors.New(fmt.Sprintf("auth_mode = '%s' is not valid value. Use either 'none', 'basic' or 'oidc'.", authMode))
	}

	return cfg, nil
}

func err(result *oxc.Result, e error) error {
	if result.Error {
		return errors.New(result.Message)
	}
	return e
}
