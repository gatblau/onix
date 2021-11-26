/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package deploy

import (
	"fmt"
)

// AppManifest the application manifest that is made up of one or more service manifests
type AppManifest []SvcRef

type SvcRef struct {
	// the name of the service
	Name string `yaml:"name"`
	// the uri of the service manifest
	URI string `yaml:"uri,omitempty"`
	// the URI of the service image containing the service manifest
	Image string `yaml:"image,omitempty"`
	// whether this service should not be publicly exposed, by default is false
	Private bool `yaml:"private,omitempty"`
	// the service port, if not specified, the application port (in the service manifest) is used
	Port string `yaml:"port,omitempty"`
	// the service manifest loaded from remote image
	Service *SvcManifest `yaml:"service,omitempty"`
}

// NewAppMan creates a new application manifest from an URI (supported schemes are http(s):// and file://
func NewAppMan(uri string) (*AppManifest, error) {
	if ok, path := isFile(uri); ok {
		return loadFromFile(path)
	}
	if isURL(uri) {
		return loadFromURL(uri)
	}
	return nil, fmt.Errorf("invalid URI value '%s': should start with either file://, http:// or https://\n", uri)
}

// Explode augments an app manifest that has remote references to service manifests
func (m *AppManifest) Explode() (*AppManifest, error) {
	var err error
	// create a copy of the passed in light manifest to become the exploded version
	appMan := m.deepCopy()
	// validate the app manifest
	if err = m.validate(); err != nil {
		return nil, err
	}
	// loop through
	var svcMan *SvcManifest
	for i, svc := range *m {
		// image only
		if len(svc.Image) > 0 && len(svc.URI) == 0 {
			svcMan, err = loadSvcManFromImage(svc)
			if err != nil {
				return nil, fmt.Errorf("cannot load service manifest for '%s': %s\n", svc.Image, err)
			}
		} else if len(svc.Image) > 0 && len(svc.URI) > 0 {
			svcMan, err = loadSvcManFromURI(svc)
			if err != nil {
				return nil, fmt.Errorf("cannot load service manifest for '%s': %s\n", svc.Image, err)
			}
		}
		appMan[i].Service = svcMan
	}
	return &appMan, nil
}

func (m *AppManifest) validate() error {
	for _, svc := range *m {
		// check that the manifest has named services
		if len(svc.Name) == 0 {
			return fmt.Errorf("service manifest name is required for image %s\n", svc.Image)
		}
		// case of manifest embedded in docker image then no URI is needed (image only)
		// case of manifest in git repo (uri + image required)
		// so cases to avoid is uri only
		if len(svc.Image) == 0 && len(svc.URI) > 0 {
			return fmt.Errorf("invalid entry for service '%s' manifest in application manifest: only one of Image or URI attributes must be specified\n", svc.Name)
		}
		// or uri & image not provided
		if len(svc.Image) == 0 && len(svc.URI) == 0 {
			return fmt.Errorf("invalid entry for service '%s' manifest in application manifest: either one of Image or URI attributes must be specified\n", svc.Name)
		}
	}
	return nil
}

func (m *AppManifest) deepCopy() AppManifest {
	result := AppManifest{}
	for _, svc := range *m {
		result = append(result, svc)
	}
	return result
}
