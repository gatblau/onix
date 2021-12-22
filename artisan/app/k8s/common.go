/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package k8s

const (
	AppsVersion = "apps/v1"
	CoreVersion = "v1"
)

type Labels struct {
	App string `yaml:"app"`
}

type Metadata struct {
	Annotations Annotations `yaml:"annotations"`
	Labels      Labels      `yaml:"labels,omitempty"`
	Name        string      `yaml:"name"`
}
