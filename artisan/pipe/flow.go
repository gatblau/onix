/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package pipe

import "github.com/gatblau/onix/artisan/data"

type Flow struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Steps       []Step `yaml:"steps"`
}

type Secret struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Path        string `yaml:"path,omitempty"`
	Value       string `yaml:"value,omitempty"`
}

type Step struct {
	Name        string        `yaml:"name"`
	Description string        `yaml:"description"`
	Runtime     string        `yaml:"runtime"`
	Function    string        `yaml:"function"`
	Package     string        `yaml:"package,omitempty"`
	Secret      []*Secret     `yaml:"secret,omitempty"`
	Input       []*data.Input `yaml:"input,omitempty"`
}
