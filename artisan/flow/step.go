package flow

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"github.com/gatblau/onix/artisan/data"
	"strings"
)

type Step struct {
	Name            string      `yaml:"name" json:"name"`
	Description     string      `yaml:"description,omitempty" json:"description,omitempty"`
	Runtime         string      `yaml:"runtime" json:"runtime"`
	RuntimeManifest string      `yaml:"runtime_manifest,omitempty" json:"runtime_manifest,omitempty"`
	Function        string      `yaml:"function,omitempty" json:"function,omitempty"`
	Package         string      `yaml:"package,omitempty" json:"package,omitempty"`
	PackageSource   string      `yaml:"source,omitempty" json:"source,omitempty"`
	Input           *data.Input `yaml:"input,omitempty" json:"input,omitempty"`
	Privileged      bool        `yaml:"privileged" json:"privileged"`
}

func (s *Step) surveyBuildfile(requiresGitSource bool) bool {
	// requires a git source, it defines a function without package
	return requiresGitSource && len(s.Function) > 0 && len(s.Package) == 0
}

func (s *Step) surveyManifest() bool {
	// defines a function and a package or in the case of a package merge a function is not required but package and source = merge exist
	return (len(s.Function) > 0 && len(s.Package) > 0) || (len(s.Package) > 0 && len(s.Function) == 0 && strings.ToLower(s.PackageSource) == "merge")
}

func (s *Step) surveyRuntime() bool {
	// defines a runtime manifest location, but no package or function
	return len(s.RuntimeManifest) > 0 && len(s.Package) == 0 && len(s.Function) == 0
}
