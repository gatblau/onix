/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package server

// @title Artisan Package Registry
// @version 0.0.4
// @description Registry for Artisan packages
// @contact.name gatblau
// @contact.url http://onix.gatblau.org/
// @contact.email onix@gatblau.org
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gatblau/onix/artisan/artreg/backend"
	_ "github.com/gatblau/onix/artisan/artreg/docs"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/registry"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggo/http-swagger" // http-swagger middleware
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type Server struct {
	lock  *lock
	conf  *ServerConfig
	start time.Time
}

func NewServer() *Server {
	return &Server{
		// the server configuration
		conf: new(ServerConfig),
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
		fmt.Printf("? Download API available at /api\n")
		router.PathPrefix("/api").Handler(httpSwagger.WrapHandler)
	}

	if s.conf.MetricsEnabled() {
		// prometheus metrics
		fmt.Printf("? /metrics endpoint is enabled\n")
		router.Handle("/metrics", promhttp.Handler()).Methods("GET")
	}

	// push an package using a tag
	router.HandleFunc("/package/{repository-group}/{repository-name}/tag/{package-tag}", s.packageUploadHandler).Methods("POST")

	// update package information by id
	router.HandleFunc("/package/{repository-group}/{repository-name}/id/{package-id}", s.packageInfoUpdateHandler).Methods("PUT")
	router.HandleFunc("/package/{repository-group}/{repository-name}/id/{package-id}", s.packageInfoGetHandler).Methods("GET")

	// get repository information
	router.HandleFunc("/repository/{repository-group}/{repository-name}", s.repositoryInfoHandler).Methods("GET")

	// files download
	router.HandleFunc("/file/{repository-group}/{repository-name}/{filename}", s.fileDownloadHandler).Methods("GET")

	// create a webhook
	router.HandleFunc("/webhook/{repository-group}/{repository-name}", s.webhookCreateHandler).Methods("POST")
	// delete a webhook
	router.HandleFunc("/webhook/{repository-group}/{repository-name}/{webhook-id}", s.webhookDeleteHandler).Methods("DELETE")
	// retrieve webhooks
	router.HandleFunc("/webhook/{repository-group}/{repository-name}", s.webhookGetHandler).Methods("GET")

	fmt.Printf("? backend => %s\n", GetBackend().Name())

	// starts the server
	s.listen(router)
}

// @Summary Check that the registry HTTP API is live
// @Description Checks that the registry HTTP server is listening on the required port.
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

// @Summary Download a file from the registry
// @Description
// @Tags Files
// @Produce octet-stream
// @Router /file/{repository-group}/{repository-name}/{filename} [get]
// @Param repository-group path string true "the package repository group name"
// @Param repository-name path string true "the package repository name"
// @Param filename path string true "the filename to download"
// @Success 200 {file} package has been downloaded successfully
// @Failure 500 {string} internal server error
func (s *Server) fileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	// get request variables
	vars := mux.Vars(r)
	group := vars["repository-group"]
	name := vars["repository-name"]
	filename := vars["filename"]

	// get the backend to use
	back := GetBackend()

	file, _ := back.Download(group, name, filename, s.conf.HttpUser(), s.conf.HttpPwd())
	defer file.Close()

	fileHeader := make([]byte, 512)
	file.Read(fileHeader)

	fileStat, _ := file.Stat()

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))

	_, err := file.Seek(0, 0)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	_, err = io.Copy(w, file)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	return
}

