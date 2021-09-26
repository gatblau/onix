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
	_ "github.com/gatblau/onix/pilotctl/docs"
	. "github.com/gatblau/onix/pilotctl/types"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// pingHandler excluded from swagger as it is accessed by pilot with a special time-bound access token
func pingHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("cannot read ping request body: %s\n", err)
		http.Error(w, "cannot read ping request body, check the server logs\n", http.StatusBadRequest)
		return
	}
	if len(body) > 0 {
		pingRequest := &PingRequest{}
		err = json.Unmarshal(body, pingRequest)
		if err != nil {
			log.Printf("cannot unmarshal ping request body: %s\n", err)
			http.Error(w, "cannot unmarshal ping request body, check the server logs\n", http.StatusBadRequest)
			return
		}
		// if the ping request contains a job result
		if pingRequest.Result != nil {
			// persist the result of the job
			err = api.CompleteJob(pingRequest.Result)
			if err != nil {
				log.Printf("cannot set job status: %s\n", err)
				http.Error(w, "set job status, check the server logs\n", http.StatusBadRequest)
				return
			}
		}
		// if the ping request contains syslog events
		if pingRequest.Events != nil && len(pingRequest.Events.Events) > 0 {
			// publish those events to registered sources
			api.PublishEvents(pingRequest.Events)
		}
	}
	// todo: add support for fx version
	jobId, fxKey, _, err := api.Ping()
	if err != nil {
		log.Printf("can't record ping time: %v\n", err)
		http.Error(w, "can't record ping time, check the server logs\n", http.StatusInternalServerError)
		return
	}
	// create a command with no job
	var cmdValue = &CmdInfo{
		JobId: jobId,
	}
	// if we have a job to execute
	if jobId > 0 {
		// fetches the definition for the job function to run from Onix
		cmdValue, err = api.GetCommandValue(fxKey)
		if err != nil {
			log.Printf("can't retrieve Artisan function definition from Onix: %v\n", err)
			http.Error(w, "can't retrieve Artisan function definition from Onix, check server logs\n", http.StatusInternalServerError)
			return
		}
		// set the job reference
		cmdValue.JobId = jobId
	}
	cr, err := NewPingResponse(*cmdValue, api.PingInterval())
	if err != nil {
		log.Printf("can't sign ping response: %v\n", err)
		http.Error(w, "can't sign ping response, check the server logs\n", http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(cr)
	if err != nil {
		log.Printf("can't marshal ping response: %s\n", err)
		http.Error(w, "can't marshal ping response, check the server logs\n", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(bytes)
}

// @Summary Get All Hosts
// @Description Returns a list of remote hosts
// @Tags Host
// @Router /host [get]
// @Param og query string false "the organisation group key to filter the query"
// @Param or query string false "the organisation key to filter the query"
// @Param ar query string false "the area key to filter the query"
// @Param lo query string false "the location key to filter the query"
// @Param label query string false "a pipe | separated list of labels associated to the host(s) to retrieve"
// @Produce json
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func hostQueryHandler(w http.ResponseWriter, r *http.Request) {
	orgGroup := r.FormValue("og")
	org := r.FormValue("or")
	area := r.FormValue("ar")
	location := r.FormValue("lo")
	labels := r.FormValue("label")
	var label []string
	if len(labels) > 0 {
		label = strings.Split(labels, "|")
	}
	hosts, err := api.GetHosts(orgGroup, org, area, location, label)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	server.Write(w, r, hosts)
}

// registerHandler excluded from swagger as it is accessed by pilot with a special time-bound access token
func registerHandler(w http.ResponseWriter, r *http.Request) {
	// get http body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body, check the server logs for more details", http.StatusBadRequest)
		return
	}
	// unmarshal body
	reg := &RegistrationRequest{}
	err = json.Unmarshal(body, reg)
	if err != nil {
		log.Printf("Error unmarshalling body: %v", err)
		http.Error(w, "can't unmarshal body, check the server logs for more details", http.StatusBadRequest)
		return
	}
	regInfo, err := api.Register(reg)
	if err != nil {
		log.Printf("Failed to register host, Onix responded with an error: %v", err)
		http.Error(w, "Failed to register host, Onix responded with an error, check the server logs for more details", http.StatusInternalServerError)
		return
	}
	log.Printf("host %s - %s registered", reg.Hostname, reg.MachineId)
	bytes, err := json.Marshal(regInfo)
	if err != nil {
		log.Printf("Failed to marshal registration configuration: %v", err)
		http.Error(w, "Failed to marshal registration configuration, check the server logs for more details", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(bytes)
}

// @Summary Create or Update a Command
// @Description creates a new or updates an existing command definition
// @Tags Command
// @Router /cmd [put]
// @Param command body types.Cmd true "the command definition"
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
	cmd := new(Cmd)
	err = json.Unmarshal(bytes, cmd)
	if err != nil {
		log.Printf("failed to unmarshal request: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = api.PutCommand(cmd)
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
	cmd, err := api.GetCommand(name)
	if err != nil {
		log.Printf("can't query command with name '%s': %v\n", name, err)
		http.Error(w, fmt.Sprintf("can't query command with name '%s': %v\n", name, err), http.StatusInternalServerError)
		return
	}
	server.Write(w, r, cmd)
}

// @Summary Delete a Command definition
// @Description deletes a specific a command definition using the command name
// @Tags Command
// @Router /cmd/{name} [delete]
// @Param name path string true "the unique name for the command to delete"
// @Accepts json
// @Produce plain
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func deleteCmdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	resultingOperation, err := api.DeleteCommand(name)
	if err != nil {
		log.Printf("can't delete command with name '%s': %v\n", name, err)
		http.Error(w, fmt.Sprintf("can't delete command with name '%s', check server logs for more details\n", name), http.StatusInternalServerError)
		return
	}
	server.Write(w, r, resultingOperation)
}

// @Summary Get all Command definitions
// @Description gets a list of all command definitions
// @Tags Command
// @Router /cmd [get]
// @Accepts json
// @Produce plain
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func getAllCmdHandler(w http.ResponseWriter, r *http.Request) {
	cmds, err := api.GetAllCommands()
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
// @Param command body types.JobBatchInfo true "the information required to create a new job"
// @Accepts json
// @Produce plain
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func newJobHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("can't read http body: %v\n", err)
		http.Error(w, fmt.Sprintf("can't read http body: check server logs\n"), http.StatusInternalServerError)
		return
	}
	var batch = new(JobBatchInfo)
	err = json.Unmarshal(bytes, batch)
	if err != nil {
		log.Printf("can't unmarshal http body: %v\n", err)
		http.Error(w, fmt.Sprintf("can't unmarshal http body, check the server logs\n"), http.StatusInternalServerError)
		return
	}
	jobBatchId, err := api.CreateJobBatch(*batch)
	if err != nil {
		log.Printf("can't create job batch: %v\n", err)
		http.Error(w, fmt.Sprintf("can't create job batch, check the server logs\n"), http.StatusInternalServerError)
		return
	}
	// return the batch ID
	w.Write([]byte(strconv.FormatInt(jobBatchId, 10)))
}

// @Summary Get Jobs
// @Description Returns a list of jobs filtered by the specified logistics tags
// @Tags Job
// @Router /job [get]
// @Param bid query int64 false "the unique identifier (number) of the job batch to retrieve"
// @Param og query string false "the organisation group key to filter the query"
// @Param or query string false "the organisation key to filter the query"
// @Param ar query string false "the area key to filter the query"
// @Param lo query string false "the location key to filter the query"
// @Produce json
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func getJobsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		bid *int64
		id  int64
		err error
	)
	batchId := r.FormValue("bid")
	// if a batch id was provided
	if len(batchId) > 0 {
		// try and parse to int64
		id, err = strconv.ParseInt(batchId, 10, 64)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// update pointer variable
		bid = &id
	}
	orgGroup := r.FormValue("og")
	org := r.FormValue("or")
	area := r.FormValue("ar")
	location := r.FormValue("lo")

	jobs, err := api.GetJobs(orgGroup, org, area, location, bid)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	server.Write(w, r, jobs)
}

// @Summary Get Job Batches
// @Description Returns a list of jobs batches with various filters
// @Tags Job
// @Router /job/batch [get]
// @Param name query string false "the name of the batch as in name% format"
// @Param owner query string false "the creator of the batch"
// @Param label query string false "a pipe | separated list of labels associated to the batch"
// @Param from query string false "the time from which to get batches (format should be dd-MM-yyyy)"
// @Param to query string false "the time to which to get batches (format should be dd-MM-yyyy)"
// @Produce json
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func getJobBatchHandler(w http.ResponseWriter, r *http.Request) {
	nameParam := r.FormValue("name")
	ownerParam := r.FormValue("owner")
	labelParam := r.FormValue("label")
	fromParam := r.FormValue("from")
	toParam := r.FormValue("to")

	var fromTime *time.Time
	if len(fromParam) > 0 {
		from, err := time.Parse("02-01-2006", fromParam)
		if err != nil {
			log.Printf("failed to parse FROM date '%s': %s\n", fromParam, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fromTime = &from
	}

	var toTime *time.Time
	if len(toParam) > 0 {
		to, err := time.Parse("02-01-2006", toParam)
		if err != nil {
			log.Printf("failed to parse TO date '%s': %s\n", toParam, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		toTime = &to
	}

	var label []string
	if len(labelParam) > 0 {
		label = strings.Split(labelParam, "|")
	}

	var name, owner *string
	if len(nameParam) > 0 {
		name = &nameParam
	}
	if len(ownerParam) > 0 {
		owner = &ownerParam
	}
	batches, err := api.GetJobBatches(name, owner, fromTime, toTime, &label)
	if err != nil {
		log.Printf("failed to retrieve job batches: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	server.Write(w, r, batches)
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
	areas, err := api.GetAreas(orgGroup)
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
	areas, err := api.GetOrgGroups()
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
	areas, err := api.GetOrgs(orgGroup)
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
	areas, err := api.GetLocations(area)
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
// @Param command body []types.Admission true "the required admission information"
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
	var admissions []Admission
	err = json.Unmarshal(bytes, &admissions)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, admission := range admissions {
		err = api.SetAdmission(admission)
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
	packages, err := api.GetPackages()
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
	api, err := api.GetPackageAPI(n)
	if err != nil {
		log.Printf("failed to get package API from registry: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	server.Write(w, r, api)
}
