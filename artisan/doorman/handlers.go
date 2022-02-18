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
// @description Transfer (pull, verify, scan, resign and push) artefacts between repositories
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
	"github.com/gatblau/onix/oxlib/httpserver"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
)

// @Summary Creates or updates a cryptographic key
// @Description creates or updates a cryptographic key used by either inbound or outbound routes to verify or sign
// @Description packages respectively
// @Tags Keys
// @Router /key [put]
// @Param key body types.Key true "the data for the key to persist"
// @Accept application/yaml, application/json
// @Produce plain
// @Failure 400 {string} bad request: the server cannot or will not process the request due to something that is perceived to be a client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing)
// @Failure 500 {string} internal server error: the server encountered an unexpected condition that prevented it from fulfilling the request.
// @Success 200 {string} object has been updated
// @Success 201 {string} object has been created
func upsertKeyHandler(w http.ResponseWriter, r *http.Request) {
	key := new(types.Key)
	err := unmarshal(r, key)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	// validate the data in the key struct
	if isErr(w, key.Valid(), http.StatusBadRequest, "invalid key data") {
		return
	}
	db := core.NewDb()
	var resultCode int
	_, err, resultCode = db.UpsertObject(types.KeysCollection, key)
	if isErr(w, err, resultCode, "cannot update key in database") {
		return
	}
	w.WriteHeader(resultCode)
}

// @Summary Creates or updates a command
// @Description creates or updates a command
// @Tags Commands
// @Router /command [put]
// @Param key body types.Command true "the data for the command to persist"
// @Accept application/yaml, application/json
// @Produce plain
// @Failure 400 {string} bad request: the server cannot or will not process the request due to something that is perceived to be a client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing)
// @Failure 500 {string} internal server error: the server encountered an unexpected condition that prevented it from fulfilling the request.
// @Success 200 {string} object has been updated
// @Success 201 {string} object has been created
func upsertCommandHandler(w http.ResponseWriter, r *http.Request) {
	cmd := new(types.Command)
	err := unmarshal(r, cmd)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if isErr(w, err, http.StatusBadRequest, "cannot unmarshal command data") {
		return
	}
	// validate the data in the key struct
	if isErr(w, cmd.Valid(), http.StatusBadRequest, "invalid command data") {
		return
	}
	db := core.NewDb()
	var resultCode int
	_, err, resultCode = db.UpsertObject(types.CommandsCollection, cmd)
	if isErr(w, err, resultCode, "cannot update command in database") {
		return
	}
	w.WriteHeader(resultCode)
}

// @Summary Creates or updates an inbound route
// @Description creates or updates an inbound route
// @Tags Routes
// @Router /route/in [put]
// @Param key body types.InRoute true "the data for the inbound route to persist"
// @Accept application/yaml, application/json
// @Produce plain
// @Failure 400 {string} bad request: the server cannot or will not process the request due to something that is perceived to be a client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing)
// @Failure 500 {string} internal server error: the server encountered an unexpected condition that prevented it from fulfilling the request.
// @Success 200 {string} object has been updated
// @Success 201 {string} object has been created
func upsertInboundRouteHandler(w http.ResponseWriter, r *http.Request) {
	inRoute := new(types.InRoute)
	err := unmarshal(r, inRoute)
	if isErr(w, err, http.StatusBadRequest, "cannot unmarshal inbound route data") {
		return
	}
	// validate the data in the key struct
	if isErr(w, inRoute.Valid(), http.StatusBadRequest, "invalid inbound route data") {
		return
	}
	db := core.NewDb()
	var resultCode int
	_, err, resultCode = db.UpsertObject(types.InRouteCollection, inRoute)
	if isErr(w, err, resultCode, "cannot update inbound route in database") {
		return
	}
	w.WriteHeader(resultCode)
}

// @Summary Creates or updates an inbound route
// @Description creates or updates an inbound route
// @Tags Routes
// @Router /route/out [put]
// @Param key body types.OutRoute true "the data for the outbound route to persist"
// @Accept application/yaml, application/json
// @Produce plain
// @Failure 400 {string} bad request: the server cannot or will not process the request due to something that is perceived to be a client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing)
// @Failure 500 {string} internal server error: the server encountered an unexpected condition that prevented it from fulfilling the request.
// @Success 200 {string} object has been updated
// @Success 201 {string} object has been created
func upsertOutboundRouteHandler(w http.ResponseWriter, r *http.Request) {
	outRoute := new(types.OutRoute)
	err := unmarshal(r, outRoute)
	if isErr(w, err, http.StatusBadRequest, "cannot unmarshal outbound route data") {
		return
	}
	// validate the data in the key struct
	if isErr(w, outRoute.Valid(), http.StatusBadRequest, "invalid outbound route data") {
		return
	}
	db := core.NewDb()
	var resultCode int
	_, err, resultCode = db.UpsertObject(types.OutRouteCollection, outRoute)
	if isErr(w, err, resultCode, "cannot update outbound route in database") {
		return
	}
	w.WriteHeader(resultCode)
}

