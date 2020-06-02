//   Onix Config Db - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package util

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gatblau/oxc"
	"net/http"
)

// the source of database scripts
type ScriptManager struct {
	client    *oxc.Client
	plan      *Plan
	manifests []Release
	cfg       *AppCfg
}

// factory function
func NewScriptManager(cfg *AppCfg, client *oxc.Client) (*ScriptManager, error) {
	// creates a new struct
	source := new(ScriptManager)
	// setup attributes
	source.cfg = cfg
	source.client = client
	return source, nil
}

// new oxc configuration
func NewOxClientConf(cfg *AppCfg) *oxc.ClientConf {
	return &oxc.ClientConf{
		BaseURI:            cfg.Get(SchemaURI),
		InsecureSkipVerify: false,
		AuthMode:           oxc.None,
	}
}

// get database initialisation information
func (s *ScriptManager) fetchInit() (*DbInit, error) {
	response, err := s.client.Get(fmt.Sprintf("%s/init/init.json", s.get(SchemaURI)), s.addHttpHeaders)
	if err != nil {
		return nil, err
	}
	init := &DbInit{}
	init, err = init.decode(response)
	defer func() {
		if ferr := response.Body.Close(); ferr != nil {
			err = ferr
		}
	}()
	return init, err
}

// access a cached plan reference
func (s *ScriptManager) getPlan() *Plan {
	// if the plan is not fetched
	if s.plan == nil {
		// fetches it
		err := s.loadPlan()
		if err != nil {
			fmt.Sprintf("cannot retrieve plan, %v", err)
		}
	}
	return s.plan
}

// (re)loads the internal plan reference
func (s *ScriptManager) loadPlan() error {
	p, err := s.fetchPlan()
	s.plan = p
	return err
}

// fetches the release plan
func (s *ScriptManager) fetchPlan() (*Plan, error) {
	if s.cfg == nil {
		return nil, errors.New("configuration object not initialised when fetching release getPlan")
	}
	response, err := s.client.Get(fmt.Sprintf("%s/plan.json", s.get(SchemaURI)), s.addHttpHeaders)
	if err != nil {
		return nil, err
	}
	p := &Plan{}
	p, err = p.decode(response)
	defer func() {
		if ferr := response.Body.Close(); ferr != nil {
			err = ferr
		}
	}()
	return p, err
}

// fetches the scripts for a database release
func (s *ScriptManager) fetchRelease(appVersion string) (*Release, error) {
	// if cfg not initialised, no point in continuing
	if s.cfg == nil {
		return nil, errors.New("configuration object not initialised when calling fetching release")
	}
	// get the release information based on the
	ri, err := s.release(appVersion)
	if err != nil {
		// could not find release information in the release plan
		return nil, err
	}
	// builds a uri to fetch the specific release manifest
	uri := fmt.Sprintf("%s/%s/release.json", s.get(SchemaURI), ri.Path)
	// fetch the release.json manifest
	response, err := s.client.Get(uri, s.addHttpHeaders)
	// if the request was unsuccessful then return the error
	if err != nil {
		return nil, err
	}
	// request was good so construct a release manifest reference
	r := &Release{}
	r, err = r.decode(response)
	defer func() {
		if ferr := response.Body.Close(); ferr != nil {
			err = ferr
		}
	}()
	return r, err
}

// get the release information for a given application version
func (s *ScriptManager) release(appVersion string) (*Info, error) {
	for _, release := range s.getPlan().Releases {
		if release.AppVersion == appVersion {
			return &release, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("information for application version '%s' does not exist in the release plan", appVersion))
}

// add http headers to the request object
func (s *ScriptManager) addHttpHeaders(req *http.Request, payload oxc.Serializable) error {
	// add headers to disable caching
	req.Header.Add("Cache-Control", `no-cache"`)
	req.Header.Add("Pragma", "no-cache")
	// if there is an access token defined
	if len(s.get(SchemaUsername)) > 0 && len(s.get(SchemaToken)) > 0 {
		credentials := base64.StdEncoding.EncodeToString([]byte(
			fmt.Sprintf("%s:%s", s.get(SchemaUsername), s.get(SchemaToken))))
		req.Header.Add("Authorization", credentials)
	}
	return nil
}

func (s *ScriptManager) get(key string) string {
	return s.cfg.Get(key)
}
