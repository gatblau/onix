package main

import (
	"github.com/gatblau/onix/artisan/flow"
	"github.com/gatblau/onix/artisan/server"
	"github.com/gorilla/mux"
	"testing"
)

func Test(t *testing.T) {
	f := new(flow.Flow)
	z := f.Labels["aaa"]
	print(z)
	// creates a generic http server
	s := server.New("onix/artisan-runner")
	// add handlers
	s.Serve(func(router *mux.Router) {
		router.HandleFunc("/flow", createFlowFromPayloadHandler).Methods("POST")
		router.HandleFunc("/list", listHandler).Methods("GET")
	})
}
