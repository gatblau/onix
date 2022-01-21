//   Onix Config DatabaseProvider - Dbman
//   Copyright (c) 2018-Present by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.

package core

import (
	"errors"
	"fmt"
	. "github.com/gatblau/onix/dbman/plugin"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ScriptManager the source of database scripts
type ScriptManager struct {
	cfg *Config
}

// NewScriptManager factory function
func NewScriptManager(cfg *Config) (*ScriptManager, error) {
	// creates a new struct
	source := new(ScriptManager)
	// setup attributes
	source.cfg = cfg
	// source.client = client
	return source, nil
}

// fetchPlan fetches the getReleaseInfo plan
func (s *ScriptManager) fetchPlan() (*Plan, error) {
	// get the base uri to retrieve scripts (includes credentials if set)
	baseUri, err := s.getRepoUri()
	if err != nil {
		return nil, err
	}
	// note: passing http headers results in 503 on gitlab.com when using credentials on URI
	content, err := s.readFile(fmt.Sprintf("%s/plan.json", baseUri))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("! cannot retrieve release plan: %v", err))
	}
	p := &Plan{}
	p, err = p.decode(content)
	return p, err
}

// FetchFile reads a file stored in the remote repository
// path: the relative path to the file in the repository
func (s *ScriptManager) FetchFile(path string) (*string, error) {
	// get the base uri to retrieve scripts (includes credentials if set)
	baseUri, err := s.getRepoUri()
	if err != nil {
		return nil, err
	}
	// builds a uri to the specific file
	uri := fmt.Sprintf("%s%s", baseUri, path)
	// fetch the file
	// note: passing http headers results in 503 on gitlab.com when using credentials on URI
	content, err := s.readFile(uri)
	// if the request was unsuccessful then return the error
	if err != nil {
		return nil, err
	}
	result := string(content)
	// if the result is not an empty string
	if len(strings.Trim(result, " ")) > 0 {
		// return the result
		return &result, nil
	} else {
		// otherwise, return nil
		return nil, nil
	}
}

// fetchManifest a release manifest
// - appVersion: the version of the application release to fetchManifest
// - contentTypes: list of content type content to fetchManifest
func (s *ScriptManager) fetchManifest(appVersion string) (*Info, *Manifest, error) {
	// get the base uri to retrieve scripts (includes credentials if set)
	baseUri, err := s.getRepoUri()
	if err != nil {
		return nil, nil, err
	}
	// get the ReleaseInfo information based on the
	release, err := s.getReleaseInfo(appVersion)
	if err != nil {
		// could not find ReleaseInfo information in the getReleaseInfo plan
		return nil, nil, err
	}
	// builds a uri to fetchManifest the specific release manifest
	uri := fmt.Sprintf("%s/%s/manifest.json", baseUri, release.Path)
	// fetchManifest the manifest
	// note: passing http headers results in 503 on gitlab.com when using credentials on URI
	content, err := s.readFile(uri)
	// if the request was unsuccessful then return the error
	if err != nil {
		return nil, nil, err
	}
	// request was good so construct a release manifest reference
	man := &Manifest{}
	man, err = man.Decode(content)
	return release, man, nil
}

