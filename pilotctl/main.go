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
	"github.com/reugn/go-quartz/quartz"
	"os"
	"time"
)

var (
	rem *core.ReMan
)

func init() {
	var err error
	rem, err = core.NewReMan()
	if err != nil {
		fmt.Printf("ERROR: fail to create remote manager: %s", err)
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
		router.HandleFunc("/admission", getAdmissionsHandler).Methods("GET")
		router.HandleFunc("/admission", setAdmissionHandler).Methods("PUT")
		router.HandleFunc("/package", getPackagesHandler).Methods("GET")
		router.HandleFunc("/package/{name}/api", getPackagesApiHandler).Methods("GET")
	}
	// add asynchronous jobs
	// starts a job to record events if host connection status changes
	s.Jobs = func() error {
		conf := core.NewConf()
		interval := time.Duration(conf.GetPingInterval())
		// creates a job to check for changes in the base image
		updateConnStatusJob, err := core.NewUpdateConnStatusJob(rem)
		if err != nil {
			return fmt.Errorf("cannot create connection status update job: %s", err)
		}
		// create a new scheduler
		sched := quartz.NewStdScheduler()
		// start the scheduler
		sched.Start()
		// schedule the job
		err = sched.ScheduleJob(updateConnStatusJob, quartz.NewSimpleTrigger(time.Duration(interval*time.Second)))
		if err != nil {
			return fmt.Errorf("cannot schedule connection status update job: %s", err)
		}
		return nil
	}
	// set up specific authentication for host pilot agents
	s.Auth = map[string]func(string) bool{
		"/register": pilotAuth,
		"/ping/.*":  pilotAuth,
	}
	s.Serve()
}

var pilotAuth = func(token string) bool {
	return rem.Authenticate(token)
}
