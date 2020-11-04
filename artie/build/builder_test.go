/*
  Onix Config Manager - Artie
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package build

import (
	"github.com/gatblau/onix/artie/core"
	"testing"
)

func TestBuild(t *testing.T) {
	p := NewBuilder()
	// p.Build("https://github.com/gatblau/boot", "", core.ParseName("boot"), "")
	p.Build("https://gitlab.com/sap7/sapg-plugin.git", "", "hcPVXRkULPTbUWG2gxV5", core.ParseName("sap-nexusrepo-amosonline-aws-01-sapgatewaycd.apps.amosds.amosonline.io/sap/sap-jco"), "")
	// p.Build("/Users/andresalos/artie-test/boot", "", "", core.ParseName("localhost:8081/gatblau/boot"), "")
}
