/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package main

import (
	"github.com/gatblau/onix/warden/server"
	"github.com/gorilla/mux"
)

func main() {
	// creates a generic http server
	s := server.New("onix/warden")
	// add handlers
	s.Serve(func(router *mux.Router) {
		router.HandleFunc("/list", listHandler).Methods("GET")
	})
}
