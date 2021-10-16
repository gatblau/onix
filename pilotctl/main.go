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
	"github.com/gatblau/onix/oxlib/httpserver"
	"github.com/gatblau/onix/oxlib/oxc"
	"github.com/gatblau/onix/pilotctl/core"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

var (
	api *core.API
)

func init() {
	var err error
	api, err = core.NewAPI(new(core.Conf))
	if err != nil {
		fmt.Printf("ERROR: fail to create backend services API: %s", err)
		os.Exit(1)
	}
}

func main() {
	// creates a generic http server
	s := httpserver.New("pilotctl")
	// add handlers
	s.Http = func(router *mux.Router) {
		// enable encoded path  vars
		router.UseEncodedPath()
		// middleware
		// router.Use(s.LoggingMiddleware)
		router.Use(s.AuthenticationMiddleware)

		// pilot http handlers
		router.HandleFunc("/ping", pingHandler).Methods("POST")
		router.HandleFunc("/register", registerHandler).Methods("POST")

		// apply authorisation to admin user http handlers
		router.Handle("/host", s.Authorise(hostQueryHandler)).Methods("GET")
		router.Handle("/cmd", s.Authorise(updateCmdHandler)).Methods("PUT")
		router.Handle("/cmd", s.Authorise(getAllCmdHandler)).Methods("GET")
		router.Handle("/cmd/{name}", s.Authorise(getCmdHandler)).Methods("GET")
		router.Handle("/cmd/{name}", s.Authorise(deleteCmdHandler)).Methods("DELETE")
		router.Handle("/org-group", s.Authorise(getOrgGroupsHandler)).Methods("GET")
		router.Handle("/org-group/{org-group}/area", s.Authorise(getAreasHandler)).Methods("GET")
		router.Handle("/org-group/{org-group}/org", s.Authorise(getOrgHandler)).Methods("GET")
		router.Handle("/area/{area}/location", s.Authorise(getLocationsHandler)).Methods("GET")
		router.Handle("/admission", s.Authorise(setAdmissionHandler)).Methods("PUT")
		router.Handle("/package", s.Authorise(getPackagesHandler)).Methods("GET")
		router.Handle("/package/{name}/api", s.Authorise(getPackagesApiHandler)).Methods("GET")
		router.Handle("/job", s.Authorise(newJobHandler)).Methods("POST")
		router.Handle("/job", s.Authorise(getJobsHandler)).Methods("GET")
		router.Handle("/job/batch", s.Authorise(getJobBatchHandler)).Methods("GET")
		router.Handle("/user", s.Authorise(getUserHandler)).Methods("GET")

		router.HandleFunc("/pub", getKeyHandler).Methods("GET")
	}
	// set up specific authentication for host pilot agents
	s.Auth = map[string]func(http.Request) *oxc.UserPrincipal{
		"/register": pilotAuth,
		"/ping":     pilotAuth,
	}
	s.DefaultAuth = defaultAuth
	s.Serve()
}

// the overridden authentication mechanism used by the authentication middleware for specific routes
// specified in server.Auth map
var pilotAuth = func(r http.Request) *oxc.UserPrincipal {
	token := r.Header.Get("Authorization")
	return api.AuthenticatePilot(token)
}

// the default authentication mechanism user by the authentication middleware
var defaultAuth = func(r http.Request) *oxc.UserPrincipal {
	return api.AuthenticateUser(r)
}
