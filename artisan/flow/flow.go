/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package flow

import (
	"github.com/gatblau/onix/artisan/data"
)

type Flow struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Steps       []*Step `yaml:"steps"`
}

func (f *Flow) StepByFx(fxName string) *Step {
	for _, step := range f.Steps {
		if step.Function == fxName {
			return step
		}
	}
	return nil
}

type Secret struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Path        string `yaml:"path,omitempty"`
	Value       string `yaml:"value,omitempty"`
}

type Step struct {
	Name        string      `yaml:"name"`
	Description string      `yaml:"description,omitempty"`
	Runtime     string      `yaml:"runtime"`
	Function    string      `yaml:"function,omitempty"`
	Package     string      `yaml:"package,omitempty"`
	Input       *data.Input `yaml:"input,omitempty"`
}
