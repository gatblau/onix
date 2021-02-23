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
	// get a tekton builder
	builder := tkn.NewBuilder(f)

	resources := builder.BuildSlice()

	ctx := context.Background()
	k8s, err := NewK8S()
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot create kubernetes client: %s", err), http.StatusInternalServerError)
		return
	}
	for _, resource := range resources {
		err = k8s.Apply(string(resource), ctx)
		http.Error(w, fmt.Sprintf("cannot apply kubernetes resources: %s", err), http.StatusInternalServerError)
		return
	}
}
