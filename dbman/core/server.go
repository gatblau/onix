//   Onix Config DatabaseProvider - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package core

// @title Onix DbMan
// @version 0.0.4
// @description Call DbMan's commands using HTTP operations from anywhere.
// @contact.name Gatblau
// @contact.url http://onix.gatblau.org/
// @contact.email onix@gatblau.org
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	_ "github.com/gatblau/onix/dbman/docs" // documentation needed for swagger
	"github.com/gatblau/onix/dbman/plugin"
	"github.com/ghodss/yaml"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggo/http-swagger" // http-swagger middleware
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

type Server struct {
	cfg   *Config
	start time.Time
}

func NewServer(cfg *Config) *Server {
	return &Server{
		cfg: cfg,
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
	router.HandleFunc("/ready", s.readyHandler).Methods("GET")
	router.HandleFunc("/conf", s.showConfigHandler).Methods("GET")
	router.HandleFunc("/conf/check", s.checkConfigHandler).Methods("GET")
	router.HandleFunc("/db/info/server", s.dbServerHandler).Methods("GET")
	router.HandleFunc("/db/info/queries", s.queriesHandler).Methods("GET")
	router.HandleFunc("/db/query/{name}", s.queryHandler).Methods("GET")
	router.HandleFunc("/db/create", s.createHandler).Methods("POST")
	router.HandleFunc("/db/deploy", s.deployHandler).Methods("POST")
	router.HandleFunc("/db/upgrade", s.upgradeHandler).Methods("POST")

	// swagger configuration
	router.PathPrefix("/api").Handler(httpSwagger.WrapHandler)

	if s.cfg.GetBool(HttpMetrics) {
		// prometheus metrics
		fmt.Printf("? /metrics endpoint is enabled\n")
		router.Handle("/metrics", promhttp.Handler()).Methods("GET")
	}

	// starts the server
	s.listen(router)
}

// a liveliness probe to prove the http service is listening
func (s *Server) liveHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("OK"))
	if err != nil {
		fmt.Printf("!!! I cannot write response: %v", err)
	}
}

// @Summary Check that DbMan is Ready
// @Description Checks that DbMan is ready to accept calls
// @Tags General
// @Produce  plain
// @Success 200 {string} OK
// @Failure 500 {string} error message
// @Router /ready [get]
func (s *Server) readyHandler(w http.ResponseWriter, r *http.Request) {
	ready, err := DM.CheckReady()
	if !ready {
		fmt.Printf("! I am not ready: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
		}
	} else {
		_, _ = w.Write([]byte("OK"))
	}
}

// @Summary Retrieves database server information
// @Description Gets specific information about the database server to which DbMan is configured to connect.
// @Tags Database
// @Produce  application/json, application/yaml
// @Success 200 {json} database server information
// @Failure 500 {string} error message
// @Router /db/info/server [get]
func (s *Server) dbServerHandler(w http.ResponseWriter, r *http.Request) {
	info, err := DM.GetDbInfo()
	if err != nil {
		s.writeError(w, errors.New(fmt.Sprintf("!!! I cannot get database server information: %v\n", err)))
		return
	}
	s.write(w, r, info)
}

// @Summary Gets a list of available queries.
// @Description Lists all of the queries declared in the current release manifest.
// @Tags Database
// @Produce  application/json, application/yaml
// @Success 200 {string} configuration variables
// @Failure 500 {string} error message
// @Router /db/info/queries [get]
func (s *Server) queriesHandler(writer http.ResponseWriter, request *http.Request) {
	// get the release manifest for the current application version
	_, manifest, err := DM.GetReleaseInfo(DM.Cfg.GetString(AppVersion))
	if err != nil {
		s.writeError(writer, errors.New(fmt.Sprintf("!!! I cannot fetch release information: %v\n", err)))
		return
	}
	s.write(writer, request, manifest.Queries)
}

