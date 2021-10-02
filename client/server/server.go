package server

/*
  Onix Config Manager - Http Client
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
	"github.com/gatblau/onix/client"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
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
	Auth map[string]func(http.Request) *client.UserPrincipal
	// default authentication function
	DefaultAuth func(http.Request) *client.UserPrincipal
}

func New(realm string) *Server {
	conf := new(ServerConfig)
	return &Server{
		// the server configuration
		Conf:  conf,
		realm: realm,
		// defines a default authentication function using Basic Authentication
		// can be overridden to change the behaviour or made nil to have an unauthenticated service or endpoint
		DefaultAuth: func(r http.Request) *client.UserPrincipal {
			requestToken := r.Header.Get("Authorization")
			// authenticates if the http request token matches the configured basic authentication token
			if requestToken == conf.BasicToken() {
				user, _ := ParseBasicToken(r)
				// return the user principal
				return &client.UserPrincipal{
					Username: user,
					Rights:   nil,
					Created:  time.Now(),
				}
			}
			// otherwise, return nil meaning that authentication has failed
			return nil
		},
	}
}

// Serve starts the server
func (s *Server) Serve() {
	// compute the time the server is called
	s.start = time.Now()

	router := mux.NewRouter()

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

// @Summary Check that the HTTP API is live
// @Description Checks that the HTTP server is listening on the required port.
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

// LoggingMiddleware log http requests to stdout
func (s *Server) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path, _ := url.PathUnescape(r.URL.Path)
		fmt.Printf("request from: %s %s %s\n", r.RemoteAddr, r.Method, path)
		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Println(err)
		}
		log.Println(string(requestDump))
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// AuthenticationMiddleware determines if the request is authenticated
func (s *Server) AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// holds user principal
		var user *client.UserPrincipal
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
				user = authenticate(*r)
				// if authentication fails the there is no user principal returned
				if user == nil {
					// Write an error and stop the handler chain
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				} else {
					// exit loop
					break
				}
			}
		}
		// if not authenticated by a custom handler then use default handler
		if user == nil {
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
					user = s.DefaultAuth(*r)
					if user == nil {
						// if the authentication failed, write an error and stop the handler chain
						http.Error(w, "Forbidden", http.StatusForbidden)
						return
					}
				}
			}
		}
		// create a user context containing the user principal
		userContext := context.WithValue(r.Context(), "User", user)
		// create a shallow copy of the request with the user context added to it
		req := r.WithContext(userContext)
		// pass down the request to the next middleware (or final handler)
		next.ServeHTTP(w, req)
	})
}

// AuthorisationMiddleware authorises the http request based on the rights in user principal in the request context
func (s *Server) AuthorisationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUserPrincipal(r)
		// if no principal is found reject the request
		if user == nil || !user.Rights.RequestAllowed(s.realm, r) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
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
