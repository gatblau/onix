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
	"gopkg.in/yaml.v2"
	"testing"
)

func TestSpec_SaveSpec(t *testing.T) {
	s, err := NewSpec(".", "")
	if err != nil {
		t.Fatal(err)
	}
	err = ExportSpec(*s, "", "minioadmin:minioadmin", "", "")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSpec_ImportSpec(t *testing.T) {
	err := ImportSpec("s3://localhost:9000/app1/v1", "minioadmin:minioadmin", "", "", false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSpec_ExportSpec(t *testing.T) {
	s, err := NewSpec("../../deploy/1.0", "")
	if err != nil {
		t.Fatal(err)
	}
	err = ExportSpec(*s, "", "minioadmin:minioadmin", "", "")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSpec_Unmarshal(t *testing.T) {
	specBytes, err := core.ReadFile("spec.yaml", "")
	if err != nil {
		t.Fatal(fmt.Errorf("cannot read spec.yaml: %s", err))
	}
	spec := new(Spec)
	err = yaml.Unmarshal(specBytes, spec)
	if err != nil {
		t.Fatal(fmt.Errorf("cannot unmarshal spec.yaml: %s", err))
	}
}

func TestDownloadSpec(t *testing.T) {
	DownloadSpec("s3://s3.cmsee.cloud/cmsee/1.0.2", " a667925:YWUwM2ExNDhiMmUy", "./app1-spec")
}

func TestSpec_Push(t *testing.T) {
	PushSpec("./", "localhost:5000", "", "", "", true, false, false)
}
