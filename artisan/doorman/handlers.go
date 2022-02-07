/*
  Onix Config Manager - Artisan's Doorman
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package main

// @title Artisan's Doorman
// @version 0.0.4
// @description Transfer (pull, verify, scan, resign and push) artefacts between networks
// @contact.name gatblau
// @contact.url http://onix.gatblau.org/
// @contact.email onix@gatblau.org
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

import (
	"encoding/json"
	"fmt"
	"github.com/gatblau/onix/artisan/doorman/core"
	_ "github.com/gatblau/onix/artisan/doorman/docs"
	"github.com/gatblau/onix/artisan/doorman/types"
	"io/ioutil"
	"log"
	"net/http"
)

// @Summary Upload a new key
// @Description uploads a new key used by doorman for cryptographic operations
// @Tags Keys
// @Router /key [post]
// @Param key body types.Key true "the data for the key to persist"
// @Produce plain
// @Failure 400 {string} bad request: the server cannot or will not process the request due to something that is perceived to be a client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing)
// @Failure 500 {string} internal server error: the server encountered an unexpected condition that prevented it from fulfilling the request.
// @Success 201 {string} created
func newKeyHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if isErr(w, err, http.StatusBadRequest, "cannot read request body") {
		return
	}
	key := new(types.Key)
	err = json.Unmarshal(body, key)
	if isErr(w, err, http.StatusBadRequest, "cannot unmarshal request body") {
		return
	}
	// validate the data in the key struct
	if isErr(w, key.Valid(), http.StatusBadRequest, "invalid payload") {
		return
	}
	db := core.NewDb()
	if db.ObjectExists(types.KeysColl, key.Name) {
		isErr(w, fmt.Errorf("key with name %s already exist\n", key.Name), http.StatusBadRequest, "")
		return
	}
	_, err = db.InsertObject(types.KeysColl, key)
	if isErr(w, err, http.StatusInternalServerError, "cannot insert key in database") {
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// @Summary Create a new command
// @Description creates  a new command
// @Tags Commands
// @Router /command [post]
// @Param key body types.Command true "the data for the command to persist"
// @Produce plain
// @Failure 400 {string} bad request: the server cannot or will not process the request due to something that is perceived to be a client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing)
// @Failure 500 {string} internal server error: the server encountered an unexpected condition that prevented it from fulfilling the request.
// @Success 201 {string} created
func newCommandHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if isErr(w, err, http.StatusBadRequest, "cannot read request body") {
		return
	}
	cmd := new(types.Command)
	err = json.Unmarshal(body, cmd)
	if isErr(w, err, http.StatusBadRequest, "cannot unmarshal request body") {
		return
	}
	// validate the data in the key struct
	if isErr(w, cmd.Valid(), http.StatusBadRequest, "invalid payload") {
		return
	}
	db := core.NewDb()
	if db.ObjectExists(types.CommandsColl, cmd.Name) {
		isErr(w, fmt.Errorf("command with name %s already exists\n", cmd.Name), http.StatusBadRequest, "")
		return
	}
	_, err = db.InsertObject(types.CommandsColl, cmd)
	if isErr(w, err, http.StatusInternalServerError, "cannot insert command in database") {
		return
	}
	w.WriteHeader(http.StatusCreated)
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
