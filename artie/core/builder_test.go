/*
  Onix Config Manager - Artie
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

import "testing"

func TestPack(t *testing.T) {
	p := NewBuilder()
	p.Build("https://github.com/gatblau/boot", "", "nexus.io/gatblau/boot:linux-01", "linux")
	//p.Build("/Users/andresalos/test-pak/boot", "", "nexus.io/gatblau/boot")
}
