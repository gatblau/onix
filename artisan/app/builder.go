/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Builder the contract for application deployment configuration builders
type Builder interface {
	Build() ([]DeploymentRsx, error)
}

func NewBuilder(builder BuilderType, appManifest Manifest) (Builder, error) {
	switch builder {
	case DockerCompose:
		return newComposeBuilder(appManifest), nil
	case Kubernetes:
		return newKubeBuilder(appManifest), nil
	}
	return nil, fmt.Errorf("builder type not supported\n")
}

type DeploymentRsxType int

const (
	ComposeProject DeploymentRsxType = iota
	K8SResource
	EnvironmentFile
	ConfigurationFile
	SvcInitScript
	DbInitScript
	DeployScript
	BuildFile
)

type DeploymentRsx struct {
	Name    string
	Content []byte
	Type    DeploymentRsxType
}

type BuilderType int

const (
	DockerCompose BuilderType = iota
	Kubernetes
)

// GenerateResources generates application deployment resources for a particular platform specified by the format parameter
// uri: the uri of the application manifest to use to generate resources
// format: the resources platform format (e.g. compose, k8s)
// profile: the application profile name, that describes the services to generate from the application manifest
// creds: credentials to get the app manifest from git service in the format uname:pwd - if not required pass empty string
// path: the file path where the resources should be saved
func GenerateResources(uri, format, profile, creds, path, artHome string) error {
	// create an application manifest
	manifest, err := NewAppMan(uri, profile, creds, artHome)
	if err != nil {
		return err
	}
	// create a builder
	var builderType BuilderType
	switch strings.ToLower(format) {
	case "compose":
		builderType = DockerCompose
	case "k8s":
		builderType = Kubernetes
	default:
		return fmt.Errorf("invalid format, valid formats are compose or k8s")
	}
	builder, err := NewBuilder(builderType, *manifest)
	if err != nil {
		return err
	}
	// build the app deployment resources
	files, err := builder.Build()
	if err != nil {
		return err
	}
	// work out a target path
	path, err = filepath.Abs(path)
	if err != nil {
		return err
	}
	// ensure path exists
	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("cannot create folder '%s': %s\n", path, err)
		}
	}
	// write files to disk
	for _, file := range files {
		// if the path contains a directory
		if !isFilename(file.Name) {
			// if it is absolute
			if isAbs(file.Name) {
				// makes the directory relative (removes leading slash)
				file.Name = file.Name[1:]
			}
			// creates the relative directory
			dir, err2 := filepath.Abs(filepath.Join(path, filepath.Dir(file.Name)))
			if err2 != nil {
				return err
			}
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return err
			}
		}
		fpath := filepath.Join(path, file.Name)
		err = os.WriteFile(fpath, file.Content, os.ModePerm)
		if err != nil {
			return fmt.Errorf("cannot write configuration file %s: %s\n", fpath, err)
		}
	}
	return nil
}

func isAbs(path string) bool {
	// the path is considered absolute if:
	// a. starts with a forward slash
	// b. contains at least two forward slashes
	return len(strings.Split(path, "/")) > 1 && path[0] == '/'
}

func isFilename(path string) bool {
	// the path is considered a filename only if it does not contain any forward slashes
	return strings.Index(path, "/") == -1
}
