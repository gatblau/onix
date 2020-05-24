//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package util

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gatblau/oxc"
	"github.com/rs/zerolog/log"
	"net/http"
)

// database release information
type RInfo struct {
	client    *oxc.Client
	index     *Index
	manifests []Release
	cfg       *Config
}

// factory function
func NewRInfo(cfg *Config, client *oxc.Client) (*RInfo, error) {
	// creates a new struct
	script := new(RInfo)
	// setup attributes
	script.cfg = cfg
	script.client = client
	return script, nil
}

// new oxc configuration
func NewClientConf(cfg *Config) *oxc.ClientConf {
	return &oxc.ClientConf{
		BaseURI:            cfg.Path,
		InsecureSkipVerify: false,
		AuthMode:           oxc.None,
	}
}

// access a cached index reference
func (info *RInfo) ix() *Index {
	// if the index is not fetched
	if info.index == nil {
		// fetches it
		err := info.loadIx()
		if err != nil {
			log.Error().Msgf("cannot retrieve index, %info", err)
		}
	}
	return info.index
}

// (re)loads the internal index reference
func (info *RInfo) loadIx() error {
	ix, err := info.fetchIndex()
	info.index = ix
	return err
}

// fetches the release index
func (info *RInfo) fetchIndex() (*Index, error) {
	if info.cfg == nil {
		return nil, errors.New("configuration object not initialised when fetching release index")
	}
	response, err := info.client.Get(fmt.Sprintf("%s/index.json", info.cfg.SchemaURI), info.addHttpHeaders)
	if err != nil {
		return nil, err
	}
	i := &Index{}
	i, err = i.decode(response)
	defer func() {
		if ferr := response.Body.Close(); ferr != nil {
			err = ferr
		}
	}()
	return i, err
}

// fetches the scripts for a database release
func (info *RInfo) fetchRelease(appVersion string) (*Release, error) {
	// if cfg not initialised, no point in continuing
	if info.cfg == nil {
		return nil, errors.New("configuration object not initialised when calling fetching release")
	}
	// get the release information based on the
	ri, err := info.release(appVersion)
	if err != nil {
		// could not find release information in the release index
		return nil, err
	}
	// builds a uri to fetch the specific release manifest
	uri := fmt.Sprintf("%s/%s/release.json", info.cfg.SchemaURI, ri.Path)
	// fetch the release.json manifest
	response, err := info.client.Get(uri, info.addHttpHeaders)
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
	// result, err := read(response, &Release{})
	// return result.(*Release), err
	return nil, nil
}

func read(response *http.Response, obj *interface{}) (*interface{}, error) {
	var err error
	defer func() {
		if ferr := response.Body.Close(); ferr != nil {
			err = ferr
		}
	}()
	err = json.NewDecoder(response.Body).Decode(obj)
	return obj, err
}

// get the release information for a given application version
func (info *RInfo) release(appVersion string) (*ReleaseInfo, error) {
	for _, release := range info.ix().Releases {
		if release.AppVersion == appVersion {
			return &release, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("information for application version '%s' does not exist in the release index", appVersion))
}

// add http headers to the request object
func (info *RInfo) addHttpHeaders(req *http.Request, payload oxc.Serializable) error {
	// add headers to disable caching
	req.Header.Add("Cache-Control", `no-cache"`)
	req.Header.Add("Pragma", "no-cache")
	// if there is an access token defined
	if len(info.cfg.SchemaUsername) > 0 && len(info.cfg.SchemaToken) > 0 {
		credentials := base64.StdEncoding.EncodeToString([]byte(
			fmt.Sprintf("%s:%s", info.cfg.SchemaUsername, info.cfg.SchemaToken)))
		req.Header.Add("Authorization", credentials)
	}
	return nil
}
