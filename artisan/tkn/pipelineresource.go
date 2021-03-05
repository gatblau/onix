/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package tkn

type PipelineResource struct {
	APIVersion string    `yaml:"apiVersion,omitempty"`
	Kind       string    `yaml:"kind,omitempty"`
	Metadata   *Metadata `yaml:"metadata,omitempty"`
	Spec       *Spec     `yaml:"spec,omitempty"`
}
