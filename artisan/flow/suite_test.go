/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package flow

import (
	"github.com/gatblau/onix/artisan/core"
	"testing"
)

func TestSample(t *testing.T) {
	// creates a generator from a file
	g, err := NewFromPath("ci-flow.yaml", "build.yaml")
	core.CheckErr(err, "failsed to create generator")
	g.FillIn()

}
