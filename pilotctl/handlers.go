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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gatblau/onix/oxlib/httpserver"
	"github.com/gatblau/onix/pilotctl/core"
	_ "github.com/gatblau/onix/pilotctl/docs"
	. "github.com/gatblau/onix/pilotctl/types"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
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
			err = core.Api().CompleteJob(pingRequest.Result)
			if err != nil {
				log.Printf("cannot set job status: %s\n", err)
				http.Error(w, "set job status, check the server logs\n", http.StatusBadRequest)
				return
			}
		}
		// if the ping request contains syslog events
		if pingRequest.Events != nil && len(pingRequest.Events.Events) > 0 {
			// adds extra information to the events
			events, err := core.Api().Augment(pingRequest.Events)
			// if the augmentation fails
			if err != nil && len(events.Events) > 0 {
				// log a warning but continue
				log.Printf("WARNING: event augmentation failed for host uuid '%s': %s", events.Events[0].HostUUID, err)
			}
			// publish those events to registered sources
			core.Api().PublishEvents(events)
		}
	}
	// todo: add support for fx version
	jobId, fxKey, _, err := core.Api().Ping()
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
		cmdValue, err = core.Api().GetCommandValue(fxKey)
		if err != nil {
			log.Printf("can't retrieve Artisan function definition from Onix: %v\n", err)
			http.Error(w, "can't retrieve Artisan function definition from Onix, check server logs\n", http.StatusInternalServerError)
			return
		}
		// set the job reference
		cmdValue.JobId = jobId
	}
	cr, err := NewPingResponse(*cmdValue, core.Api().PingInterval())
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
	hosts, err := core.Api().GetHosts(orgGroup, org, area, location, label)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	httpserver.Write(w, r, hosts)
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
	regInfo, err := core.Api().Register(reg)
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
	err = core.Api().PutCommand(cmd)
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
	cmd, err := core.Api().GetCommand(name)
	if err != nil {
		log.Printf("can't query command with name '%s': %v\n", name, err)
		http.Error(w, fmt.Sprintf("can't query command with name '%s': %v\n", name, err), http.StatusInternalServerError)
		return
	}
	httpserver.Write(w, r, cmd)
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
	resultingOperation, err := core.Api().DeleteCommand(name)
	if err != nil {
		log.Printf("can't delete command with name '%s': %v\n", name, err)
		http.Error(w, fmt.Sprintf("can't delete command with name '%s', check server logs for more details\n", name), http.StatusInternalServerError)
		return
	}
	httpserver.Write(w, r, resultingOperation)
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
	cmds, err := core.Api().GetAllCommands()
	if err != nil {
		log.Printf("can't query list of commands: %v\n", err)
		http.Error(w, fmt.Sprintf("can't query list of commands: %s\n", err), http.StatusInternalServerError)
		return
	}
	httpserver.Write(w, r, cmds)
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
	jobBatchId, err := core.Api().CreateJobBatch(*batch)
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
	var bid *int64
	batchId := r.FormValue("bid")
	// if a batch id was provided
	if len(batchId) > 0 {
		id, parseErr := strconv.ParseInt(batchId, 10, 64)
		if isErr(w, parseErr, http.StatusBadRequest, "cannot parse batch Id") {
			return
		}
		bid = &id
	}
	orgGroup := r.FormValue("og")
	org := r.FormValue("or")
	area := r.FormValue("ar")
	location := r.FormValue("lo")

	jobs, err := core.Api().GetJobs(orgGroup, org, area, location, bid)
	if isErr(w, err, http.StatusBadRequest, "cannot retrieve jobs from database") {
		return
	}
	httpserver.Write(w, r, jobs)
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
	batches, err := core.Api().GetJobBatches(name, owner, fromTime, toTime, &label)
	if err != nil {
		log.Printf("failed to retrieve job batches: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	httpserver.Write(w, r, batches)
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
	areas, err := core.Api().GetAreas(orgGroup)
	if err != nil {
		log.Printf("failed to retrieve areas: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	httpserver.Write(w, r, areas)
}

// @Summary Get Organisation Groups
// @Description Get a list of organisation groups
// @Tags Logistics
// @Router /org-group [get]
// @Produce json
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func getOrgGroupsHandler(w http.ResponseWriter, r *http.Request) {
	areas, err := core.Api().GetOrgGroups()
	if err != nil {
		log.Printf("failed to retrieve org groups: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	httpserver.Write(w, r, areas)
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
	areas, err := core.Api().GetOrgs(orgGroup)
	if err != nil {
		log.Printf("failed to retrieve organisations: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	httpserver.Write(w, r, areas)
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
	areas, err := core.Api().GetLocations(area)
	if err != nil {
		log.Printf("failed to retrieve organisations: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	httpserver.Write(w, r, areas)
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
		err = core.Api().SetAdmission(admission)
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
	packages, err := core.Api().GetPackages()
	if err != nil {
		log.Printf("failed to retrieve package list from Artisan Registry: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	httpserver.Write(w, r, packages)
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
	api, err := core.Api().GetPackageAPI(n)
	if err != nil {
		log.Printf("failed to get package API from registry: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	httpserver.Write(w, r, api)
}

// @Summary Retrieve the logged user principal
// @Description Retrieve the logged user principal containing a list of access controls granted to the user
// @Description use it primarily to log in user interface services and retrieve a list of access controls to inform which
// @Description operations are available to the user via the user interface
// @Tags Access Control
// @Router /user [get]
// @Accepts json
// @Produce plain
// @Failure 401 {string} the request failed to authenticate the user
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	// if we got this far is because the user is authenticated
	// then return the access controls for the user
	user := httpserver.GetUserPrincipal(r)
	if user != nil {
		httpserver.Write(w, r, user)
	}
}

// @Summary Retrieve the service public PGP key
// @Description Retrieve the service public PGP key used to verify the authenticity of the service by pilot agents
// @Tags PGP
// @Router /pub [get]
// @Accepts json
// @Produce plain
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func getKeyHandler(w http.ResponseWriter, r *http.Request) {
	// load the verification key
	path, err := KeyFilePath("verify")
	if err != nil {
		log.Printf("cannot find public PGP key: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	key, err := os.ReadFile(path)
	if err != nil {
		log.Printf("cannot read public PGP key file: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(hex.EncodeToString(key)))
}

// @Summary Registers a Host so that it can be activated
// @Description requests the activation service to reserve an activation for a host of the specified mac-address
// @Tags Activation
// @Router /registration [post]
// @Param command body []types.Registration true "the required registration information"
// @Accepts json
// @Produce plain
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 201 {string} OK
func registrationHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("cannot read payload: %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var registrations []Registration
	err = json.Unmarshal(bytes, &registrations)
	if err != nil {
		log.Printf("cannot unmarshal payload: %s\nthe payload was: '%s'\n", err, string(bytes[:]))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c := core.NewConf()
	for _, registration := range registrations {
		// call activation service and reserve the mac-address for this tenant
		_, err = core.HttpRequest(&http.Client{Timeout: 60 * time.Second}, fmt.Sprintf("%s/provision/%s/%s", c.GetActivationURI(), c.GetTenant(), registration.MacAddress), "POST", c.GetActivationUser(), c.GetActivationPwd(), 201)
		if err != nil {
			log.Printf("cannot provision mac-address %s with activation service: %s\n", registration.MacAddress, err)
			http.Error(w, fmt.Sprintf("cannot provision mac-address %s with activation service, check the server logs for more information\n", registration.MacAddress), http.StatusInternalServerError)
			return
		}
		// as the provisioning of the mac-address has been successful records the host in pilot-ctl db
		err = core.Api().SetRegistration(registration)
		if err != nil {
			log.Printf("cannot record registration information in database: %s\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
}

// activationHandler notifies PilotCtl of a Host Activation
// used by the activation service to notify pilot control that a host has been issued with an activation key
// not in swagger as it authenticates with activation service credentials
func activationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	macAddr := vars["macAddress"]
	ma, err := url.PathUnescape(macAddr)
	if err != nil {
		log.Printf("failed to unescape mac-address '%s': %s\n", macAddr, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	uuid := vars["uuid"]
	id, err := url.PathUnescape(uuid)
	if err != nil {
		log.Printf("failed to unescape host UUID '%s': %s\n", uuid, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = core.Api().AdmitRegistered(ma, id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func isErr(w http.ResponseWriter, err error, statusCode int, msg string) bool {
	if err != nil {
		msg = fmt.Sprintf("%s: %s\n", msg, err)
		log.Printf(msg)
		w.WriteHeader(statusCode)
		w.Write([]byte(msg))
		return true
	}
	return false
}
