/*
  Onix Config Manager - Artisan Host Runner
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package main

// @title Artisan Host Runner
// @version 0.0.4
// @description Run Artisan packages with in a host
// @contact.name gatblau
// @contact.url http://onix.gatblau.org/
// @contact.email onix@gatblau.org
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/gatblau/onix/artisan/build"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/merge"
	_ "github.com/gatblau/onix/artisan/runner/host/docs"
	o "github.com/gatblau/onix/artisan/runner/host/onix"
	"github.com/gorilla/mux"
)

// @Summary Build patching artisan package
// @Description Trigger a new build to create artisan package from the vulnerabilty scanned csv report passed in the payload.
// @Tags Runners
// @Router /host/{cmd-key} [post]
// @Param cmd-key path string true "the key of the command to retrieve"
// @Produce plain
// @Param flow body flow.Flow true "the artisan flow to run"
// @Failure 500 {string} there was an error in the server, error the server logs
// @Failure 422 {string} command-key was not found in database, error the server logs
// @Success 200 {string} OK

func executeCommandHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	api := o.Api()
	cmdkey := vars["cmd-key"]
	cmd, err := api.GetCommand(cmdkey)
	if checkErr(w, fmt.Sprintf("%s: [ %s ]\n", "Error while getting command using cmd key ", cmdkey), err) {
		return
	}
	if cmd == nil {
		msg := fmt.Sprintf("No command item for item type ART_FX found in database for cmd key [ %s ] , please check if this item exists ", cmdkey)
		fmt.Printf(msg)
		http.Error(w, msg, http.StatusUnprocessableEntity)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if checkErr(w, "Error while reading http request body ", err) {
		return
	}

	t, err := core.NewTempDir()
	if checkErr(w, "Error while creating temp folder ", err) {
		return
	}

	d := path.Join(t, "context")
	err = core.WriteFile(body, d, "")
	if checkErr(w, fmt.Sprintf("%s: [ %s ]\n", "Error while writing request body to temp path ", d), err) {
		return
	}

	// get the variables in the host environment
	hostEnv := merge.NewEnVarFromSlice(os.Environ())
	// get the variables in the command
	cmdEnv := merge.NewEnVarFromSlice(cmd.Env())
	// if not containerised add PATH to execution environment
	hostEnv.Merge(cmdEnv)
	cmdEnv = hostEnv
	// if running in verbose mode
	if cmd.Verbose {
		// add ARTISAN_DEBUG to execution environment
		cmdEnv.Vars["ARTISAN_DEBUG"] = "true"
	}

	cmdString := fmt.Sprintf("art %s -u %s:%s %s %s --path=%s", "exe", cmd.User, cmd.Pwd, cmd.Package, cmd.Function, t)
	// run and return
	out, err := build.ExeAsync(cmdString, ".", cmdEnv, false)
	if checkErr(w, fmt.Sprintf("Error while executing artisan package function using command [ %s ]", cmdString), err) {
		return
	} else {
		msg := fmt.Sprintf("%s [%s %s ] : [ %s ] \n", "Result of executing artisan package function using command", cmd.Package, cmd.Function, out)
		fmt.Printf(msg)
	}
}
