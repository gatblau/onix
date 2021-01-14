/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package flow

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/registry"
	"testing"
)

func TestSample(t *testing.T) {
	// creates a generator from a file
	m, err := NewFromPath("/Users/andresalos/go/src/github.com/gatblau/onix/artisan/test/art/flow/p2i-flow-merged.yaml", "id_rsa_pub.pgp", "")
	core.CheckErr(err, "failed to create generator")
	m.FillIn(registry.NewLocalRegistry())
	fmt.Println(m.flow.Steps)
}
