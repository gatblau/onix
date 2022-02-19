/*
  Onix Config Manager - Artisan's Doorman
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
	s := httpserver.New("doorman")
	// add handlers
	s.Http = func(router *mux.Router) {
		router.HandleFunc("/key", upsertKeyHandler).Methods("PUT")
		router.HandleFunc("/command", upsertCommandHandler).Methods("PUT")
		router.HandleFunc("/route/in", upsertInboundRouteHandler).Methods("PUT")
		router.HandleFunc("/route/out", upsertOutboundRouteHandler).Methods("PUT")
		router.HandleFunc("/notification", upsertNotificationHandler).Methods("PUT")
		router.HandleFunc("/notification", getAllNotificationsHandler).Methods("GET")
		router.HandleFunc("/notification-template", upsertNotificationTemplateHandler).Methods("PUT")
		router.HandleFunc("/notification-template", getAllNotificationTemplatesHandler).Methods("GET")
		router.HandleFunc("/pipe", upsertPipelineHandler).Methods("PUT")
		router.HandleFunc("/pipe/{name}", getPipelineHandler).Methods("GET")
		router.HandleFunc("/pipe", getAllPipelinesHandler).Methods("GET")
		router.HandleFunc("/event/{uri}", eventHandler).Methods("POST")
	}
	s.Serve()
}
