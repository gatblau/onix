/*
  Onix Config Manager - Art
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

// structure of package.yaml file
type BuildFile struct {
	Type     string            `yaml:"type"`
	Env      map[string]string `yaml:"env"`
	License  string            `yaml:"license"`
	Labels   map[string]string `yaml:"labels"`
	Profiles []Profile         `yaml:"profiles"`
}

func (b *BuildFile) getEnv() []string {
	var vars []string
	// adds the environment variables defined in the profile
	for key, value := range b.Env {
		vars = append(vars, fmt.Sprintf("%s=%s", key, value))
	}
	return vars
}

type Profile struct {
	Name   string            `yaml:"name"`
	Labels map[string]string `yaml:"labels"`
	Env    map[string]string `yaml:"env"`
	Run    []string          `yaml:"run"`
	Target string            `yaml:"target"`
}

// gets a slice of string with each element containing key=value
func (p *Profile) getEnv() []string {
	var vars []string
	// adds the environment variables defined in the profile
	for key, value := range p.Env {
		vars = append(vars, fmt.Sprintf("%s=%s", key, value))
	}
	return vars
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
