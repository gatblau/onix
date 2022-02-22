package main

import (
	"net/http"

	"github.com/gatblau/onix/artisan/runner/host/handlers"
)

func createOSPatchingHandler(w http.ResponseWriter, r *http.Request) {

	osph := handlers.OSpatchingHandler{}
	osph.HandleEvent(w, r)
}
