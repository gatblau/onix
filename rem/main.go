package main

/*
  Onix Config Manager - REMote Host Service
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
	s := server.New("onix/rem")
	// add handlers
	s.Serve(func(router *mux.Router) {
		router.HandleFunc("/ping/{host-key}", pingHandler).Methods("POST")
		router.HandleFunc("/host", hostQueryHandler).Methods("GET")
		router.HandleFunc("/register", registerHandler).Methods("POST")
		router.HandleFunc("/cmd", updateCmdHandler).Methods("POST")
		router.HandleFunc("/cmd/{id}", getCmdHandler).Methods("GET")
	})
}