// @Summary Runs a query.
// @Description Execute a query defined in the release manifest and return the result as a generic serializable table.
// @Tags Database
// @Produce  application/json, application/yaml, application/xml, text/csv
// @Success 200 {Table} a generic table
// @Failure 500 {string} error message
// @Param name path string true "the name of the query as defined in the release manifest"
// @Param params query string false "a string of parameters to be passed to the query in the format 'key1=value1,...,keyN=valueN'"
// @Router /db/query/{name} [get]
func (s *Server) queryHandler(writer http.ResponseWriter, request *http.Request) {
	// get request variables
	vars := mux.Vars(request)
	queryName := vars["name"]
	// if no query name has been specified it cannot continue
	if len(queryName) == 0 {
		s.writeError(writer, errors.New(fmt.Sprintf("!!! I cannot run the query as a query name has not been provided\n")))
		return
	}
	// now check the query has parameters
	queryParams := request.URL.Query()["params"]
	params := make(map[string]string)
	if len(queryParams) > 0 {
		parts := strings.Split(queryParams[0], ",")
		for _, part := range parts {
			subPart := strings.Split(part, "=")
			if len(subPart) != 2 {
				fmt.Printf("!!! I cannot break down query parameter '%s': format should be 'key=value'\n", subPart)
				return
			}
			params[strings.Trim(subPart[0], " ")] = strings.Trim(subPart[1], " ")
		}
	}
	table, _, err := DM.Query(queryName, params)
	if err != nil {
		s.writeError(writer, errors.New(fmt.Sprintf("!!! I cannot execute the query: %v\n", err)))
		return
	}
	s.write(writer, request, *table)
}

// @Summary Creates a new database
// @Description When the database does not already exists, this operation executes the manifest commands required to create the new database.
// @Tags Database
// @Produce  plain
// @Success 200 {string} execution logs
// @Failure 500 {string} error message
// @Router /db/create [post]
func (s *Server) createHandler(w http.ResponseWriter, r *http.Request) {
	// deploy the schema and functions
	output, err, elapsed := DM.Create()
	w.Write([]byte(output.String()))
	// return an error if failed
	if err != nil {
		s.writeError(w, err)
	} else {
		_, err = w.Write([]byte(fmt.Sprintf("? I have completed the action in %v\n", elapsed)))
		if err != nil {
			fmt.Printf("!!! I failed to write error to response: %v", err)
		}
	}
}

// @Summary Deploys the schema and objects in an empty database.
// @Description When the database is empty, this operation executes the manifest commands required to deploy the  database schema and objects.
// @Tags Database
// @Produce  plain
// @Success 200 {string} execution logs
// @Failure 500 {string} error message
// @Router /db/deploy [post]
func (s *Server) deployHandler(w http.ResponseWriter, r *http.Request) {
	// deploy the schema and functions
	output, err, elapsed := DM.Deploy()
	w.Write([]byte(output.String()))
	// return an error if failed
	if err != nil {
		s.writeError(w, err)
	} else {
		_, err = w.Write([]byte(fmt.Sprintf("? I have completed the action in %v\n", elapsed)))
		if err != nil {
			fmt.Printf("!!! I failed to write error to response: %v", err)
		}
	}
}

// @Summary Upgrade a database to a specific version.
// @Description This operation executes the manifest commands required to upgrade an existing database schema and objects to a new version. The target version is defined by DbMan's configuration value "AppVersion". This operation support rolling upgrades.
// @Tags Database
// @Produce  plain
// @Success 200 {string} execution logs
// @Failure 500 {string} error message
// @Router /db/upgrade [post]
func (s *Server) upgradeHandler(w http.ResponseWriter, r *http.Request) {
	// deploy the schema and functions
	output, err, elapsed := DM.Upgrade()
	w.Write([]byte(output.String()))
	// return an error if failed
	if err != nil {
		s.writeError(w, err)
	} else {
		_, err = w.Write([]byte(fmt.Sprintf("? I have completed the action in %v\n", elapsed)))
		if err != nil {
			fmt.Printf("!!! I failed to write error to response: %v", err)
		}
	}
}

