/*
   Onix Kube - Copyright (c) 2019 by www.gatblau.org

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
   Unless required by applicable law or agreed to in writing, software distributed under
   the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
   either express or implied.
   See the License for the specific language governing permissions and limitations under the License.

   Contributors to this project, hereby assign copyright in this code to the project,
   to be licensed under the same terms as the rest of the code.
*/

package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

type Webhook struct {
	log    *logrus.Entry
	config WebhookConf
	ox     *Client
	ready  bool
}

// launch a webhook on a TCP port listening for events
func (c *Webhook) Start(client *Client) {
	// set the ox client
	c.ox = client

	// for security reasons, avoid using the DefaultServeMux
	// and instead use a locally-scoped ServeMux
	mux := http.NewServeMux()

	// registers web handlers
	c.log.Tracef("Registering web root / and livelyness probe /live handlers")
	mux.HandleFunc("/", c.rootHandler)
	mux.HandleFunc("/live", c.rootHandler)

	c.log.Tracef("Registering readyness probe handler /ready")
	mux.HandleFunc("/ready", c.readyHandler)

	c.log.Tracef("Registering handler for web path /%s.", c.config.Path)
	mux.HandleFunc(c.config.Path, c.webhookHandler)

	if c.config.Metrics {
		// prometheus metrics
		c.log.Tracef("Metrics is enabled, registering handler for endpoint /metrics.")
		http.Handle("/metrics", promhttp.Handler())
	}

	// creates an http server listening on the specified TCP port
	server := &http.Server{Addr: fmt.Sprintf(":%s", c.config.Port), Handler: mux}

	// runs the server asynchronously
	go func() {
		c.log.Println(fmt.Sprintf("OxKube listening on :%s", c.config.Port))
		if err := server.ListenAndServe(); err != nil {
			c.log.Fatal(err)
		}
	}()

	// creates a channel to pass a SIGINT (ctrl+C) kernel signal with buffer capacity 1
	stop := make(chan os.Signal, 1)

	// sends any SIGINT signal to the stop channel
	signal.Notify(stop, os.Interrupt)

	// waits for the SIGINT signal to be raised (pkill -2)
	<-stop

	// gets a context with some delay to shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	// releases resources if main completes before the delay period elapses
	defer cancel()

	// on error shutdown
	if err := server.Shutdown(ctx); err != nil {
		c.log.Fatal(err)
		c.log.Println("Shutting down Webhook consumer.")
	}
}

func (c *Webhook) readyHandler(w http.ResponseWriter, r *http.Request) {
	if !c.ready {
		c.log.Warnf("Webhook is not ready")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Webhook is not ready"))
		if err != nil {
			c.log.Error(err)
		}
	} else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}
}

func (c *Webhook) webhookHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// get the request data
	event, err := c.readRequestBody(r)

	if err != nil {
		c.log.Errorf("failed to get request data: %s", err)
		return
	}

	c.log.Tracef("request: %s", event)

	// if basic auth enabled
	if strings.ToLower(c.config.AuthMode) == "basic" {
		if r.Header.Get("Authorization") == "" {
			// if no authorisation header is passed, then it prompts a client browser to authenticate
			w.Header().Set("WWW-Authenticate", `Basic realm="oxkube"`)
			w.WriteHeader(http.StatusUnauthorized)
			c.log.Tracef("Unauthorised request.")
			return
		} else {
			// authenticate the request
			requiredToken := NewBasicToken(c.config.Username, c.config.Password)
			providedToken := r.Header.Get("Authorization")
			// if the authentication fails
			if !strings.Contains(providedToken, requiredToken) {
				// returns an unauthorised request
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}
	}

	switch r.Method {
	case "GET":
		_, _ = io.WriteString(w, "OxKube webhook is ready: use an HTTP POST to send events.")
	case "PUT":
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = io.WriteString(w, "OxKube webhook only supports HTTP POST to send events.")
	case "POST":
		fallthrough
	case "DELETE":
		result, err := c.process(event)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			msg := fmt.Sprintf("Error whilst processing request: %s", err)
			_, _ = w.Write([]byte(msg))
			c.log.Error(msg)
			return
		}
		// protective code instead process does not return a result
		if result == nil {
			c.log.Error("No result whilst processing request to %s with payload %s", r.URL.Path, r.Body)
			return
		}
		if result.Error {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(result.Message))
			c.log.Errorf("Error whilst processing request: %s", result.Message)
			return
		}
		if result.Changed {
			if result.Operation == "I" {
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte("created"))
			} else if result.Operation == "U" {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("updated"))
			}
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("nothing to update"))
		}
	}
}

func (c *Webhook) rootHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	switch r.Method {
	case "GET":
		_, _ = io.WriteString(w, fmt.Sprintf("OxKube is ready.\n"+
			"POST events to webhook: /%s.", c.config.Path))
	case "POST":
	case "PUT":
	case "DELETE":
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("405 - Method Not Allowed"))
	}
}

func (c *Webhook) process(event []byte) (*Result, error) {
	// get the kind of K8S object
	chgKind := gjson.GetBytes(event, "Change.kind")
	// get the type of change
	chgType := gjson.GetBytes(event, "Change.type")

	switch strings.ToLower(chgKind.String()) {
	case "namespace":
		switch strings.ToLower(chgType.String()) {
		case "create":
			fallthrough
		case "update":
			return c.ox.putNamespace(event)
		case "delete":
			return c.ox.deleteResource("item", NS(event))
		}
	case "pod":
		switch strings.ToLower(chgType.String()) {
		case "create":
			fallthrough
		case "update":
			return c.ox.putPod(event)
		case "delete":
			return c.ox.deleteResource("item", itemKey(event, PodNameTag))
		}
	case "service":
		switch strings.ToLower(chgType.String()) {
		case "create":
			fallthrough
		case "update":
			return c.ox.putService(event)
		case "delete":
			return c.ox.deleteResource("item", itemKey(event, ServiceNameTag))
		}
	case "persistent_volume_claim":
		switch strings.ToLower(chgType.String()) {
		case "create":
			fallthrough
		case "update":
			return c.ox.putPersistentVolumeClaim(event)
		case "delete":
			return c.ox.deleteResource("item", itemKey(event, PersistentVolumeClaimNameTag))
		}
	case "replication_controller":
		switch strings.ToLower(chgType.String()) {
		case "create":
			fallthrough
		case "update":
			return c.ox.putReplicationController(event)
		case "delete":
			return c.ox.deleteResource("item", itemKey(event, ReplicationControllerNameTag))
		}
	case "ingress":
		switch strings.ToLower(chgType.String()) {
		case "create":
			fallthrough
		case "update":
			fallthrough
		case "delete":
			c.log.Trace("ingress recording not implemented")
		}
	case "resourcequota":
		switch strings.ToLower(chgType.String()) {
		case "create":
			fallthrough
		case "update":
			return c.ox.putResourceQuota(event)
		case "delete":
			return c.ox.deleteResource("item", itemKey(event, ResourceQuotaNameTag))
		}
	}
	return nil, nil
}

// unmarshal the http request into a json like structure
func (c *Webhook) readRequestBody(r *http.Request) ([]byte, error) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (c *Webhook) handleOther(path string, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != path {
		http.NotFound(w, r)
		return
	}
}
