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
	"log"
	"os"

	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/oxlib/httpserver"
	m "github.com/gatblau/onix/oxlib/mqtt"
	"github.com/gorilla/mux"
)

func main() {

	app_port := os.Getenv("HOST_RUNNER_PORT")
	if len(app_port) > 0 {
		os.Setenv("OX_HTTP_PORT", app_port)
	}

	// creates a generic http server
	s := httpserver.New("art-host-runner")
	// add handlers
	s.Http = func(router *mux.Router) {
		router.HandleFunc("/host/{cmd-key}", executeCommandHandler).Methods("POST")
		router.HandleFunc("/flow", executeFlowFromPayloadHandler).Methods("POST")
		router.HandleFunc("/webhook/{flow-key}/push", executeWebhookFlowHandler).Methods("POST")
		core.Debug("new handler is registered...")
	}

	connstatus := make(chan error, 1)
	go func() {
		fmt.Println("launching broker")
		er := launchBroker()
		connstatus <- er
	}()
	// not using s.Jobs to executive above go routine, because it gets into deadlock condition
	// why because, below select statement wait for data from connstatus channel and data in
	// connstatus channel is set when above function is called from s.Jobs, and s.Jobs is called
	// from s.Serve(), where the code never reached

	select {
	case err := <-connstatus:
		{
			if err != nil {
				log.Fatalf("ERROR: mqtt client failed to connect broker : %s \n", err)
			}
		}
	}
	core.Debug("starting http server")
	s.Serve()
}

func launchBroker() error {
	mqc := m.BuildClient()
	err := mqc.Start(30)
	if err != nil {
		return err
	}
	mqc.Subscribe(eventMessageHandler)
	return err
}