// @Summary Validates the current DbMan's configuration.
// @Description Checks that the information in the current configuration set is ok to connect to backend services and the format of manifest is correct.
// @Tags Configuration
// @Produce  plain
// @Success 200 {string} execution logs
// @Failure 500 {string} error message
// @Router /conf/check [get]
func (s *Server) checkConfigHandler(w http.ResponseWriter, r *http.Request) {
	results := DM.CheckConfigSet()
	for check, result := range results {
		_, _ = w.Write([]byte(fmt.Sprintf("[%v] => %v\n", check, result)))
	}
}

// @Summary Show DbMan's current configuration.
// @Description Lists all variables in DbMan's configuration.
// @Tags Configuration
// @Produce  plain
// @Success 200 {string} configuration variables
// @Failure 500 {string} error message
// @Router /conf [get]
func (s *Server) showConfigHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(DM.ConfigSetAsString()))
}

// creates a new Basic Authentication Token
func (s *Server) newBasicToken(user string, pwd string) string {
	return fmt.Sprintf(
		"Basic %s",
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, pwd))))
}

func (s *Server) listen(handler http.Handler) {
	// creates an http server listening on the specified TCP port
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", s.get(HttpPort)),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  time.Second * 60,
		Handler:      handler,
	}

	// creates a channel to pass a SIGINT (ctrl+C) kernel signal with buffer capacity 1
	stop := make(chan os.Signal, 1)

	// runs the server asynchronously
	go func() {
		fmt.Printf("? I am listening on :%s\n", s.get(HttpPort))
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
	w.WriteHeader(http.StatusInternalServerError)
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
		// gets the authentication mode
		authMode := s.get(HttpAuthMode)

		// if there is a username and password
		if len(authMode) > 0 && strings.ToLower(authMode) == "basic" {
			if r.Header.Get("Authorization") == "" {
				// if no authorisation header is passed, then it prompts a client browser to authenticate
				w.Header().Set("WWW-Authenticate", `Basic realm="dbman"`)
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Printf("? I have received an (unauthorised) http request from: '%v'\n", r.RemoteAddr)
			} else {
				// authenticate the request
				requiredToken := s.newBasicToken(s.get(HttpUsername), s.get(HttpPassword))
				providedToken := r.Header.Get("Authorization")
				// if the authentication fails
				if providedToken != requiredToken {
					// Write an error and stop the handler chain
					http.Error(w, "Forbidden", http.StatusForbidden)
				}
			}
		}
		// Pass down the request to the next middleware (or final handler)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) get(key string) string {
	return s.cfg.GetString(key)
}

// writes the content of an object using the response writer in the format specified by the accept http header
// supporting content negotiation for json, yaml, xml and csv formats
func (s *Server) write(w http.ResponseWriter, r *http.Request, obj interface{}) {
	var (
		bytes []byte
		err   error
	)
	// gets the accept http header
	accept := r.Header.Get("Accept")
	switch accept {
	case "application/json":
		{
			w.Header().Set("Content-Type", "application/json")
			bytes, err = json.Marshal(obj)
		}
	case "application/yaml":
		{
			w.Header().Set("Content-Type", "application/yaml")
			bytes, err = yaml.Marshal(obj)
		}
	case "application/xml":
		{
			w.Header().Set("Content-Type", "application/xml")
			bytes, err = xml.Marshal(obj)
		}
	case "text/csv":
		{
			// only support serialization of a Table struct
			if table, ok := obj.(plugin.Table); ok {
				w.Header().Set("Content-Type", "text/csv")
				bytes = []byte(table.AsCSV())
			} else {
				err = errors.New(fmt.Sprintf("!!! I cannot convert to CSV an object that is not a table"))
			}
		}
	default:
		err = errors.New(fmt.Sprintf("!!! I do not support the accept content type '%s'", accept))
	}
	if err != nil {
		s.writeError(w, err)
	}
	w.Write(bytes)
}
