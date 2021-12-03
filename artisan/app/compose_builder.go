/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package app

import (
	"fmt"
	"github.com/compose-spec/compose-go/types"
	"gopkg.in/yaml.v2"
	"strings"
)

type ComposeBuilder struct {
	Manifest Manifest
}

// newComposeBuilder called internally by NewBuilder()
func newComposeBuilder(appMan Manifest) Builder {
	return &ComposeBuilder{Manifest: appMan}
}

func (b *ComposeBuilder) Build() ([]DeploymentRsx, error) {
	p := b.buildProject()
	composeProject, err := yaml.Marshal(p)
	if err != nil {
		return nil, err
	}
	return []DeploymentRsx{
		{
			Name:    "docker-compose.yml",
			Content: composeProject,
			Type:    ComposeProject,
		},
	}, nil
}

func (b *ComposeBuilder) buildProject() types.Project {
	p := types.Project{}
	p.Name = fmt.Sprintf("Docker Compose Project for %s", strings.ToUpper(b.Manifest.Name))
	return p
}
