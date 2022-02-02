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

func NewSpec(path, creds string) (*Spec, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("cannot get absolute path: %s", err)
	}
	specFile := filepath.Join(path, "spec.yaml")
	content, err := core.ReadFile(specFile, creds)
	if err != nil {
		return nil, fmt.Errorf("cannot read spec file %s: %s", specFile, err)
	}
	spec := new(Spec)
	err = yaml.Unmarshal(content, spec)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal spec file: %s", err)
	}
	spec.content = content
	return spec, nil
}

func (s *Spec) Export(targetUri, sourceCreds, targetCreds string) error {
	// first, save the spec to the target location
	uri := fmt.Sprintf("%s/spec.yaml", targetUri)
	err := core.WriteFile(s.content, uri, targetCreds)
	if err != nil {
		return fmt.Errorf("cannot save spec file: %s", err)
	}
	core.InfoLogger.Printf("writing spec.yaml to %s", targetUri)
	// save packages first
	l := registry.NewLocalRegistry()
	for _, value := range s.Packages {
		name, err2 := core.ParseName(value)
		if err2 != nil {
			return fmt.Errorf("invalid package name: %s", err)
		}
		uri = fmt.Sprintf("%s/%s.tar", targetUri, pkgName(value))
		err = l.ExportPackage([]core.PackageName{*name}, sourceCreds, uri, targetCreds)
		if err != nil {
			return fmt.Errorf("cannot save package %s: %s", value, err)
		}
	}
	// save images
	for _, value := range s.Images {
		// note: the package is saved with a name exactly the same as the container image
		// to avoid the art package name parsing from failing, any images with no host or user/group in the name should be avoided
		// e.g. docker.io/mongo-express:latest will fail so use docker.io/library/mongo-express:latest instead
		err = ExportImage(value, value, targetUri, targetCreds)
		if err != nil {
			return fmt.Errorf("cannot save image %s: %s", value, err)
		}
	}
	return nil
}

func ImportSpec(targetUri, targetCreds string) error {
	r := registry.NewLocalRegistry()
	uri := fmt.Sprintf("%s/spec.yaml", targetUri)
	core.InfoLogger.Printf("retrieving %s\n", uri)
	specBytes, err := core.ReadFile(uri, targetCreds)
	if err != nil {
		return fmt.Errorf("cannot read spec.yaml: %s", err)
	}
	spec := new(Spec)
	err = yaml.Unmarshal(specBytes, spec)
	if err != nil {
		return fmt.Errorf("cannot unmarshal spec.yaml: %s", err)
	}
	// import packages
	for _, pkName := range spec.Packages {
		name := fmt.Sprintf("%s/%s.tar", targetUri, pkgName(pkName))
		err2 := r.Import([]string{name}, targetCreds)
		if err2 != nil {
			return fmt.Errorf("cannot read %s.tar: %s", pkgName(pkName), err2)
		}
	}
	// import images
	for _, image := range spec.Images {
		name := fmt.Sprintf("%s/%s.tar", targetUri, pkgName(image))
		err2 := r.Import([]string{name}, targetCreds)
		if err2 != nil {
			return fmt.Errorf("cannot read %s.tar: %s", pkgName(image), err)
		}
		core.InfoLogger.Printf("loading => %s\n", image)
		_, err2 = build.Exe(fmt.Sprintf("art exe %s import", image), ".", merge.NewEnVarFromSlice([]string{}), false)
		if err2 != nil {
			return fmt.Errorf("cannot import image %s: %s", image, err2)
		}
	}
	return nil
}

