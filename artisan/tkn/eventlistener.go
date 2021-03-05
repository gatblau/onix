/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package tkn

type EventListener struct {
	APIVersion string    `yaml:"apiVersion"`
	Kind       string    `yaml:"kind"`
	Metadata   *Metadata `yaml:"metadata"`
	Spec       *Spec     `yaml:"spec"`
}

type Bindings struct {
	Name string `yaml:"name"`
}

type Template struct {
	Name string `yaml:"name"`
}

type Triggers struct {
	Bindings []*Bindings `yaml:"bindings"`
	Template *Template   `yaml:"template"`
}