// @Summary Creates or updates an inbound route
// @Description creates or updates an inbound route
// @Tags Pipelines
// @Router /pipe [put]
// @Param key body types.PipelineConf true "the data for the pipeline to persist"
// @Accept application/yaml, application/json
// @Produce plain
// @Failure 400 {string} bad request: the server cannot or will not process the request due to something that is perceived to be a client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing)
// @Failure 500 {string} internal server error: the server encountered an unexpected condition that prevented it from fulfilling the request.
// @Success 200 {string} object has been updated
// @Success 201 {string} object has been created
func upsertPipelineHandler(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		code int
	)
	pipe := new(types.PipelineConf)
	err = unmarshal(r, pipe)
	if isErr(w, err, http.StatusBadRequest, "cannot unmarshal pipeline data") {
		return
	}
	// validate the data in the key struct
	if isErr(w, pipe.Valid(), http.StatusBadRequest, "invalid pipeline data") {
		return
	}
	err, code = core.UpsertPipeline(*pipe)
	if isErr(w, err, http.StatusBadRequest, "cannot create or update pipeline configuration") {
		return
	}
	w.WriteHeader(code)
}

// @Summary Creates or updates a notification template
// @Description creates or updates a notification template
// @Tags Notifications
// @Router /notification-template [put]
// @Param key body types.NotificationTemplate true "the data for the notification template to persist"
// @Accept application/yaml, application/json
// @Produce plain
// @Failure 400 {string} bad request: the server cannot or will not process the request due to something that is perceived to be a client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing)
// @Failure 500 {string} internal server error: the server encountered an unexpected condition that prevented it from fulfilling the request.
// @Success 200 {string} object has been updated
// @Success 201 {string} object has been created
func upsertNotificationTemplateHandler(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		code int
	)
	template := new(types.NotificationTemplate)
	err = unmarshal(r, template)
	if isErr(w, err, http.StatusBadRequest, "cannot unmarshal notification template data") {
		return
	}
	// validate the data in the key struct
	if isErr(w, template.Valid(), http.StatusBadRequest, "invalid notification template data") {
		return
	}
	db := core.NewDb()
	_, err, code = db.UpsertObject(types.NotificationTemplatesCollection, *template)
	if isErr(w, err, code, "cannot create or update notification template configuration") {
		return
	}
	w.WriteHeader(code)
}

// @Summary Creates or updates a notification
// @Description creates or updates a notification
// @Tags Notifications
// @Router /notification [put]
// @Param key body types.Notification true "the data for the notification to persist"
// @Accept application/yaml, application/json
// @Produce plain
// @Failure 400 {string} bad request: the server cannot or will not process the request due to something that is perceived to be a client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing)
// @Failure 500 {string} internal server error: the server encountered an unexpected condition that prevented it from fulfilling the request.
// @Success 200 {string} object has been updated
// @Success 201 {string} object has been created
func upsertNotificationHandler(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		code int
	)
	notification := new(types.Notification)
	err = unmarshal(r, notification)
	if isErr(w, err, http.StatusBadRequest, "cannot unmarshal notification data") {
		return
	}
	// validate the data in the key struct
	if isErr(w, notification.Valid(), http.StatusBadRequest, "invalid notification data") {
		return
	}
	err, code = core.UpsertNotification(*notification)
	if isErr(w, err, code, "cannot create or update notification configuration") {
		return
	}
	w.WriteHeader(code)
}

// @Summary Gets a pipeline
// @Description gets a pipeline
// @Tags Pipelines
// @Router /pipe/{name} [get]
// @Param name path string true "the name of the pipeline to retrieve"
// @Produce application/json, application/yaml, application/xml
// @Failure 400 {string} bad request: the server cannot or will not process the request due to something that is perceived to be a client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing)
// @Failure 404 {string} not found: the requested object does not exist
// @Failure 500 {string} internal server error: the server encountered an unexpected condition that prevented it from fulfilling the request.
// @Success 200 {string} success
func getPipelineHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pipeName := vars["name"]
	pipe, err := core.FindPipeline(pipeName)
	if isErr(w, err, http.StatusInternalServerError, fmt.Sprintf("cannot retrieve pipeline %s: %s", pipeName, err)) {
		return
	}
	for i := 0; i < len(pipe.OutboundRoutes); i++ {
		pipe.OutboundRoutes[i].PackageRegistry.PrivateKey = "*******"
	}
	httpserver.Write(w, r, pipe)
}

