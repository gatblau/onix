/*
  Onix Config Manager - Artisan's Doorman
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package main

import (
	artCore "github.com/gatblau/onix/artisan/core"
	"log"
	"path/filepath"
)

var D *Doorman

func main() {
	if err := checkDoormanHome(); err != nil {
		log.Fatalf("cannot launch  doorman, cannot write to file system: %s", err)
	}
	D = NewDoorman(NewDefaultProcFactory())
	D.RegisterHandlers()
	D.Start()
}

func checkDoormanHome() error {
	path := filepath.Join(artCore.HomeDir(), ".doorman")
	return artCore.EnsureRegistryPath(path)
}
