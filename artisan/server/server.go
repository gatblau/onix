package server

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"time"
)

// Server generic http server
type Server struct {
	// the start time of the server
	start time.Time
	// the server configuration
	Conf *ServerConfig
	// basic auth realm
	realm string
	// http function to register http handlers
	Http func(*mux.Router)
	// jobs function to register async jobs
	Jobs func() error
	// map of authentication handlers
	Auth map[string]func(string) bool
	// default authentication function
	DefaultAuth func(string) bool
}

func New(realm string) *Server {
	conf := new(ServerConfig)
	return &Server{
		// the server configuration
		Conf:  conf,
		realm: realm,
		// defines a default authentication function using Basic Authentication
		// can be overridden to change the behaviour or made nil to have an unauthenticated service or endpoint
		DefaultAuth: func(requestToken string) bool {
			// authenticates if the http request token matches the configured basic authentication token
			return requestToken == conf.BasicToken()
		},
	}
}

// Serve starts the server
func (s *Server) Serve() {
	// compute the time the server is called
	s.start = time.Now()

	router := mux.NewRouter()
	router.Use(s.loggingMiddleware)
	router.Use(s.authenticationMiddleware)

	// registers web handlers
	router.HandleFunc("/", s.liveHandler).Methods("GET")

	// swagger configuration
	if s.Conf.SwaggerEnabled() {
		fmt.Printf("? OpenAPI available at /api\n")
		router.PathPrefix("/api").Handler(httpSwagger.WrapHandler)
	}

	// Prometheus endpoint
	if s.Conf.MetricsEnabled() {
		// prometheus metrics
		fmt.Printf("? /metrics endpoint is enabled\n")
		router.Handle("/metrics", promhttp.Handler()).Methods("GET")
	}

	// add the http handlers to the router if a registering function has been dclared
	if s.Http != nil {
		s.Http(router)
	} else {
		// warn that no handler has been provided
		log.Printf("WARNING: no http handler has been registered, no application specific endpoints will be available\n" +
			"have you forgotten to set the server Http function?\n")
	}

	// run jobs if there are any
	if s.Jobs != nil {
		err := s.Jobs()
		if err != nil {
			log.Printf(err.Error())
		}
	}

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
func (s *Server) liveHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("OK"))
	if err != nil {
		fmt.Printf("error: cannot write response: %v", err)
	}
}

func (s *Server) listen(handler http.Handler) {
	// creates an http server listening on the specified TCP port
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", s.Conf.HttpPort()),
		WriteTimeout: 180 * time.Second,
		ReadTimeout:  180 * time.Second,
		IdleTimeout:  time.Second * 180,
		Handler:      handler,
	}

	// creates a channel to pass a SIGINT (ctrl+C) kernel signal with buffer capacity 1
	stop := make(chan os.Signal, 1)

	// runs the server asynchronously
	go func() {
		fmt.Printf("server listening on :%s\n", s.Conf.HttpPort())
		fmt.Printf("server started in %v\n", time.Since(s.start))
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("server stopping: %v\n", err)
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
		// get the authentication token in the request
		providedToken := r.Header.Get("Authorization")
		// not authenticated yet
		authenticated := false
		// loop through specific authentication by URL path
		for urlPattern, authenticate := range s.Auth {
			// if the request URL match the authentication function pattern
			matched, err := regexp.Match(urlPattern, []byte(r.URL.Path))
			// regex error?
			if err != nil {
				// Write an error and stop the handler chain
				log.Printf("authentication function error: %s\n", err)
				http.Error(w, "Authentication Error", http.StatusInternalServerError)
			}
			// if the regex matched the URL path
			if matched {
				// then try and authenticate using the specified function
				// if authentication fails
				if !authenticate(providedToken) {
					// Write an error and stop the handler chain
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				} else {
					// authentication succeeded, set flag to true
					authenticated = true
					// exit loop
					break
				}
			}
		}
		// if not authenticated by a custom handler then use default handler
		if !authenticated {
			// no specific authentication function matched the request URL, so tries
			// the default authentication function if it has been defined
			// if no function has been defined then do not authenticate the request
			if s.DefaultAuth != nil {
				// if no Authorization header is found
				if r.Header.Get("Authorization") == "" {
					// prompts a client to authenticate by setting WWW-Authenticate response header
					w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, s.realm))
					w.WriteHeader(http.StatusUnauthorized)
					fmt.Printf("! unauthorised http request from: '%v'\n", r.RemoteAddr)
				} else {
					// authenticate the request using the default handler
					if !s.DefaultAuth(providedToken) {
						// if the authentication failed, write an error and stop the handler chain
						http.Error(w, "Forbidden", http.StatusForbidden)
						return
					}
				}
			}
		}
		// Pass down the request to the next middleware (or final handler)
		next.ServeHTTP(w, r)
	})
}

// writes the content of an object using the response writer in the format specified by the accept http header
// supporting content negotiation for json, yaml, and xml formats
func (s *Server) Write(w http.ResponseWriter, r *http.Request, obj interface{}) {
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

func (s *Server) writeError(w http.ResponseWriter, err error, errorCode int) {
	fmt.Printf(fmt.Sprintf("%s\n", err))
	w.WriteHeader(errorCode)
}
