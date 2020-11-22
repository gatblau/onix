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
	"encoding/json"
	"fmt"
	"github.com/gatblau/onix/artie/core"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

// Artie's HTTP API
type GenericApi struct {
	https  bool
	domain string
	client *http.Client
}

func NewGenericAPI(domain string, useTls bool) *GenericApi {
	return &GenericApi{
		https:  useTls,
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

func (r *GenericApi) UploadArtefact(name *core.ArtieName, artefactRef string, zipfile multipart.File, jsonFile multipart.File, metaInfo *Artefact, user string, pwd string) error {
	// ensure files are properly closed
	defer zipfile.Close()
	defer jsonFile.Close()

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	info, err := metaInfo.ToJson()
	core.CheckErr(err, "cannot marshall artefact info")
	err = r.addField(writer, "artefact-meta", info)
	core.CheckErr(err, "cannot add artefact meta file")
	err = r.addFile(writer, "artefact-seal", fmt.Sprintf("%s.json", artefactRef), jsonFile)
	core.CheckErr(err, "cannot add seal file")
	err = r.addFile(writer, "artefact-file", fmt.Sprintf("%s.zip", artefactRef), zipfile)
	core.CheckErr(err, "cannot add artefact file")
	// don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	err = writer.Close()
	core.CheckErr(err, "cannot close writer")
	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", r.artefactTagURI(name.Group, name.Name, name.Tag), &b)
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
	switch res.StatusCode {
	case http.StatusCreated:
		core.Msg("artefact pushed")
	case http.StatusOK:
		core.Msg("nothing to do")
	}
	return nil
}

func (r *GenericApi) UpdateArtefactInfo(name *core.ArtieName, artefact *Artefact, user string, pwd string) error {
	str, err := artefact.ToJson()
	if err != nil {
		return err
	}
	body := bytes.NewReader([]byte(str))
	req, err := http.NewRequest("PUT", r.artefactIdURI(name.Group, name.Name, artefact.Id), body)
	if err != nil {
		return err
	}
	req.Header.Set("accept", "application/json")
	if len(user) > 0 && len(pwd) > 0 {
		req.Header.Add("authorization", basicToken(user, pwd))
	}
	// Submit the request
	res, err := r.client.Do(req)
	core.CheckErr(err, "cannot post to backend")
	// Check the response
	if res.StatusCode > 299 {
		return fmt.Errorf("failed to update artefact info, the remote server responded with: %s", res.Status)
	}
	return nil
}

func (r *GenericApi) GetRepositoryInfo(group, name, user, pwd string) (*Repository, error) {
	req, err := http.NewRequest("GET", r.repoURI(group, name), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	if len(user) > 0 && len(pwd) > 0 {
		req.Header.Add("authorization", basicToken(user, pwd))
	}
	// Submit the request
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	repo := new(Repository)
	err = json.Unmarshal(bytes, repo)
	return repo, err
}

func (r *GenericApi) GetArtefactInfo(group, name, id, user, pwd string) (*Artefact, error) {
	req, err := http.NewRequest("GET", r.artefactIdURI(group, name, id), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	if len(user) > 0 && len(pwd) > 0 {
		req.Header.Add("authorization", basicToken(user, pwd))
	}
	// Submit the request
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// if the body contains a nil response
	if string(bytes) == "null" {
		return nil, nil
	}
	artefact := new(Artefact)
	err = json.Unmarshal(bytes, artefact)
	return artefact, err
}

func (r *GenericApi) DownloadArtefact() {
	panic("implement me")
}

func (r *GenericApi) repoURI(group, name string) string {
	scheme := "http"
	if r.https {
		scheme = fmt.Sprintf("%ss", scheme)
	}
	// {scheme}://{domain}/repository/{repository-group}/{repository-name}
	return fmt.Sprintf("%s://%s/repository/%s/%s", scheme, r.domain, group, name)
}

func (r *GenericApi) artefactURI(group, name string) string {
	scheme := "http"
	if r.https {
		scheme = fmt.Sprintf("%ss", scheme)
	}
	// {scheme}://{domain}/registry/{repository-group}/{repository-name}/{tag}
	return fmt.Sprintf("%s://%s/artefact/%s/%s", scheme, r.domain, group, name)
}

func (r *GenericApi) artefactTagURI(group, name, tag string) string {
	return fmt.Sprintf("%s/tag/%s", r.artefactURI(group, name), tag)
}

func (r *GenericApi) artefactIdURI(group, name, id string) string {
	return fmt.Sprintf("%s/id/%s", r.artefactURI(group, name), id)
}

// add a field to a multipart form
func (r *GenericApi) addField(writer *multipart.Writer, fieldName, fieldValue string) error {
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
func (r *GenericApi) addFile(writer *multipart.Writer, fieldName, fileName string, file multipart.File) error {
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
