package main

import (
	"github.com/gatblau/onix/artisan/server"
	"github.com/gorilla/mux"
)

func main() {
	// creates a generic http server
	s := server.New("onix/rem")
	// add handlers
	s.Serve(func(router *mux.Router) {
		router.HandleFunc("/bear", beatHandler).Methods("POST")
	})
}
