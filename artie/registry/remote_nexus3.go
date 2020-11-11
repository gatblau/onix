/*
  Onix Config Manager - Artie
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package registry

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/gatblau/onix/artie/core"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

// a Nexus3 implementation of a remote registry
type Nexus3Registry struct {
	domain string
	client *http.Client
}

func NewNexus3Registry(domain string) Remote {
	return &Nexus3Registry{
		domain: domain,
		client: &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

func (r *Nexus3Registry) CreateRepository() {}
func (r *Nexus3Registry) DeleteRepository() {}

// Upload an artefact
func (r *Nexus3Registry) UploadArtefact(name *core.ArtieName, artefactRef string, zipfile multipart.File, jsonFile multipart.File, user string, pwd string) error {
	// ensure files are properly closed
	defer zipfile.Close()
	defer jsonFile.Close()

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	err := r.addField(writer, "raw.directory", name.Repository())
	if err != nil {
		return err
	}
	err = r.addField(writer, "raw.asset1.filename", fmt.Sprintf("%s.json", artefactRef))
	if err != nil {
		return err
	}
	err = r.addFile(writer, "raw.asset1", fmt.Sprintf("%s.json", artefactRef), jsonFile)
	if err != nil {
		return err
	}
	err = r.addField(writer, "raw.asset2.filename", fmt.Sprintf("%s.zip", artefactRef))
	if err != nil {
		return err
	}
	err = r.addFile(writer, "raw.asset2", fmt.Sprintf("%s.zip", artefactRef), zipfile)
	if err != nil {
		return err
	}
	// don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	writer.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", r.uploadURI(), &b)
	if err != nil {
		return err
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("accept", "application/json")
	if len(user) > 0 && len(pwd) > 0 {
		req.Header.Add("authorization", basicToken(user, pwd))
	}
	// Submit the request
	res, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// Check the response
	if res.StatusCode > 299 {
		return fmt.Errorf("failed to push, the remote server responded with status code %d: %s", res.StatusCode, res.Status)
	}
	return nil
}

// add a field to a multipart form
func (r *Nexus3Registry) addField(writer *multipart.Writer, fieldName, fieldValue string) error {
	// create a writer with the mime header for the field
	formWriter, err := writer.CreateFormField(fieldName)
	if err != nil {
		return err
	}
	// writes the field value
	_, err = io.Copy(formWriter, strings.NewReader(fieldValue))
	return err
}

// add a file to a multipart form
func (r *Nexus3Registry) addFile(writer *multipart.Writer, fieldName, fileName string, file multipart.File) error {
	// create a writer with the mime header for the field
	formWriter, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return err
	}
	// writes the field value
	_, err = io.Copy(formWriter, file)
	file.Close()
	return err
}

func (r *Nexus3Registry) uploadURI() string {
	return fmt.Sprintf("%s/service/rest/v1/components?repository=artie", r.domain)
}

func (r *Nexus3Registry) DownloadArtefact() {}
