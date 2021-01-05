/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package registry

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Artie's HTTP Registry API
type Api struct {
	https  bool
	domain string
	client *http.Client
	tmp    string
}

func NewGenericAPI(domain string, noTLS bool) *Api {
	core.TmpExists()
	return &Api{
		https:  !noTLS,
		domain: domain,
		tmp:    core.TmpPath(),
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

func (r *Api) UploadArtefact(name *core.PackageName, artefactRef string, zipfile multipart.File, jsonFile multipart.File, metaInfo *Artefact, user string, pwd string) error {
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

func (r *Api) UpdateArtefactInfo(name *core.PackageName, artefact *Artefact, user string, pwd string) error {
	b, err := json.Marshal(artefact)
	if err != nil {
		return err
	}
	body := bytes.NewReader([]byte(b))
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

func (r *Api) GetRepositoryInfo(group, name, user, pwd string) (*Repository, error) {
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
	switch resp.StatusCode {
	case http.StatusNotFound:
		// if repository is nil then the client is not talking to the proper artefact registry
		return nil, fmt.Errorf("\"%s\" does not conform to the artefact registry api, are you sure the artefact domain is correct", r.domain)
	case http.StatusForbidden:
		return nil, fmt.Errorf("invalid credentials, access to the registry is forbidden")
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// if the result body is not in JSON format is likely that the domain of artefact does not exist
	if !isJSON(string(b)) {
		return nil, fmt.Errorf("the artefact was not found: its domain/group/name is likely to be incorrect")
	}
	// if not response then return an empty repository
	if len(b) == 0 {
		return &Repository{
			Repository: fmt.Sprintf("%s/%s", group, name),
			Artefacts:  make([]*Artefact, 0),
		}, nil
	}
	repo := new(Repository)
	err = json.Unmarshal(b, repo)
	return repo, err
}

func (r *Api) GetArtefactInfo(group, name, id, user, pwd string) (*Artefact, error) {
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
	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, nil
	case http.StatusForbidden:
		return nil, fmt.Errorf("invalid credentials, access to the registry is forbidden")
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if !isJSON(string(b)) {
		return nil, fmt.Errorf("invalid artefact name: %s/%s/%s", r.domain, group, name)
	}
	artefact := new(Artefact)
	err = json.Unmarshal(b, artefact)
	return artefact, err
}

func (r *Api) Download(group, name, filename, user, pwd string) (string, error) {
	req, err := http.NewRequest("GET", r.fileURI(group, name, filename), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("accept", "application/json")
	if len(user) > 0 && len(pwd) > 0 {
		req.Header.Add("authorization", basicToken(user, pwd))
	}
	// Submit the request
	res, err := r.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case http.StatusNotFound:
		return "", fmt.Errorf("file '%s' not found in registry", filename)
	case http.StatusForbidden:
		return "", fmt.Errorf("invalid credentials, access to the registry is forbidden")
	}
	// write response to a temp file
	var b bytes.Buffer
	out := bufio.NewWriter(&b)
	_, err = io.Copy(out, res.Body)
	if err != nil {
		return "", err
	}
	err = out.Flush()
	if err != nil {
		return "", err
	}
	file, err := os.Create(filepath.Join(r.tmp, filename))
	if err != nil {
		return "", err
	}
	_, err = file.Write(b.Bytes())
	file.Close()
	return file.Name(), err
}

func (r *Api) repoURI(group, name string) string {
	scheme := "http"
	if r.https {
		scheme = fmt.Sprintf("%ss", scheme)
	}
	// {scheme}://{domain}/repository/{repository-group}/{repository-name}
	return fmt.Sprintf("%s://%s/repository/%s/%s", scheme, r.domain, group, name)
}

func (r *Api) artefactURI(group, name string) string {
	scheme := "http"
	if r.https {
		scheme = fmt.Sprintf("%ss", scheme)
	}
	// {scheme}://{domain}/artefact/{repository-group}/{repository-name}/{tag}
	return fmt.Sprintf("%s://%s/artefact/%s/%s", scheme, r.domain, group, name)
}

func (r *Api) artefactTagURI(group, name, tag string) string {
	return fmt.Sprintf("%s/tag/%s", r.artefactURI(group, name), tag)
}

func (r *Api) artefactIdURI(group, name, id string) string {
	return fmt.Sprintf("%s/id/%s", r.artefactURI(group, name), id)
}

func (r *Api) fileURI(group, name, filename string) string {
	scheme := "http"
	if r.https {
		scheme = fmt.Sprintf("%ss", scheme)
	}
	return fmt.Sprintf("%s://%s/file/%s/%s/%s", scheme, r.domain, group, name, filename)
}

// add a field to a multipart form
func (r *Api) addField(writer *multipart.Writer, fieldName, fieldValue string) error {
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
func (r *Api) addFile(writer *multipart.Writer, fieldName, fileName string, file multipart.File) error {
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
