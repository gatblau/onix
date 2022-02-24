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
	"net/url"
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
// @Accept application/json, application/yaml
// @Produce plain
// @Failure 400 {string} bad request: the server cannot or will not process the request due to something that is perceived to be a client error (e.g., malformed request syntax, invalid request message framing, or deceptive request routing)
// @Failure 500 {string} internal server error: the server encountered an unexpected condition that prevented it from fulfilling the request.
// @Success 201 {string} event has been processed
func minioEventsHandler(w http.ResponseWriter, r *http.Request) {
	event := new(MinioS3Event)
	err := util.Unmarshal(r, event)
	if util.IsErr(w, err, http.StatusBadRequest, "cannot unmarshal webhook payload") {
		return
	}
	if event.Records == nil {
		util.Err(w, http.StatusBadRequest, "incorrect webhook payload, missing Records, cannot continue")
		return
	}
	object := event.Records[0].S3.Object
	if !strings.HasSuffix(object.Key, "spec.yaml") {
		util.Err(w, http.StatusBadRequest, fmt.Sprintf("invalid event, changed object was %s but required spec.yaml", object.Key))
		return
	}
	key, err := url.PathUnescape(object.Key)
	if util.IsErr(w, err, http.StatusBadRequest, fmt.Sprintf("cannot unescape object key %s", object.Key)) {
		return
	}
	cut := strings.LastIndex(key, "/")
	// get the path within the bucket
	folderName := key[:cut]
	// get the unique identifier for the bucket
	deploymentId := event.Records[0].ResponseElements.XMinioDeploymentID
	// get the bucket name
	bucketName := event.Records[0].S3.Bucket.Name
	// constructs the URI of the object that changed
	doormanBaseURI, err := getDoormanBaseURI()
	if util.IsErr(w, err, http.StatusInternalServerError, "missing configuration") {
		return
	}
	fmt.Printf("︎⚡️ new release:\n")
	fmt.Printf("  ✔ from   = %s\n", event.Records[0].ResponseElements.XMinioOriginEndpoint)
	fmt.Printf("  ✔ id     = %s\n", deploymentId)
	fmt.Printf("  ✔ bucket = %s\n", bucketName)
	fmt.Printf("  ✔ folder = %s\n", folderName)
	fmt.Println("  ✔ type   = minio compatible")
	requestURI := fmt.Sprintf("%s/event/%s/%s/%s", doormanBaseURI, deploymentId, bucketName, folderName)
	if _, postErr, code := newRequest("POST", requestURI); postErr != nil {
		w.WriteHeader(code)
		w.Write([]byte(postErr.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
}
