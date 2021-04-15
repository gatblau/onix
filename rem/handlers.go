package main

/*
  Onix Config Manager - REMote Host Service
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

// @title Onix Remote Host
// @version 0.0.4
// @description Remote Ctrl Service for Onix Pilot
// @contact.name gatblau
// @contact.url http://onix.gatblau.org/
// @contact.email onix@gatblau.org
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

import (
	"github.com/gatblau/onix/artisan/server"
	"github.com/gatblau/onix/rem/core"
	_ "github.com/gatblau/onix/rem/docs"
	"github.com/gorilla/mux"
	"net/http"
)

// @Summary Host Ping
// @Description receives a periodic ping request from Onix Pilot
// @Tags Pilot
// @Router /ping/{host-key} [post]
// @Produce json
// @Param host-key path string true "the unique key for the host"
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func pingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_ = vars["host-key"]
	var commands []core.Command
	server.Write(w, r, commands)
}

// @Summary Host Query
// @Description Returns a list of managed hosts
// @Tags Admin
// @Router /host [get]
// @Produce json
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func hostQueryHandler(w http.ResponseWriter, r *http.Request) {
	hosts := []core.Host{
		core.Host{
			Name:      "HOST-001",
			Customer:  "CUST-01",
			Region:    "UK-North-West",
			Location:  "Manchester",
			Connected: true,
			Up:        true,
		},
		core.Host{
			Name:      "HOST-002",
			Customer:  "CUST-01",
			Region:    "UK-North-West",
			Location:  "Manchester",
			Connected: false,
			Up:        false,
		},
		core.Host{
			Name:      "HOST-003",
			Customer:  "CUST-01",
			Region:    "UK-South-West",
			Location:  "Devon",
			Connected: false,
			Up:        true,
		},
	}
	server.Write(w, r, hosts)
}
