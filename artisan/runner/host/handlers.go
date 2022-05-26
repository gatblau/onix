/*
  Onix Config Manager - Artisan Host Runner
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package main

// @title Artisan Host Runner
// @version 0.0.4
// @description Run Artisan packages with in a host
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
	"path"
	"path/filepath"
	"strings"

	"github.com/gatblau/onix/artisan/build"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/flow"
	"github.com/gatblau/onix/artisan/merge"
	"github.com/gatblau/onix/artisan/runner"
	_ "github.com/gatblau/onix/artisan/runner/host/docs"
	g "github.com/gatblau/onix/artisan/runner/host/git"
	o "github.com/gatblau/onix/artisan/runner/host/onix"
	"github.com/gatblau/onix/oxlib/resx"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	h "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/gorilla/mux"
)

type flowExecution func(path string, f *flow.Flow, w http.ResponseWriter) error

// @Summary Build patching artisan package
// @Description Trigger a new build to create artisan package from the vulnerabilty scanned csv report passed in the payload.
// @Tags Runner
// @Router /host/{cmd-key} [post]
// @Param cmd-key path string true "the unique key of the command to retrieve"
// @Produce plain
// @Failure 500 {string} there was an error in the server, error the server logs
// @Failure 422 {string} command-key was not found in database, error the server logs
// @Success 200 {string} OK
func executeCommandHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	api := o.Api()
	cmdkey := vars["cmd-key"]
	cmd, err := api.GetCommand(cmdkey)
	if checkErr(w, fmt.Sprintf("%s: [ %s ]\n", "Error while getting command using cmd key ", cmdkey), err) {
		return
	}
	if cmd == nil {
		msg := fmt.Sprintf("No command item for item type ART_FX found in database for cmd key [ %s ] , please check if this item exists ", cmdkey)
		fmt.Printf(msg)
		http.Error(w, msg, http.StatusUnprocessableEntity)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if checkErr(w, "Error while reading http request body ", err) {
		return
	}

	t, err := core.NewTempDir()
	if checkErr(w, "Error while creating temp folder ", err) {
		return
	}

	d := path.Join(t, "context")
	err = resx.WriteFile(body, d, "")
	if checkErr(w, fmt.Sprintf("%s: [ %s ]\n", "Error while writing request body to temp path ", d), err) {
		os.RemoveAll(t)
		return
	}

	// get the variables in the host environment
	hostEnv := merge.NewEnVarFromSlice(os.Environ())
	// get the variables in the command
	cmdEnv := merge.NewEnVarFromSlice(cmd.Env())
	// if not containerised add PATH to execution environment
	hostEnv.Merge(cmdEnv)
	cmdEnv = hostEnv
	// if running in verbose mode
	if cmd.Verbose {
		// add ARTISAN_DEBUG to execution environment
		cmdEnv.Vars["ARTISAN_DEBUG"] = "true"
	}

	cmdString := fmt.Sprintf("art %s -u %s:%s %s %s --path=%s", "exe", cmd.User, cmd.Pwd, cmd.Package, cmd.Function, t)
	// run and return
	out, err := build.ExeAsync(cmdString, ".", cmdEnv, false)
	if checkErr(w, fmt.Sprintf("Error while executing artisan package function using command [ %s ]", cmdString), err) {
		os.RemoveAll(t)
		return
	} else {
		msg := fmt.Sprintf("%s [%s %s ] : [ %s ] \n", "Result of executing artisan package function using command", cmd.Package, cmd.Function, out)
		fmt.Printf(msg)
	}

	os.RemoveAll(t)
}

// @Summary Execute an Artisan flow
// @Description Execute a flow from the definition passed in the payload.
// @Tags Runner
// @Router /flow [post]
// @Produce plain
// @Failure 500 {string} there was an error in the server, error the server logs
// @Success 200 {string} OK
func executeFlowFromPayloadHandler(w http.ResponseWriter, r *http.Request) {

	core.Debug("reading request body ....")
	body, err := ioutil.ReadAll(r.Body)
	if checkErr(w, "cannot read request payload ", err) {
		return
	}
	// unmarshal the flow bytes
	core.Debug("creating flow from request ....")
	f, err := flow.NewFlow(body)
	if checkErr(w, "cannot read flow ", err) {
		return
	}

	core.Debug("creating new temp folder ....")
	path, err := core.NewTempDir()
	if checkErr(w, "Error while creating temp folder ", err) {
		return
	}

	if f.RequiresGitSource() {
		fmt.Printf(" Git content %+v\n", f.Git)
		err = gitClone(path, f.Git)
		if checkErr(w, fmt.Sprintf("Error while cloning git uri  [%s]", f.Git.Uri), err) {
			os.RemoveAll(path)
			return
		}
	}
	err = executeRunCFlow(path, f, w)
	if checkErr(w, "error while executing flow spec ", err) {
		os.RemoveAll(path)
		return
	}
	os.RemoveAll(path)
}

// @Summary Retrieve a configured flow from CMDB and execute it.
// @Description Connect to CMDB and retrieves a flow using configuration item natural key passed in flow-key from CMDB
// @Tags Runner
// @Router /webhook/{flow-key}/push [post]
// @Produce plain
// @Param flow-key path string true "the unique key of the flow specification in cmdb"
// @Failure 500 {string} the health check failed with an error, check server logs for details
// @Success 200 {string} OK, the health check succeeded
func executeWebhookFlowHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	api := o.Api()
	flowkey := vars["flow-key"]

	body, err := ioutil.ReadAll(r.Body)
	if checkErr(w, "Error while reading http request body ", err) {
		return
	}
	p, err := g.NewPushEvent(body)
	if checkErr(w, "Error while unmarshalling http request push event body ", err) {
		return
	}

	//get the flow spec json from cmdb using flow-key
	fl, err := api.GetFlow(flowkey)
	if checkErr(w, fmt.Sprintf("%s: [ %s ]\n", "Error while getting flow spec using flow key ", flowkey), err) {
		return
	}

	// build flow using flow spec json
	f, err := flow.NewFlow(fl)
	if checkErr(w, fmt.Sprintf("failed to build flow using flow spec obtained by config item key %s ", flowkey), err) {
		return
	}
	core.Debug("validating flow ....")
	err = f.IsValid()
	if checkErr(w, "Invalid flow ", err) {
		return
	}

	core.Debug("creating new temp folder ....")
	path, err := core.NewTempDir()
	if checkErr(w, "Error while creating temp folder ", err) {
		return
	}

	if f.RequiresGitSource() {
		err = p.IsValidUri(f.Git.Uri)
		if checkErr(w, "Push event git uri validation failed ", err) {
			os.RemoveAll(path)
			return
		}
		err = p.IsValidBranch(f.Git.Branch)
		if checkErr(w, "Push event git branch validation failed ", err) {
			os.RemoveAll(path)
			return
		}

		err = gitClone(path, f.Git)
		if checkErr(w, fmt.Sprintf("Error while cloning git uri  [%s]", f.Git.Uri), err) {
			os.RemoveAll(path)
			return
		}
	}

	if *f.UseRuntimes == true {
		err = executeRunCFlow(path, f, w)
		if checkErr(w, "error while executing flow spec ", err) {
			os.RemoveAll(path)
			return
		}
	} else {
		executeRunFlow(path, f, w)
	}

	os.RemoveAll(path)
}

func executeRunCFlow(path string, f *flow.Flow, w http.ResponseWriter) error {

	var r *runner.Runner
	var err error
	if f.RequiresGitSource() {
		r, err = runner.NewFromPath(path)
		if checkErr(w, fmt.Sprintf("Error while creating runner using current path as [%s] for build file ", path), err) {
			return err
		}
	}

	var env *merge.Envar
	core.Debug("Executing steps ", len(f.Steps))
	if len(f.Steps) > 0 {
		for _, s := range f.Steps {
			i := s.Input
			if i != nil {
				env = i.Env(true)
			}
			// for surce type 'create' delete the folder contents
			if strings.EqualFold(s.PackageSource, "create") {
				err := deleteFolderContents(filepath.Join(path, "*"))
				if checkErr(w, fmt.Sprintf("Error while deleting content of folder path [%s] ", path), err) {
					return err
				}
			}

			//for package source as create/merge, open the package at the give location
			if strings.EqualFold(s.PackageSource, "create") || strings.EqualFold(s.PackageSource, "merge") {
				err := openArtisanPackage(path, s)
				if checkErr(w, fmt.Sprintf("Error while opening artisan package [ %s ]", s.Package), err) {
					return err
				}

				/*build new runner after art open so that new build.yaml is loaded into the runner*/
				r, err = runner.NewFromPath(path)
				if checkErr(w, fmt.Sprintf("Error while creating runner using current path as [%s] for build file ", path), err) {
					return err
				}
			}

			err := r.RunC(s.Function, false, env, "host")
			if checkErr(w, fmt.Sprintf("Error while executing function [%s] using runc command ", s.Function), err) {
				return err
			}
		}
	}
	return nil
}

