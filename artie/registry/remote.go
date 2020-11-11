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

// the interface implemented by a remote registry
type Remote interface {
	// create an artefact repository
	CreateRepository()
	// delete a repository
	DeleteRepository()
	// upload an artefact
	UploadArtefact(name *core.ArtieName, artefactRef string, zipfile multipart.File, jsonFile multipart.File, user string, pwd string) error
	// download an artefact
	DownloadArtefact()
}
