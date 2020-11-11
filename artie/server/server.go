/*
  Onix Config Manager - Artie
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package server

// @title Artie: the generic artefact manager API
// @version 0.0.4
// @description Artie's HTTP API for artefact backends.
// @contact.name gatblau
// @contact.url http://onix.gatblau.org/
// @contact.email onix@gatblau.org
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

import (
	"context"
	"fmt"
	"github.com/gatblau/onix/artie/core"
	_ "github.com/gatblau/onix/artie/docs" // documentation needed for swagger
	"github.com/gatblau/onix/artie/registry"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggo/http-swagger" // http-swagger middleware
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	lock  *lock
	conf  *core.ServerConfig
	start time.Time
}

func NewServer() *Server {
	return &Server{
		// the server configuration
		conf: new(core.ServerConfig),
		// a read/write lock
		lock: new(lock),
	}
}

func (s *Server) Serve() {
	// compute the time the server is called
	s.start = time.Now()
	// ensure the locks path is created
	s.lock.ensurePath()

	router := mux.NewRouter()
	router.Use(s.loggingMiddleware)
	router.Use(s.authenticationMiddleware)

	// registers web handlers
	fmt.Printf("? I am registering http handlers\n")
	router.HandleFunc("/", s.liveHandler).Methods("GET")
	// router.HandleFunc("/ready", s.readyHandler).Methods("GET")

	// swagger configuration
	if s.conf.SwaggerEnabled() {
		fmt.Printf("? Open API available at /api\n")
		router.PathPrefix("/api").Handler(httpSwagger.WrapHandler)
	}

	if s.conf.MetricsEnabled() {
		// prometheus metrics
		fmt.Printf("? /metrics endpoint is enabled\n")
		router.Handle("/metrics", promhttp.Handler()).Methods("GET")
	}

	// push artefacts
	router.HandleFunc("/artefact", s.uploadHandler).Methods("POST")

	fmt.Printf("? using %s backend @ %s\n", s.conf.Backend(), s.conf.BackendDomain())

	// starts the server
	s.listen(router)
}

// @Summary Check that Artie's HTTP API is live
// @Description Checks that Artie's HTTP server is listening on the required port.
// @Description Use a liveliness probe.
// @Description It does not guarantee the server is ready to accept calls.
// @Tags General
// @Produce  plain
// @Success 200 {string} OK
// @Router / [get]
func (s *Server) liveHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("OK"))
	if err != nil {
		fmt.Printf("!!! I cannot write response: %v", err)
	}
}

// @Summary Push an artefact to the configured backend
// @Description uploads the artefact file and its seal to the pre-configured backend (e.g. Nexus, etc)
// @Tags Artefacts
// @Produce  plain
// @Success 204 {string} artefact has been uploaded successfully. the server has nothing to respond.
// @Failure 423 {string} the artefact is locked (pessimistic locking)
// @Router /artefact [post]
// @Param artefact.fileRef formData string true "the artefact file reference"
// @Param artefact.repository formData string true "the artefact repository"
// @Param artefact.file formData file true "the artefact file part of the multipart message"
// @Param seal.file formData file true "the seal file part of the multipart message"
func (s *Server) uploadHandler(w http.ResponseWriter, r *http.Request) {
	// file limit 50 MB
	r.ParseMultipartForm(50 << 20)
	artefactRef := r.FormValue("artefact.fileRef")
	repositoryName := r.FormValue("artefact.repository")
	zipfile, _, err := r.FormFile("artefact.file")
	if err != nil {
		log.Printf("error retrieving artefact file: %s", err)
		s.writeError(w, err)
	}
	jsonFile, _, err := r.FormFile("seal.file")
	if err != nil {
		log.Printf("error retrieving seal file: %s", err)
		s.writeError(w, err)
		return
	}
	// try and upload checking the resource is not locked
	isLocked := s.upload(w, repositoryName, artefactRef, zipfile, jsonFile)
	// if the resource was locked
	if isLocked {
		// try again
		isLocked = s.upload(w, repositoryName, artefactRef, zipfile, jsonFile)
		// if locked again error
		w.WriteHeader(http.StatusLocked)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// // @Summary Check that Artie's HTTP API is Ready
// // @Description Checks that the HTTP API is ready to accept calls
// // @Tags General
// // @Produce  plain
// // @Success 200 {string} OK
// // @Failure 500 {string} error message
// // @Router /ready [get]
// func (s *Server) readyHandler(w http.ResponseWriter, r *http.Request) {
// 	ready, err := checkReady()
// 	if !ready {
// 		fmt.Printf("! I am not ready: %v", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		_, err := w.Write([]byte(err.Error()))
// 		if err != nil {
// 		}
// 	} else {
// 		_, _ = w.Write([]byte("OK"))
// 	}
// }

func (s *Server) listen(handler http.Handler) {
	// creates an http server listening on the specified TCP port
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", s.conf.HttpPort()),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  time.Second * 60,
		Handler:      handler,
	}

	// creates a channel to pass a SIGINT (ctrl+C) kernel signal with buffer capacity 1
	stop := make(chan os.Signal, 1)

	// runs the server asynchronously
	go func() {
		fmt.Printf("? I am listening on :%s\n", s.conf.HttpPort())
		fmt.Printf("? I have taken %v to start\n", time.Since(s.start))
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("! Stopping the server: %v\n", err)
			os.Exit(1)
		}
	}()

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
		fmt.Printf("? I am shutting down due to an error: %v\n", err)
	}
}

func (s *Server) writeError(w http.ResponseWriter, err error) {
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		fmt.Printf("!!! I failed to write error to response: %v\n", err)
	}
}

// log http requests to stdout
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("? I received an http request from: %s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// determines if the request is authenticated
func (s *Server) authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			// if no authorisation header is passed, then it prompts a client browser to authenticate
			w.Header().Set("WWW-Authenticate", `Basic realm="onix/artie"`)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Printf("? I have received an (unauthorised) http request from: '%v'\n", r.RemoteAddr)
		} else {
			// authenticate the request
			requiredToken := s.conf.BasicToken()
			providedToken := r.Header.Get("Authorization")
			// if the authentication fails
			if providedToken != requiredToken {
				// Write an error and stop the handler chain
				http.Error(w, "Forbidden", http.StatusForbidden)
			}
		}
		// Pass down the request to the next middleware (or final handler)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) upload(w http.ResponseWriter, repositoryName string, artefactRef string, zipfile multipart.File, jsonFile multipart.File) (isLocked bool) {
	locked, err := s.lock.acquire(repositoryName)
	if err != nil {
		if err != nil {
			log.Printf("error trying to release lock: %s", err)
			s.writeError(w, err)
			return true
		}
	}
	if locked > 0 {
		artieName := core.ParseName(repositoryName)
		backend := registry.NewBackendFactory().Get()
		err := backend.UploadArtefact(artieName, artefactRef, zipfile, jsonFile, s.conf.HttpUser(), s.conf.HttpPwd())
		s.lock.release(repositoryName)
		if err != nil {
			log.Printf("error whilst pushing to %s backend: %s", s.conf.Backend(), err)
			s.writeError(w, err)
		}
		return false
	} else {
		err := s.lock.tryRelease(repositoryName, 15)
		if err != nil {
			log.Printf("error trying to release lock: %s", err)
			s.writeError(w, err)
			return true
		}
	}
	// if we are here it is because the resource was locked and a time out occurred
	return true
}
