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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

// the provider configuration
var cfg Config

// Configuration information for the Terraform provider
type Config struct {
	URI                string
	User               string
	Pwd                string
	Client             *oxc.Client
	Token              string
	TokenURI           string
	ClientId           string
	AppSecret          string
	AuthMode           string
	InsecureSkipVerify bool
}

// determined if the configuration has been loaded
func (cfg *Config) loaded() bool {
	return len(cfg.URI) > 0
}

// load the provider configuration from environment first
// then from terraform resource data
// NOTE: prefer loading credentials from environment variables
func (cfg *Config) load(data *schema.ResourceData) error {
	// set time format to UNIX Time as it is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// gets the web API URI
	cfg.URI = os.Getenv("TF_PROVIDER_OX_URI")
	if len(cfg.URI) == 0 {
		cfg.URI = data.Get("uri").(string)
		log.Warn().Msg("Loading 'uri' from resource data. Consider setting the TF_PROVIDER_OX_URI environment variable instead.")
	}
	// gets the flag to skip TLS cert verification
	skipVerify := os.Getenv("TF_PROVIDER_INSECURE_SKIP_VERIFY")
	if len(skipVerify) > 0 {
		skip, err := strconv.ParseBool(skipVerify)
		if err != nil {
			log.Warn().Msgf("failed to parse boolean '%s' from TF_PROVIDER_INSECURE_SKIP_VERIFY. Ignoring.", skipVerify)
		}
		cfg.InsecureSkipVerify = skip
	} else {
		cfg.InsecureSkipVerify = data.Get("insecure_skip_verify").(bool)
	}
	// gets the authentication mode to use
	cfg.AuthMode = os.Getenv("TF_PROVIDER_OX_AUTH_MODE")
	if len(cfg.AuthMode) == 0 {
		cfg.AuthMode = data.Get("auth_mode").(string)
	}
	// if the auth mode is not set then it defaults to basic authentication
	if len(cfg.AuthMode) == 0 {
		log.Info().Msg("auth_mode not found, defaulting to basic authentication")
		cfg.AuthMode = "basic"
	}
	// if basic authentication selected
	if cfg.AuthMode == "basic" {
		cfg.User = os.Getenv("TF_PROVIDER_OX_USER")
		if len(cfg.User) == 0 {
			cfg.User = data.Get("user").(string)
			log.Warn().Msg("Loading 'user' from resource data. Consider setting the TF_PROVIDER_OX_USER environment variable instead.")
		}
		cfg.Pwd = os.Getenv("TF_PROVIDER_OX_PWD")
		if len(cfg.Pwd) == 0 {
			cfg.Pwd = data.Get("pwd").(string)
			log.Warn().Msg("Loading 'pwd' from resource data. Consider setting the TF_PROVIDER_OX_PWD environment variable instead.")
		}
	}
	// if open id connect selected
	if cfg.AuthMode == "oidc" {
		cfg.TokenURI = os.Getenv("TF_PROVIDER_OX_TOKEN_URI")
		if len(cfg.TokenURI) == 0 {
			cfg.TokenURI = data.Get("token_uri").(string)
			log.Warn().Msg("Loading 'token_uri' from resource data. Consider setting the TF_PROVIDER_OX_TOKEN_URI environment variable instead.")
		} else {
			return errors.New("TF_PROVIDER_OX_TOKEN_URI is required if TF_PROVIDER_OX_AUTH_MODE=oidc")
		}
		cfg.ClientId = os.Getenv("TF_PROVIDER_OX_APP_CLIENT_ID")
		if len(cfg.ClientId) == 0 {
			val := data.Get("app_client_id")
			if val != nil {
				cfg.ClientId = val.(string)
				log.Warn().Msg("Loading 'client_secret' from resource data. Consider setting the TF_PROVIDER_OX_CLIENT_ID environment variable instead.")
			} else {
				return errors.New("TF_PROVIDER_OX_CLIENT_ID is required if TF_PROVIDER_OX_AUTH_MODE=oidc")
			}
		}
		cfg.AppSecret = os.Getenv("TF_PROVIDER_OX_APP_SECRET")
		if len(cfg.AppSecret) == 0 {
			val := data.Get("app_secret")
			if val != nil {
				cfg.AppSecret = val.(string)
				log.Warn().Msg("Loading 'app_secret' from resource data. Consider setting the TF_PROVIDER_OX_APP_SECRET environment variable instead.")
			} else {
				return errors.New("TF_PROVIDER_OX_APP_SECRET is required if TF_PROVIDER_OX_AUTH_MODE=oidc")
			}
		}
	}

	// build the web api client configuration
	conf := &oxc.ClientConf{
		BaseURI:            cfg.URI,
		InsecureSkipVerify: cfg.InsecureSkipVerify,
		Username:           cfg.User,
		Password:           cfg.Pwd,
		TokenURI:           cfg.TokenURI,
		ClientId:           cfg.ClientId,
		AppSecret:          cfg.AppSecret,
	}
	// set the authentication mode
	conf.SetAuthMode(cfg.AuthMode)

	// gets a new client instance
	client, err := oxc.NewClient(conf)
	if err != nil {
		return err
	}

	// assigns the client to the config object
	cfg.Client = client

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
			"insecure_skip_verify": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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

// read the version from the resource data and returns a 0 if the value is nil
func getVersion(data *schema.ResourceData) int64 {
	i := data.Get("version")
	if i != nil {
		return int64(i.(int))
	}
	return 0
}
