/*
  Onix Config Manager - K(ubernetes Cluster Ops Manager)
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package main

// @title K - Kubernetes Cluster Ops Manager
// @version 0.0.4
// @description Onix Config Manager Service to Deploy and Continuously Configure Kubernetes
// @contact.name gatblau
// @contact.url http://onix.gatblau.org/
// @contact.email onix@gatblau.org
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

import (
	_ "github.com/gatblau/onix/k/docs"
	"net/http"
)

func testHandler(w http.ResponseWriter, r *http.Request) {

}
