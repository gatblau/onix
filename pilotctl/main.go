package main

/*
  Onix Config Manager - Pilot Control
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"fmt"
	"github.com/gatblau/onix/artisan/server"
	"github.com/gatblau/onix/pilotctl/core"
	"github.com/gorilla/mux"
	"os"
)

var (
	api *core.API
)

func init() {
	var err error
	api, err = core.NewAPI(new(core.Conf))
	if err != nil {
		fmt.Printf("ERROR: fail to create backedn services API: %s", err)
		os.Exit(1)
	}
}

func main() {
	// creates a generic http server
	s := server.New("onix/pilotctl")
	// add handlers
	s.Http = func(router *mux.Router) {
		// enable encoded path  vars
		router.UseEncodedPath()
		// add http handlers
		router.HandleFunc("/ping/{machine-id}", pingHandler).Methods("POST")
		router.HandleFunc("/host", hostQueryHandler).Methods("GET")
		router.HandleFunc("/register", registerHandler).Methods("POST")
		router.HandleFunc("/cmd", updateCmdHandler).Methods("PUT")
		router.HandleFunc("/cmd", getAllCmdHandler).Methods("GET")
		router.HandleFunc("/cmd/{name}", getCmdHandler).Methods("GET")
		router.HandleFunc("/org-group", getOrgGroupsHandler).Methods("GET")
		router.HandleFunc("/org-group/{org-group}/area", getAreasHandler).Methods("GET")
		router.HandleFunc("/org-group/{org-group}/org", getOrgHandler).Methods("GET")
		router.HandleFunc("/area/{area}/location", getLocationsHandler).Methods("GET")
		router.HandleFunc("/admission", setAdmissionHandler).Methods("PUT")
		router.HandleFunc("/package", getPackagesHandler).Methods("GET")
		router.HandleFunc("/package/{name}/api", getPackagesApiHandler).Methods("GET")
		router.HandleFunc("/job", newJobHandler).Methods("POST")
		router.HandleFunc("/job", getJobsHandler).Methods("GET")
	}
	// set up specific authentication for host pilot agents
	s.Auth = map[string]func(string) bool{
		"/register": pilotAuth,
		"/ping/.*":  pilotAuth,
	}
	s.Serve()
}

var pilotAuth = func(token string) bool {
	return api.Authenticate(token)
}
