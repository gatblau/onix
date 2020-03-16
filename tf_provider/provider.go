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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"log"
	"os"
)

// the provider configuration
var cfg Config

// Configuration information for the Terraform provider
type Config struct {
	URI       string
	User      string
	Pwd       string
	Client    oxc.Client
	Token     string
	TokenURI  string
	ClientId  string
	AppSecret string
	AuthMode  string
}

// determined if the configuration has been loaded
func (cfg *Config) loaded() bool {
	return len(cfg.URI) > 0
}

// load the provider configuration from environment first
// then from terraform resource data
// NOTE: prefer loading credentials from environment variables
func (cfg *Config) load(data *schema.ResourceData) error {
	uri := os.Getenv("TF_PROVIDER_OX_URI")
	if len(uri) == 0 {
		cfg.URI = data.Get("uri").(string)
		log.Printf("WARNING: Loading URI from resource data. Consider setting the TF_PROVIDER_OX_URI env var instead.")
	}
	user := os.Getenv("TF_PROVIDER_OX_USER")
	if len(user) == 0 {
		cfg.User = data.Get("user").(string)
	}
	pwd := os.Getenv("TF_PROVIDER_OX_PWD")
	if len(pwd) == 0 {
		cfg.Pwd = data.Get("pwd").(string)
	}
	authMode := os.Getenv("TF_PROVIDER_OX_AUTH_MODE")
	if len(authMode) == 0 {
		cfg.AuthMode = data.Get("auth_mode").(string)
	}
	tokenURI := os.Getenv("TF_PROVIDER_OX_TOKEN_URI")
	if len(tokenURI) == 0 {
		cfg.TokenURI = data.Get("token_uri").(string)
	}
	clientId := os.Getenv("TF_PROVIDER_OX_CLIENT_ID")
	if len(clientId) == 0 {
		cfg.ClientId = data.Get("client_id").(string)
	}
	clientSecret := os.Getenv("TF_PROVIDER_OX_CLIENT_SECRET")
	if len(clientSecret) == 0 {
		val := data.Get("app_secret")
		if val != nil {
			cfg.AppSecret = val.(string)
		}
	}

	cfg.Client = oxc.Client{BaseURL: cfg.URI}

	switch cfg.AuthMode {
	case "none":
		// ensure no token value is specified
		cfg.Client.SetAuthToken("")

	case "basic":
		// sets a basic authentication token
		cfg.Client.SetAuthToken(cfg.Client.NewBasicToken(cfg.User, cfg.Pwd))

	case "oidc":
		// sets an OAuth Bearer token
		bearerToken, err := cfg.Client.GetBearerToken(tokenURI, clientId, clientSecret, user, pwd)
		if err != nil {
			return err
		}
		cfg.Client.SetAuthToken(bearerToken)

	default:
		// can't recognise the auth_mode provided
		return errors.New(fmt.Sprintf("auth_mode = '%s' is not valid value. Use either 'none', 'basic' or 'oidc'.", authMode))
	}
	return nil
}

// return a provider for production use
func provider() terraform.ResourceProvider {
	// pass in isTest = false
	return newProvider(false)
}

// create a provider instance
func newProvider(isTest bool) terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"uri": {
				Type:     schema.TypeString,
				Required: !isTest,
				Optional: isTest,
			},
			"user": {
				Type:     schema.TypeString,
				Required: !isTest,
				Optional: isTest,
			},
			"pwd": {
				Type:      schema.TypeString,
				Required:  !isTest,
				Optional:  isTest,
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
			"ox_model":          ModelDataSource(),
			"ox_item_type":      ItemTypeDataSource(),
			"ox_item_type_attr": ItemTypeAttributeDataSource(),
			"ox_link_type":      LinkTypeDataSource(),
			"ox_link_type_attr": LinkTypeAttributeDataSource(),
			"ox_link_rule":      LinkRuleDataSource(),
			"ox_partition":      PartitionDataSource(),
			"ox_role":           RoleDataSource(),
			"ox_privilege":      PrivilegeDataSource(),
			"ox_item":           ItemDataSource(),
			"ox_link":           LinkDataSource(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(data *schema.ResourceData) (interface{}, error) {
	if !cfg.loaded() {
		cfg.load(data)
	}
	return cfg, nil
}

func err(result *oxc.Result, e error) error {
	if result.Error {
		return errors.New(result.Message)
	}
	return e
}
