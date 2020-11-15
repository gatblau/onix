/*
  Onix Config Manager - Artie
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package build

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

// structure of build.yaml file
type BuildFile struct {
	// the type of technology used by the application that can be used to determine the tool chain to use
	// e.g. java, nodejs, golang, python, php, etc
	Type string `yaml:"type"`
	// the environment variables that apply to the build
	// any variables defined at this level will be available to all build profiles
	// in addition, the defined variables are added on top of the existing environment
	Env map[string]string `yaml:"env"`
	// the type of license used by the application
	// if not empty, it is added to the artefact seal
	License string `yaml:"license"`
	// a list of labels to be added to the artefact seal
	// they should be used to document key aspects of the artefact in a generic way
	Labels map[string]string `yaml:"labels"`
	// a list of build configurations in the form of labels, commands to run and environment variables
	Profiles []Profile `yaml:"profiles"`
	// a list of functions containing a list of commands to execute
	Functions []Function `yaml:"functions"`
}

func (b *BuildFile) getEnv() map[string]string {
	return b.Env
}

// return the default profile if exists
func (b *BuildFile) defaultProfile() *Profile {
	for _, profile := range b.Profiles {
		if profile.Default {
			return &profile
		}
	}
	return nil
}

// return the function in the build file specified by its name
func (b *BuildFile) fx(name string) *Function {
	for _, fx := range b.Functions {
		if fx.Name == name {
			return &fx
		}
	}
	return nil
}

type Profile struct {
	// the name of the profile
	Name string `yaml:"name"`
	// whether this is the default profile
	Default bool `yaml:"default"`
	// a set of labels associated with the profile
	Labels map[string]string `yaml:"labels"`
	// a set of environment variables required by the run commands
	Env map[string]string `yaml:"env"`
	// the commands to be executed to build the application
	Run []string `yaml:"run"`
	// the output of the build process, namely either a file or a folder, that has to be compressed
	// as part of the packaging process
	Target string `yaml:"target"`
}

// gets a slice of string with each element containing key=value
func (p *Profile) getEnv() map[string]string {
	return p.Env
}

func LoadBuildFile(path string) *BuildFile {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	buildFile := &BuildFile{}
	err = yaml.Unmarshal(bytes, buildFile)
	if err != nil {
		log.Fatal(err)
	}
	return buildFile
}

type Function struct {
	// the name of the function
	Name string `yaml:"name"`
	// the description for the function
	Description string `yaml:"description"`
	// a set of environment variables required by the run commands
	Env map[string]string `yaml:"env"`
	// the commands to be executed by the function
	Run []string `yaml:"run"`
}

// gets a slice of string with each element containing key=value
func (p *Function) getEnv() map[string]string {
	return p.Env
}
