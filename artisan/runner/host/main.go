/*
  Onix Config Manager - Artisan Runner
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package main

import (
	"fmt"

	"github.com/gatblau/onix/artisan/runner/host/handlers"
	"github.com/gatblau/onix/oxlib/httpserver"
	"github.com/gorilla/mux"
)

func main() {
	/*
		cmd := exec.Command("art", "pull", "aps-edge-registry.amosonline.io/aps/keys/pk-registry:latest", "-u", "admin:n3xU5@APS")
		var outb, errb bytes.Buffer
		cmd.Stdout = &outb
		cmd.Stderr = &errb
		err := cmd.Run()
		if err != nil {
			fmt.Println("ERROR  occurred out -:", outb.String(), "err -:", errb.String())
			log.Fatal(err)
		}
		fmt.Println("out ===", outb.String(), "err ==", errb.String())
	*/

	// creates a generic http server
	handlerMgr := handlers.NewHandlerManager()
	s := httpserver.New("art-host-runner")
	// add handlers
	s.Http = func(router *mux.Router) {
		fmt.Printf("handler is registered...")
		router.Handle("/host/{package}/{function}", handlerMgr).Methods("POST")
	}
	s.Serve()

	/*
		// creates a generic http server
		s := httpserver.New("art-host-runner")
		// add handlers
		s.Http = func(router *mux.Router) {
			fmt.Printf("new handler is registered...\n")
			router.HandleFunc("/host/{package}/{function}", createOSPatchingHandler).Methods("POST")
		}
		s.Serve()

	*/

}
