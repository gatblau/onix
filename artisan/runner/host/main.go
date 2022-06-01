/*
  Onix Config Manager - Artisan Host Runner
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package main

import (
	"fmt"
	"os"

	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/oxlib/httpserver"
	"github.com/gorilla/mux"
)

func main() {

	fmt.Println(" setting port")
	app_port := os.Getenv("OX_HTTP_PORT")
	app_port1 := os.Getenv("HOST_RUNNER_PORT")
	if len(app_port) > 0 {
		os.Setenv("OX_HTTP_PORT", app_port)
	}
	core.Debug("host-runner listening at port ", app_port)
	fmt.Println("host-runner listening at port ", app_port)
	fmt.Println("host-runner listening at port1 ", app_port1)
	// creates a generic http server
	s := httpserver.New("art-host-runner")
	// add handlers
	s.Http = func(router *mux.Router) {
		router.HandleFunc("/host/{cmd-key}", executeCommandHandler).Methods("POST")
		router.HandleFunc("/flow", executeFlowFromPayloadHandler).Methods("POST")
		router.HandleFunc("/webhook/{flow-key}/push", executeWebhookFlowHandler).Methods("POST")
		fmt.Printf("new handler is registered...")
	}
	s.Serve()
}
