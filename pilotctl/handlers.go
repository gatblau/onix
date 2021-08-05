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

// @Summary Ping
// @Description submits a ping from a host to the control plane
// @Tags Host
// @Router /ping/{machine-id} [post]
// @Produce json
// @Param machine-id path string true "the machine Id of the host"
// @Param cmd-result body string false "the result of the execution of the last command or nil if no result is available"
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func pingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	machineId := vars["machine-id"]
	if len(machineId) == 0 {
		log.Printf("missing machine Id")
		http.Error(w, "missing machine Id", http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("cannot read ping request body: %s\n", err)
		http.Error(w, "cannot read ping request body, check the server logs\n", http.StatusBadRequest)
		return
	}
	if len(body) > 0 {
		result := &core.Result{}
		err = json.Unmarshal(body, result)
		if err != nil {
			log.Printf("cannot unmarshal ping request body: %s\n", err)
			http.Error(w, "cannot unmarshal ping request body, check the server logs\n", http.StatusBadRequest)
			return
		}
		err = rem.CompleteJob(result)
		if err != nil {
			log.Printf("cannot set job status: %s\n", err)
			http.Error(w, "set job status, check the server logs\n", http.StatusBadRequest)
			return
		}
	}
	// todo: add support for fx version
	jobId, fxKey, _, err := rem.Beat(machineId)
	if err != nil {
		log.Printf("can't record ping: %v\n", err)
		http.Error(w, "can't record ping, check server logs\n", http.StatusInternalServerError)
		return
	}
	// create a command with no job
	var cmdValue = &core.CmdValue{
		JobId: jobId,
	}
	// if we have a job to execute
	if jobId > 0 {
		// fetches the definition for the job function to run from Onix
		cmdValue, err = rem.GetCommandValue(fxKey)
		if err != nil {
			log.Printf("can't retrieve Artisan function definition from Onix: %v\n", err)
			http.Error(w, "can't retrieve Artisan function definition from Onix, check server logs\n", http.StatusInternalServerError)
			return
		}
		// set the job reference
		cmdValue.JobId = jobId
	}
	cr, err := core.NewCmdRequest(*cmdValue)
	if err != nil {
		log.Printf("can't sign command request: %v\n", err)
		http.Error(w, "can't sign command request, check server logs\n", http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(cr)
	if err != nil {
		log.Printf("can't marshal command request: %s\n", err)
		http.Error(w, "can't marshal command request, check server logs\n", http.StatusInternalServerError)
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
	err = rem.PutCommand(cmd)
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

// @Summary Get Areas in Organisation Group
// @Description Get a list of areas setup in an organisation group
// @Tags Logistics
// @Router /org-group/{org-group}/area [get]
// @Param org-group path string true "the unique id for organisation group under which areas are defined"
// @Produce json
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func getAreasHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgGroup := vars["org-group"]
	areas, err := rem.GetAreas(orgGroup)
	if err != nil {
		log.Printf("failed to retrieve areas: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	server.Write(w, r, areas)
}

// @Summary Get Organisation Groups
// @Description Get a list of organisation groups
// @Tags Logistics
// @Router /org-group [get]
// @Produce json
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func getOrgGroupsHandler(w http.ResponseWriter, r *http.Request) {
	areas, err := rem.GetOrgGroups()
	if err != nil {
		log.Printf("failed to retrieve org groups: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	server.Write(w, r, areas)
}

// @Summary Get Organisations in Organisation Group
// @Description Get a list of organisations setup in an organisation group
// @Tags Logistics
// @Router /org-group/{org-group}/org [get]
// @Param org-group path string true "the unique id for organisation group under which organisations are defined"
// @Produce json
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func getOrgHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgGroup := vars["org-group"]
	areas, err := rem.GetOrgs(orgGroup)
	if err != nil {
		log.Printf("failed to retrieve organisations: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	server.Write(w, r, areas)
}

// @Summary Get Locations in an Area
// @Description Get a list of locations setup in an area
// @Tags Logistics
// @Router /area/{area}/location [get]
// @Param area path string true "the unique id for area under which locations are defined"
// @Produce json
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func getLocationsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	area := vars["area"]
	areas, err := rem.GetLocations(area)
	if err != nil {
		log.Printf("failed to retrieve organisations: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	server.Write(w, r, areas)
}

// @Summary Admits a host into service
// @Description inform pilotctl to accept management connections coming from a host pilot agent
// @Description admitting a host also requires associating the relevant logistic information such as org, area and location for the host
// @Tags Admission
// @Router /admission [put]
// @Param command body []core.Admission true "the required admission information"
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
	var admissions []core.Admission
	err = json.Unmarshal(bytes, &admissions)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, admission := range admissions {
		err = rem.SetAdmission(admission)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