func executeRunFlow(path string, f *flow.Flow, w http.ResponseWriter) error {

	b := build.NewBuilder()
	var env *merge.Envar
	core.Debug("Executing steps")
	bl := false
	if len(f.Steps) > 0 {
		for _, s := range f.Steps {
			i := s.Input
			if i != nil {
				env = i.Env(true)
			}
			// for surce type 'create' delete the folder contents
			if strings.EqualFold(s.PackageSource, "create") {
				err := deleteFolderContents(filepath.Join(path, "*"))
				if checkErr(w, fmt.Sprintf("Error while deleting content of folder path [%s] ", path), err) {
					return err
				}
			}

			//for package source as create/merge, open the package at the give location
			if strings.EqualFold(s.PackageSource, "create") || strings.EqualFold(s.PackageSource, "merge") {
				err := openArtisanPackage(path, s)
				if checkErr(w, fmt.Sprintf("Error while opening artisan package [ %s ]", s.Package), err) {
					return err
				}
			}

			b.Run(s.Function, path, bl, env)
		}
	}
	return nil
}

func deleteFolderContents(path string) error {
	contents, err := filepath.Glob(path)
	if err != nil {
		return err
	}
	for _, item := range contents {
		err = os.RemoveAll(item)
		if err != nil {
			return err
		}
	}
	return nil
}

