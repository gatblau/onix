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
	"github.com/gorilla/mux"
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
// @Failure 500 {string} there was an error in the server, error the server logs
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
	resources, pr, _ := builder.Build()
	ctx := context.Background()
	k8s, err := NewK8S()
	if checkErr(w, "cannot create kubernetes client", err) {
		return
	}
	// create the pipeline resources
	for _, resource := range resources {
		err = k8s.Patch(resource, ctx)
		if checkErr(w, "cannot apply kubernetes resources", err) {
			return
		}
	}
	msg := fmt.Sprintf("? starting pipeline run %s\n", pr)
	log.Print(msg)
	w.Write([]byte(msg))

	// // stream the execution logs to the client using the tkn cli:
	// // tkn pr logs pipelinerun-name -a -f -n namespace
	// err = execute("tkn", []string{"pr", "logs", pr, "-a", "-f", "-n", f.Labels["namespace"]}, w)
	// if err != nil {
	// 	msg := "cannot retrieve pipeline logs, pipeline might still be running"
	// 	log.Printf(msg)
	// 	http.Error(w, msg, http.StatusInternalServerError)
	// 	return
	// }
	//
	// // if got source is not required, the it is not a CI pipeline listening for source code changes
	// // and therefore it is assume to be transient, when the job is done it should be deleted
	// if !requiresGit {
	// 	// now can delete all resources
	// 	err = k8s.DeleteAll(resources, ctx)
	// 	if err != nil {
	// 		log.Printf(err.Error())
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// }
}

// @Summary Launch an existing flow from a Git commit
// @Description starts a flow execution from a commit in a Git repository
// @Tags Flows
// @Router /webhook/{namespace}/{flow-name} [post]
// @Param namespace path string true "the kubernetes namespace where the pipeline run is created"
// @Param flow-name path string true "the name of the flow to run"
// @Produce plain
// @Failure 500 {string} there was an error in the server, error the server logs
// @Success 200 {string} OK
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	flowName := vars["flow-name"]
	namespace := vars["namespace"]

	builder := new(tkn.Builder)
	pr := builder.NewNamedPipelineRun(flowName, namespace)

	ctx := context.Background()
	k8s, err := NewK8S()
	if checkErr(w, "cannot create kubernetes client", err) {
		return
	}

	err = k8s.Patch(tkn.ToYaml(pr, "pipelinerun"), ctx)
	if checkErr(w, "cannot create pipelinerun", err) {
		return
	}
}

func checkErr(w http.ResponseWriter, msg string, err error) bool {
	if err != nil {
		msg := fmt.Sprintf("%s: %s\n", msg, err)
		fmt.Printf(msg)
		http.Error(w, msg, http.StatusInternalServerError)
	}
	return err != nil
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	execute("tkn", []string{"resource", "list"}, w)
}
