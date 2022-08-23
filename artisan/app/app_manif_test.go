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
	"testing"
)

func TestAppManifest_Explode(t *testing.T) {
	m, err := NewAppMan("./test/onix.yaml", "", "", "")
	if err != nil {
		t.Fatalf("cannot create app manifest: %s\n", err)
	}
	fmt.Println(len(m.Services))
}

func TestAppManifest_ExplodeHTTP(t *testing.T) {
	m, err := NewAppMan("https://raw.githubusercontent.com/gatblau/onix/dev/deploy/onix.yaml", "full", "", "")
	if err != nil {
		t.Fatalf("cannot create app manifest: %s\n", err)
	}
	fmt.Println(len(m.Services))
}

func TestNewAppMan(t *testing.T) {
	m, err := NewAppMan("test2/app.yaml", "config-db", "", "")
	if err != nil {
		t.Fatalf("cannot create app manifest: %s\n", err)
	}
	fmt.Println(len(m.Services))
}
