/*
   Onix Config Manager - SeS - Onix Webhook Receiver for AlertManager
   Copyright (c) 2020 by www.gatblau.org

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
package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gatblau/oxc"
	"github.com/gorilla/mux"
	"github.com/prometheus/alertmanager/template"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"time"
)

// the Service Status web server
type SeS struct {
	ox    *oxc.Client // client to connect to Onix WAPI
	conf  *Config
	ready bool
}

// creates a new instance of the Service Status
func NewSeS() (*SeS, error) {
	// load the configuration
	conf, err := NewConfig()
	if err != nil {
		return nil, err
	}
	// creates an Onix client
	client, err := oxc.NewClient(conf.Ox)
	if err != nil {
		return nil, err
	}
	// create an instance of the service
	ses := &SeS{
		ox:   client,
		conf: conf,
	}
	// launch a go routine to check for and create the meta-model
	go ses.checkModel()
	// returns the service instance
	return ses, nil
}

// check the metamodel exist and creates it if not
func (s *SeS) checkModel() {
	// checks if a meta model for Terraform is defined in Onix
	model := NewSeSModel(s.ox)
	err := model.create()
	if err != nil {
		log.Error().Msgf("cannot create SeS model in Onix: %s", err)
	} else {
		// if no error then set the ready state to true
		s.ready = true
	}
}

// starts the http service
func (s *SeS) Start() {
	// gets a new router
	mux := mux.NewRouter()
	// logs incoming calls
	mux.Use(loggingMiddleware)
	// registers web handlers
	mux.HandleFunc("/live", s.liveHandler).Methods("GET")
	mux.HandleFunc("/ready", s.readyHandler).Methods("GET")
	mux.HandleFunc(fmt.Sprintf("/%s", s.conf.Path), s.svcHandler).Methods("POST")
	if s.conf.Metrics {
		// prometheus metrics
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
		log.Trace().Msgf("SeS listening on :%s", s.conf.Port)
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
		log.Fatal().Err(err)
	}

	log.Info().Msg("shutting down SeS")
}

// liveliness probe handler
func (s *SeS) liveHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

// readyness probe handler
func (s *SeS) readyHandler(w http.ResponseWriter, r *http.Request) {
	if !s.ready {
		s.notReady(w)
	} else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}
}

// responds with a not ready error message
func (s *SeS) notReady(w http.ResponseWriter) {
	notReadyMsg := "Service Status WebHook Receiver is not ready"
	log.Warn().Msg(notReadyMsg)
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte(notReadyMsg))
	if err != nil {
		log.Error().Err(err)
	}
}

// main service handler
func (s *SeS) svcHandler(w http.ResponseWriter, r *http.Request) {
	// if not ready then return an error
	if !s.ready {
		s.notReady(w)
	}

	// continues only if the request is authenticated
	if !s.authenticate(w, r) {
		return
	}

	var payload template.Data
	// attempts to de-serialise the payload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Error().Msgf("cannot read the alerts in the payload: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// sort the incoming alerts by StartsAt time
	alerts := NewTimeSortedAlerts(payload.Alerts)
	sort.Sort(alerts)

	for _, alert := range alerts {
		ok, service := s.contain(alert.Annotations, "service")
		if !ok {
			http.Error(w, fmt.Sprintf("cannot find 'service' annotation in alert '%s'", alert), http.StatusBadRequest)
			return
		}
		ok, status := s.contain(alert.Annotations, "status")
		if !ok {
			http.Error(w, fmt.Sprintf("cannot find 'status' annotation in alert '%s'", alert), http.StatusBadRequest)
			return
		}
		ok, uri := s.contain(alert.Annotations, "uri")
		if !ok {
			http.Error(w, fmt.Sprintf("cannot find 'uri' annotation in alert '%s'", alert), http.StatusBadRequest)
			return
		}
		ok, description := s.contain(alert.Annotations, "description")
		if !ok {
			http.Error(w, fmt.Sprintf("cannot find 'description' annotation in alert '%s'", alert), http.StatusBadRequest)
			return
		}
		log.Info().Msgf("service: %s:%s is %s", service, uri, status)
		result, err := s.ox.PutItem(&oxc.Item{
			Key:         fmt.Sprintf("%s_%s", service, strings.Replace(strings.Replace(uri, ":", "_", -1), ".", "_", -1)),
			Name:        fmt.Sprintf("%s Service", service),
			Description: description,
			Type:        SeSServiceItemType,
			Attribute: map[string]interface{}{
				"name":        service,
				"status":      status,
				"description": description,
				"uri":         uri,
			},
		})
		if err != nil {
			log.Error().Msgf("cannot update service status due to a technical issue: %s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if result.Error {
			log.Error().Msgf("cannot update service status: %s", err)
			http.Error(w, result.Message, http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

// logs incoming http requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Trace().Msgf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// authenticates an incoming request
func (s *SeS) authenticate(w http.ResponseWriter, r *http.Request) bool {
	// if there is a username and password
	if len(s.conf.AuthMode) > 0 && strings.ToLower(s.conf.AuthMode) == "basic" {
		if r.Header.Get("Authorization") == "" {
			// if no authorisation header is passed, then it prompts a client browser to authenticate
			w.Header().Set("WWW-Authenticate", `Basic realm="OxSeS"`)
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
func (s *SeS) newBasicToken(user string, pwd string) string {
	return fmt.Sprintf("Basic %s",
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, pwd))))
}

// check if the passed-in map contains the specified key and returns its value
func (s *SeS) contain(values template.KV, key string) (bool, string) {
	for k, v := range values {
		if strings.ToLower(key) == strings.ToLower(k) {
			return true, v
		}
	}
	return false, ""
}

// process the received
func (s *SeS) processAlerts(data template.Alerts) error {
	// sort the incoming alerts by StartsAt time
	alerts := NewTimeSortedAlerts(data)
	sort.Sort(alerts)

	for _, alert := range alerts {
		ok, service := s.contain(alert.Annotations, "service")
		if !ok {
			return errors.New(fmt.Sprintf("cannot find 'service' annotation in alert '%s'", alert))
		}
		ok, status := s.contain(alert.Annotations, "status")
		if !ok {
			return errors.New(fmt.Sprintf("cannot find 'status' annotation in alert '%s'", alert))
		}
		ok, uri := s.contain(alert.Annotations, "uri")
		if !ok {
			return errors.New(fmt.Sprintf("cannot find 'uri' annotation in alert '%s'", alert))
		}
		ok, description := s.contain(alert.Annotations, "description")
		if !ok {
			return errors.New(fmt.Sprintf("cannot find 'description' annotation in alert '%s'", alert))
		}
		log.Info().Msgf("service: %s:%s is %s", service, uri, status)
		result, err := s.ox.PutItem(&oxc.Item{
			Key:         fmt.Sprintf("%s_%s", service, strings.Replace(strings.Replace(uri, ":", "_", -1), ".", "_", -1)),
			Name:        fmt.Sprintf("%s Service", service),
			Description: description,
			Type:        SeSServiceItemType,
			Attribute: map[string]interface{}{
				"name":        service,
				"status":      status,
				"description": description,
				"uri":         uri,
				"time":        alert.StartsAt.String(),
			},
		})
		if err != nil {
			return err
		}
		if result.Error {
			return errors.New(fmt.Sprintf("cannot update service status: %s", result.Message))
		}
	}
	return nil
}
