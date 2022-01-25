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
	"github.com/gatblau/onix/artisan/build"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/merge"
	"github.com/gatblau/onix/artisan/registry"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
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
		name, err2 := core.ParseName(value)
		if err2 != nil {
			return fmt.Errorf("invalid package name: %s", err)
		}
		uri = fmt.Sprintf("%s/%s.tar", targetUri, key)
		err = l.Save([]core.PackageName{*name}, sourceCreds, uri, targetCreds)
		if err != nil {
			return fmt.Errorf("cannot save package %s: %s", value, err)
		}
		core.InfoLogger.Println(value)
	}
	// save images
	for key, value := range s.Images {
		uri = fmt.Sprintf("%s/%s.tar", targetUri, key)
		err = SaveImage(value, value, uri, targetCreds)
		if err != nil {
			return fmt.Errorf("cannot save image %s: %s", value, err)
		}
		core.InfoLogger.Println(value)
	}
	return nil
}

func ImportSpec(targetUri, targetCreds, localPath string) error {
	r := registry.NewLocalRegistry()
	uri := fmt.Sprintf("%s/spec.yaml", targetUri)
	specBytes, err := core.ReadFile(uri, targetCreds)
	if err != nil {
		return fmt.Errorf("cannot read spec.yaml: %s", err)
	}
	spec := new(Spec)
	err = yaml.Unmarshal(specBytes, spec)
	if err != nil {
		return fmt.Errorf("cannot unmarshal spec.yaml: %s", err)
	}
	// if the uri is s3 allows using localPath
	if strings.HasPrefix(targetUri, "s3") && len(localPath) > 0 {
		path, err2 := filepath.Abs(localPath)
		if err2 != nil {
			return err2
		}
		// if the path does not exist
		if _, err = os.Stat(path); os.IsNotExist(err) {
			// creates it
			err = os.MkdirAll(path, 0755)
			if err != nil {
				return err
			}
		}
		localPath = path
		err = os.WriteFile(filepath.Join(localPath, "spec.yaml"), specBytes, 0755)
		if err != nil {
			return err
		}
	} else {
		// otherwise, return error
		return fmt.Errorf("local path cannot be specified if URI is not s3")
	}
	// import packages
	for k, _ := range spec.Packages {
		name := fmt.Sprintf("%s/%s.tar", targetUri, k)
		err2 := r.Import([]string{name}, targetCreds, localPath)
		if err2 != nil {
			return fmt.Errorf("cannot read %s.tar: %s", k, err2)
		}
		core.InfoLogger.Println(name)
	}
	// import images
	for k, _ := range spec.Images {
		name := fmt.Sprintf("%s/%s.tar", targetUri, k)
		err2 := r.Import([]string{name}, targetCreds, localPath)
		if err2 != nil {
			return fmt.Errorf("cannot read %s.tar: %s", k, err)
		}
		core.InfoLogger.Println(name)
	}
	// import images
	for _, name := range spec.Images {
		_, err2 := build.Exe(fmt.Sprintf("art exe %s import", name), ".", merge.NewEnVarFromSlice([]string{}), false)
		if err2 != nil {
			return fmt.Errorf("cannot import image %s: %s", name, err)
		}
		core.InfoLogger.Println(name)
	}
	return nil
}
