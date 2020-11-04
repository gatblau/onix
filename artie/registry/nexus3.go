/*
  Onix Config Manager - Artie
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package registry

import (
	"errors"
	"fmt"
	"github.com/gatblau/onix/artie/core"
	"io"
	"log"
	"net/http"
	"strings"
)

// a Nexus3 implementation of a remote registry
type Nexus3Registry struct {
}

func NewNexus3Registry() Remote {
	return &Nexus3Registry{}
}

func (r *Nexus3Registry) CreateRepository() {}
func (r *Nexus3Registry) DeleteRepository() {}

// Upload an artefact
func (r *Nexus3Registry) UploadArtefact(client *http.Client, name *core.ArtieName, localPath string, fileReference string, credentials string) error {
	// prepare the reader instances to encode
	values := map[string]io.Reader{
		"raw.directory": strings.NewReader(name.Repository()),
		// the json filename
		"raw.asset1.filename": strings.NewReader(fmt.Sprintf("%s.json", fileReference)),
		// the json file (seal)
		"raw.asset1": mustOpen(fmt.Sprintf("%s/%s.json", localPath, fileReference)),
		// the zip filename
		"raw.asset2.filename": strings.NewReader(fmt.Sprintf("%s.zip", fileReference)),
		// the zip file (artefact)
		"raw.asset2": mustOpen(fmt.Sprintf("%s/%s.zip", localPath, fileReference)),
	}
	user, pwd := userPwd(credentials)
	return upload(client, r.uploadURI(name), values, user, pwd)
}

func (r *Nexus3Registry) uploadURI(name *core.ArtieName) string {
	return fmt.Sprintf("http://%s/service/rest/v1/components?repository=artie", name.Domain)
}

func (r *Nexus3Registry) DownloadArtefact() {}

func userPwd(creds string) (user, pwd string) {
	// if credentials not provided then no user / pwd
	if len(creds) == 0 {
		return "", ""
	}
	parts := strings.Split(creds, ":")
	if len(parts) != 2 {
		log.Fatal(errors.New("credentials in incorrect format, they should be USER:PWD"))
	}
	return parts[0], parts[1]
}