func getCredentials(e *merge.Envar) (string, error) {
	usr := e.Vars["ART_REG_USER"]
	pwd := e.Vars["ART_REG_PWD"]
	if len(usr) == 0 && len(pwd) == 0 {
		return "", errors.New("artisan registry credentials missing, credentials must be defined through environment variable ART_REG_USER, ART_REG_PWD ")
	} else if len(usr) == 0 {
		return "", errors.New("artisan registry user missing, user must be defined through environment variable ART_REG_USER ")
	} else if len(pwd) == 0 {
		return "", errors.New("artisan registry password missing, password must be defined through environment variable ART_REG_PWD ")
	} else {
		return fmt.Sprintf("%s:%s", usr, pwd), nil
	}
}

func gitClone(path string, g *flow.Git) error {
	fmt.Printf("git struts content is \n %+v\n", g)
	var opts *git.CloneOptions
	if len(g.Branch) > 0 {
		opts = &git.CloneOptions{
			URL:           g.Uri,
			Progress:      os.Stdout,
			ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", g.Branch)),
		}
	} else {
		opts = &git.CloneOptions{
			URL:      g.Uri,
			Progress: os.Stdout,
		}
	}

	// if authentication token has been provided
	core.Debug("auth credentials is there ", ((len(g.Login) > 0) && len(g.Password) > 0))
	if (len(g.Login) > 0) && len(g.Password) > 0 {
		// The intended use of a GitHub personal access token is in replace of your password
		// because access tokens can easily be revoked.
		// https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
		opts.Auth = &h.BasicAuth{
			Username: g.Login, // yes, this can be anything except an empty string
			Password: g.Password,
		}
	}
	core.Debug("Cloning git repo....")
	_, err := git.PlainClone(path, false, opts)
	if err != nil {
		return err
	}
	core.Debug("Cloning git repo completed at path ", path)

	return err
}

func openArtisanPackage(p string, s *flow.Step) error {
	i := s.Input
	var env *merge.Envar
	if i != nil {
		env = i.Env(true)
	}
	crd, err := getCredentials(env)

	cmdString := fmt.Sprintf("art %s %s -u %s ", "open", s.Package, crd)
	// run and return
	out, err := build.Exe(cmdString, p, env, false)

	if err != nil {
		return err
	}
	cmdStringErr := fmt.Sprintf("art %s %s -u %s ", "open", s.Package, "******:******")
	msg := fmt.Sprintf("opened package using command [%s] at path [%s] with message", cmdStringErr, p, out)
	fmt.Printf(msg)
	return nil
}
