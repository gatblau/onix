/*
   Terraform Http Backend - Onix - Copyright (c) 2018 by www.gatblau.org

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
	. "gatblau.org/onix/wapic"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Service struct {
	log   *logrus.Entry
	ready bool
	ox    *Client  // client to connect to Onix WAPI
	conf  *SvcConf // configuration for the http service enpoint
}

// creates a new http backend service
func NewService(backend Backend) *Service {
	svc := new(Service)
	svc.conf = backend.config.Service
	svc.log = backend.log
	svc.ox = backend.client
	svc.ready = backend.ready
	return svc
}

// launch the http backend on a TCP port
func (s *Service) Start() {
	mux := mux.NewRouter()

	// registers web handlers
	s.log.Tracef("Registering web root / and liveliness probe /live handlers")
	pattern := fmt.Sprintf("/%s/{key}", s.conf.Path)
	mux.HandleFunc(pattern, s.rootHandler)
	mux.HandleFunc("/live", s.liveHandler)

	s.log.Tracef("Registering readiness probe handler /ready")
	mux.HandleFunc("/ready", s.readyHandler)

	if s.conf.Metrics {
		// prometheus metrics
		s.log.Tracef("Metrics is enabled, registering handler for endpoint /metrics.")
		mux.Handle("/metrics", promhttp.Handler())
	}

	// creates an http server listening on the specified TCP port
	server := &http.Server{Addr: fmt.Sprintf(":%s", s.conf.Port), Handler: mux}

	// runs the server asynchronously
	go func() {
		s.log.Println(fmt.Sprintf("Terraform Http Backend listening on :%s", s.conf.Port))
		if err := server.ListenAndServe(); err != nil {
			s.log.Fatal(err)
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
		s.log.Fatal(err)
		s.log.Println("Shutting down Terraform Http Backend service.")
	}
}

func (c *Service) rootHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	defer r.Body.Close()
	switch r.Method {
	case "GET":
		_, _ = io.WriteString(w, fmt.Sprintf("Retrieving state for %s", vars["key"]))
	case "POST":
	case "PUT":
	case "DELETE":
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("405 - Method Not Allowed"))
	}
}

func (s *Service) readyHandler(w http.ResponseWriter, r *http.Request) {
	if !s.ready {
		s.log.Warnf("Terraform HTTP Backend service is not ready")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Terraform HTTP Backend service is not ready"))
		if err != nil {
			s.log.Error(err)
		}
	} else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}
}

func (s *Service) liveHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
