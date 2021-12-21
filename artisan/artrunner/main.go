/*
  Onix Config Manager - Artisan Runner
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package main

import (
	"github.com/gatblau/onix/oxlib/httpserver"
	"github.com/gorilla/mux"
)

func main() {
	// creates a generic http server
	s := httpserver.New("art-runner")
	// add handlers
	s.Http = func(router *mux.Router) {
		router.HandleFunc("/flow", createFlowFromPayloadHandler).Methods("POST")
		router.HandleFunc("/flow/name/{flow-name}/ns/{namespace}", runFlowHandler).Methods("POST")
		router.HandleFunc("/flow/key/{flow-key}/ns/{namespace}", createFlowFromConfigHandler).Methods("POST")
		router.HandleFunc("/flow/key/{flow-key}", getFlowHandler).Methods("GET")
	}
	s.Serve()
}
