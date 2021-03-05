/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package flow

import (
	"encoding/json"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/i18n"
	"github.com/gatblau/onix/artisan/registry"
)

// a set of authentication credentials for a package registry
type Credential struct {
	User     string
	Password string
	Domain   string
}

type Flow struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	// a list of labels to document key aspects of the flow execution
	// for example using a target namespace if running in Kubernetes
	Labels  map[string]string `yaml:"labels"`
	GitURI  string            `yaml:"git_uri,omitempty"`
	AppIcon string            `yaml:"app_icon,omitempty"`
	Steps   []*Step           `yaml:"steps"`
	Input   *data.Input       `yaml:"input,omitempty"`
}

func (f *Flow) StepByFx(fxName string) *Step {
	for _, step := range f.Steps {
		if step.Function == fxName {
			return step
		}
	}
	return nil
}

func (f *Flow) RequiresGitSource() bool {
	var useGitSource, usePackageSource bool
	for _, step := range f.Steps {
		// function only - needs a git source
		if len(step.Package) == 0 && len(step.Function) > 0 {
			useGitSource = true
		}
		// requires a package source
		if len(step.Package) > 0 && len(step.PackageSource) > 0 && len(step.Function) > 0 {
			usePackageSource = true
		}
	}
	// git source is required if it is not using a package source
	return useGitSource && !usePackageSource
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

// retrieve all input data required by the flow without values
// interactive mode is off - gets definition only
func (f *Flow) GetInputDefinition(b *data.BuildFile, env *core.Envar) *data.Input {
	result := &data.Input{
		Key:    make([]*data.Key, 0),
		Secret: make([]*data.Secret, 0),
		Var:    make([]*data.Var, 0),
	}
	local := registry.NewLocalRegistry()

	for _, step := range f.Steps {
		// if a function is defined without a package and the source is not a package
		if step.surveyBuildfile(f.RequiresGitSource()) {
			// check a build file has been specified
			if b == nil {
				core.RaiseErr("flow '%s' requires a build.yaml", f.Name)
			}
			// surveys the build.yaml for variables
			i := data.SurveyInputFromBuildFile(step.Function, b, false, true, env)
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
		} else if step.surveyManifest() {
			// surveys the package manifest for variables
			name, err := core.ParseName(step.Package)
			i18n.Err(err, i18n.ERR_INVALID_PACKAGE_NAME)
			manif := local.GetManifest(name)
			if manif == nil {
				core.RaiseErr("manifest for package '%s' not found", name)
			}
			i := data.SurveyInputFromManifest(f.Name, step.Name, step.PackageSource, name.Domain, step.Function, manif, false, true, env)
			i.SurveyRegistryCreds(f.Name, step.Name, step.PackageSource, name.Domain, false, true, env)
			result.Merge(i)
		} else if step.surveyRuntime() {
			// surveys runtime manifest for variables
			i := data.SurveyInputFromURI(step.RuntimeManifest, false, true, env)
			result.Merge(i)
		} else {
			flowHealthCheck(f, step)
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

func (f *Flow) JsonBytes() ([]byte, error) {
	data, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (f *Flow) Step(name string) *Step {
	for _, step := range f.Steps {
		if step.Name == name {
			return step
		}
	}
	return nil
}
