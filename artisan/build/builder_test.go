/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package build

import (
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/registry"
	"testing"
)

func TestBuild(t *testing.T) {
	p := NewBuilder()
	p.Build(".", "", "", core.ParseName("artisan"), "linux", false, false, "")
	l := registry.NewLocalRegistry()
	l.Open(core.ParseName("artisan"), "", false, "test", "", true)
}

func TestRun(t *testing.T) {
	p := NewBuilder()
	p.Run("test", "/Users/andresalos/go/src/github.com/gatblau/onix/artisan/build", false)
}

func TestExec(t *testing.T) {
	p := NewBuilder()
	p.Execute(
		core.ParseName("localhost:8082/gatblau/art-buildah"),
		"release-image",
		"",
		false,
		"",
		false,
		true)
}

func TestExec2(t *testing.T) {
	p := NewBuilder()
	p.Execute(
		core.ParseName("artisan-registry-amosonline-aws-01-sapgatewaycd.apps.amosds.amosonline.io/sap/sap-equip-jvm-ctx"),
		"build-image",
		"admin:nxrpsap",
		true,
		"",
		false,
		false)
}

func TestRunInContainer(t *testing.T) {
	err := RunInContainer("quay.io/gatblau/buildah", "localhost:8082/gatblau/art-buildah", "build-image")
	if err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
}
