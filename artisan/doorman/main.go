/*
  Onix Config Manager - Artisan's Doorman
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package main

import (
	"fmt"
	"github.com/gatblau/onix/artisan/doorman/core"
	"github.com/gatblau/onix/oxlib/httpserver"
	"github.com/gatblau/onix/oxlib/oxc"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"time"
)

var defaultAuth func(r http.Request) *oxc.UserPrincipal

func main() {
	// creates a generic http server
	s := httpserver.New("doorman")
	// add handlers
	s.Http = func(router *mux.Router) {
		// enable encoded path  vars
		router.UseEncodedPath()
		// conditionally enable middleware
		if len(os.Getenv("DOORMAN_LOGGING")) > 0 {
			router.Use(s.LoggingMiddleware)
		}
		// apply authentication
		router.Use(s.AuthenticationMiddleware)

		// admin facing endpoints
		router.HandleFunc("/key", upsertKeyHandler).Methods("PUT")
		router.HandleFunc("/command", upsertCommandHandler).Methods("PUT")
		router.HandleFunc("/route/in", upsertInboundRouteHandler).Methods("PUT")
		router.HandleFunc("/route/out", upsertOutboundRouteHandler).Methods("PUT")
		router.HandleFunc("/notification", upsertNotificationHandler).Methods("PUT")
		router.HandleFunc("/notification", getAllNotificationsHandler).Methods("GET")
		router.HandleFunc("/notification-template", upsertNotificationTemplateHandler).Methods("PUT")
		router.HandleFunc("/notification-template", getAllNotificationTemplatesHandler).Methods("GET")
		router.HandleFunc("/pipe", upsertPipelineHandler).Methods("PUT")
		router.HandleFunc("/pipe/{name}", getPipelineHandler).Methods("GET")
		router.HandleFunc("/pipe", getAllPipelinesHandler).Methods("GET")
		router.HandleFunc("/job", getTopJobsHandler).Methods("GET")

		// doorman proxy facing endpoints
		router.HandleFunc("/event/{service-id}/{bucket-name}/{folder-name}", eventHandler).Methods("POST")
		router.HandleFunc("/token/{token-value}", getWebhookAuthInfoHandler).Methods("GET")
		router.HandleFunc("/token", getWebhookAllAuthInfoHandler).Methods("GET")
	}
	// grab a reference to default auth to use it in the proxy override below
	defaultAuth = s.DefaultAuth
	// set up specific authentication for doorman proxy
	s.Auth = map[string]func(http.Request) *oxc.UserPrincipal{
		"^/token.*":  dProxyAuth,
		"^/event/.*": dProxyAuth,
	}
	// https://textkool.com/en/ascii-art-generator?hl=default&vl=default&font=Broadway%20KB&text=dproxy%0A
	fmt.Print(`
++++++++++++++| ONIX CONFIG MANAGER |+++++++++++++++
|    ___   ___   ___   ___   _       __    _       |
|   | | \ / / \ / / \ | |_) | |\/|  / /\  | |\ |   |
|   |_|_/ \_\_/ \_\_/ |_| \ |_|  | /_/--\ |_| \|   |
|                                                  |
+++++++++++|  the artisan's doorman  |++++++++++++++
`)
	s.Serve()
}

// dProxyAuth authenticates doorman's proxy requests using either proxy specific or admin credentials
func dProxyAuth(r http.Request) *oxc.UserPrincipal {
	user, userErr := core.GetProxyUser()
	if userErr != nil {
		fmt.Printf("cannot authenticate proxy: %s", userErr)
		return nil
	}
	pwd, pwdErr := core.GetProxyPwd()
	if pwdErr != nil {
		fmt.Printf("cannot authenticate proxy: %s", pwdErr)
		return nil
	}
	// try proxy specific credentials first
	if r.Header.Get("Authorization") == httpserver.BasicToken(user, pwd) {
		return &oxc.UserPrincipal{
			Username: user,
			Created:  time.Now(),
		}
	} else if defaultAuth != nil {
		// try admin credentials
		if principal := defaultAuth(r); principal != nil {
			return principal
		}
	}
	// otherwise, fail authentication
	return nil
}
