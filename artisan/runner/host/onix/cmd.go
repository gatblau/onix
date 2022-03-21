/*
  Onix Config Manager - Artisan Host Runner
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package onix

import (
	"fmt"

	"github.com/gatblau/onix/artisan/data"
)

// Cmd command information for remote host execution
type Cmd struct {
	// the natural key uniquely identifying the command
	Key string `json:"key"`
	// description of the command
	Description string `json:"description"`
	// the package to use
	Package string `json:"package"`
	// the function in the package to call
	Function string `json:"function"`
	// the function input information
	Input *data.Input `json:"input"`
	// the package registry user
	User string `json:"user"`
	// the package registry password
	Pwd string `json:"pwd"`
	// enables verbose output
	Verbose bool `json:"verbose"`
}

func (c *Cmd) Env() []string {
	var vars []string
	// append vars
	for _, v := range c.Input.Var {
		vars = append(vars, fmt.Sprintf("%s=%s", v.Name, v.Value))
	}
	// append secrets
	for _, s := range c.Input.Secret {
		vars = append(vars, fmt.Sprintf("%s=%s", s.Name, s.Value))
	}
	return vars
}

func (c *Cmd) GetVarValue(varName string) string {
	if c.Input.Var != nil {
		for _, v := range c.Input.Var {
			if v.Name == varName {
				return v.Value
			}
		}
	}
	return ""
}
