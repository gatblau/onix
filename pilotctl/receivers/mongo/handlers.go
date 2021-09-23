package main

/*
  Onix Config Manager - MongoDb event receiver for Pilot Control
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"encoding/json"
	"github.com/gatblau/onix/pilotctl/receivers/mongo/core"
	_ "github.com/gatblau/onix/pilotctl/receivers/mongo/docs"
	"github.com/gatblau/onix/pilotctl/types"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	db  *core.Db
	err error
)

func init() {
	db, err = core.NewDb()
	if err != nil {
		// TODO: add retry
		log.Printf("cannot connect to database: '%s'\n", err)
	}
}

// @title MongoDB Event Receiver for Pilot Control
// @version 0.0.4
// @description Onix Config Manager Event Receiver for Pilot Control using MongoDb
// @contact.name gatblau
// @contact.url http://onix.gatblau.org/
// @contact.email onix@gatblau.org
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @Summary Submit Syslog Events
// @Description submits syslog events to be persisted for further use
// @Tags Receiver
// @Router /events [post]
// @Param command body types.Events true "the events to submit"
// @Accepts json
// @Produce plain
// @Failure 400 {string} there was an error in the server trying to read or unmarshal the http request body, check the server logs
// @Failure 500 {string} there was an error in the server trying to persist events to the database, check the server logs
// @Success 200 {string} OK
func eventReceiverHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read request body: %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var events types.Events
	err = json.Unmarshal(bytes, &events)
	if err != nil {
		log.Printf("failed to unmarshal events: %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = db.Insert(&events)
	if err != nil {
		log.Printf("failed to insert events: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
