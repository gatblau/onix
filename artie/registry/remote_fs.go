package registry

import (
	"github.com/gatblau/onix/artie/core"
	"mime/multipart"
)

type RemoteFs struct {
	root string
}

// create an artefact repository
func (f *RemoteFs) CreateRepository() {
}

// delete a repository
func (f *RemoteFs) DeleteRepository() {}

// upload an artefact
func (f *RemoteFs) UploadArtefact(name *core.ArtieName, artefactRef string, zipfile multipart.File, jsonFile multipart.File, user string, pwd string) error {

	return nil
}

// download an artefact
func (f *RemoteFs) DownloadArtefact() {}
