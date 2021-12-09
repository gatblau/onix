/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package app

// SvcManifest the manifest that describes how a software service should be configured
type SvcManifest struct {
	// Name the name of the service
	Name string `yaml:"name"`
	// Description describes what the service is all about
	Description string `yaml:"description"`
	// the port used by the http service
	Port string `yaml:"port"`
	// the URI to determine if the service is ready to use
	ReadyURI string `yaml:"ready_uri,omitempty"`
	// the variables passed to the service (either ordinary or secret)
	Var []Var `yaml:"var,omitempty"`
	// the files used by the service (either ordinary or secret)
	File []File `yaml:"file,omitempty"`
	// one or more persistent volumes
	Volume []Volume `yaml:"volume,omitempty"`
}

// Var describes a variable used by a service
type Var struct {
	// the variable name
	Name string `yaml:"name"`
	// a human-readable description for the variable
	Description string `yaml:"description,omitempty"`
	// if defined, the fix value for the variable
	Value string `yaml:"value,omitempty"`
	// whether the variable should be treated as a secret
	Secret bool `yaml:"secret,omitempty"`
	// a default value for the variable
	Default string `yaml:"default,omitempty"`
}

// File describes a file used by a service
type File struct {
	// the file path
	Path string `yaml:"path"`
	// a human-readable description for the file
	Description string `yaml:"description,omitempty"`
	// whether the file should be treated as a secret
	Secret bool `yaml:"secret,omitempty"`
	// the template to use to create the file
	Template string `yaml:"template,omitempty"`
}

type Volume struct {
	// the name of the volume
	Name string `yaml:"name"`
	// the volume use description
	Description string `yaml:"description,omitempty"`
	// the volume source path
	Path string `yaml:"path"`
}
