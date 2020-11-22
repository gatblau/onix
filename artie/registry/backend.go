/*
  Onix Config Manager - Artie
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package registry

import (
	"github.com/gatblau/onix/artie/core"
	"mime/multipart"
)

// the interface implemented by a backend
type Backend interface {
	// upload an artefact to the remote repository
	UploadArtefact(name *core.ArtieName, artefactRef string, zipfile multipart.File, jsonFile multipart.File, repo multipart.File, user string, pwd string) error
	// get repository information
	GetRepositoryInfo(group, name, user, pwd string) (*Repository, error)
	// get artefact information
	GetArtefactInfo(group, name, id, user, pwd string) (*Artefact, error)
	// update artefact information
	UpdateArtefactInfo(group string, name string, artefact *Artefact, user string, pwd string) error
}

func GetBackend() Backend {
	conf := new(core.ServerConfig)
	// get the configured factory
	switch conf.Backend() {
	case core.Nexus3:
		return NewNexus3Backend(
			conf.BackendDomain(), // the nexus scheme://domain:port
		)
	}
	core.RaiseErr("backend not recognised")
	return nil
}
