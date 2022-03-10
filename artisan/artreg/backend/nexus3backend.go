package backend

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
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
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/registry"
	"github.com/gatblau/onix/oxlib/httpserver"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Nexus3Backend a Nexus3 implementation of a remote registry
type Nexus3Backend struct {
	domain string
	client *http.Client
	tmp    string
}

func NewNexus3Backend(domain string) Backend {
	core.TmpExists()
	return &Nexus3Backend{
		domain: domain,
		tmp:    core.TmpPath(),
		client: &http.Client{
			Timeout: time.Second * 120,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

func (r *Nexus3Backend) Name() string {
	return fmt.Sprintf("NEXUS3 @ %s\n", r.domain)
}

func (r *Nexus3Backend) Download(repoGroup, repoName, fileName, user, pwd string) (*os.File, error) {
	// get the file download URI
	downloadURI := r.fileDownloadURI(repoGroup, repoName, fileName)
	req, err := http.NewRequest("GET", downloadURI, nil)
	if err != nil {
		return nil, err
	}
	if len(user) > 0 && len(pwd) > 0 {
		req.Header.Add("authorization", httpserver.BasicToken(user, pwd))
	}
	// Submit the request
	res, err := r.client.Do(req)
	// must close the body
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	out := bufio.NewWriter(&b)
	// Write the body to file
	_, err = io.Copy(out, res.Body)
	if err != nil {
		return nil, err
	}
	err = out.Flush()
	if err != nil {
		return nil, err
	}
	file, err := os.Create(filepath.Join(r.tmp, fileName))
	if err != nil {
		return nil, err
	}
	_, err = file.Write(b.Bytes())
	if err != nil {
		return nil, err
	}
	file.Close()
	f, err := os.Open(file.Name())
	if err != nil {
		return nil, err
	}
	err = os.Remove(file.Name())
	if err != nil {
		return nil, err
	}
	return f, err
}

func (r *Nexus3Backend) UpdatePackageInfo(group, name string, packageInfo *registry.Package, user string, pwd string) error {
	// get the repository info
	repo, err := r.GetRepositoryInfo(group, name, user, pwd)
	if err != nil {
		return err
	}
	// update the repository
	updated := repo.UpdatePackage(packageInfo)
	if !updated {
		return fmt.Errorf("package not found in remote repository, no update was made")
	}
	// turn the repository into a file to upload
	// create a repository file
	repoFile, err := core.ToJsonFile(repo)
	if err != nil {
		return err
	}
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	err = r.addField(writer, "raw.directory", fmt.Sprintf("%s/%s", group, name))
	if err != nil {
		return err
	}
	err = r.addField(writer, "raw.asset1.filename", "repository.json")
	if err != nil {
		return err
	}
	err = r.addFile(writer, "raw.asset1", "repository.json", repoFile)
	if err != nil {
		return err
	}
	writer.Close()
	return r.postMultipart(b, writer, user, pwd)
}

// UploadPackage upload a new package
func (r *Nexus3Backend) UploadPackage(group, name, packageRef string, zipfile multipart.File, jsonFile multipart.File, repo multipart.File, user string, pwd string) error {
	// ensure files are properly closed
	defer zipfile.Close()
	defer jsonFile.Close()
	defer repo.Close()

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	err := r.addField(writer, "raw.directory", fmt.Sprintf("%s/%s", group, name))
	if err != nil {
		return err
	}
	err = r.addField(writer, "raw.asset1.filename", fmt.Sprintf("%s.json", packageRef))
	if err != nil {
		return err
	}
	err = r.addFile(writer, "raw.asset1", fmt.Sprintf("%s.json", packageRef), jsonFile)
	if err != nil {
		return err
	}
	err = r.addField(writer, "raw.asset2.filename", fmt.Sprintf("%s.zip", packageRef))
	if err != nil {
		return err
	}
	err = r.addFile(writer, "raw.asset2", fmt.Sprintf("%s.zip", packageRef), zipfile)
	if err != nil {
		return err
	}
	err = r.addField(writer, "raw.asset3.filename", "repository.json")
	if err != nil {
		return err
	}
	err = r.addFile(writer, "raw.asset3", "repository.json", repo)
	if err != nil {
		return err
	}
	// don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	writer.Close()
	return r.postMultipart(b, writer, user, pwd)
}

// DeletePackage delete a specific package
func (r *Nexus3Backend) DeletePackage(group, name, packageRef, user, pwd string) error {
	assetList, err := r.getAssets(user, pwd)
	if err != nil {
		return fmt.Errorf("cannot get list of assets in Nexus: %s", err)
	}
	var ids []string
	for _, asset := range assetList.Items {
		if strings.HasPrefix(asset.Path, fmt.Sprintf("%s/%s/%s", group, name, packageRef)) {
			ids = append(ids, asset.ID)
		}
	}
	for _, id := range ids {
		if err = r.deleteComponent(id, user, pwd); err != nil {
			return err
		}
	}
	return nil
}

func (r *Nexus3Backend) postMultipart(b bytes.Buffer, writer *multipart.Writer, user string, pwd string) error {
	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", r.componentsURI(), &b)
	if err != nil {
		return err
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("accept", "application/json")
	if len(user) > 0 && len(pwd) > 0 {
		req.Header.Add("authorization", httpserver.BasicToken(user, pwd))
	}
	// Submit the request
	res, err := r.client.Do(req)
	// must close the body
	if res.Body != nil {
		res.Body.Close()
	}
	if err != nil {
		return err
	}
	// Check the response
	if res.StatusCode > 299 {
		return fmt.Errorf("failed to push, the remote server responded with status code %d: %s", res.StatusCode, res.Status)
	}
	return nil
}

func (r *Nexus3Backend) deleteComponent(id, user, pwd string) error {
	req, err := http.NewRequest("DELETE", r.componentsWithIdURI(id), nil)
	if err != nil {
		return err
	}
	req.Header.Set("accept", "application/json")
	if len(user) > 0 && len(pwd) > 0 {
		req.Header.Add("authorization", httpserver.BasicToken(user, pwd))
	}
	res, err := r.client.Do(req)
	// must close the body
	if res.Body != nil {
		res.Body.Close()
	}
	if err != nil {
		return err
	}
	// success
	switch res.StatusCode {
	case 204:
		// success
		return nil
	case 403:
		return fmt.Errorf("insufficient permissions to delete component")
	case 404:
		return fmt.Errorf("component not found")
	case 422:
		return fmt.Errorf("malformed id")
	default:
		return fmt.Errorf("unknown response status: %s", res.Status)
	}
}

func (r *Nexus3Backend) GetRepositoryInfo(group, name, user, pwd string) (*registry.Repository, error) {
	// NOTE: the commented out validation below is not working in some cases, as sometime the metadata in Nexus seems to get corrupted
	// and although the files are uploaded successfully the components meta data is not being updated accordingly

	// // check the repository.json file exists in nexus
	// components, err := r.getComponents(user, pwd)
	// if err != nil {
	// 	return nil, err
	// }
	// // if the repository file does not exist
	// if !components.containsFile(group, name, "repository.json") {
	// 	// returns an empty repository
	// 	return &Repository{
	// 		Repository: fmt.Sprintf("%s/%s", group, name),
	// 		Artefacts:  make([]*Package, 0),
	// 	}, nil
	// }

	// otherwise fetches the content and returns it
	b, err := r.getFile(group, name, "repository.json", user, pwd)
	if err != nil {
		return nil, err
	}
	// if the file is not in JSON format then
	if !core.IsJSON(string(b)) {
		// assume file not found (404 HTML page)
		// returns an empty repository
		return &registry.Repository{
			Repository: fmt.Sprintf("%s/%s", group, name),
			Packages:   make([]*registry.Package, 0),
		}, nil
	}
	repo := new(registry.Repository)
	err = json.Unmarshal(b, repo)
	return repo, err
}

func (r *Nexus3Backend) GetAllRepositoryInfo(user, pwd string) ([]*registry.Repository, error) {
	assets, err := r.getAssets(user, pwd)
	if err != nil {
		return nil, fmt.Errorf("cannot get list of assets in Nexus: %s", err)
	}
	infos := make([]*registry.Repository, 0)
	// loop through the assets
	for _, asset := range assets.Items {
		// find repository descriptors
		if strings.HasSuffix(asset.Path, "repository.json") {
			// get the group
			path := asset.Path[0 : len(asset.Path)-len("repository.json")-1]
			group := path[:strings.LastIndex(path, "/")]
			name := path[strings.LastIndex(path, "/")+1:]
			repositoryInfo, err := r.GetRepositoryInfo(group, name, user, pwd)
			if err != nil {
				return nil, fmt.Errorf("cannot get information for repository %s: %s", asset.Path, err)
			}
			infos = append(infos, repositoryInfo)
		}
	}
	return infos, nil
}

func (r *Nexus3Backend) GetPackageInfo(group, name, id, user, pwd string) (*registry.Package, error) {
	repo, err := r.GetRepositoryInfo(group, name, user, pwd)
	if err != nil {
		return nil, err
	}
	if repo != nil {
		return repo.FindPackage(id), nil
	}
	return nil, nil
}

func (r *Nexus3Backend) GetPackageManifest(group, name, tag, user, pwd string) (*data.Manifest, error) {
	repo, err := r.GetRepositoryInfo(group, name, user, pwd)
	if err != nil {
		return nil, err
	}
	if len(tag) == 0 {
		tag = "latest"
	}
	ref, err := getFileRef(repo, tag)
	if err != nil {
		return nil, err
	}
	manifestFile, err := r.Download(group, name, fmt.Sprintf("%s.json", ref), user, pwd)
	if err != nil {
		return nil, err
	}
	defer manifestFile.Close()
	bytes, err := ioutil.ReadAll(manifestFile)
	if err != nil {
		return nil, err
	}
	var seal data.Seal
	err = json.Unmarshal(bytes, &seal)
	if err != nil {
		return nil, err
	}
	return seal.Manifest, nil
}

func getFileRef(repo *registry.Repository, tag string) (string, error) {
	for _, pack := range repo.Packages {
		for _, t := range pack.Tags {
			if t == tag {
				return pack.FileRef, nil
			}
		}
	}
	return "", fmt.Errorf("package not found using tag %s\n", tag)
}

// add a field to a multipart form
func (r *Nexus3Backend) addField(writer *multipart.Writer, fieldName, fieldValue string) error {
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
func (r *Nexus3Backend) addFile(writer *multipart.Writer, fieldName, fileName string, file multipart.File) error {
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

func (r *Nexus3Backend) componentsURI() string {
	return fmt.Sprintf("%s/service/rest/v1/components?repository=artisan", r.domain)
}

func (r *Nexus3Backend) componentsWithIdURI(id string) string {
	return fmt.Sprintf("%s/service/rest/v1/components/%s", r.domain, id)
}

func (r *Nexus3Backend) assetsURI(continuationToken string) string {
	if len(continuationToken) > 0 {
		return fmt.Sprintf("%s/service/rest/v1/assets?repository=artisan&continuationToken=%s", r.domain, continuationToken)
	}
	return fmt.Sprintf("%s/service/rest/v1/assets?repository=artisan", r.domain)
}

func (r *Nexus3Backend) fileDownloadURI(group, name, filename string) string {
	return fmt.Sprintf("%s/repository/artisan/%s/%s/%s", r.domain, group, name, filename)
}

func (r *Nexus3Backend) downloadURI(repoGroup, repoName, filename string) string {
	return fmt.Sprintf("%s/repository/artisan/%s/%s/%s", r.domain, repoGroup, repoName, filename)
}

// get the content of a file from Nexus
func (r *Nexus3Backend) getFile(repoGroup, repoName, filename, user, pwd string) ([]byte, error) {
	req, err := http.NewRequest("GET", r.downloadURI(repoGroup, repoName, filename), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	if len(user) > 0 && len(pwd) > 0 {
		req.Header.Add("authorization", httpserver.BasicToken(user, pwd))
	}
	// Submit the request
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (r *Nexus3Backend) getAssetsPages(user, pwd, continuationToken string) (*assets, error) {
	result, err := r.getMeta(user, pwd, r.assetsURI(continuationToken), new(assets))
	if err != nil {
		return nil, err
	}
	return result.(*assets), nil
}

func (r *Nexus3Backend) getAssets(user, pwd string) (*assets, error) {
	// first time no continuation token no required
	token := ""
	a, err := r.getAssetsPages(user, pwd, token)
	if err != nil {
		return nil, err
	}
	// if a continuation token exists then there are more pages
	if a.Continuationtoken != nil {
		// get the token
		token = a.Continuationtoken.(string)
	}
	// starts a loop to fetch all pages
	for len(token) > 0 {
		// query the next page using the continuation token
		a2, err2 := r.getAssetsPages(user, pwd, token)
		if err2 != nil {
			return nil, err
		}
		// append the items in second page to the items in first page
		a.Items = append(a.Items, a2.Items...)
		// if there is a continuation token
		if a2.Continuationtoken != nil {
			// set the token to the value of the continuation token in the second page
			// to allow the loop to fetch the next page
			token = a2.Continuationtoken.(string)
		} else {
			// otherwise set the token to empty to finish the loop
			token = ""
		}
	}
	return a, nil
}

func (r *Nexus3Backend) getComponents(user, pwd string) (*components, error) {
	result, err := r.getMeta(user, pwd, r.componentsURI(), new(components))
	if err != nil {
		return nil, err
	}
	return result.(*components), nil
}

// get metadata from Nexus
func (r *Nexus3Backend) getMeta(user, pwd, uri string, result interface{}) (interface{}, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	if len(user) > 0 && len(pwd) > 0 {
		req.Header.Add("authorization", httpserver.BasicToken(user, pwd))
	}
	// Submit the request
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// if the result body is not in JSON format is likely that the domain of the backend does not exist
	if !core.IsJSON(string(b)) {
		return nil, fmt.Errorf("the response body was in an incorrect format, which suggests \nthe backend URI '%s' is not correct, \nor the server responsed with a bogus payload", r.domain)
	}
	err = json.Unmarshal(b, result)
	return result, err
}

// Nexus3 assets in a Repository
type assets struct {
	Items []struct {
		Downloadurl string `json:"downloadUrl"`
		Path        string `json:"path"`
		ID          string `json:"id"`
		Repository  string `json:"repository"`
		Format      string `json:"format"`
		Checksum    struct {
			Sha1   string `json:"sha1"`
			Sha512 string `json:"sha512"`
			Sha256 string `json:"sha256"`
			Md5    string `json:"md5"`
		} `json:"checksum"`
		Contenttype    string      `json:"contentType"`
		Lastmodified   time.Time   `json:"lastModified"`
		Blobcreated    time.Time   `json:"blobCreated"`
		Lastdownloaded interface{} `json:"lastDownloaded"`
	} `json:"items"`
	Continuationtoken interface{} `json:"continuationToken"`
}

// Nexus3 components in a Repository
type components struct {
	Items []struct {
		ID         string      `json:"id"`
		Repository string      `json:"repository"`
		Format     string      `json:"format"`
		Group      string      `json:"group"`
		Name       string      `json:"name"`
		Version    interface{} `json:"version"`
		Assets     []struct {
			DownloadURL string `json:"downloadUrl"`
			Path        string `json:"path"`
			ID          string `json:"id"`
			Repository  string `json:"repository"`
			Format      string `json:"format"`
			Checksum    struct {
				Sha1   string `json:"sha1"`
				Sha512 string `json:"sha512"`
				Sha256 string `json:"sha256"`
				Md5    string `json:"md5"`
			} `json:"checksum"`
		} `json:"assets"`
	} `json:"items"`
	ContinuationToken interface{} `json:"continuationToken"`
}

// determines if the file is in the nexus Repository
func (c *components) containsFile(group, name, filename string) bool {
	for _, item := range c.Items {
		if item.Name == fmt.Sprintf("%s/%s/%s", group, name, filename) {
			return true
		}
	}
	return false
}
