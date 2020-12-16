/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package build

import (
	"testing"
)

func TestBuild(t *testing.T) {
	// p := NewBuilder()
	// p.Build("/Users/andresalos/go/src/github.com/gatblau/onix", "artie", "", core.ParseName("gatblau/artie"), "", false)
}

func TestRun(t *testing.T) {
	p := NewBuilder()
	p.Run("test", "/Users/andresalos/go/src/github.com/gatblau/onix/artisan/test/test", false)
}
