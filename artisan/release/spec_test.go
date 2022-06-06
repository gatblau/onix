/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package release

import (
	"fmt"
	"github.com/gatblau/onix/artisan/data"
	"gopkg.in/yaml.v2"
	"os"
	"testing"
)

func TestSpec_SaveSpec(t *testing.T) {
	s, err := NewSpec(".", "")
	if err != nil {
		t.Fatal(err)
	}
	err = ExportSpec(ExportOptions{s, "", "minioadmin:minioadmin", "", "", ""})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSpec_ImportSpec(t *testing.T) {
	_, err := ImportSpec(ImportOptions{"s3://localhost:9000/app1/v1", "minioadmin:minioadmin", "", "", nil, ""})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSpec_ExportSpec(t *testing.T) {
	s, err := NewSpec("../../deploy/1.0", "")
	if err != nil {
		t.Fatal(err)
	}
	err = ExportSpec(ExportOptions{s, "", "minioadmin:minioadmin", "", "", ""})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSpec_Unmarshal(t *testing.T) {
	specBytes, err := os.ReadFile("spec.yaml")
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
	_, _ = DownloadSpec(UpDownOptions{"s3://s3.cmsee.cloud/cmsee/1.0.2", " a667925:YWUwM2ExNDhiMmUy", "./app1-spec"})
}

func TestSpec_Push(t *testing.T) {
	_ = PushSpec(PushOptions{"./", "localhost:5000", "", "", "", true, false, false, ""})
}

func TestSerialise(t *testing.T) {
	spec := &Spec{
		Name:        "App 1",
		Description: "Description of App 1",
		Version:     "1.0",
		Info:        "This is the first release",
		Author:      "ACME Ltd",
		Images: map[string]string{
			"IMAGE1": "quay.io/images/image-1:latest",
			"IMAGE2": "quay.io/images/image-2:latest",
		},
		Packages: map[string]string{
			"PKG1": "localhost:8082/pkg/pkg-1:latest",
			"PKG2": "localhost:8082/pkg/pkg-2:latest",
		},
		Run: []Run{
			{
				Package:  "pk1",
				Function: "fx1",
				Input: &data.Input{
					Var: data.Vars{
						{
							Name:  "Key1",
							Value: "Value1",
						},
					},
				},
				Event: ReleaseDeploy,
			},
		},
	}
	bytes, _ := yaml.Marshal(spec)
	fmt.Println(string(bytes[:]))
}
