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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
	"testing"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProviders = map[string]terraform.ResourceProvider{
		"ox": testProvider(),
	}
}

// returns a provider for testing purposes
func testProvider() terraform.ResourceProvider {
	testProvider := newProvider(true)
	testProvider.Configure(testProviderCfg())
	return testProvider
}

// verify the structure of the provider and all of the resources,
// and reports an error if it is invalid.
func TestProvider(t *testing.T) {
	if err := testProvider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

// get the configuration of the test provider
// requires environment variables or defaults to preset values
// NOTE: only Basic Authentication is currently implemented for testing purposes
func testProviderCfg() *terraform.ResourceConfig {
	uri := getVar("TF_PROVIDER_OX_URI", "http://localhost:8080")
	user := getVar("TF_PROVIDER_OX_USER", "admin")
	pwd := getVar("TF_PROVIDER_OX_PWD", "0n1x")
	return &terraform.ResourceConfig{
		Config: map[string]interface{}{
			"uri":  uri,
			"user": user,
			"pwd":  pwd,
		},
	}
}

// gets a configuration variable from the environment or applies a default value
func getVar(name string, defValue string) string {
	v := os.Getenv(name)
	if len(v) == 0 {
		return defValue
	}
	return v
}

// check the attribute in TF state matches the one in the database
func checkEntityAttr(rs *terraform.ResourceState, attrName string, targetValue string) error {
	if rs.Primary.Attributes[attrName] != targetValue {
		return fmt.Errorf("attribute '%s' expected value %s, but found %s", attrName, rs.Primary.Attributes[attrName], targetValue)
	}
	return nil
}
