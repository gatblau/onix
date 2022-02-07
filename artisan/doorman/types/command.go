/*
  Onix Config Manager - Artisan's Doorman
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package types

import "fmt"

// Command an instruction to be executed by a pipeline
// @Description an instruction to be executed by a pipeline
type Command struct {
	// a unique name for the command
	Name string `bson:"_id" json:"name" example:"clamscan"`
	// the command description
	Description string `bson:"description" json:"description" example:"scan files in specified path"`
	// the value of the command
	Value string `bson:"value" json:"value" example:"freshclam && clamscan -r ${path}"`
	// a regex used to determine if the command execution has errored
	ErrorRegex string `bson:"error_regex" json:"errorRegex" example:".*Infected files: [^0].*"`
	// determines if the process should stop on a command execution error
	StopOnError bool `bson:"stop_on_error" json:"stopOnError" example:"true"`
}

func (c Command) GetName() string {
	return c.Name
}

func (c Command) Valid() error {
	if len(c.Value) == 0 {
		return fmt.Errorf("command must have a value")
	}
	return nil
}
