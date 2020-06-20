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
	. "github.com/gatblau/onix/dbman/plugins"
	"github.com/gatblau/oxc"
	"io/ioutil"
	"net/http"
	"strings"
)

// the source of database scripts
type ScriptManager struct {
	client *oxc.Client
	cfg    *Config
}

// factory function
func NewScriptManager(cfg *Config, client *oxc.Client) (*ScriptManager, error) {
	// creates a new struct
	source := new(ScriptManager)
	// setup attributes
	source.cfg = cfg
	source.client = client
	return source, nil
}

// new oxc configuration
func NewOxClientConf(cfg *Config) *oxc.ClientConf {
	return &oxc.ClientConf{
		BaseURI:            cfg.Get(SchemaURI),
		InsecureSkipVerify: false,
		AuthMode:           oxc.None,
	}
}

// fetches the getReleaseInfo plan
func (s *ScriptManager) fetchPlan() (*Plan, error) {
	// get the base uri to retrieve scripts (includes credentials if set)
	baseUri, err := s.getSchemaUri()
	if err != nil {
		return nil, err
	}
	response, err := s.client.Get(fmt.Sprintf("%s/plan.json", baseUri), s.addHttpHeaders)
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

// fetchManifest a release manifest
// - appVersion: the version of the application release to fetchManifest
// - contentTypes: list of content type content to fetchManifest
func (s *ScriptManager) fetchManifest(appVersion string) (*Manifest, error) {
	// get the base uri to retrieve scripts (includes credentials if set)
	baseUri, err := s.getSchemaUri()
	if err != nil {
		return nil, err
	}
	// get the ReleaseInfo information based on the
	release, err := s.getReleaseInfo(appVersion)
	if err != nil {
		// could not find ReleaseInfo information in the getReleaseInfo plan
		return nil, err
	}
	// builds a uri to fetchManifest the specific release manifest
	uri := fmt.Sprintf("%s/%s/manifest.json", baseUri, release.Path)
	// fetchManifest the manifest
	response, err := s.client.Get(uri, s.addHttpHeaders)
	// if the request was unsuccessful then return the error
	if err != nil {
		return nil, err
	}
	// request was good so construct a release manifest reference
	man := &Manifest{}
	man, err = man.Decode(response)
	return man, nil
}

func (s *ScriptManager) fetchCommandContent(appVersion string, subPath string, command Command) (*Command, error) {
	// get the base uri to retrieve scripts (includes credentials if set)
	baseUri, err := s.getSchemaUri()
	if err != nil {
		return nil, err
	}
	// get the ReleaseInfo information based on the
	release, err := s.getReleaseInfo(appVersion)
	if err != nil {
		// could not find ReleaseInfo information in the getReleaseInfo plan
		return nil, err
	}
	command.Scripts, err = s.addScriptsContent(baseUri, fmt.Sprintf("%s/%s", release.Path, subPath), command.Scripts)
	if err != nil {
		return nil, err
	}
	return &command, nil
}

func (s *ScriptManager) fetchQueryContent(appVersion string, subPath string, query Query, params []string) (*Query, error) {
	// get the base uri to retrieve scripts (includes credentials if set)
	baseUri, err := s.getSchemaUri()
	if err != nil {
		return nil, err
	}
	// get the ReleaseInfo information based on the
	release, err := s.getReleaseInfo(appVersion)
	if err != nil {
		// could not find ReleaseInfo information in the getReleaseInfo plan
		return nil, err
	}
	query, err = s.addQueryContent(baseUri, fmt.Sprintf("%s/%s", release.Path, subPath), query, params)
	if err != nil {
		return nil, err
	}
	return &query, nil
}

// get the release information for a given application version
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
	if len(s.get(SchemaUsername)) > 0 && len(s.get(SchemaPassword)) > 0 {
		credentials := base64.StdEncoding.EncodeToString([]byte(
			fmt.Sprintf("%s:%s", s.get(SchemaUsername), s.get(SchemaPassword))))
		req.Header.Add("Authorization", credentials)
	}
	return nil
}

func (s *ScriptManager) get(key string) string {
	return s.cfg.Get(key)
}

// add the content form the remote repository to the passed-in scripts
func (s *ScriptManager) addScriptsContent(baseUri string, path string, scripts []Script) ([]Script, error) {
	var result []Script
	for _, script := range scripts {
		content, err := s.getContent(baseUri, path, script.File)
		if err != nil {
			return nil, err
		}
		script.Content = content
		mergedScript, err := s.merge(script.Content, script.Vars, nil)
		if err != nil {
			return nil, err
		}
		script.Content = mergedScript
		result = append(result, script)
	}
	return result, nil
}

func (s *ScriptManager) addQueryContent(baseUri string, path string, query Query, params []string) (Query, error) {
	// retrieve content from the remote repository
	content, err := s.getContent(baseUri, path, query.File)
	if err != nil {
		return Query{}, err
	}
	// assign the content to the query
	query.Content = content
	// merge vars
	mergedQuery, err := s.merge(query.Content, query.Vars, params)
	if err != nil {
		return query, err
	}
	query.Content = mergedQuery
	return query, nil
}

// get the content of a particular script on a git repo via http client
func (s *ScriptManager) getContent(baseUri string, path string, file string) (string, error) {
	// get the uri of the script
	uri := fmt.Sprintf("%v/%v/%v", baseUri, path, file)
	// issue an http request for the content
	response, err := s.client.Get(uri, s.addHttpHeaders)
	if err != nil {
		return "", err
	}
	// decode response into a string
	if response.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			response.Body.Close()
			return "", err
		}
		return string(bodyBytes), err

	}
	return "", err
}

// merges the passed-in script with the values in of the script vars
func (s *ScriptManager) merge(script string, vars []Var, params []string) (string, error) {
	// merge vars if any
	for ix, variable := range vars {
		var value string
		// if variable is in configuration
		if len(variable.FromConf) > 0 {
			// get the value from the configuration set
			value = s.cfg.Get(variable.FromConf)
		} else
		// if variable has a value in the manifest
		if len(variable.FromValue) > 0 {
			// get the value from the manifest
			value = variable.FromValue
		} else
		// if a variable has a value in the release
		if len(params) > 0 {
			value = params[ix]
		}
		// merge the variable value
		script = strings.Replace(script, variable.Name, value, -1)
	}
	return script, nil
}

func (s ScriptManager) getSchemaUri() (string, error) {
	uri := s.cfg.Get(SchemaURI)
	if len(uri) == 0 {
		return "", errors.New(fmt.Sprintf("!!! The SchemaURI is not defined"))
	}
	if !strings.HasPrefix(strings.ToLower(uri), "http") {
		return "", errors.New(fmt.Sprintf("!!! The SchemaURI must be an http(s) address"))
	}
	// if the username and password have been set
	if len(s.cfg.Get(SchemaUsername)) > 0 && len(s.cfg.Get(SchemaPassword)) > 0 {
		uriParts := strings.Split(uri, "//")
		return fmt.Sprintf("%s//%s:%s@%s", uriParts[0], s.cfg.Get(SchemaUsername), s.cfg.Get(SchemaPassword), uriParts[1]), nil
	}
	return uri, nil
}