// @Summary Push an package to the configured backend
// @Description uploads the package file and its seal to the pre-configured backend (e.g. Nexus, etc)
// @Tags Packages
// @Produce  plain
// @Success 204 {string} package has been uploaded successfully. the server has nothing to respond.
// @Failure 423 {string} the package is locked (pessimistic locking)
// @Router /package/{repository-group}/{repository-name}/tag/{package-tag} [post]
// @Param repository-group path string true "the package repository group name"
// @Param repository-name path string true "the package repository name"
// @Param tag path string true "the package reference name"
// @Param package-meta formData string true "the package metadata in JSON base64 encoded string format"
// @Param package-file formData file true "the package file part of the multipart message"
// @Param package-seal formData file true "the seal file part of the multipart message"
func (s *Server) packageUploadHandler(w http.ResponseWriter, r *http.Request) {
	// get request variables
	vars := mux.Vars(r)
	repoGroup := vars["repository-group"]
	repoName := vars["repository-name"]
	packageTag := vars["package-tag"]
	repoGroup, _ = url.PathUnescape(repoGroup)
	name := &core.PackageName{
		Group: repoGroup,
		Name:  repoName,
		Tag:   packageTag,
	}
	// file upload limit in MB
	err := r.ParseMultipartForm(s.conf.HttpUploadLimit() << 20)
	if err != nil {
		s.writeError(w, fmt.Errorf("error parsing multipart form: %s", err), http.StatusBadRequest)
		return
	}
	info := r.FormValue("package-meta")
	meta, err := base64.StdEncoding.DecodeString(info)
	if err != nil {
		core.CheckErr(err, "failed to base64 decode package information")
	}
	jsonFile, _, err := r.FormFile("package-seal")
	if err != nil {
		s.writeError(w, fmt.Errorf("error retrieving seal file: %s", err), http.StatusBadRequest)
		return
	}
	zipFile, _, err := r.FormFile("package-file")
	if err != nil {
		s.writeError(w, fmt.Errorf("error retrieving package zip file: %s", err), http.StatusBadRequest)
		return
	}
	// convert the meta file into an package
	packageMeta := new(registry.Package)
	err = json.Unmarshal(meta, packageMeta)
	if err != nil {
		s.writeError(w, fmt.Errorf("cannot unmashall package metadata: %s", err), http.StatusBadRequest)
		return
	}
	// try and upload checking the resource is not locked
	repoPath := fmt.Sprintf("%s/%s", repoGroup, repoName)
	// get the backend to use
	back := GetBackend()
	// retrieve the repository meta data
	repo, err := back.GetRepositoryInfo(repoGroup, repoName, s.conf.HttpUser(), s.conf.HttpPwd())
	if err != nil {
		s.writeError(w, fmt.Errorf("cannot retrieve repository information form the backend: %s", err), http.StatusInternalServerError)
		return
	}
	// try to find the package being pushed in the remote backend
	backendPackage := repo.FindPackage(packageMeta.Id)
	// if the package exists
	if backendPackage != nil {
		// check the tag does not exist
		if backendPackage.HasTag(packageTag) {
			// package already exist
			fmt.Printf("package already exist, nothing to push")
			// returns ok but not created to indicate there is nothing to do
			w.WriteHeader(http.StatusOK)
			return
		}
		// if the tag does not exist then add the tag to the backend package
		backendPackage.Tags = append(backendPackage.Tags, packageTag)
		// update the package information
		if !repo.UpdatePackage(backendPackage) {
			s.writeError(w, fmt.Errorf("cannot update repository information: %s", backendPackage.Id), http.StatusInternalServerError)
			return
		}
		err = back.UpdatePackageInfo(name, backendPackage, s.conf.HttpUser(), s.conf.HttpPwd())
		if err != nil {
			s.writeError(w, fmt.Errorf("cannot update package information in Nexus backend: %s", err), http.StatusInternalServerError)
			return
		}
		// returns a 201 to indicate the metadata (tag) was added
		w.WriteHeader(http.StatusCreated)
		return
	}
	// if the package does not exist
	// add it to the repository
	repo.Packages = append(repo.Packages, packageMeta)
	// create a repository file
	repoFile, err := core.ToJsonFile(repo)
	if err != nil {
		s.writeError(w, fmt.Errorf("cannot create repository file: %s", err), http.StatusBadRequest)
		return
	}
	// try and acquire a lock
	locked, err := s.lock.acquire(repoPath)
	if err != nil {
		s.writeError(w, fmt.Errorf("cannot acquire lock as it already exists: %s", err), http.StatusBadRequest)
		return
	}
	if locked > 0 {
		err = back.UploadPackage(name, packageMeta.FileRef, zipFile, jsonFile, repoFile, s.conf.HttpUser(), s.conf.HttpPwd())
		_, e := s.lock.release(repoPath)
		if err != nil {
			log.Printf("error whilst pushing to %s backend: %s", s.conf.Backend(), err)
			s.writeError(w, fmt.Errorf("error whilst pushing to %s backend: %s", s.conf.Backend(), err), http.StatusInternalServerError)
			return
		}
		if e != nil {
			s.writeError(w, fmt.Errorf("cannot release lock on repository: %s, %s", repoPath, err), http.StatusInternalServerError)
			return
		}
		// returns a created code to indicate the package was added
		w.WriteHeader(http.StatusCreated)
	} else {
		err := s.lock.tryRelease(repoPath, 15)
		if err != nil {
			s.writeError(w, fmt.Errorf("error trying to release lock: %s", err), http.StatusLocked)
		}
	}
}

