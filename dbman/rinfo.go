//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package main

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
)

// database release information
type RInfo struct {
	client    *Source
	index     *Index
	manifests []Release
	cfg       *Config
}

// factory function
func NewRInfo(cfg *Config, client *Source) (*RInfo, error) {
	// creates a new struct
	script := new(RInfo)
	// setup attributes
	script.cfg = cfg
	script.client = client
	return script, nil
}

// access a cached index reference
func (s *RInfo) ix() *Index {
	// if the index is not fetched
	if s.index == nil {
		// fetches it
		err := s.loadIx()
		if err != nil {
			log.Error().Msgf("cannot retrieve index, %s", err)
		}
	}
	return s.index
}

// (re)loads the internal index reference
func (s *RInfo) loadIx() error {
	ix, err := s.fetchIndex()
	s.index = ix
	return err
}

// fetches the release index
func (s *RInfo) fetchIndex() (*Index, error) {
	if s.cfg == nil {
		return nil, errors.New("configuration object not initialised when fetching release index")
	}
	response, err := s.client.get(fmt.Sprintf("%s/index.json", s.cfg.SchemaURI))
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
func (s *RInfo) fetchRelease(appVersion string) (*Release, error) {
	// if cfg not initialised, no point in continuing
	if s.cfg == nil {
		return nil, errors.New("configuration object not initialised when calling fetching release")
	}
	// get the release information based on the
	info, err := s.release(appVersion)
	if err != nil {
		// could not find release information in the release index
		return nil, err
	}
	// builds a uri to fetch the specific release manifest
	uri := fmt.Sprintf("%s/%s/release.json", s.cfg.SchemaURI, info.Path)
	// fetch the release.json manifest
	response, err := s.client.get(uri)
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
func (s *RInfo) release(appVersion string) (*ReleaseInfo, error) {
	for _, release := range s.ix().Releases {
		if release.AppVersion == appVersion {
			return &release, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("information for application version '%s' does not exist in the release index", appVersion))
}
