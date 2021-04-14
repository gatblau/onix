/*
  Onix Config Manager - REMote Ctrl Service
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package main

// @title Onix Remote
// @version 0.0.4
// @description Remote Ctrl Service for Onix Pilot
// @contact.name gatblau
// @contact.url http://onix.gatblau.org/
// @contact.email onix@gatblau.org
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

import (
	"fmt"
	_ "github.com/gatblau/onix/rem/docs"
	"io/ioutil"
	"net/http"
)

// @Summary heart bit
// @Description receives a heart bit from a host agent via HTTP
// @Tags Agent
// @Router /beat [post]
// @Produce plain
// @Param host-key body string true "the unique key identifying the host sending the heart beat"
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func beatHandler(w http.ResponseWriter, r *http.Request) {
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot read request payload: %s", err), http.StatusInternalServerError)
		return
	}
}
