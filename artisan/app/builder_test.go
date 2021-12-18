/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package app

import (
	"testing"
)

func TestGenerateResources(t *testing.T) {
	// if err := GenerateResources("./test/onix.yaml", "compose", "micro", "", "onix"); err != nil {
	// 	t.Fatalf(err.Error())
	// }
	if err := GenerateResources("./test/artisan_registry.yaml", "compose", "nexus", "", "nexus"); err != nil {
		t.Fatalf(err.Error())
	}
}
