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
	"testing"
)

func TestComposeBuilder(t *testing.T) {
	m, err := NewAppMan("https://raw.githubusercontent.com/gatblau/onix/dev/deploy/onix.yaml", "micro", "")
	if err != nil {
		t.Fatalf("cannot create app manifest: %s\n", err)
	}
	builder, err := NewBuilder(DockerCompose, *m)
	if err != nil {
		t.Fatalf("cannot create resource builder: %s\n", err)
	}
	resx, err := builder.Build()
	if err != nil {
		t.Fatalf("cannot build docker compose resources: %s\n", err)
	}
	fmt.Println(string(resx[0].Content[:]))
	fmt.Println(string(resx[1].Content[:]))
}
