package main

/*
  Onix Config Manager - Pilot Control Discovery
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

import (
	_ "github.com/gatblau/onix/pilotdisco/docs"
	"net/http"
)

// @title Pilot Control Discovery
// @version 0.0.4
// @description Onix Config Manager Discovery Service for Pilot Control
// @contact.name gatblau
// @contact.url http://onix.gatblau.org/
// @contact.email onix@gatblau.org
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @Summary Discover Pilot Control service
// @Description discovers the Pilot Control service allocated to the Pilot making the request and if successful,
// @Description admits the Pilot into service
// @Tags Discovery
// @Router /disco [post]
// @Accepts json
// @Produce json
// @Failure 401 {string} authentication failed
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 201 {string} OK
func discoveryHandler(writer http.ResponseWriter, request *http.Request) {

}
