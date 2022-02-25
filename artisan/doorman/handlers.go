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
	"fmt"
	"github.com/gatblau/onix/artisan/doorman/core"
	_ "github.com/gatblau/onix/artisan/doorman/docs"
	"github.com/gatblau/onix/artisan/doorman/types"
	util "github.com/gatblau/onix/oxlib/httpserver"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"strconv"
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
	err := util.Unmarshal(r, key)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	// validate the data in the key struct
	if util.IsErr(w, key.Valid(), http.StatusBadRequest, "invalid key data") {
		return
	}
	db := core.NewDb()
	var resultCode int
	err, resultCode = db.UpsertKey(key)
	if util.IsErr(w, err, resultCode, "cannot update key in database") {
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
	err := util.Unmarshal(r, cmd)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if util.IsErr(w, err, http.StatusBadRequest, "cannot unmarshal command data") {
		return
	}
	// validate the data in the key struct
	if util.IsErr(w, cmd.Valid(), http.StatusBadRequest, "invalid command data") {
		return
	}
	db := core.NewDb()
	var resultCode int
	_, err, resultCode = db.UpsertObject(types.CommandsCollection, cmd)
	if util.IsErr(w, err, resultCode, "cannot update command in database") {
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
	err := util.Unmarshal(r, inRoute)
	if util.IsErr(w, err, http.StatusBadRequest, "cannot unmarshal inbound route data") {
		return
	}
	// validate the data in the key struct
	if util.IsErr(w, inRoute.Valid(), http.StatusBadRequest, "invalid inbound route data") {
		return
	}
	db := core.NewDb()
	var resultCode int
	_, err, resultCode = db.UpsertObject(types.InRouteCollection, inRoute)
	if util.IsErr(w, err, resultCode, "cannot update inbound route in database") {
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
	err := util.Unmarshal(r, outRoute)
	if util.IsErr(w, err, http.StatusBadRequest, "cannot unmarshal outbound route data") {
		return
	}
	// validate the data in the key struct
	if util.IsErr(w, outRoute.Valid(), http.StatusBadRequest, "invalid outbound route data") {
		return
	}
	db := core.NewDb()
	var resultCode int
	_, err, resultCode = db.UpsertObject(types.OutRouteCollection, outRoute)
	if util.IsErr(w, err, resultCode, "cannot update outbound route in database") {
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
	err = util.Unmarshal(r, pipe)
	if util.IsErr(w, err, http.StatusBadRequest, "cannot unmarshal pipeline data") {
		return
	}
	// validate the data in the key struct
	if util.IsErr(w, pipe.Valid(), http.StatusBadRequest, "invalid pipeline data") {
		return
	}
	db := core.NewDb()
	err, code = db.UpsertPipeline(*pipe)
	if util.IsErr(w, err, http.StatusBadRequest, "cannot create or update pipeline configuration") {
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
	err = util.Unmarshal(r, template)
	if util.IsErr(w, err, http.StatusBadRequest, "cannot unmarshal notification template data") {
		return
	}
	// validate the data in the key struct
	if util.IsErr(w, template.Valid(), http.StatusBadRequest, "invalid notification template data") {
		return
	}
	db := core.NewDb()
	_, err, code = db.UpsertObject(types.NotificationTemplatesCollection, *template)
	if util.IsErr(w, err, code, "cannot create or update notification template configuration") {
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
	err = util.Unmarshal(r, notification)
	if util.IsErr(w, err, http.StatusBadRequest, "cannot unmarshal notification data") {
		return
	}
	// validate the data in the key struct
	if util.IsErr(w, notification.Valid(), http.StatusBadRequest, "invalid notification data") {
		return
	}
	db := core.NewDb()
	err, code = db.UpsertNotification(*notification)
	if util.IsErr(w, err, code, "cannot create or update notification configuration") {
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
	db := core.NewDb()
	pipe, err := db.FindPipeline(pipeName)
	if util.IsErr(w, err, http.StatusInternalServerError, fmt.Sprintf("cannot retrieve pipeline %s: %s", pipeName, err)) {
		return
	}
	for i := 0; i < len(pipe.OutboundRoutes); i++ {
		pipe.OutboundRoutes[i].PackageRegistry.PrivateKey = "*******"
	}
	util.Write(w, r, pipe)
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
	db := core.NewDb()
	pipelines, err := db.FindAllPipelines()
	if util.IsErr(w, err, http.StatusInternalServerError, fmt.Sprintf("cannot retrieve pipelines: %s", err)) {
		return
	}
	util.Write(w, r, pipelines)
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
	db := core.NewDb()
	pipelines, err := db.FindAllNotifications()
	if util.IsErr(w, err, http.StatusInternalServerError, fmt.Sprintf("cannot retrieve notifications: %s", err)) {
		return
	}
	util.Write(w, r, pipelines)
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
	db := core.NewDb()
	notificationTemplates, err := db.FindAllNotificationTemplates()
	if util.IsErr(w, err, http.StatusInternalServerError, fmt.Sprintf("cannot retrieve notification templates: %s", err)) {
		return
	}
	util.Write(w, r, notificationTemplates)
}

// @Summary Triggers the ingestion of an artisan spec artefacts
// @Description Triggers the ingestion of a specification
// @Tags Webhook
// @Router /event/{service-id}/{bucket-name}/{folder-name} [post]
// @Param deployment-id path string true "a unique identifier for the bucket endpoint (e.g. x-minio-deployment-id for MinIO)"
// @Param bucket-name path string true "the name of the bucket that contains the uploaded files"
// @Param folder-name path string true "the name of the folder within the bucket that contains the uploaded files"
// @Produce plain
// @Failure 400 {string} bad request: the server cannot or will not process the request due to something that is perceived to be a client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing)
// @Failure 500 {string} internal server error: the server encountered an unexpected condition that prevented it from fulfilling the request.
// @Success 201 {string} ingestion process has started
func eventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["service-id"]
	id, err := url.PathUnescape(id)
	if util.IsErr(w, err, http.StatusInternalServerError, "cannot unescape deployment-id") {
		return
	}
	bucketName := vars["bucket-name"]
	bucketName, err = url.PathUnescape(bucketName)
	if util.IsErr(w, err, http.StatusInternalServerError, "cannot unescape bucket-name") {
		return
	}
	folderName := vars["folder-name"]
	folderName, err = url.PathUnescape(folderName)
	if util.IsErr(w, err, http.StatusInternalServerError, "cannot unescape folderName") {
		return
	}
	processor := core.NewProcessor(id, bucketName, folderName)
	processor.Start()
	w.WriteHeader(http.StatusCreated)
}

// @Summary Get authentication information for specified webhook token
// @Description checks that an opaque string / authentication token sent to a webhook has been defined for
// @Description an inbound route, returning required referrer URL and IP white list
// @Description NOTE: this endpoint is called by the proxy to authenticate its webhook
// @Tags Webhook
// @Router /token/{token-value} [get]
// @Param token-value path string true "the authentication token presented to doorman proxy webhook"
// @Produce application/json, application/yaml, application/xml
// @Failure 400 {string} bad request: the server cannot or will not process the request due to something that is perceived to be a client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing)
// @Failure 404 {string} not found: the specified token has not been found
// @Failure 500 {string} internal server error: the server encountered an unexpected condition that prevented it from fulfilling the request.
// @Success 200 {string} success
func getWebhookAuthInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token-value"]
	if len(token) == 0 {
		util.Err(w, http.StatusBadRequest, "token is required")
		return
	}
	db := core.NewDb()
	routes, err := db.FindInboundRoutesByWebHookToken(token)
	if util.IsErr(w, err, http.StatusInternalServerError, fmt.Sprintf("cannot retrieve inbound routes for token %s: %s", token, err)) {
		return
	}
	// if no routes exists then return not found
	if len(routes) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var info []types.WebhookAuthInfo
	for _, route := range routes {
		info = append(info, types.WebhookAuthInfo{
			WebhookToken: route.WebhookToken,
			ReferrerURL:  route.ServiceHost,
			Whitelist:    route.WebhookWhitelist,
			Filter:       route.Filter,
		})
	}
	util.Write(w, r, info)
}

// @Summary get authentication information for all webhook tokens
// @Description get authentication information for all webhook tokens
// @Description an inbound route, returning required referrer URL and IP white list
// @Description NOTE: this endpoint is called by the proxy to authenticate its webhook
// @Tags Webhook
// @Router /token [get]
// @Produce application/json, application/yaml, application/xml
// @Failure 400 {string} bad request: the server cannot or will not process the request due to something that is perceived to be a client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing)
// @Failure 404 {string} not found: the specified token has not been found
// @Failure 500 {string} internal server error: the server encountered an unexpected condition that prevented it from fulfilling the request.
// @Success 200 {string} success
func getWebhookAllAuthInfoHandler(w http.ResponseWriter, r *http.Request) {
	db := core.NewDb()
	routes, err := db.FindAllInRoutes()
	if util.IsErr(w, err, http.StatusInternalServerError, fmt.Sprintf("cannot retrieve inbound routes: %s", err)) {
		return
	}
	var info []types.WebhookAuthInfo
	for _, route := range routes {
		info = append(info, types.WebhookAuthInfo{
			WebhookToken: route.WebhookToken,
			ReferrerURL:  route.ServiceHost,
			Whitelist:    route.WebhookWhitelist,
			Filter:       route.Filter,
		})
	}
	util.Write(w, r, info)
}

// @Summary Gets top jobs
// @Description Gets a list of top job information logs
// @Tags Jobs
// @Router /job [get]
// @Param count query int false "the number of top jobs to retrieve, if not specified returns the top 1 job"
// @Produce application/json, application/yaml, application/xml
// @Failure 400 {string} bad request: the server cannot or will not process the request due to something that is perceived to be a client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing)
// @Failure 404 {string} not found: the requested object does not exist
// @Failure 500 {string} internal server error: the server encountered an unexpected condition that prevented it from fulfilling the request.
// @Success 200 {string} success
func getTopJobsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		count = 1
		err   error
	)
	countStr := r.URL.Query().Get("count")
	if len(countStr) > 0 {
		count, err = strconv.Atoi(countStr)
		if err != nil {
			util.Err(w, http.StatusBadRequest, "job count is invalid")
			return
		}
	}
	db := core.NewDb()
	jobs, err := db.FindTopJobs(count)
	if util.IsErr(w, err, http.StatusInternalServerError, fmt.Sprintf("cannot retrieve jobs: %s", err)) {
		return
	}
	util.Write(w, r, jobs)
}
