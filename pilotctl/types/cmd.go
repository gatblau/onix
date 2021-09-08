package types

/*
  Onix Config Manager - Pilot Control
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

import "github.com/gatblau/onix/artisan/data"

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
	// run command in runtime
	Containerised bool `json:"containerised"`
}
