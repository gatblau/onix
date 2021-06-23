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
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/flow"
	"github.com/gatblau/onix/artisan/tkn"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

// @Summary Creates an Artisan flow
// @Description creates a new flow from the definition passed in the payload and starts its execution
// @Tags Flows
// @Router /flow [post]
// @Produce plain
// @Param flow body flow.Flow true "the artisan flow to run"
// @Failure 500 {string} there was an error in the server, error the server logs
// @Success 200 {string} OK
func createFlowFromPayloadHandler(w http.ResponseWriter, r *http.Request) {
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

	// // this logic requires tkn cli in the container image
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

// @Summary Creates an Artisan flow from a flow spec stored as an Onix configuration item
// @Description creates a new flow from the definition passed in the payload and starts its execution
// @Tags Flows
// @Router /flow/key/{flow-key}/ns/{namespace} [post]
// @Produce plain
// @Param namespace path string true "the kubernetes namespace where the flow is created"
// @Param flow-key path string true "the unique key of the flow specification in Onix configuration database"
// @Param file body string false "any configuration information sent by the client to the execution context"
// @Failure 500 {string} there was an error in the server, error the server logs
// @Success 200 {string} OK
func createFlowFromConfigHandler(w http.ResponseWriter, r *http.Request) {
	// read path variables
	vars := mux.Vars(r)
	flowKey := vars["flow-key"]
	namespace := vars["namespace"]
	// read the payload body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := fmt.Sprintf("cannot read request payload: %s", err)
		log.Printf("%s\n", msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	// fetches the flow specification from onix
	flowSpec, err := readFlow(flowKey)
	if err != nil {
		msg := fmt.Sprintf("cannot retrieve Artisan Flow specification from Onix: %s", err)
		log.Printf("%s\n", msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	// unmarshal the flow bytes
	f, err := flow.NewFlow(flowSpec)
	if err != nil {
		msg := fmt.Sprintf("cannot create Artisan flow: %s", err)
		log.Printf("%s\n", msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	// set the namespace label in the flow (required by the transpiler to create namespaced resources)
	f.Labels["namespace"] = namespace

	// add any payload in the request to the flow as a file on each step
	// if a payload exist
	if body != nil {
		context := &data.File{
			Name:        "context",
			Description: "the webhook payload data persisted as a context file in Tekton",
			Path:        "context",
			Content:     string(body),
		}
		// add the context file to every step
		for _, step := range f.Steps {
			if step.Input == nil {
				step.Input = &data.Input{File: data.Files{context}}
			} else {
				if step.Input.File == nil {
					step.Input.File = data.Files{context}
				} else {
					step.Input.File = append(step.Input.File, context)
				}
			}
		}
	}
	// get a tekton builder
	builder := tkn.NewBuilder(f)
	resources, pr, _ := builder.Build()
	ctx := context.Background()
	k8s, err := NewK8S()
	if err != nil {
		msg := fmt.Sprintf("cannot create K8S client: %s", err)
		log.Printf("%s\n", msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	// create the pipeline resources
	for _, resource := range resources {
		err = k8s.Patch(resource, ctx)
		if err != nil {
			msg := fmt.Sprintf("cannot apply K8S resources: %s", err)
			log.Printf("%s\n", msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
	}
	msg := fmt.Sprintf("? starting pipeline run %s\n", pr)
	log.Print(msg)
	w.Write([]byte(msg))
}

// @Summary Launch an existing flow (typically, from a Git commit hook)
// @Description starts the execution of a pre-existing flow based on its name and the namespace where is located
// @Tags Flows
// @Router /flow/name/{flow-name}/ns/{namespace} [post]
// @Param namespace path string true "the kubernetes namespace where the pipeline run is created"
// @Param flow-name path string true "the name of the flow to run"
// @Produce plain
// @Failure 500 {string} there was an error in the server, error the server logs
// @Success 200 {string} OK
func runFlowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	flowName := vars["flow-name"]
	namespace := vars["namespace"]

	builder := new(tkn.Builder)
	pr := builder.NewNamedPipelineRun(flowName, namespace)

	// need to add git repo in resources of pipeline run
	pr.Spec.Resources = []*tkn.Resources{
		{
			Name: fmt.Sprintf("%s-code-repo", flowName),
			ResourceRef: &tkn.ResourceRef{
				Name: fmt.Sprintf("%s-code-repo", flowName),
			},
		},
	}

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

// @Summary Retrieve a configuration flow from Onix
// @Description connect to Onix and retrieves a flow using its configuration item natural key in Onix
// @Tags Flows
// @Router /flow/key/{flow-key} [get]
// @Produce plain
// @Param flow-key path string true "the unique key of the flow specification in Onix configuration database"
// @Failure 500 {string} the health check failed with an error, check server logs for details
// @Success 200 {string} OK, the health check succeeded
func getFlowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	flowKey := vars["flow-key"]
	// fetches the flow specification from onix
	flowSpec, err := readFlow(flowKey)
	if err != nil {
		msg := fmt.Sprintf("cannot retrieve Artisan Flow specification from Onix: %s", err)
		log.Printf("%s\n", msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	// unmarshal the flow bytes
	f, err := flow.NewFlow(flowSpec)
	if err != nil {
		msg := fmt.Sprintf("cannot create Artisan flow: %s", err)
		log.Printf("%s\n", msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	bytes, err := f.JsonBytes()
	if err != nil {
		msg := fmt.Sprintf("cannot serialise Artisan flow before sending response: %s", err)
		log.Printf("%s\n", msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

func checkErr(w http.ResponseWriter, msg string, err error) bool {
	if err != nil {
		msg := fmt.Sprintf("%s: %s\n", msg, err)
		fmt.Printf(msg)
		http.Error(w, msg, http.StatusInternalServerError)
	}
	return err != nil
}
