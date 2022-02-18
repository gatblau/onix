/*
  Onix Config Manager - Artisan's Doorman Proxy
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package main

import (
	"fmt"
	_ "github.com/gatblau/onix/artisan/doorman/proxy/docs"
	util "github.com/gatblau/onix/oxlib/httpserver"
	"net/http"
	"strings"
)

// @title Artisan's Doorman Proxy
// @version 0.0.4
// @description Notifications & Event Sources for Doorman
// @contact.name gatblau
// @contact.url http://onix.gatblau.org/
// @contact.email onix@gatblau.org
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @Summary Sends a new notification
// @Description sends a notification of the specified type
// @Tags Notifications
// @Router /notify [post]
// @Param notification body Notification true "the notification information to send"
// @Accept application/yaml, application/json
// @Produce plain
// @Failure 400 {string} bad request: the server cannot or will not process the request due to something that is perceived to be a client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing)
// @Failure 500 {string} internal server error: the server encountered an unexpected condition that prevented it from fulfilling the request.
// @Success 201 {string} notification has been sent
func notifyHandler(w http.ResponseWriter, r *http.Request) {
	notification := new(Notification)
	err := util.Unmarshal(r, notification)
	if util.IsErr(w, err, http.StatusBadRequest, "cannot unmarshal notification") {
		return
	}
	if util.IsErr(w, notification.Valid(), http.StatusBadRequest, "invalid notification") {
		return
	}
	switch strings.ToUpper(notification.Type) {
	case "EMAIL":
		err = sendMail(*notification)
		if util.IsErr(w, err, http.StatusBadRequest, "cannot email notification") {
			return
		}
	default:
		util.Err(w, http.StatusBadRequest, fmt.Sprintf("notification type '%s' is not supported", strings.ToUpper(notification.Type)))
	}
}

// @Summary A Webhook for MinIO compatible event sources
// @Description receives a s3:ObjectCreated:Put event sent by a MinIO format compatible source
// @Tags Event Sources
// @Router /events/minio [post]
// @Param event body MinioS3Event true "the notification information to send"
// @Accept application/yaml, application/json
// @Produce plain
// @Failure 400 {string} bad request: the server cannot or will not process the request due to something that is perceived to be a client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing)
// @Failure 500 {string} internal server error: the server encountered an unexpected condition that prevented it from fulfilling the request.
// @Success 201 {string} event has been processed
func minioEventsHandler(w http.ResponseWriter, r *http.Request) {

}
