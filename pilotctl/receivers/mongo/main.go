package main

/*
  Onix Config Manager - MongoDb event receiver for Pilot Control
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"github.com/gatblau/onix/artisan/server"
	"github.com/gorilla/mux"
)

func main() {
	// creates a generic http server
	s := server.New("onix/pilotctl/receivers/mongo")

	s.Http = func(router *mux.Router) {
		// enable encoded path  vars
		router.UseEncodedPath()
		// add http handlers
		router.HandleFunc("/events", eventReceiverHandler).Methods("POST")
		router.HandleFunc("/events", eventQueryHandler).Methods("GET")
	}

	s.Serve()
}
