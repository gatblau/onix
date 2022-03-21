/*
  Onix Config Manager - Artisan Host Runner
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package main

// @title Artisan Host Runner
// @version 0.0.4
// @description Run Artisan packages with in a host
// @contact.name gatblau
// @contact.url http://onix.gatblau.org/
// @contact.email onix@gatblau.org
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

import (
	"net/http"

	_ "github.com/gatblau/onix/artisan/runner/host/docs"
	"github.com/gatblau/onix/artisan/runner/host/handlers"
)

// @Summary Build patching artisan package
// @Description Trigger a new build to create artisan package from the vulnerabilty scanned csv report passed in the payload.
// @Tags Runners
// @Router /host/{cmd-key} [post]
// @Param cmd-key path string true "the key of the command to retrieve"
// @Produce plain
// @Param flow body flow.Flow true "the artisan flow to run"
// @Failure 500 {string} there was an error in the server, error the server logs
// @Success 200 {string} OK

func createOSPatchingHandler(w http.ResponseWriter, r *http.Request) {

	osph := handlers.OSpatchingHandler{}
	osph.HandleEvent(w, r)
}
