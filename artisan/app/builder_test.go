/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package app

import (
	"testing"
)

func TestGenerateResources(t *testing.T) {
	if err := GenerateResources("./test/onix.yaml", "compose", "full-test-data", "", "onix-compose"); err != nil {
		t.Fatalf(err.Error())
	}
	// if err := GenerateResources("./test/artisan_registry.yaml", "compose", "full", "", "art-reg"); err != nil {
	// 	t.Fatalf(err.Error())
	// }
	// if err := GenerateResources("https://raw.githubusercontent.com/gatblau/onix/dev/deploy/onix.yaml", "k8s", "full", "", "onix-k8s"); err != nil {
	// 	t.Fatalf(err.Error())
	// }
}
