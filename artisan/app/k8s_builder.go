/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package app

type K8SBuilder struct {
	manifest Manifest
}

// newComposeBuilder called internally by NewBuilder()
func newK8SBuilder(appMan Manifest) Builder {
	return &K8SBuilder{manifest: appMan}
}

func (b *K8SBuilder) Build() ([]DeploymentRsx, error) {
	rsx := make([]DeploymentRsx, 0)
	return rsx, nil
}
