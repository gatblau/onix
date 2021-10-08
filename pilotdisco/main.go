package main

/*
  Onix Config Manager - Pilot Control Discovery
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"github.com/gatblau/onix/oxlib/httpserver"
	"github.com/gorilla/mux"
)

func main() {
	// creates a generic http server
	s := httpserver.New("pilot-disco")
	// add handlers
	s.Http = func(router *mux.Router) {
		// enable encoded path  vars
		router.UseEncodedPath()
		// middleware
		// router.Use(s.LoggingMiddleware)
		router.Use(s.AuthenticationMiddleware)

		// pilot http handlers
		router.HandleFunc("/disco", discoveryHandler).Methods("POST")
	}
	s.Serve()
}
