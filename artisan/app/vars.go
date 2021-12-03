/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package app

import "sort"

// Vars the list of variables used by an application
type Vars struct {
	Items []AppVar
}

// AppVar describes a variable used by the app as a whole
type AppVar struct {
	// the variable name
	Name string `yaml:"name"`
	// a human-readable description for the variable
	Description string `yaml:"description,omitempty"`
	// if defined, the fix value for the variable
	Value string `yaml:"value,omitempty"`
	// whether the variable should be treated as a secret
	Secret bool `yaml:"secret,omitempty"`
	// the name of the service that originated the variable
	Service string `yaml:"service,omitempty"`
}

func (a *Vars) SortByService() {
	sort.Slice(a, func(i, j int) bool {
		return a.Items[i].Service < a.Items[j].Service
	})
}
