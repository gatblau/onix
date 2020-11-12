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

// Artie's HTTP API
type RemoteApi struct {
	https  bool
	client *http.Client
}

func NewRemoteAPI(useTls bool) Remote {
	return &RemoteApi{
		https: useTls,
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

func (r *RemoteApi) CreateRepository() {
	panic("implement me")
}

func (r *RemoteApi) DeleteRepository() {
	panic("implement me")
}

func (r *RemoteApi) UploadArtefact(name *core.ArtieName, artefactRef string, zipfile multipart.File, jsonFile multipart.File, user string, pwd string) error {
	// ensure files are properly closed
	defer zipfile.Close()
	defer jsonFile.Close()

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	// err := r.addField(writer, "artefact.repository", name.Repository())
	// core.CheckErr(err, "cannot add field artefact.repository")
	// err = r.addField(writer, "artefact.fileRef", artefactRef)
	// core.CheckErr(err, "cannot add field artefact.fileRef")
	err := r.addFile(writer, "artefact-seal", fmt.Sprintf("%s.json", artefactRef), jsonFile)
	core.CheckErr(err, "cannot add seal file")
	err = r.addFile(writer, "artefact-file", fmt.Sprintf("%s.zip", artefactRef), zipfile)
	core.CheckErr(err, "cannot add artefact file")
	// don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	writer.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", r.uploadURI(name, artefactRef), &b)
	core.CheckErr(err, "cannot create http request")
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("accept", "application/json")
	if len(user) > 0 && len(pwd) > 0 {
		req.Header.Add("authorization", basicToken(user, pwd))
	}
	// Submit the request
	res, err := r.client.Do(req)
	core.CheckErr(err, "cannot post to backend")

	// Check the response
	if res.StatusCode > 299 {
		return fmt.Errorf("failed to push, the remote server responded with status code %d: %s", res.StatusCode, res.Status)
	}
	return nil
}

func (r *RemoteApi) DownloadArtefact() {
	panic("implement me")
}

func (r *RemoteApi) uploadURI(name *core.ArtieName, artefactRef string) string {
	scheme := "http"
	if r.https {
		scheme = fmt.Sprintf("%ss", scheme)
	}
	// {scheme}://{domain}/artefact/{repository-group}/{repository-name}/{artefact-ref}
	return fmt.Sprintf("%s://%s/registry/%s/%s/%s", scheme, name.Domain, name.Group, name.Name, artefactRef)
}

// add a field to a multipart form
func (r *RemoteApi) addField(writer *multipart.Writer, fieldName, fieldValue string) error {
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
func (r *RemoteApi) addFile(writer *multipart.Writer, fieldName, fileName string, file multipart.File) error {
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
