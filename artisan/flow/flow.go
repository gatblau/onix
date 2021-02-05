/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package flow

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/registry"
)

// a set of authentication credentials for a package registry
type Credential struct {
	User     string
	Password string
	Domain   string
}

type Flow struct {
	Name        string      `yaml:"name"`
	Description string      `yaml:"description"`
	GitURI      string      `yaml:"git_uri,omitempty"`
	AppIcon     string      `yaml:"app_icon,omitempty"`
	Steps       []*Step     `yaml:"steps"`
	Input       *data.Input `yaml:"input,omitempty"`
}

func (f *Flow) StepByFx(fxName string) *Step {
	for _, step := range f.Steps {
		if step.Function == fxName {
			return step
		}
	}
	return nil
}

func (f *Flow) RequiresSource() bool {
	for _, step := range f.Steps {
		if len(step.Package) == 0 && len(step.Function) > 0 {
			return true
		}
	}
	return false
}

func (f *Flow) RequiresKey() bool {
	for _, step := range f.Steps {
		if step.Input != nil && step.Input.Key != nil {
			return true
		}
	}
	return false
}

func (f *Flow) RequiresSecrets() bool {
	for _, step := range f.Steps {
		if step.Input != nil && step.Input.Secret != nil {
			return true
		}
	}
	return false
}

type Step struct {
	Name            string      `yaml:"name"`
	Description     string      `yaml:"description,omitempty"`
	Runtime         string      `yaml:"runtime"`
	RuntimeManifest string      `yaml:"runtime_manifest,omitempty"`
	Function        string      `yaml:"function,omitempty"`
	Package         string      `yaml:"package,omitempty"`
	Input           *data.Input `yaml:"input,omitempty"`
}

// retrieve all input data required by the flow without values
// interactive mode is off
func (f *Flow) GetCombinedInput(b *data.BuildFile) *data.Input {
	result := &data.Input{
		Key:    make([]*data.Key, 0),
		Secret: make([]*data.Secret, 0),
		Var:    make([]*data.Var, 0),
	}
	local := registry.NewLocalRegistry()
	for _, step := range f.Steps {
		// if a function is defined without a package
		if len(step.Function) > 0 && len(step.Package) == 0 {
			// check a build file has been specified
			if b == nil {
				core.RaiseErr("flow '%s' requires a build.yaml", f.Name)
			}
			// surveys the build.yaml for variables
			i := data.SurveyInputFromBuildFile(step.Function, b, false)
			// add GIT_URI if not already added
			if i == nil || !result.VarExist("GIT_URI") {
				i.Var = append(i.Var, &data.Var{
					Name:        "GIT_URI",
					Description: "the URI of the project GIT repository",
					Required:    true,
					Type:        "url",
				})
			}
			result.Merge(i)
		} else if len(step.Function) > 0 && len(step.Package) > 0 {
			// surveys the package manifest for variables
			name, err := core.ParseName(step.Package)
			core.CheckErr(err, "package name is invalid")
			manif := local.GetManifest(name)
			if manif == nil {
				core.RaiseErr("manifest for package '%s' not found", name)
			}
			i := data.SurveyInputFromManifest(name, step.Function, manif, false)
			// add package registry credentials
			userKey := fmt.Sprintf("ART_REG_USER_%s", name.Domain)
			pwdKey := fmt.Sprintf("ART_REG_PWD_%s", name.Domain)
			if i == nil || !result.SecretExist(userKey) {
				i.Secret = append(i.Secret, &data.Secret{
					Name:        userKey,
					Description: fmt.Sprintf("the username to authenticate with the artefact registry at '%s'", name.Domain),
				})
			}
			if i == nil || !result.SecretExist(pwdKey) {
				i.Secret = append(i.Secret, &data.Secret{
					Name:        pwdKey,
					Description: fmt.Sprintf("the password to authenticate with the artefact registry at '%s'", name.Domain),
				})
			}
			result.Merge(i)
		} else {
			// surveys runtime manifest for variables
			i := data.SurveyInputFromURI(step.RuntimeManifest, false)
			result.Merge(i)
		}
		// try augment the result with default values in the build.yaml
		if b != nil {
			for _, v := range b.Input.Var {
				for _, v2 := range result.Var {
					if v.Name == v2.Name && len(v.Default) > 0 {
						v2.Default = v.Default
					}
				}
			}
		}
	}
	return result
}

func (f *Flow) HasVar(name string) bool {
	for _, v := range f.Input.Var {
		if v.Name == name {
			return true
		}
	}
	return false
}

func (f *Flow) HasSecret(name string) bool {
	for _, s := range f.Input.Secret {
		if s.Name == name {
			return true
		}
	}
	return false
}

func (f *Flow) HasKey(name string) bool {
	for _, k := range f.Input.Key {
		if k.Name == name {
			return true
		}
	}
	return false
}
