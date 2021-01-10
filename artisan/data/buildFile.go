/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package data

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

// structure of build.yaml file
type BuildFile struct {
	// the environment variables that apply to the build
	// any variables defined at this level will be available to all build profiles
	// in addition, the defined variables are added on top of the existing environment
	Env map[string]string `yaml:"env"`
	// a list of labels to be added to the artefact seal
	// they should be used to document key aspects of the artefact in a generic way
	Labels map[string]string `yaml:"labels"`
	// a list of build configurations in the form of labels, commands to run and environment variables
	Profiles []*Profile `yaml:"profiles"`
	// a list of functions containing a list of commands to execute
	Functions []*Function `yaml:"functions"`
}

func (b *BuildFile) GetEnv() map[string]string {
	return b.Env
}

// return the default profile if exists
func (b *BuildFile) DefaultProfile() *Profile {
	for _, profile := range b.Profiles {
		if profile.Default {
			return profile
		}
	}
	return nil
}

// return the function in the build file specified by its name
func (b *BuildFile) Fx(name string) *Function {
	for _, fx := range b.Functions {
		if fx.Name == name {
			return fx
		}
	}
	return nil
}

type Profile struct {
	// the name of the profile
	Name string `yaml:"name"`
	// whether this is the default profile
	Default bool `yaml:"default"`
	// the name of the application
	Application string `yaml:"application"`
	// the type of license used by the application
	// if not empty, it is added to the artefact seal
	License string `yaml:"license"`
	// the type of technology used by the application that can be used to determine the tool chain to use
	// e.g. java, nodejs, golang, python, php, etc
	Type string `yaml:"type"`
	// the artefact name
	Artefact string `yaml:"artefact"`
	// the pipeline Icon
	Icon string `yaml:"icon"`
	// runtime image that should be used to execute functions in the package
	Runtime string `yaml:"runtime,omitempty"`
	// Sonar configuration
	Sonar *Sonar `yaml:"sonar,omitempty"`
	// a set of labels associated with the profile
	Labels map[string]string `yaml:"labels"`
	// a set of environment variables required by the run commands
	Env map[string]string `yaml:"env"`
	// the commands to be executed to build the application
	Run []string `yaml:"run"`
	// the output of the build process, namely either a file or a folder, that has to be compressed
	// as part of the packaging process
	Target string `yaml:"target"`
	// merged target if exist, internal use only
	MergedTarget string
}

// gets a slice of string with each element containing key=value
func (p *Profile) GetEnv() map[string]string {
	return p.Env
}

// return the build profile specified by its name
func (b *BuildFile) Profile(name string) *Profile {
	for _, profile := range b.Profiles {
		if profile.Name == name {
			return profile
		}
	}
	return nil
}

// survey all missing variables in the profile
func (p *Profile) Survey(bf *BuildFile) map[string]string {
	env := bf.Env
	// merges the profile environment with the passed in environment
	for k, v := range p.Env {
		env[k] = v
	}
	// attempt to merge any environment variable in the profile run commands
	// run the merge in interactive mode so that any variables not available in the build file environment are surveyed
	_, updatedEnvironment := core.MergeEnvironmentVars(p.Run, env, true)
	// attempt to merge any environment variable in the functions run commands
	for _, run := range p.Run {
		// if the run line has a function
		if ok, fxName := core.HasFunction(run); ok {
			// merge any variables on the function
			env = bf.Fx(fxName).Survey(env)
		}
	}
	return updatedEnvironment
}

func LoadBuildFile(path string) (*BuildFile, error) {
	if !filepath.IsAbs(path) {
		abs, err := filepath.Abs(path)
		core.CheckErr(err, "cannot get absolute path for %s", path)
		path = abs
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot load build file from %s: %s", path, err)
	}
	buildFile := &BuildFile{}
	err = yaml.Unmarshal(bytes, buildFile)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal build file %s: %s", path, err)
	}
	return buildFile, nil
}

type Sonar struct {
	URI        string `yaml:"uri"`
	ProjectKey string `yaml:"project_key"`
	Sources    string `yaml:"sources"`
	Binaries   string `yaml:"binaries"`
}
