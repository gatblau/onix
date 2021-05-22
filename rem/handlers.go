package main

/*
  Onix Config Manager - REMote Host Service
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

// @title Onix Remote Host
// @version 0.0.4
// @description Remote Ctrl Service for Onix Pilot
// @contact.name gatblau
// @contact.url http://onix.gatblau.org/
// @contact.email onix@gatblau.org
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

import (
	"encoding/json"
	"fmt"
	"github.com/gatblau/onix/artisan/server"
	"github.com/gatblau/onix/rem/core"
	_ "github.com/gatblau/onix/rem/docs"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	rem *core.ReMan
	err error
)

func init() {
	rem, err = core.NewReMan()
	if err != nil {
		fmt.Printf("ERROR: fail to create remote manager: %s", err)
		os.Exit(1)
	}
}

// @Summary Ping
// @Description submits a ping from a host to the control plane
// @Tags Host
// @Router /ping/{host-key} [post]
// @Produce json
// @Param host-key path string true "the unique key for the host"
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func pingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	host := vars["host-key"]
	if len(host) == 0 {
		log.Printf("invalid host value: '%v'", host)
		http.Error(w, fmt.Sprintf("invalid host value: '%s'", host), http.StatusBadRequest)
		return
	}
	err = rem.Beat(host)
	if err != nil {
		log.Printf("Error recording ping: %v", err)
		http.Error(w, fmt.Sprintf("can't record ping: %s", err), http.StatusInternalServerError)
		return
	}
	log.Printf("host '%s' ping\n", host)
	// return an empty command list for now
	cmds := make([]core.CmdRequest, 0)
	bytes, err := json.Marshal(cmds)
	if err != nil {
		fmt.Printf("error: cant marshal commands: %s", err)
	}
	w.Write(bytes)
}

// @Summary Get All Hosts
// @Description Returns a list of remote hosts
// @Tags Host
// @Router /host [get]
// @Produce json
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func hostQueryHandler(w http.ResponseWriter, r *http.Request) {
	// lastSeen, _ := time.Parse("2006-Jan-02 Monday 03:04:05", "2020-Jan-29 Wednesday 12:19:25")
	// hosts := []core.Host{
	// 	{
	// 		Name:      "HOST-001",
	// 		Customer:  "CUST-01",
	// 		Region:    "UK-North-West",
	// 		Location:  "Manchester",
	// 		Connected: true,
	// 		LastSeen:  lastSeen,
	// 	},
	// 	{
	// 		Name:      "HOST-002",
	// 		Customer:  "CUST-01",
	// 		Region:    "UK-North-West",
	// 		Location:  "Manchester",
	// 		Connected: false,
	// 		LastSeen:  lastSeen,
	// 	},
	// 	{
	// 		Name:      "HOST-003",
	// 		Customer:  "CUST-01",
	// 		Region:    "UK-South-West",
	// 		Location:  "Devon",
	// 		Connected: false,
	// 		LastSeen:  lastSeen,
	// 	},
	// }
	hosts, err := rem.GetHostStatus()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	server.Write(w, r, hosts)
}

// @Summary Register a Host
// @Description registers a new host and its technical details with the service
// @Tags Host
// @Router /register [post]
// @Param registration-info body core.Registration true "the host registration configuration"
// @Accepts json
// @Produce plain
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func registerHandler(w http.ResponseWriter, r *http.Request) {
	// get http body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	// unmarshal body
	reg := &core.Registration{}
	err = json.Unmarshal(body, reg)
	if err != nil {
		log.Printf("Error unmarshalling body: %v", err)
		http.Error(w, "can't unmarshal body", http.StatusBadRequest)
		return
	}
	err = rem.Register(reg)
	if err != nil {
		log.Printf("Error recording ping: %v", err)
		http.Error(w, fmt.Sprintf("can't record ping: %s", err), http.StatusInternalServerError)
		return
	}
	log.Printf("host %s - %s registered", reg.Hostname, reg.MachineId)
}

// @Summary Log Events
// @Description log host events (e.g. up, down, connected, disconnected)
// @Tags Host
// @Router /log [post]
// @Param logs body core.Events true "the host logs to post"
// @Accepts json
// @Produce plain
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func newLogHandler(w http.ResponseWriter, r *http.Request) {
}

// @Summary Get Events by Host
// @Description get log host events (e.g. up, down, connected, disconnected) by specific host
// @Tags Host
// @Router /log/{host-id} [get]
// @Param host-key path string true "the unique key for the host"
// @Accepts json
// @Produce plain
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func geLogHandler(w http.ResponseWriter, r *http.Request) {
}

// @Summary Create or Update a Command
// @Description creates a new or updates an existing command definition
// @Tags Command
// @Router /cmd [put]
// @Param command body core.Cmd true "the command definition"
// @Accepts json
// @Produce plain
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func updateCmdHandler(w http.ResponseWriter, r *http.Request) {
}

// @Summary Get a Command definition
// @Description get a specific a command definition
// @Tags Command
// @Router /cmd/{id} [get]
// @Param id path string true "the unique id for the command to retrieve"
// @Accepts json
// @Produce plain
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func getCmdHandler(w http.ResponseWriter, r *http.Request) {
}

// @Summary Get All Command definitions
// @Description get all command definitions
// @Tags Command
// @Router /cmd [get]
// @Accepts json
// @Produce plain
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func getAllCmdHandler(w http.ResponseWriter, r *http.Request) {
}

// @Summary Create a Job
// @Description create a new job for execution on one or more remote hosts
// @Tags Job
// @Router /job [post]
// @Param command body core.Cmd true "the job definition"
// @Accepts json
// @Produce plain
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func newJobHandler(w http.ResponseWriter, r *http.Request) {
}

// @Summary Get Job Information
// @Description get a specific a job information
// @Tags Job
// @Router /job/{id} [get]
// @Param id path string true "the unique id for the job to retrieve"
// @Accepts json
// @Produce plain
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func getJobHandler(w http.ResponseWriter, r *http.Request) {
}

// @Summary Get All Jobs Information
// @Description get all jobs
// @Tags Job
// @Router /job [get]
// @Accepts json
// @Produce plain
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func getJobsHandler(w http.ResponseWriter, r *http.Request) {
}

// @Summary Submit a Vulnerability Scan Report
// @Description submits a vulnerability report for a specific host
func uploadVulnerabilityReportHandler(w http.ResponseWriter, r *http.Request) {
}

// @Summary Get Regions
// @Description get a list of regions where hosts are deployed
// @Tags Region
// @Router /region [get]
// @Produce json
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func getRegionsHandler(w http.ResponseWriter, r *http.Request) {
	regions := []core.Region{
		{
			Key:  "NE",
			Name: "North East",
		},
		{
			Key:  "NW",
			Name: "North West",
		},
		{
			Key:  "NE",
			Name: "North East",
		},
		{
			Key:  "YH",
			Name: "Yorkshire & The Humber",
		},
		{
			Key:  "WM",
			Name: "West Midlands",
		},
		{
			Key:  "EM",
			Name: "East Midlands",
		},
		{
			Key:  "EE",
			Name: "East of England",
		},
		{
			Key:  "LO",
			Name: "London",
		},
		{
			Key:  "SE",
			Name: "South East",
		},
		{
			Key:  "SW",
			Name: "South West",
		},
	}
	server.Write(w, r, regions)
}

// @Summary Get Locations by Region
// @Description get a list of locations within a particular region
// @Tags Region
// @Router /region/{region-key}/location [get]
// @Produce json
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func geLocationsByRegionHandler(w http.ResponseWriter, r *http.Request) {
	regions := []core.Location{
		{
			Key:       "CHE",
			Name:      "Cheshire",
			RegionKey: "NW",
		},
		{
			Key:       "GM",
			Name:      "Greater Manchester",
			RegionKey: "NW",
		},
		{
			Key:       "CU",
			Name:      "Cumbria",
			RegionKey: "NW",
		},
		{
			Key:       "LANC",
			Name:      "Lancashire",
			RegionKey: "NW",
		},
		{
			Key:       "MER",
			Name:      "Merseyside",
			RegionKey: "NW",
		},
		// london
		{
			Key:       "CITY",
			Name:      "London City",
			RegionKey: "LO",
		},
		{
			Key:       "BX",
			Name:      "Brixton",
			RegionKey: "LO",
		},
		{
			Key:       "CR",
			Name:      "Croydon",
			RegionKey: "LO",
		},
		{
			Key:       "CA",
			Name:      "Camden",
			RegionKey: "LO",
		},
		{
			Key:       "GRE",
			Name:      "Greenwich",
			RegionKey: "LO",
		},
	}
	server.Write(w, r, regions)
}
