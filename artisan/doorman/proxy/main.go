/*
  Onix Config Manager - Artisan's Doorman Proxy
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gatblau/onix/artisan/doorman/types"
	"github.com/gatblau/onix/oxlib/httpserver"
	"github.com/gatblau/onix/oxlib/oxc"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	WhInfo      []types.WebhookAuthInfo
	defaultAuth func(r http.Request) *oxc.UserPrincipal
)

func main() {
	// load web hook information from doorman (requires doorman to be up and running)
	WhInfo = loadWhInfo()
	// creates a generic http server
	s := httpserver.New("doorman")
	// add handlers
	s.Http = func(router *mux.Router) {
		// apply authentication
		router.Use(s.AuthenticationMiddleware)

		router.HandleFunc("/events/minio", minioEventsHandler).Methods("POST")
		router.HandleFunc("/notify", notifyHandler).Methods("POST")
	}
	// grab a reference to default auth to use it in the proxy override below
	defaultAuth = s.DefaultAuth
	// set up specific authentication for doorman proxy
	s.Auth = map[string]func(http.Request) *oxc.UserPrincipal{
		"^/events/.*": whAuth,
	}
	s.Serve()
}

// whAuth authenticates web hook requests using opaque string (bearer token)
func whAuth(r http.Request) *oxc.UserPrincipal {
	ip := httpserver.FindRealIP(&r)
	token := r.Header.Get("Authorization")
	for _, info := range WhInfo {
		// if the bearer token is ok
		if strings.HasSuffix(token, info.WebhookToken) {
			var whitelisted bool
			// if a whitelist has been set up for the webhook
			if info.Whitelist != nil && len(info.Whitelist) > 0 {
				// check that the requester real IP is in the whitelist
				for _, listedIp := range info.Whitelist {
					if ip == listedIp {
						whitelisted = true
					}
				}
				// if it is not then block the request
				if !whitelisted {
					return nil
				}
			}
			return &oxc.UserPrincipal{
				Username: "webhook-user",
				Created:  time.Now(),
				Context:  token,
			}
		}
	}
	// try with a admin credentials
	if defaultAuth != nil {
		if principal := defaultAuth(r); principal != nil {
			return principal
		}
	}
	// otherwise, fail authentication
	return nil
}

func loadWhInfo() []types.WebhookAuthInfo {
	fmt.Printf("INFO: contacting doorman to load webhook configuration\n")
	doormanBaseURI, err := getDoormanBaseURI()
	if err != nil {
		fmt.Printf("ERROR: cannot retrieve configuration: %s, exiting\n", err)
		os.Exit(1)
	}
	requestURI := fmt.Sprintf("%s/token", doormanBaseURI)
	resp, err, _ := newRequest("GET", requestURI)
	if err != nil {
		fmt.Printf("ERROR: cannot contact doorman: %s, exiting\n", err)
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ERROR: cannot read doorman's response: %s, exiting\n", err)
		os.Exit(1)
	}
	var info []types.WebhookAuthInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		fmt.Printf("ERROR: cannot unmarshal doorman's response body: %s, exiting\n", err)
		os.Exit(1)
	}
	return info
}