// @Summary Gets all pipelines
// @Description gets all pipelines
// @Tags Pipelines
// @Router /pipe [get]
// @Produce application/json, application/yaml, application/xml
// @Failure 400 {string} bad request: the server cannot or will not process the request due to something that is perceived to be a client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing)
// @Failure 404 {string} not found: the requested object does not exist
// @Failure 500 {string} internal server error: the server encountered an unexpected condition that prevented it from fulfilling the request.
// @Success 200 {string} success
func getAllPipelinesHandler(w http.ResponseWriter, r *http.Request) {
	pipelines, err := core.FindAllPipelines()
	if isErr(w, err, http.StatusInternalServerError, fmt.Sprintf("cannot retrieve pipelines: %s", err)) {
		return
	}
	httpserver.Write(w, r, pipelines)
}

// @Summary Gets all notifications
// @Description gets all notifications
// @Tags Notifications
// @Router /notification [get]
// @Produce application/json, application/yaml, application/xml
// @Failure 400 {string} bad request: the server cannot or will not process the request due to something that is perceived to be a client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing)
// @Failure 404 {string} not found: the requested object does not exist
// @Failure 500 {string} internal server error: the server encountered an unexpected condition that prevented it from fulfilling the request.
// @Success 200 {string} success
func getAllNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	pipelines, err := core.FindAllNotifications()
	if isErr(w, err, http.StatusInternalServerError, fmt.Sprintf("cannot retrieve notifications: %s", err)) {
		return
	}
	httpserver.Write(w, r, pipelines)
}

// @Summary Gets all notification templates
// @Description gets all notification templates
// @Tags Notifications
// @Router /notification-template [get]
// @Produce application/json, application/yaml, application/xml
// @Failure 400 {string} bad request: the server cannot or will not process the request due to something that is perceived to be a client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing)
// @Failure 404 {string} not found: the requested object does not exist
// @Failure 500 {string} internal server error: the server encountered an unexpected condition that prevented it from fulfilling the request.
// @Success 200 {string} success
func getAllNotificationTemplatesHandler(w http.ResponseWriter, r *http.Request) {
	notificationTemplates, err := core.FindAllNotificationTemplates()
	if isErr(w, err, http.StatusInternalServerError, fmt.Sprintf("cannot retrieve notification templates: %s", err)) {
		return
	}
	httpserver.Write(w, r, notificationTemplates)
}

// @Summary Triggers the ingestion of an artisan spec artefacts from a MinIO bucket notification
// @Description Triggers the ingestion of a specification from a MinIO bucket notification
// @Tags Events
// @Router /event/minio [post]
// @Param key body types.S3MinioEvent true "the payload sent by minio notification"
// @Require application/json
// @Produce plain
// @Failure 400 {string} bad request: the server cannot or will not process the request due to something that is perceived to be a client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing)
// @Failure 500 {string} internal server error: the server encountered an unexpected condition that prevented it from fulfilling the request.
// @Success 201 {string} ingestion process has started
func eventMinioHandler(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if isErr(w, err, http.StatusBadRequest, "cannot read MinIO event") {
		return
	}
	minIoEvent := new(types.S3MinioEvent)
	err = json.Unmarshal(content, minIoEvent)
	if isErr(w, err, http.StatusBadRequest, "cannot unmarshal MinIO event") {
		return
	}
	ev, err := types.NewSpecEventFromMinio(*minIoEvent)
	if isErr(w, err, http.StatusBadRequest, "invalid MinIO event") {
		return
	}
	core.ProcessAsync(*ev)
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

func unmarshal(r *http.Request, v interface{}) error {
	contentType := r.Header.Get("Content-Type")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("cannot read key data: %s\n", err)
	}
	switch contentType {
	case "application/json":
		err = json.Unmarshal(body, v)
		if err != nil {
			return fmt.Errorf("cannot unmarshal input data: %s\n", err)
		}
	case "application/yaml":
		err = yaml.Unmarshal(body, v)
		if err != nil {
			return fmt.Errorf("cannot unmarshal input data: %s\n", err)
		}
	default:
		return fmt.Errorf("invalid Content-Type %s\n", contentType)
	}
	return nil
}
