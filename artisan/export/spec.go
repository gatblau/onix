/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package export

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/registry"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

// Spec the specification for artisan artefacts to be exported
type Spec struct {
	Version  string            `yaml:"version"`
	Images   map[string]string `yaml:"images"`
	Packages map[string]string `yaml:"packages"`

	content []byte
}

func NewSpec(path string) (*Spec, error) {
	// finds the absolute path
	specFile, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("cannot get absolute path: %s", err)
	}
	// appends spec filename
	specFile = filepath.Join(specFile, "spec.yaml")
	// reads spec.yaml
	content, err := os.ReadFile(specFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read spec file: %s", err)
	}
	spec := new(Spec)
	// unmarshal yaml
	err = yaml.Unmarshal(content, spec)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal spec file: %s", err)
	}
	// set the content of the spec file for later use
	spec.content = content
	return spec, nil
}

func (s *Spec) Save(targetUri, sourceCreds, targetCreds string) error {
	// first, save the spec to the target location
	uri := fmt.Sprintf("%s/spec.yaml", targetUri)
	err := core.WriteFile(s.content, uri, targetCreds)
	if err != nil {
		return fmt.Errorf("cannot save spec file: %s", err)
	}
	core.InfoLogger.Println("spec.yaml")
	// target uri should not have a file extension
	if len(filepath.Ext(targetUri)) > 0 {
		return fmt.Errorf("target URI must to have a file extension")
	}
	// save packages first
	l := registry.NewLocalRegistry()
	for key, value := range s.Packages {
		name, err := core.ParseName(value)
		if err != nil {
			return fmt.Errorf("invalid package name: %s", err)
		}
		uri := fmt.Sprintf("%s/%s.tar", targetUri, key)
		err = l.Save([]core.PackageName{*name}, sourceCreds, uri, targetCreds)
		if err != nil {
			fmt.Errorf("cannot save package %s: %s", value, err)
		}
		core.InfoLogger.Println(value)
	}
	// save images
	for key, value := range s.Images {
		uri := fmt.Sprintf("%s/%s.tar", targetUri, key)
		err := SaveImage(value, value, uri, targetCreds)
		if err != nil {
			return fmt.Errorf("cannot save image %s: %s", value, err)
		}
		core.InfoLogger.Println(value)
	}
	return nil
}
