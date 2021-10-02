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
	"github.com/gatblau/onix/client/server"
	"github.com/gatblau/onix/pilotctl/receivers/mongo/core"
	_ "github.com/gatblau/onix/pilotctl/receivers/mongo/docs"
	"github.com/gatblau/onix/pilotctl/types"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	db  *core.Db
	err error
)

func init() {
	db, err = core.NewDb()
	if err != nil {
		// TODO: add retry
		log.Printf("ERROR: cannot connect to database: '%s'\n", err)
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

// @Summary Get filtered events
// @Description Returns a list of syslog entries following the specified filter
// @Tags Query
// @Router /events [get]
// @Param uuid query string false "the host UUID of the entries to retrieve"
// @Param og query string false "the organisation of the device where the syslog entry was created"
// @Param or query string false "the organisation of the device where the syslog entry was created"
// @Param ar query string false "the area of the device where the syslog entry was created"
// @Param lo query string false "the location of the device where the syslog entry was created"
// @Param tag query string false "syslog entry tag"
// @Param pri query string false "the syslog entry priority"
// @Param sev query string false "the syslog entry severity"
// @Param from query string false "the time FROM which syslog entries are shown (time format must be ddMMyyyyHHmmSS)"
// @Param to query string false "the time TO which syslog entries are shown (time format must be ddMMyyyyHHmmSS)"
// @Produce json
// @Failure 500 {string} there was an error in the server, check the server logs
// @Success 200 {string} OK
func eventQueryHandler(w http.ResponseWriter, r *http.Request) {
	filter := make(map[string]interface{})
	filter = appendString(r, filter, "uuid", "host_uuid")
	filter = appendString(r, filter, "og", "org_group")
	filter = appendString(r, filter, "or", "org")
	filter = appendString(r, filter, "ar", "area")
	filter = appendString(r, filter, "lo", "location")
	filter = appendString(r, filter, "tag", "tag")
	filter = appendInt(r, filter, "pri", "priority")
	filter = appendInt(r, filter, "sev", "severity")
	filter = appendDateRange(r, filter, "from", "to", "time")

	events, err := db.Query(filter)
	if err != nil {
		log.Printf("failed to query events: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	server.Write(w, r, events)
}

func appendString(r *http.Request, filter map[string]interface{}, formKey, keyName string) map[string]interface{} {
	keyValue := r.FormValue(formKey)
	if len(keyValue) > 0 {
		filter[keyName] = keyValue
	}
	return filter
}

func appendInt(r *http.Request, filter map[string]interface{}, formKey, keyName string) map[string]interface{} {
	keyValue := r.FormValue(formKey)
	if len(keyValue) > 0 {
		intValue, err := strconv.Atoi(keyValue)
		if err != nil {
			log.Printf("WARNING: filter '%s' discarded: %s\n", formKey, err)
			return filter
		}
		filter[keyName] = intValue
	}
	return filter
}

func appendDateRange(r *http.Request, filter map[string]interface{}, fromKey, toKey, filterKey string) map[string]interface{} {
	// date format to use
	layout := "02012006030405"
	fromValue := r.FormValue(fromKey)
	toValue := r.FormValue(toKey)
	if len(fromValue) > 0 && len(toValue) == 0 {
		dateValue, err := time.Parse(layout, fromValue)
		if err != nil {
			log.Printf("WARNING: filter '%s' discarded: %s\n", fromKey, err)
			return filter
		}
		filter[filterKey] = bson.M{"$gte": dateValue}
		return filter
	}
	if len(fromValue) == 0 && len(toValue) > 0 {
		dateValue, err := time.Parse(layout, toValue)
		if err != nil {
			log.Printf("WARNING: filter '%s' discarded: %s\n", toKey, err)
			return filter
		}
		filter[filterKey] = bson.M{"$lte": dateValue}
		return filter
	}
	dateFromValue, err := time.Parse(layout, fromValue)
	if err != nil {
		log.Printf("WARNING: filter '%s' discarded: %s\n", fromKey, err)
		return filter
	}
	dateToValue, err := time.Parse(layout, toValue)
	if err != nil {
		log.Printf("WARNING: filter '%s' discarded: %s\n", toKey, err)
		return filter
	}
	filter["$and"] = []bson.M{
		{filterKey: bson.M{"$gte": dateFromValue}},
		{filterKey: bson.M{"$lte": dateToValue}},
	}
	return filter
}

func appendLteDate(r *http.Request, filter map[string]interface{}, formKey, keyName string) map[string]interface{} {
	// date format to use
	layout := "02012006030405"
	keyValue := r.FormValue(formKey)
	if len(keyValue) > 0 {
		dateValue, err := time.Parse(layout, keyValue)
		if err != nil {
			log.Printf("WARNING: filter '%s' discarded: %s\n", formKey, err)
			return filter
		}
		filter[keyName] = bson.M{"$lte": dateValue}
	}
	return filter
}