func (s *ScriptManager) fetchCommandContent(appVersion string, subPath string, command Command) (*Command, error) {
	// get the base uri to retrieve scripts (includes credentials if set)
	baseUri, err := s.getRepoUri()
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

func (s *ScriptManager) fetchQueryContent(appVersion string, subPath string, query Query, params map[string]string) (*Query, error) {
	// get the base uri to retrieve scripts (includes credentials if set)
	baseUri, err := s.getRepoUri()
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
		return &query, err
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

func (s *ScriptManager) get(key string) string {
	return s.cfg.GetString(key)
}

// add the content from the remote repository to the passed-in scripts
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

// add the content of the query from the remote repository
func (s *ScriptManager) addQueryContent(baseUri string, path string, query Query, params map[string]string) (Query, error) {
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
	// note: passing http headers results in 503 on gitlab.com when using credentials on URI
	content, err := s.readFile(uri)
	if err != nil {
		return "", err
	}
	return string(content[:]), err
}

// merges the passed-in script with the values in of the script vars
func (s *ScriptManager) merge(script string, vars []Var, params map[string]string) (string, error) {
	// merge vars if any
	for _, variable := range vars {
		var value string
		// if variable is in configuration
		if len(variable.FromConf) > 0 {
			// get the value from the configuration set
			value = s.get(variable.FromConf)
		} else
		// if variable has a value in the manifest
		if len(variable.FromValue) > 0 {
			// get the value from the manifest
			value = variable.FromValue
		} else
		// if a variable has a value passed-in as an input from the CLI or http URI
		if len(variable.FromInput) > 0 {
			value = params[variable.FromInput]
		}
		// validate for suspicious values
		if s.suspicious(value) {
			return "", errors.New(fmt.Sprintf("!!! I found suspicious content for variable '%s'", variable.Name))
		}
		// merge the variable value
		script = strings.Replace(script, fmt.Sprintf("{{%s}}", variable.Name), value, -1)
	}
	return script, nil
}

func (s *ScriptManager) suspicious(value string) bool {
	// containing blank spaces
	return len(strings.Split(value, " ")) > 1 ||
		// containing line breaks
		len(strings.Split(value, "\n")) > 1
}

func (s ScriptManager) getRepoUri() (string, error) {
	uri := s.get(RepoURI)
	if len(uri) == 0 {
		return "", errors.New(fmt.Sprintf("!!! The Repo.URI is not defined"))
	}
	// if the username and password have been set
	if len(s.get(RepoUsername)) > 0 && len(s.get(RepoPassword)) > 0 {
		uriParts := strings.Split(uri, "//")
		return fmt.Sprintf("%s//%s:%s@%s", uriParts[0], s.get(RepoUsername), s.get(RepoPassword), uriParts[1]), nil
	}
	return uri, nil
}

func (s *ScriptManager) readFile(uri string) ([]byte, error) {
	if strings.HasPrefix(uri, "http") {
		return getHttpFile(uri, fmt.Sprintf("%s:%s", s.get(RepoUsername), s.get(RepoPassword)))
	} else {
		return getFsFile(uri)
	}
	return nil, nil
}

// getFsFile reads a file from the file system
func getFsFile(uri string) ([]byte, error) {
	path, err := filepath.Abs(uri)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(path)
}

// getHttpFile reads a file from a http endpoint
func getHttpFile(uri, creds string) ([]byte, error) {
	// if credentials are provided
	if len(creds) > 0 {
		// add them to the uri scheme
		u, err := addCredsToHttpURI(uri, creds)
		if err != nil {
			return nil, err
		}
		uri = u
	}
	// create an http client with defined timeout
	client := http.Client{
		Timeout: 60 * time.Second,
	}
	// create a new http request
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	// add headers to disable caching
	req.Header.Add("Cache-Control", `no-cache"`)
	req.Header.Add("Pragma", "no-cache")
	// execute the request
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp != nil {
		switch resp.StatusCode {
		case 401:
			fallthrough
		case 403:
			return nil, errors.New(fmt.Sprintf("!!! I do not have permission to get the content of %s\n", uri))
		case 404:
			return nil, errors.New(fmt.Sprintf("!!! I cannot find %s\n", uri))
		case 408:
			return nil, errors.New(fmt.Sprintf("!!! I cannot retrieve content of %s as the server is not responding\n", uri))
		case 500:
			return nil, errors.New(fmt.Sprintf("!!! I cannot retrieve content of %s as the server responded with an internal error\n", uri))
		}
		// return the byte content in the response body
		return ioutil.ReadAll(resp.Body)
	}
	return nil, fmt.Errorf("server sent no response: %s", err)
}

// addCredsToHttpURI add credentials to http(s) URI
func addCredsToHttpURI(uri string, creds string) (string, error) {
	// if there are no credentials or the uri is a file path
	if len(creds) == 0 || strings.HasPrefix(uri, "http") {
		// skip and return as is
		return uri, nil
	}
	parts := strings.Split(uri, "/")
	if !strings.HasPrefix(parts[0], "http") {
		return uri, fmt.Errorf("invalid URI scheme, http(s) expected when specifying credentials\n")
	}
	parts[2] = fmt.Sprintf("%s@%s", creds, parts[2])
	return strings.Join(parts, "/"), nil
}
