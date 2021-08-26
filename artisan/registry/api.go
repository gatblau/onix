package registry

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/i18n"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Api HTTP Registry API
type Api struct {
	domain string
	client *http.Client
	tmp    string
}

func NewGenericAPI(domain string) *Api {
	core.TmpExists()
	return &Api{
		domain: domain,
		tmp:    core.TmpPath(),
		client: &http.Client{
			Timeout: time.Second * 60,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

func (r *Api) UploadPackage(name *core.PackageName, packageRef string, zipfile multipart.File, jsonFile multipart.File, metaInfo *Package, user string, pwd string, https bool) error {
	// ensure files are properly closed
	defer zipfile.Close()
	defer jsonFile.Close()

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	info, err := metaInfo.ToJson()
	core.CheckErr(err, "cannot marshall package info")
	err = r.addField(writer, "package-meta", info)
	core.CheckErr(err, "cannot add package meta file")
	err = r.addFile(writer, "package-seal", fmt.Sprintf("%s.json", packageRef), jsonFile)
	core.CheckErr(err, "cannot add seal file")
	err = r.addFile(writer, "package-file", fmt.Sprintf("%s.zip", packageRef), zipfile)
	core.CheckErr(err, "cannot add package file")
	// don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	err = writer.Close()
	core.CheckErr(err, "cannot close writer")
	// create and start bar
	bar := pb.Simple.New(b.Len()).Start()
	bar.Set("prefix", "package > ")
	bar.SetWriter(os.Stdout)
	// create proxy reader
	reader := bar.NewProxyReader(&b)
	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", r.packageTagURI(name.Group, name.Name, name.Tag, https), reader)
	core.CheckErr(err, "cannot create http request")
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("accept", "application/json")
	if len(user) > 0 && len(pwd) > 0 {
		req.Header.Add("authorization", core.BasicToken(user, pwd))
	}
	// Submit the request
	res, err := r.client.Do(req)
	core.CheckErr(err, "cannot post to backend")
	switch res.StatusCode {
	case http.StatusCreated:
		i18n.Printf(i18n.INFO_PUSHED, name.String())
	case http.StatusOK:
		i18n.Printf(i18n.INFO_NOTHING_TO_PUSH)
	default:
		if res.StatusCode > 299 {
			return fmt.Errorf("failed to push, the remote server responded with status code %d: %s", res.StatusCode, res.Status)
		}
	}
	return nil
}

func (r *Api) UpdatePackageInfo(name *core.PackageName, pack *Package, user string, pwd string, https bool) error {
	b, err := json.Marshal(pack)
	if err != nil {
		return err
	}
	body := bytes.NewReader([]byte(b))
	req, err := http.NewRequest("PUT", r.packageIdURI(name.Group, name.Name, pack.Id, https), body)
	if err != nil {
		return err
	}
	req.Header.Set("accept", "application/json")
	if len(user) > 0 && len(pwd) > 0 {
		req.Header.Add("authorization", core.BasicToken(user, pwd))
	}
	// Submit the request
	res, err := r.client.Do(req)
	core.CheckErr(err, "cannot post to backend")
	// Check the response
	if res.StatusCode > 299 {
		return fmt.Errorf("failed to update package info, the remote server responded with: %s", res.Status)
	}
	return nil
}

func (r *Api) GetRepositoryInfo(group, name, user, pwd string, https bool) (*Repository, error) {
	// note: repoURI() escape the group
	req, err := http.NewRequest("GET", r.repoURI(group, name, https), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	if len(user) > 0 && len(pwd) > 0 {
		req.Header.Add("authorization", core.BasicToken(user, pwd))
	}
	// Submit the request
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusNotFound:
		// if repository is nil then the client is not talking to the proper package registry
		return nil, fmt.Errorf("\"%s\" does not conform to the Package Registry API, are you sure the package domain is correct?", r.domain)
	case http.StatusForbidden:
		return nil, fmt.Errorf("access to the registry is forbidden")
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("invalid credentials, access to the registry is unauthorised")
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// if the result body is not in JSON format is likely that the domain of package does not exist
	if !core.IsJSON(string(b)) {
		return nil, fmt.Errorf("the package was not found: its domain/group/name is likely to be incorrect")
	}
	// if not response then return an empty repository
	if len(b) == 0 {
		return &Repository{
			Repository: fmt.Sprintf("%s/%s", group, name),
			Packages:   make([]*Package, 0),
		}, nil
	}
	repo := new(Repository)
	err = json.Unmarshal(b, repo)
	return repo, err
}

func (r *Api) GetPackageInfo(group, name, id, user, pwd string, https bool) (*Package, error) {
	req, err := http.NewRequest("GET", r.packageIdURI(group, name, id, https), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	if len(user) > 0 && len(pwd) > 0 {
		req.Header.Add("authorization", core.BasicToken(user, pwd))
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
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("the remote registry responded with an internal error, check the registry logs for more information")
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if !core.IsJSON(string(b)) {
		return nil, fmt.Errorf("invalid package name: %s/%s/%s", r.domain, group, name)
	}
	pack := new(Package)
	err = json.Unmarshal(b, pack)
	return pack, err
}

func (r *Api) Download(group, name, filename, user, pwd string, https bool) (string, error) {
	req, err := http.NewRequest("GET", r.fileURI(group, name, filename, https), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("accept", "application/json")
	if len(user) > 0 && len(pwd) > 0 {
		req.Header.Add("authorization", core.BasicToken(user, pwd))
	}
	// Submit the request
	res, err := r.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// download progress bar
	// retrieve the content length to download from the http reader
	limit, err := strconv.ParseInt(res.Header.Get("Content-Length"), 0, 64)
	if err != nil {
		return "", err
	}
	// start simple new progress bar
	bar := pb.Simple.Start64(limit)
	// NOTE: must set to stdout as default is stderr to prevent downstream code to think there is an error when
	// the bar is writing its progress to the stream
	bar.SetWriter(os.Stdout)
	// adjust the prefix in the progress bar according to the file being downloaded
	if filepath.Ext(filename) == ".json" {
		bar.Set("prefix", "seal    > ")
	} else {
		bar.Set("prefix", "package > ")
	}
	// create proxy reader for the progress bar
	reader := bar.NewProxyReader(res.Body)

	switch res.StatusCode {
	case http.StatusNotFound:
		return "", fmt.Errorf("file '%s' not found in registry", filename)
	case http.StatusForbidden:
		return "", fmt.Errorf("invalid credentials, access to the registry is forbidden")
	}
	// write response to a temp file
	var b bytes.Buffer
	out := bufio.NewWriter(&b)
	_, err = io.Copy(out, reader)
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
	bar.Finish()
	return file.Name(), err
}

func (r *Api) repoURI(group, name string, https bool) string {
	scheme := "http"
	if https {
		scheme = fmt.Sprintf("%ss", scheme)
	}
	// {scheme}://{domain}/repository/{repository-group}/{repository-name}
	return fmt.Sprintf("%s://%s/repository/%s/%s", scheme, r.domain, Escape(group), name)
}

func (r *Api) packageURI(group, name string, https bool) string {
	scheme := "http"
	if https {
		scheme = fmt.Sprintf("%ss", scheme)
	}
	// {scheme}://{domain}/package/{repository-group}/{repository-name}/{tag}
	return fmt.Sprintf("%s://%s/package/%s/%s", scheme, r.domain, Escape(group), name)
}

func (r *Api) packageTagURI(group, name, tag string, https bool) string {
	return fmt.Sprintf("%s/tag/%s", r.packageURI(group, name, https), tag)
}

func (r *Api) packageIdURI(group, name, id string, https bool) string {
	// group escaped by packageURI()
	return fmt.Sprintf("%s/id/%s", r.packageURI(group, name, https), id)
}

func (r *Api) fileURI(group, name, filename string, https bool) string {
	scheme := "http"
	if https {
		scheme = fmt.Sprintf("%ss", scheme)
	}
	return fmt.Sprintf("%s://%s/file/%s/%s/%s", scheme, r.domain, Escape(group), name, filename)
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

// Escape slashes in path variables
func Escape(path string) string {
	// NOTE: not sure why but need to escape twice for the request to work properly!
	return url.PathEscape(url.PathEscape(path))
}
