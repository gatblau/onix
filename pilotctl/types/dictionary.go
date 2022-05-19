/*
  Onix Config Manager - Pilot Control
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package types

// Dictionary a key value pair list with name and description
type Dictionary struct {
	// Key a natural key used to uniquely identify this dictionary for the purpose of idempotent opeartions
	Key string `json:"key" yaml:"key"`
	// Name a friendly name for the dictionary
	Name string `json:"name" yaml:"name"`
	// Description describe the purpose / content of the dictionary
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Values a map containing key/value pairs that are the content held by the dictionary
	Values map[string]interface{} `json:"values,omitempty" yaml:"values,omitempty"`
	// Tags a list of string based tags used for categorising the dictionary
	Tags []interface{} `json:"tags,omitempty" yaml:"tags,omitempty"`
}
