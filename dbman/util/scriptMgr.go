//   Onix Config DatabaseProvider - Dbman
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
	"io/ioutil"
	"net/http"
)

// the source of database scripts
type ScriptManager struct {
	client *oxc.Client
	cfg    *AppCfg
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
	// fetch the db init manifest
	response, err := s.client.Get(fmt.Sprintf("%s/init/init.json", s.get(SchemaURI)), s.addHttpHeaders)
	if err != nil {
		return nil, err
	}
	init := &DbInit{}
	init, err = init.decode(response)
	if err != nil {
		err = errors.New(fmt.Sprintf("database initialisation manifest is not in the right format: %v\n", err))
	}
	err = response.Body.Close()
	if err != nil {
		return nil, err
	}
	// creates a result to hold the fetched scripts
	result := &DbInit{
		Items: make([]Item, len(init.Items)),
	}
	// for each item retrieve the underlying db script
	for ix, item := range init.Items {
		// fetch the db script
		response, err = s.client.Get(fmt.Sprintf("%s/init/%s", s.get(SchemaURI), item.Script), s.addHttpHeaders)
		// decode response into a string
		if response.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(response.Body)
			if err != nil {
				response.Body.Close()
				return nil, err
			}
			bodyString := string(bodyBytes)
			// updates the item script with the script content fetched
			item.Script = bodyString
			// assign the item to the result
			result.Items[ix] = item
		}
		if err != nil {
			return nil, err
		}
		response.Body.Close()
	}
	return result, err
}

// fetches the getReleaseInfo plan
func (s *ScriptManager) fetchPlan() (*Plan, error) {
	if s.cfg == nil {
		return nil, errors.New("configuration object not initialised when fetching getReleaseInfo getPlan")
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

// fetches the scripts for a database getReleaseInfo
func (s *ScriptManager) fetchRelease(appVersion string) (*Release, error) {
	// if cfg not initialised, no point in continuing
	if s.cfg == nil {
		return nil, errors.New("configuration object not initialised when calling fetching getReleaseInfo")
	}
	// get the getReleaseInfo information based on the
	ri, err := s.getReleaseInfo(appVersion)
	if err != nil {
		// could not find getReleaseInfo information in the getReleaseInfo plan
		return nil, err
	}
	// builds a uri to fetch the specific getReleaseInfo manifest
	uri := fmt.Sprintf("%s/%s/release.json", s.get(SchemaURI), ri.Path)
	// fetch the getReleaseInfo.json manifest
	response, err := s.client.Get(uri, s.addHttpHeaders)
	// if the request was unsuccessful then return the error
	if err != nil {
		return nil, err
	}
	// request was good so construct a release manifest reference
	r := &Release{}
	r, err = r.decode(response)
	if err != nil {
		return nil, err
	}
	err = response.Body.Close()
	if err != nil {
		return nil, err
	}
	// fetch the schema scripts
	schemas, err := s.getScripts(ri.Path, r.Schemas)
	if err != nil {
		return nil, err
	}
	r.Schemas = schemas
	// fetch function scripts
	funcs, err := s.getScripts(ri.Path, r.Functions)
	if err != nil {
		return nil, err
	}
	r.Functions = funcs
	// fetch upgrade scripts
	up, err := s.getScripts(ri.Path, r.Upgrade)
	if err != nil {
		return nil, err
	}
	r.Upgrade = up
	return r, err
}

// get the getReleaseInfo information for a given application version
func (s *ScriptManager) getReleaseInfo(appVersion string) (*Info, error) {
	plan, err := s.fetchPlan()
	if err != nil {
		return nil, err
	}
	for _, release := range plan.Releases {
		if release.AppVersion == appVersion {
			return &release, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("!!! information for application version '%s' does not exist in the release plan", appVersion))
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

// returns a slice with release scripts
// path: the release path
// the list of file names under the path to read
func (s *ScriptManager) getScripts(path string, files []string) ([]string, error) {
	result := make([]string, len(files))
	for ix, file := range files {
		uri := fmt.Sprintf("%v/%v/%v", s.cfg.Get(SchemaURI), path, file)
		response, err := s.client.Get(uri, s.addHttpHeaders)
		if err != nil {
			return nil, err
		}
		// decode response into a string
		if response.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(response.Body)
			if err != nil {
				response.Body.Close()
				return nil, err
			}
			result[ix] = string(bodyBytes)
		}
	}
	return result, nil
}
