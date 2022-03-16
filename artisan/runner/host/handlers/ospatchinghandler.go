/*
  Onix Config Manager - Artisan Runner
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package handlers

// @title Artisan Host Runner
// @version 0.0.4
// @description Run Artisan flows in host
// @contact.name gatblau
// @contact.url http://onix.gatblau.org/
// @contact.email onix@gatblau.org
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gatblau/onix/artisan/build"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/merge"
	o "github.com/gatblau/onix/artisan/runner/host/onix"
	"github.com/gatblau/onix/artisan/runner/host/parser"
	"github.com/gorilla/mux"
)

type OSpatchingHandler struct {
	// access and secret key for object store
	accessKey  string
	secretKey  string
	flowkey    string
	body       []byte
	cmd        *o.Cmd
	scanedFile string
}

var (
	api *o.API
)

func (h OSpatchingHandler) HandleEvent(w http.ResponseWriter, r *http.Request) {

	err := h.initialize(r)
	if checkErr(w, "initialization of OSpatchingHandler failed :", err) {
		return
	}

	err = h.downloadAndCopyFile()
	if checkErr(w, "error while downloading and copying object from s3 bucket :", err) {
		return
	}

	err = h.executeCommand(w)
	if checkErr(w, "error while executing artisan function in the host machine :", err) {
		return
	}
}

func (h *OSpatchingHandler) initialize(r *http.Request) error {

	vars := mux.Vars(r)
	h.flowkey = vars["flow-key"]
	api = o.Api()
	cmd, err := api.GetCommand(h.flowkey)
	if err != nil {
		msg := fmt.Sprintf("%s: %s\n", "Error while getting command using flow key ", h.flowkey, err)
		fmt.Printf(msg)
		return err
	}
	if cmd == nil {
		msg := fmt.Sprintf("%s: %s\n", "No command item for item type ART_FX found in database for flow key, please check if it was created ", h.flowkey)
		fmt.Printf(msg)
		return errors.New(msg)
	}
	h.cmd = cmd
	path := h.cmd.GetVarValue("XLXS_FILE_PATH_WITH_NAME")
	if len(strings.TrimSpace(path)) == 0 {
		msg := "Value for environment variable XLXS_FILE_PATH_WITH_NAME is missing"
		fmt.Printf(msg)
		return errors.New(msg)
	}
	h.scanedFile = path

	h.accessKey = os.Getenv("OBJECT_STORE_ACCESS_KEY")
	h.secretKey = os.Getenv("OBJECT_STORE_SECRET_KEY")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := fmt.Sprintf("%s: %s\n", "cannot read request payload:", err)
		fmt.Printf(msg)
		return errors.New(msg)
	}
	h.body = body
	return nil
}

func (h OSpatchingHandler) downloadAndCopyFile() error {

	ev, err := parser.NewS3Event(h.body)
	if err != nil {
		msg := fmt.Sprintf("%s: %s\n", "Error while parsing s3 event :", err)
		fmt.Printf(msg)
		return errors.New(msg)
	}

	link, err := ev.GetObjectDownloadURL()
	if err != nil {
		msg := fmt.Sprintf("%s: %s\n", "Error while getting object download url from s3 event :", err)
		fmt.Printf(msg)
		return errors.New(msg)
	}

	cred := fmt.Sprintf("%s:%s", h.accessKey, h.secretKey)
	f, err := core.ReadFile(link, cred)
	if err != nil {
		msg := fmt.Sprintf("%s %s: %s\n", "Error while downloading object from s3 bucket using url", link, err)
		fmt.Printf(msg)
		return errors.New(msg)
	}

	d := filepath.Dir(h.scanedFile)
	err = os.MkdirAll(d, 0755)
	if err != nil {
		return err
	}
	if err != nil {
		msg := fmt.Sprintf("Error while creating folder [%s] to copy file downloaded from object store : %s\n", d, err)
		fmt.Printf(msg)
		return errors.New(msg)
	}

	err = core.WriteFile(f, h.scanedFile, "")
	if err != nil {
		msg := fmt.Sprintf("%s [ %s ] : %s\n", "Error while copying file downloaded from S3 bucket at path ", h.scanedFile, err)
		fmt.Printf(msg)
		return errors.New(msg)
	}

	return nil
}

func (h OSpatchingHandler) executeCommand(w http.ResponseWriter) error {

	// get the variables in the host environment
	hostEnv := merge.NewEnVarFromSlice(os.Environ())
	// get the variables in the command
	cmdEnv := merge.NewEnVarFromSlice(h.cmd.Env())
	// if not containerised add PATH to execution environment
	hostEnv.Merge(cmdEnv)
	cmdEnv = hostEnv
	// if running in verbose mode
	if h.cmd.Verbose {
		// add ARTISAN_DEBUG to execution environment
		cmdEnv.Vars["ARTISAN_DEBUG"] = "true"
	}
	// create the command statement to run
	cmdString := fmt.Sprintf("art %s -u %s:%s %s %s", "exe", h.cmd.User, h.cmd.Pwd, h.cmd.Package, h.cmd.Function)
	// run and return
	out, err := build.ExeAsync(cmdString, ".", cmdEnv, false)
	if err != nil {
		msg := fmt.Sprintf("%s [%s %s ] : %s\n", "Error while executing artisan package function using command", h.cmd.Package, h.cmd.Function, err)
		fmt.Printf(msg)
		return errors.New(msg)
	} else {
		msg := fmt.Sprintf("%s [%s %s ] : [ %s ] \n", "Result of executing artisan package function using command", h.cmd.Package, h.cmd.Function, out)
		fmt.Printf(msg)
	}

	//once patching package is created rename the scan result file with current date time as backup,
	// so we know when this scanned result file was processed and corresponding artisan package was created.
	newFileName := fmt.Sprintf("%s-%s", h.scanedFile, time.Now().Format("2006-01-02-15-04-05"))
	err = os.Rename(h.scanedFile, newFileName)
	if err != nil {
		msg := fmt.Sprintf("After patching package completion error while renaming file [ %s ] to [ %s ] : %s\n", h.scanedFile, newFileName, err)
		fmt.Printf(msg)
		return errors.New(msg)
	}
	return nil
}
