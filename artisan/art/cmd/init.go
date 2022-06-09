/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	"github.com/gatblau/onix/artisan/core"
	"log"
)

func init() {
	// ensure the registry folder structure is in place
	if err := core.EnsureRegistryPath(""); err != nil {
		log.Fatal("cannot run artisan without a local registry, its creation failed: %", err)
	}
}
