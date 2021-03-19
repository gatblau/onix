/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package data

import (
	"github.com/gatblau/onix/artisan/core"
)

type Function struct {
	// the name of the function
	Name string `yaml:"name"`
	// the description for the function
	Description string `yaml:"description,omitempty"`
	// a set of environment variables required by the run commands
	Env map[string]string `yaml:"env,omitempty"`
	// the commands to be executed by the function
	Run []string `yaml:"run,omitempty"`
	// is this function to be available in the manifest
	Export *bool `yaml:"export,omitempty"`
	// defines any bindings to inputs required to run this function
	Input *InputBinding `yaml:"input,omitempty"`
	// the runtime to run this function
	Runtime string `yaml:"runtime,omitempty"`
}

// list the names of the inputs required by a function
type InputBinding struct {
	Var    []string `yaml:"var"`
	Secret []string `yaml:"secret"`
	Key    []string `yaml:"key"`
	File   []string `yaml:"file"`
}

// gets a slice of string with each element containing key=value
func (f *Function) GetEnv() map[string]string {
	return f.Env
}

// survey all missing variables in the function
// pass in any available environment variables so that they are not surveyed
func (f *Function) Survey(env map[string]string) map[string]string {
	// merges the function environment with the passed in environment
	for k, v := range f.Env {
		env[k] = v
	}
	// attempt to merge any environment variable in the run commands
	// run the merge in interactive mode so that any variables not available in the build file environment are surveyed
	_, updatedEnvironment := core.MergeEnvironmentVars(f.Run, env, true)
	return updatedEnvironment
}
