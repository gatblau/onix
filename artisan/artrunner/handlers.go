/*
  Onix Config Manager - Artisan Runner
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package main

import (
	"net/http"
)

// @Summary Executes an Artisan flow
// @Description uploads an Artisan flow and triggers the flow execution
// @Tags Flows
// @Router /flow [post]
// @Param tag path string true "the artefact reference name"
func runHandler(w http.ResponseWriter, _ *http.Request) {

}
