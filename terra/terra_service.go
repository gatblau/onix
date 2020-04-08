/*
   Onix Config Manager - OxTerra - Terraform Http Backend for Onix
   Copyright (c) 2018-2020 by www.gatblau.org

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
	"encoding/base64"
	"encoding/json"
	"fmt"
	. "github.com/gatblau/oxc"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

type TerraService struct {
	ready bool
	ox    *Client // client to connect to Onix WAPI
	conf  *Config
}

// creates a new http backend service
func NewTerraService(backend Backend) *TerraService {
	svc := new(TerraService)
	svc.conf = backend.config
	svc.ox = backend.client
	svc.ready = backend.ready
	return svc
}

// launch the http backend on a TCP port
func (s *TerraService) Start() {
	mux := mux.NewRouter()
	mux.Use(loggingMiddleware)

	// registers web handlers
	log.Trace().Msg("registering web root / and liveliness probe /live handlers")
	pattern := fmt.Sprintf("/%s/{key}", s.conf.Path)
	mux.HandleFunc(pattern, s.rootHandler)
	mux.HandleFunc("/live", s.liveHandler)

	log.Trace().Msg("registering readiness probe handler /ready")
	mux.HandleFunc("/ready", s.readyHandler)

	if s.conf.Metrics {
		// prometheus metrics
		log.Trace().Msg("metrics is enabled, registering handler for endpoint /metrics.")
		mux.Handle("/metrics", promhttp.Handler())
	}

	// creates an http server listening on the specified TCP port
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", s.conf.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  time.Second * 60,
		Handler:      mux,
	}

	// runs the server asynchronously
	go func() {
		log.Trace().Msgf("terra listening on :%s", s.conf.Port)
		if err := server.ListenAndServe(); err != nil {
			log.Err(err).Msg("Failed to start server.")
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
		log.Info().Msg("shutting down Terra")
		log.Fatal().Err(err)
	}
}

func (s *TerraService) rootHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	defer r.Body.Close()

	// continues only if the request is authenticated
	if !s.authenticate(w, r) {
		return
	}

	switch r.Method {
	case "GET":
		{
			// constructs a default terraform state object
			state := TfState{Version: 1}
			// attempts to load the state from Onix
			err := state.loadState(s.ox, vars["key"])
			if err != nil {
				if !strings.Contains(err.Error(), "404") {
					// only logs the error if it is anything other than 404 (Not Found)
					log.Error().Msg(err.Error())
				}
			}
			// writes the state to the response
			io.WriteString(w, state.toJSONString())
		}

	case "POST":
		{
			// reads the request body
			state, err := s.readRequestBody(r)
			if err != nil {
				log.Err(err).Msg("failed to read request payload")
				return
			}
			// persists the state in Onix
			err = state.save(s.ox, vars["key"])
			if err != nil {
				log.Error().Msg(err.Error())
			}
		}

	case "PUT":
		{
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte("405 - Method Not Allowed"))
		}
	case "DELETE":
		{
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte("405 - Method Not Allowed"))
		}
	}
}

func (s *TerraService) readyHandler(w http.ResponseWriter, r *http.Request) {
	if !s.ready {
		log.Warn().Msg("Terraform HTTP Backend service is not ready")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Terraform HTTP Backend service is not ready"))
		if err != nil {
			log.Error().Err(err)
		}
	} else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}
}

func (s *TerraService) liveHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

// unmarshal the http request into a TfState structure
func (s *TerraService) readRequestBody(r *http.Request) (*TfState, error) {
	var state TfState
	// read the request body into a byte array
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	// unmarshal the byte array into a TfState object
	err = json.Unmarshal(bytes, &state)
	if err != nil {
		return nil, err
	}
	// return the terraform state
	return &state, nil
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Trace().Msgf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func (s *TerraService) authenticate(w http.ResponseWriter, r *http.Request) bool {
	// if there is a username and password
	if len(s.conf.AuthMode) > 0 && strings.ToLower(s.conf.AuthMode) == "basic" {
		if r.Header.Get("Authorization") == "" {
			// if no authorisation header is passed, then it prompts a client browser to authenticate
			w.Header().Set("WWW-Authenticate", `Basic realm="oxterra"`)
			w.WriteHeader(http.StatusUnauthorized)
			log.Trace().Msg("Unauthorised request.")
			return false
		} else {
			// authenticate the request
			requiredToken := s.newBasicToken(s.conf.Username, s.conf.Password)
			providedToken := r.Header.Get("Authorization")
			// if the authentication fails
			if !strings.Contains(providedToken, requiredToken) {
				// returns an unauthorised request
				w.WriteHeader(http.StatusForbidden)
				return false
			}
		}
	}
	return true
}

// creates a new Basic Authentication Token
func (s *TerraService) newBasicToken(user string, pwd string) string {
	return fmt.Sprintf("Basic %s",
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, pwd))))
}
