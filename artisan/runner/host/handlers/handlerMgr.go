/*
  Onix Config Manager - Artisan Runner
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package handlers

// @title Artisan Host Runner
// @version 0.0.4
// @description Run Artisan flows in host
// @contact.name gatblau
// @contact.url http://onix.gatblau.org/
// @contact.email onix@gatblau.org
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type HandlerManager struct {
	handlerMapping map[string]EventHandler
}

func NewHandlerManager() *HandlerManager {
	mgr := new(HandlerManager)
	mgr.handlerMapping = make(map[string]EventHandler)
	// adding S3EventHandler to map
	osph := OSpatchingHandler{}
	var h EventHandler
	h = osph
	//map's key must be the same as the key stored in the onix db table item for item-type ART_FX
	mgr.handlerMapping["build-patching-pkg"] = h

	return mgr
}

func (h HandlerManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	flowkey := vars["flow-key"]
	eh := h.handlerMapping[flowkey]
	if eh == nil {
		msg := fmt.Sprintf("No handler is registered for flow-key %s\n", flowkey)
		fmt.Printf(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	} else {
		eh.HandleEvent(w, r)
	}
}
