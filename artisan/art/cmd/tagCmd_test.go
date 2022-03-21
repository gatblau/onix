/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	"github.com/gatblau/onix/artisan/build"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/registry"
	"testing"
)

// TestTagV1ToLatestExist test that a package with a V1 tag can be tagged to latest when a previous latest tag exists
// and the existing tag is renamed so that there is no dangling packages
func TestTagV1ToLatestExist(t *testing.T) {
	// pre-conditions
	reg := registry.NewLocalRegistry()
	// cleanup
	testLatest, _ := core.ParseName("test:latest")
	testV1, _ := core.ParseName("test:V1")
	// clear registry if previous packages exist
	// TODO: ensure all packages are removed
	reg.Remove([]string{"test:latest", "test:V1"})
	// build latest
	builder := build.NewBuilder()
	builder.Build(".", "", "", testLatest, "test1", false, false, "", false)
	// build V1
	builder.Build(".", "", "", testV1, "test1", false, false, "", false)
	// reload the registry
	reg.Load()
	testV1PId := reg.FindPackage(testV1).Id
	testLatestPId := reg.FindPackage(testLatest).Id
	// execute action tag
	reg.Tag("test:V1", "test:latest")
	// reload the registry
	reg.Load()
	// check post-conditions
	testLatestP := reg.FindPackageNamesById(testLatestPId)
	if testLatestP == nil {
		t.Fatalf("test:latest package not found")
	}
	// the old latest renamed to avoid dangling
	if len(testLatestP) != 1 {
		t.Fatalf("old test:latest package should have only one renamed tag")
	}
	testV1P := reg.FindPackageNamesById(testV1PId)
	// the new latest tag added on top of existing package with V1 tag
	if len(testV1P) != 2 {
		t.Fatalf("")
	}
}

// TestTagV1ToLatest test that a package with a V1 tag can be tagged to latest when a previous latest tag does not exist
func TestTagV1ToLatest(t *testing.T) {
	reg := registry.NewLocalRegistry()
	// cleanup
	testLatest, _ := core.ParseName("test:latest")
	testV1, _ := core.ParseName("test:V1")
	reg.Remove([]string{"test:latest", "test:V1"})
	// build latest
	builder := build.NewBuilder()
	// build V1
	builder.Build(".", "", "", testV1, "test1", false, false, "", false)
	// reload the registry
	reg.Load()
	// tag
	reg.Tag("test:V1", "test:latest")
	// check post-conditions
	if reg.FindPackage(testLatest) == nil {
		t.Fatalf("test:latest package not found")
	}
	if reg.FindPackage(testV1) == nil {
		t.Fatalf("test:V1 package not found")
	}
}
