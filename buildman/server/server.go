/*
  Onix Config Manager - Build Manager
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package server

// @title Onix - Build Manager
// @version 0.0.4
// @description Build Manager API
// @description build linux container images based on policies
// @contact.name gatblau
// @contact.url http://onix.gatblau.org/
// @contact.email onix@gatblau.org
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	_ "github.com/gatblau/onix/buildman/docs" // documentation needed for swagger
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/reugn/go-quartz/quartz"
	"github.com/swaggo/http-swagger" // http-swagger middleware
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	start time.Time
	conf  *Config
	sched *quartz.StdScheduler
}

func NewServer() *Server {
	return &Server{
		conf:  &Config{},
		sched: quartz.NewStdScheduler(),
	}
}

func (s *Server) Serve() {
	// compute the time the server is called
	s.start = time.Now()

	router := mux.NewRouter()
	router.Use(s.loggingMiddleware)
	router.Use(s.authenticationMiddleware)

	// registers web handlers
	fmt.Printf("? I am registering http handlers\n")
	router.HandleFunc("/", s.liveHandler).Methods("GET")
	// router.HandleFunc("/ready", s.readyHandler).Methods("GET")

	// swagger configuration
	if s.conf.SwaggerEnabled() {
		fmt.Printf("? Download API available at /api\n")
		router.PathPrefix("/api").Handler(httpSwagger.WrapHandler)
	}

	if s.conf.MetricsEnabled() {
		// prometheus metrics
		fmt.Printf("? /metrics endpoint is enabled\n")
		router.Handle("/metrics", promhttp.Handler()).Methods("GET")
	}

	// creates a job to check for changes in the base image
	checkBaseImageJob, err := NewCheckImageJob()
	if err != nil {
		log.Fatalf("cannot create base image check job: %s", err)
	}
	s.sched.Start()

	err = s.sched.ScheduleJob(checkBaseImageJob, quartz.NewSimpleTrigger(time.Minute*time.Duration(checkBaseImageJob.cfg.Interval)))
	if err != nil {
		log.Fatalf("cannot schedule image check job: %s", err)
	}
	// starts the server
	s.listen(router)
}

// @Summary Check that Build Manager HTTP API is live
// @Description Checks that Build Manager HTTP server is listening on the required port.
// @Description Use a liveliness probe.
// @Description It does not guarantee the server is ready to accept calls.
// @Tags General
// @Produce  plain
// @Success 200 {string} OK
// @Router / [get]
func (s *Server) liveHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("OK"))
	if err != nil {
		fmt.Printf("!!! I cannot write response: %v", err)
	}
}

func (s *Server) listen(handler http.Handler) {
	// creates an http server listening on the specified TCP port
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", s.conf.HttpPort()),
		WriteTimeout: 180 * time.Second,
		ReadTimeout:  180 * time.Second,
		IdleTimeout:  time.Second * 180,
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
			// stops the scheduler
			s.sched.Stop()
			// exit
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

func (s *Server) writeError(w http.ResponseWriter, err error, errorCode int) {
	fmt.Printf(fmt.Sprintf("%s\n", err))
	w.WriteHeader(errorCode)
}

// log http requests to stdout
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("request from: %s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		// uncomment below to dump request payload to stdout
		// requestDump, err := httputil.DumpRequest(r, true)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// fmt.Println(string(requestDump))
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// determines if the request is authenticated
func (s *Server) authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			// if no authorisation header is passed, then it prompts a client browser to authenticate
			w.Header().Set("WWW-Authenticate", `Basic realm="onix/buildman"`)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Printf("? unauthorised http request from: '%v'\n", r.RemoteAddr)
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

// writes the content of an object using the response writer in the format specified by the accept http header
// supporting content negotiation for json, yaml, and xml formats
func (s *Server) write(w http.ResponseWriter, r *http.Request, obj interface{}) {
	var (
		bs  []byte
		err error
	)
	// gets the accept http header
	accept := r.Header.Get("Accept")
	switch accept {
	case "*/*":
		fallthrough
	case "application/json":
		fallthrough
	default:
		{
			w.Header().Set("Content-Type", "application/json")
			bs, err = json.Marshal(obj)
		}
	case "application/yaml":
		{
			w.Header().Set("Content-Type", "application/yaml")
			bs, err = yaml.Marshal(obj)
		}
	case "application/xml":
		{
			w.Header().Set("Content-Type", "application/xml")
			bs, err = xml.Marshal(obj)
		}
	}
	if err != nil {
		s.writeError(w, err, 500)
	}
	_, err = w.Write(bs)
	if err != nil {
		log.Printf("error writing data to response: %s", err)
		s.writeError(w, err, 500)
	}
}