func DownloadSpec(targetUri, targetCreds, localPath string) error {
	if !strings.Contains(targetUri, "://") {
		return fmt.Errorf("invalid URI, it must have an scheme (e.g. scheme://)")
	}
	spec, err := NewSpec(targetUri, targetCreds)
	if err != nil {
		return fmt.Errorf("cannot load specification: %s", err)
	}
	if err = checkPath(localPath); err != nil {
		return fmt.Errorf("cannot create local path: %s", err)
	}
	err = os.WriteFile(filepath.Join(localPath, "spec.yaml"), spec.content, 0755)
	if err != nil {
		return err
	}
	for _, pkg := range spec.Packages {
		pkgUri := fmt.Sprintf("%s/%s.tar", targetUri, pkgName(pkg))
		core.InfoLogger.Printf("downloading => %s\n", pkgUri)
		tarBytes, err2 := core.ReadFile(pkgUri, targetCreds)
		if err2 != nil {
			return err2
		}
		pkgPath := filepath.Join(localPath, filepath.Base(targetUri))
		core.InfoLogger.Printf("writing => %s\n", pkgPath)
		err = os.WriteFile(pkgPath, tarBytes, 0755)
		if err != nil {
			return err
		}
	}
	for _, image := range spec.Images {
		imageUri := fmt.Sprintf("%s/%s.tar", targetUri, pkgName(image))
		core.InfoLogger.Printf("downloading => %s\n", imageUri)
		tarBytes, err2 := core.ReadFile(imageUri, targetCreds)
		if err2 != nil {
			return err2
		}
		targetFile := filepath.Join(localPath, fmt.Sprintf("%s.tar", pkgName(image)))
		core.InfoLogger.Printf("writing => %s\n", targetFile)
		err = os.WriteFile(targetFile, tarBytes, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func PushSpec(specPath, host, group, creds string, image bool) error {
	local := registry.NewLocalRegistry()
	if strings.Contains(specPath, "://") {
		return fmt.Errorf("spec path must be a location in the file system")
	}
	// load the spec
	spec, err := NewSpec(specPath, "")
	if err != nil {
		return fmt.Errorf("cannot load spec.yaml: %s", err)
	}
	if !image {
		for _, pac := range spec.Packages {
			tgtNameStr, tgtName, tgtNameErr := targetName(pac, group, host)
			if tgtNameErr != nil {
				return fmt.Errorf("cannot work out target name: %s", tgtNameErr)
			}
			// tag the package with the target registry name
			core.InfoLogger.Printf("tagging => '%s' to '%s'\n", pac, tgtNameStr)
			err = local.Tag(pac, tgtNameStr)
			if err != nil {
				return err
			}
			// push to remote
			core.InfoLogger.Printf("pushing => '%s'\n", tgtNameStr)
			err = local.Push(tgtName, creds)
			if err != nil {
				return err
			}
		}
	} else {
		if len(creds) > 0 {
			return fmt.Errorf("credentials specified but not used, for images ensure you are logged to the destination registry")
		}
		cli, cmdErr := containerCmd()
		if cmdErr != nil {
			return cmdErr
		}
		for _, img := range spec.Images {
			tgtNameStr, _, tgtNameErr := targetName(img, group, host)
			if tgtNameErr != nil {
				return fmt.Errorf("cannot work out target name: %s", tgtNameErr)
			}
			// tag the package with the target registry name
			core.InfoLogger.Printf("tagging => '%s' to '%s'\n", img, targetName)
			// docker tag image
			_, err = build.Exe(fmt.Sprintf("%s tag %s %s", cli, img, tgtNameStr), ".", merge.NewEnVarFromSlice([]string{}), false)
			if err != nil {
				return err
			}
			// push to remote
			core.InfoLogger.Printf("pushing => '%s'\n", tgtNameStr)
			_, err = build.Exe(fmt.Sprintf("%s push %s", cli, tgtNameStr), ".", merge.NewEnVarFromSlice([]string{}), false)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func targetName(name string, group string, host string) (string, *core.PackageName, error) {
	srcName, parseErr := core.ParseName(name)
	if parseErr != nil {
		return "", nil, parseErr
	}
	var tName string
	// if a target group has been defined
	if len(group) > 0 {
		// build a target name using host and group but keeping source package name and tag
		tName = fmt.Sprintf("%s/%s/%s:%s", host, group, srcName.Name, srcName.Tag)
	} else {
		// if a target group has not been specified, then
		// build a target name using host but keeping source package group, name and tag
		tName = fmt.Sprintf("%s/%s/%s:%s", host, srcName.Group, srcName.Name, srcName.Tag)
	}
	tgtName, parseTargetErr := core.ParseName(tName)
	if parseTargetErr != nil {
		return "", nil, parseTargetErr
	}
	return tName, tgtName, nil
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

func checkPath(path string) error {
	p, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	// if the path does not exist
	if _, err = os.Stat(p); os.IsNotExist(err) {
		// creates it
		err = os.MkdirAll(p, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
