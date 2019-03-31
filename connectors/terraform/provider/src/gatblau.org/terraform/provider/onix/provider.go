/*
   Onix CMDB - Copyright (c) 2018-2019 by www.gatblau.org

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
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"uri": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "http://localhost:8080",
			},
			"user": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "admin",
			},
			"pwd": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0n1x",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ox_item_type": ItemTypeResource(),
			"ox_item":      ItemResource(),
			"ox_link_type": LinkTypeResource(),
			"ox_link":      LinkResource(),
			"ox_link_rule": LinkRuleResource(),
			"ox_model":     ModelResource(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ox_item_type_data": ItemTypeDataSource(),
			"ox_item_data":      ItemDataSource(),
			"ox_link_type_data": ItemTypeDataSource(),
			"ox_link_data":      LinkDataSource(),
			"ox_link_rule_data": LinkRuleDataSource(),
			"ox_model_data":     ModelDataSource(),
		},
		ConfigureFunc: configureProvider,
	}
}

type Config struct {
	URI    string
	User   string
	Pwd    string
	Client Client
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	uri := d.Get("uri").(string)
	user := d.Get("user").(string)
	pwd := d.Get("pwd").(string)

	client := Client{BaseURL: uri}
	client.initBasicAuthToken(user, pwd)

	config := Config{
		URI:    uri,
		User:   user,
		Pwd:    pwd,
		Client: client,
	}

	return config, nil
}
