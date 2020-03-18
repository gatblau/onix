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
	cfg.URI = os.Getenv("TF_PROVIDER_OX_URI")
	if len(cfg.URI) == 0 {
		cfg.URI = data.Get("uri").(string)
		log.Printf("WARNING: Loading 'uri' from resource data. Consider setting the TF_PROVIDER_OX_URI environment variable instead.")
	}
	cfg.AuthMode = os.Getenv("TF_PROVIDER_OX_AUTH_MODE")
	if len(cfg.AuthMode) == 0 {
		cfg.AuthMode = data.Get("auth_mode").(string)
	}
	// if basic authentication selected
	if cfg.AuthMode == "basic" {
		cfg.User = os.Getenv("TF_PROVIDER_OX_USER")
		if len(cfg.User) == 0 {
			cfg.User = data.Get("user").(string)
			log.Printf("WARNING: Loading 'user' from resource data. Consider setting the TF_PROVIDER_OX_USER environment variable instead.")
		}
		cfg.Pwd = os.Getenv("TF_PROVIDER_OX_PWD")
		if len(cfg.Pwd) == 0 {
			cfg.Pwd = data.Get("pwd").(string)
			log.Printf("WARNING: Loading 'pwd' from resource data. Consider setting the TF_PROVIDER_OX_PWD environment variable instead.")
		}
	}
	// if open id connect selected
	if cfg.AuthMode == "oidc" {
		cfg.TokenURI = os.Getenv("TF_PROVIDER_OX_TOKEN_URI")
		if len(cfg.TokenURI) == 0 {
			cfg.TokenURI = data.Get("token_uri").(string)
			log.Printf("WARNING: Loading 'token_uri' from resource data. Consider setting the TF_PROVIDER_OX_TOKEN_URI environment variable instead.")
		} else {
			return errors.New("TF_PROVIDER_OX_TOKEN_URI is required if TF_PROVIDER_OX_AUTH_MODE=oidc")
		}
		cfg.ClientId = os.Getenv("TF_PROVIDER_OX_CLIENT_ID")
		if len(cfg.ClientId) == 0 {
			val := data.Get("client_id")
			if val != nil {
				cfg.ClientId = val.(string)
				log.Printf("WARNING: Loading 'client_secret' from resource data. Consider setting the TF_PROVIDER_OX_CLIENT_ID environment variable instead.")
			} else {
				return errors.New("TF_PROVIDER_OX_CLIENT_ID is required if TF_PROVIDER_OX_AUTH_MODE=oidc")
			}
		}
		cfg.AppSecret = os.Getenv("TF_PROVIDER_OX_APP_SECRET")
		if len(cfg.AppSecret) == 0 {
			val := data.Get("app_secret")
			if val != nil {
				cfg.AppSecret = val.(string)
				log.Printf("WARNING: Loading 'app_secret' from resource data. Consider setting the TF_PROVIDER_OX_APP_SECRET environment variable instead.")
			} else {
				return errors.New("TF_PROVIDER_OX_APP_SECRET is required if TF_PROVIDER_OX_AUTH_MODE=oidc")
			}
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
		bearerToken, err := cfg.Client.GetBearerToken(cfg.TokenURI, cfg.ClientId, cfg.AppSecret, cfg.User, cfg.Pwd)
		if err != nil {
			return err
		}
		cfg.Client.SetAuthToken(bearerToken)

	default:
		// can't recognise the auth_mode provided
		return errors.New(fmt.Sprintf("auth_mode = '%s' is not valid value. Use either 'none', 'basic' or 'oidc'.", cfg.AuthMode))
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
	if result != nil && result.Error {
		return errors.New(fmt.Sprintf("business logic error: %s", result.Message))
	}
	return e
}