// @Summary Get information about the packages in a repository
// @Description gets meta data about packages in the specified repository
// @Tags Repositories
// @Accept text/html, application/json, application/yaml, application/xml, application/xhtml+xml
// @Produce application/json, application/yaml, application/xml
// @Success 200 {string} OK
// @Router /repository/{repository-group}/{repository-name} [get]
// @Param repository-group path string true "the package repository group name"
// @Param repository-name path string true "the package repository name"
func (s *Server) repositoryInfoHandler(w http.ResponseWriter, r *http.Request) {
	// get request variables
	vars := mux.Vars(r)
	repoGroup := vars["repository-group"]
	repoName := vars["repository-name"]
	// retrieve repository metadata from the backend
	repo, err := GetBackend().GetRepositoryInfo(repoGroup, repoName, s.conf.HttpUser(), s.conf.HttpPwd())
	if err != nil {
		s.writeError(w, err, 500)
		return
	}
	s.write(w, r, repo)
}

// @Summary Get information about the specified package
// @Description gets meta data about the package identified by its id
// @Tags Packages
// @Accept text/html, application/json, application/yaml, application/xml, application/xhtml+xml
// @Produce application/json, application/yaml, application/xml
// @Success 200 {string} OK
// @Router /package/{repository-group}/{repository-name}/id/{package-id} [get]
// @Param repository-group path string true "the package repository group name"
// @Param repository-name path string true "the package repository name"
func (s *Server) packageInfoGetHandler(w http.ResponseWriter, r *http.Request) {
	// get request variables
	vars := mux.Vars(r)
	repoGroup := vars["repository-group"]
	repoName := vars["repository-name"]
	id := vars["package-id"]
	// retrieve repository metadata from the backend
	artie, err := GetBackend().GetPackageInfo(repoGroup, repoName, id, s.conf.HttpUser(), s.conf.HttpPwd())
	if err != nil {
		s.writeError(w, err, http.StatusInternalServerError)
		return
	}
	// if the package is not found send a not found error
	if artie == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	s.write(w, r, artie)
}

// @Summary Update information about the specified package
// @Description updates meta data about the package identified by its id
// @Tags Packages
// @Success 200 {string} OK
// @Router /package/{repository-group}/{repository-name}/id/{package-id} [put]
// @Param repository-group path string true "the package repository group name"
// @Param repository-name path string true "the package repository name"
// @Param package-id path string true "the package unique identifier"
// @Param package-info body interface{} true "the package information to be updated"
func (s *Server) packageInfoUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// get request variables
	vars := mux.Vars(r)
	repoGroup := vars["repository-group"]
	repoName := vars["repository-name"]
	name, err := core.ParseName(fmt.Sprintf("%s/%s", repoGroup, repoName))
	if err != nil {
		s.writeError(w, fmt.Errorf("cannot retrieve package information from request body: %s", err), 500)
		return
	}
	id := vars["package-id"]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.writeError(w, fmt.Errorf("cannot retrieve package information from request body: %s", err), 500)
		return
	}
	artie := new(registry.Package)
	err = json.Unmarshal(body, artie)
	if err != nil {
		s.writeError(w, fmt.Errorf("cannot unmarshal package information from request body: %s", err), 500)
		return
	}
	if artie.Id != id {
		s.writeError(w, fmt.Errorf("package Id in URI (%s) does not match the one provided in the payload (%s)", id, artie.Id), 500)
		return
	}
	// updates the repository metadata in Nexus
	if err = GetBackend().UpdatePackageInfo(name, artie, s.conf.HttpUser(), s.conf.HttpPwd()); err != nil {
		s.writeError(w, fmt.Errorf("cannot update repository information in Nexus backend: %s", err), http.StatusInternalServerError)
		return
	}
}

