/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package app

import "fmt"

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
