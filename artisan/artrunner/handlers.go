/*
  Onix Config Manager - Artisan Runner
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package main

// @title Artisan Flow Runner
// @version 0.0.4
// @description Run Artisan flows
// @contact.name gatblau
// @contact.url http://onix.gatblau.org/
// @contact.email onix@gatblau.org
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

import (
	"context"
	"fmt"
	_ "github.com/gatblau/onix/artisan/artrunner/docs"
	"github.com/gatblau/onix/artisan/flow"
	"github.com/gatblau/onix/artisan/tkn"
	"io/ioutil"
	"log"
	"net/http"
)

// @Summary Executes an Artisan flow
// @Description uploads an Artisan flow and triggers the flow execution
// @Tags Flows
// @Router /flow [post]
// @Produce plain
// @Param flow body flow.Flow true "the artisan flow to run"
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func runHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot read request payload: %s", err), http.StatusInternalServerError)
		return
	}
	// unmarshal the flow bytes
	f, err := flow.NewFlow(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot read flow: %s", err), http.StatusInternalServerError)
		return
	}
	// get a tekton builder
	builder := tkn.NewBuilder(f)
	resources, pr, requiresGit := builder.Build()
	ctx := context.Background()
	k8s, err := NewK8S()
	if err != nil {
		msg := fmt.Sprintf("cannot create kubernetes client: %s\n", err)
		log.Printf(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	// create the pipeline resources
	for _, resource := range resources {
		err = k8s.Patch(resource, ctx)
		if err != nil {
			msg := fmt.Sprintf("cannot apply kubernetes resources: %s\n", err)
			log.Printf(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
	}
	// stream the execution logs to the client using the tkn cli:
	// tkn pr logs pipelinerun-name -a -n namespace
	err = execute("tkn", []string{"pr", "logs", pr, "-a", "-n", f.Labels["namespace"]}, w)
	if err != nil {
		msg := "cannot retrieve pipeline logs, pipeline might still be running"
		log.Printf(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// if got source is not required, the it is not a CI pipeline listening for source code changes
	// and therefore it is assume to be transient, when the job is done it should be deleted
	if !requiresGit {
		// now can delete all resources
		err = k8s.DeleteAll(resources, ctx)
		if err != nil {
			log.Printf(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	execute("tkn", []string{"resource", "list"}, w)
}