// @Summary creates a webhook configuration
// @Description create the webhook configuration for a specified repository and url
// @Tags Webhooks
// @Accepts json
// @Success 200 {string} returns the new webhook Id
// @Failure 500 {string} internal error
// @Router /webhook/{repository-group}/{repository-name} [post]
// @Param repository-group path string true "the package repository group name"
// @Param repository-name path string true "the package repository name"
// @Param package-info body WebHookConfig true "the webhook configuration"
func (s *Server) webhookCreateHandler(w http.ResponseWriter, r *http.Request) {
	// get request variables
	vars := mux.Vars(r)
	repoGroup := vars["repository-group"]
	repoName := vars["repository-name"]
	// read the payload
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		s.writeError(w, fmt.Errorf("cannot read request payload: %s", err), http.StatusInternalServerError)
		return
	}
	// unmarshal the payload
	config := new(WebHookConfig)
	err = json.Unmarshal(body, &config)
	if err != nil {
		s.writeError(w, fmt.Errorf("cannot unmarshal request payload: %s", err), http.StatusInternalServerError)
		return
	}
	// update the  group/name
	config.Group = repoGroup
	config.Name = repoName
	// load existing configuration
	wh := NewWebHooks()
	err = wh.load()
	if err != nil {
		s.writeError(w, fmt.Errorf("cannot load webhooks configuration: %s", err), http.StatusInternalServerError)
		return
	}
	id, err := wh.Add(config)
	if err != nil {
		s.writeError(w, fmt.Errorf("cannot add webhook configuration: %s", err), 500)
		return
	}
	err = wh.save()
	if err != nil {
		s.writeError(w, fmt.Errorf("cannot save webhook configuration: %s", err), 500)
		return
	}
	// return the id
	s.write(w, r, fmt.Sprintf("{ id:\"%s\" }", id))
}

// @Summary delete a webhook configuration by Id
// @Description delete the specified webhook configuration
// @Tags Webhooks
// @Success 200 {string} successfully deleted
// @Failure 500 {string} internal error
// @Router /webhook/{repository-group}/{repository-name}/{webhook-id} [delete]
// @Param repository-group path string true "the package repository group name"
// @Param repository-name path string true "the package repository name"
// @Param webhook-id path string true "the webhook unique identifier"
func (s *Server) webhookDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// get request variables
	vars := mux.Vars(r)
	repoGroup := vars["repository-group"]
	repoName := vars["repository-name"]
	whId := vars["webhook-id"]
	wh := NewWebHooks()
	err := wh.load()
	if err != nil {
		s.writeError(w, fmt.Errorf("cannot load webhook configuration: %s", err), http.StatusInternalServerError)
		return
	}
	if wh.Remove(repoGroup, repoName, whId) {
		err := wh.save()
		if err != nil {
			s.writeError(w, fmt.Errorf("cannot update webhook configuration: %s", err), http.StatusInternalServerError)
			return
		}
	}
}

// @Summary get a list of webhooks
// @Description get a list of webhook configurations for the specified repository
// @Tags Webhooks
// @Success 200 {string} successfully deleted
// @Failure 500 {string} internal error
// @Router /webhook/{repository-group}/{repository-name} [get]
// @Param repository-group path string true "the package repository group name"
// @Param repository-name path string true "the package repository name"
func (s *Server) webhookGetHandler(w http.ResponseWriter, r *http.Request) {
	// get request variables
	vars := mux.Vars(r)
	repoGroup := vars["repository-group"]
	repoName := vars["repository-name"]
	wh := NewWebHooks()
	err := wh.load()
	if err != nil {
		s.writeError(w, fmt.Errorf("cannot load webhook configuration: %s", err), http.StatusInternalServerError)
		return
	}
	list := wh.GetList(repoGroup, repoName)
	s.write(w, r, list)
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
			w.Header().Set("WWW-Authenticate", `Basic realm="onix/artie"`)
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

func GetBackend() backend.Backend {
	conf := new(ServerConfig)
	// get the configured factory
	switch conf.Backend() {
	case backend.Nexus3:
		return backend.NewNexus3Backend(
			conf.BackendDomain(), // the nexus scheme://domain:port
		)
	case backend.FileSystem:
		return backend.NewFsBackend()
	}
	core.RaiseErr("backend not recognised")
	return nil
}
