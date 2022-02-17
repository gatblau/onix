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
	"github.com/gatblau/onix/artisan/runner/host/handlers"

)

func main() {
	// creates a generic http server
	handlerMgr := handlers.NewHandlerManager()
	s := httpserver.New("art-host-runner")
	// add handlers
	s.Http = func(router *mux.Router) {
		router.Handle("/host/{package}/{function}", handlerMgr).Methods("POST")
	}
	s.Serve()
}
