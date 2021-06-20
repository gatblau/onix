package main

/*
  Onix Config Manager - Pilot Control
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

// @title Pilot Control
// @version 0.0.4
// @description Onix Config Manager Control Service for Pilot Host agent
// @contact.name gatblau
// @contact.url http://onix.gatblau.org/
// @contact.email onix@gatblau.org
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

import (
	"encoding/json"
	"fmt"
	"github.com/gatblau/onix/artisan/server"
	"github.com/gatblau/onix/pilotctl/core"
	_ "github.com/gatblau/onix/pilotctl/docs"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var (
	err error
)

// @Summary Ping
// @Description submits a ping from a host to the control plane
// @Tags Host
// @Router /ping/{machine_id} [post]
// @Produce json
// @Param host-key path string true "the unique key for the host"
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func pingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	machineId := vars["machine_id"]
	if len(machineId) == 0 {
		log.Printf("missing machine Id")
		http.Error(w, "missing machine Id", http.StatusBadRequest)
		return
	}
	err = rem.Beat(machineId)
	if err != nil {
		log.Printf("can't record ping: %v\n", err)
		http.Error(w, fmt.Sprintf("can't record ping: %s\n", err), http.StatusInternalServerError)
		return
	}
	// return an empty command list for now
	var cmds []core.CmdRequest
	bytes, err := json.Marshal(cmds)
	if err != nil {
		log.Printf("can't marshal commands: %s\n", err)
	}
	w.WriteHeader(http.StatusCreated)
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
		log.Printf("Failed to register host, Onix responded with an error: %v", err)
		http.Error(w, "Failed to register host, Onix responded with an error", http.StatusInternalServerError)
		return
	}
	log.Printf("host %s - %s registered", reg.Hostname, reg.MachineId)
	w.WriteHeader(http.StatusCreated)
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
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read request body: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cmd := new(core.Cmd)
	err = json.Unmarshal(bytes, cmd)
	if err != nil {
		log.Printf("failed to unmarshal request: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = rem.SetCommand(cmd)
	if err != nil {
		log.Printf("failed to set command: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// @Summary Get a Command definition
// @Description get a specific a command definition
// @Tags Command
// @Router /cmd/{name} [get]
// @Param name path string true "the unique name for the command to retrieve"
// @Accepts json
// @Produce plain
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func getCmdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	cmd, err := rem.GetCommand(name)
	if err != nil {
		log.Printf("can't query command with name '%s': %v\n", name, err)
		http.Error(w, fmt.Sprintf("can't query command with name '%s': %v\n", name, err), http.StatusInternalServerError)
		return
	}
	server.Write(w, r, cmd)
}

// @Summary Get all Command definitions
// @Description get a list of all command definitions
// @Tags Command
// @Router /cmd [get]
// @Accepts json
// @Produce plain
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func getAllCmdHandler(w http.ResponseWriter, r *http.Request) {
	cmds, err := rem.GetAllCommands()
	if err != nil {
		log.Printf("can't query list of commands: %v\n", err)
		http.Error(w, fmt.Sprintf("can't query list of commands: %s\n", err), http.StatusInternalServerError)
		return
	}
	server.Write(w, r, cmds)
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

// @Summary Get Host Admissions
// @Description get a list of keys of the hosts admitted into service
// @Tags Admission
// @Router /admission [get]
// @Produce json
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func getAdmissionsHandler(w http.ResponseWriter, r *http.Request) {
	admissions, err := rem.GetAdmissions()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	server.Write(w, r, admissions)
}

// @Summary Create or Update a Host Admission
// @Description creates a new or updates an existing host admission by allowing to specify active status and search tags
// @Tags Admission
// @Router /admission [put]
// @Param command body core.Admission true "the admission to be set"
// @Accepts json
// @Produce plain
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func setAdmissionHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	admission := new(core.Admission)
	err = json.Unmarshal(bytes, admission)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = rem.SetAdmission(admission)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Get Artisan Packages
// @Description get a list of packages in the backing Artisan registry
// @Tags Registry
// @Router /package [get]
// @Produce json
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func getPackagesHandler(w http.ResponseWriter, r *http.Request) {
	packages, err := rem.GetPackages()
	if err != nil {
		log.Printf("failed to retrieve package list from Artisan Registry: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	server.Write(w, r, packages)
}

// @Summary Get the API of an Artisan Package
// @Description get a list of exported functions and inputs for the specified package
// @Tags Registry
// @Router /package/{name}/api [get]
// @Produce json
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
// @Param name path string true "the fully qualified name of the artisan package having the required API"
func getPackagesApiHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	n, err := url.PathUnescape(name)
	if err != nil {
		log.Printf("failed to unescape package name '%s': %s\n", name, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	api, err := rem.GetPackageAPI(n)
	if err != nil {
		log.Printf("failed to get package API from registry: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	server.Write(w, r, api)
}
