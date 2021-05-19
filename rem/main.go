package main

/*
  Onix Config Manager - REMote Host Service
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"fmt"
	"github.com/gatblau/onix/artisan/server"
	"github.com/gatblau/onix/rem/core"
	"github.com/gorilla/mux"
	"github.com/reugn/go-quartz/quartz"
	"time"
)

func main() {
	// creates a generic http server
	s := server.New("onix/rem")
	// add handlers
	s.ServeWithJobs(func(router *mux.Router) {
		router.HandleFunc("/ping/{host-key}", pingHandler).Methods("POST")
		router.HandleFunc("/host", hostQueryHandler).Methods("GET")
		router.HandleFunc("/register", registerHandler).Methods("POST")
		router.HandleFunc("/cmd", updateCmdHandler).Methods("POST")
		router.HandleFunc("/cmd/{id}", getCmdHandler).Methods("GET")
		router.HandleFunc("/region", getRegionsHandler).Methods("GET")
		router.HandleFunc("/region/{region-key}/location", geLocationsByRegionHandler).Methods("GET")
	}, addJobs)
}

// starts a job to record events if host connection status changes
func addJobs() error {
	conf := core.NewConf()
	interval := time.Duration(conf.GetPingInterval())
	// creates a job to check for changes in the base image
	updateConnStatusJob := core.NewUpdateConnStatusJob()
	// create a new scheduler
	sched := quartz.NewStdScheduler()
	// start the scheduler
	sched.Start()
	// schedule the job
	err := sched.ScheduleJob(updateConnStatusJob, quartz.NewSimpleTrigger(time.Duration(interval*time.Second)))
	if err != nil {
		return fmt.Errorf("cannot schedule connection status update job: %s", err)
	}
	return nil
}