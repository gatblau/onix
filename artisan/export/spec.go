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
	Images   map[string]string `yaml:"images,omitempty"`
	Packages map[string]string `yaml:"packages,omitempty"`

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
	for _, value := range s.Images {
		err = SaveImage(value, pkgName(value), targetUri, targetCreds)
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
	// if the uri is s3 allows using localPath only if local path provided
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

	}
	// import packages
	for _, pkName := range spec.Packages {
		name := fmt.Sprintf("%s/%s.tar", targetUri, pkgName(pkName))
		err2 := r.Import([]string{name}, targetCreds, localPath)
		if err2 != nil {
			return fmt.Errorf("cannot read %s.tar: %s", pkgName(pkName), err2)
		}
		core.InfoLogger.Println(name)
	}
	// import images
	for _, image := range spec.Images {
		name := fmt.Sprintf("%s/%s.tar", targetUri, pkgName(image))
		err2 := r.Import([]string{name}, targetCreds, localPath)
		if err2 != nil {
			return fmt.Errorf("cannot read %s.tar: %s", pkgName(image), err)
		}
		core.InfoLogger.Println(name)
	}
	// import images
	for _, name := range spec.Images {
		_, err2 := build.Exe(fmt.Sprintf("art exe %s import", pkgName(name)), ".", merge.NewEnVarFromSlice([]string{}), false)
		if err2 != nil {
			return fmt.Errorf("cannot import image %s: %s", name, err2)
		}
		core.InfoLogger.Println(name)
	}
	return nil
}

func (s *Spec) ContainsImage(name string) bool {
	for key, _ := range s.Images {
		if name == key {
			return true
		}
	}
	return false
}

func pkgName(name string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(name, "/", "_"), ".", "_"), "-", "_")
}
