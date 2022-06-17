/*
  Onix Config Manager - Warden
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package mode

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"github.com/gatblau/onix/warden/lib"
	"log"
	"net/http"
)

func Basic(address string, verbose bool) {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = verbose
	fmt.Println(lib.Banner())
	// lib.InfoLogger.Printf("warden starting @ %s\n", address)
	// if verbose {
	//     lib.InfoLogger.Printf("verbose output enabled\n")
	// }
	log.Fatal(http.ListenAndServe(address, proxy))
}
