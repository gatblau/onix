/*
  Onix Config Manager - Artie
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package registry

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gatblau/onix/artie/core"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"
)

// the default local registry implemented as a file system
type FileRegistry struct {
	// the reference name of the artefact corresponding to different builds
	Artefacts []*artefact `json:"artefacts"`
}

// return all the artefacts within the same repository
func (r *FileRegistry) GetArtefactsByName(name *core.ArtieName) ([]*artefact, bool) {
	var artefacts = make([]*artefact, 0)
	for _, artefact := range r.Artefacts {
		if artefact.Repository == name.FullyQualifiedRepository() {
			artefacts = append(artefacts, artefact)
		}
	}
	if len(artefacts) > 0 {
		return artefacts, true
	}
	return nil, false
}

// return the artefact that matches the specified:
// - domain/repo/name:tag or
// - artefact id substring or
// nil if not found in the FileRegistry
func (r *FileRegistry) GetArtefact(name *core.ArtieName) *artefact {
	for _, artefact := range r.Artefacts {
		// try and match against the artefact ID substring
		if strings.Contains(artefact.Id, name.Name) {
			return artefact
		}
		// if no luck use the tags
		for _, tag := range artefact.Tags {
			// try and match against the full URI
			if fmt.Sprintf("%s:%s", artefact.Repository, tag) == name.String() {
				return artefact
			}
		}
	}
	return nil
}

type artefact struct {
	// a unique identifier for the artefact calculated as the checksum of the complete seal
	Id string `json:"id"`
	// the type of application in the artefact
	Type string `json:"type"`
	// the artefact repository (name without without tag)
	Repository string `json:"repository"`
	// the artefact actual file name
	FileRef string `json:"file_ref"`
	// the list of Tags associated with the artefact
	Tags []string `json:"tags"`
	// the size
	Size string `json:"size"`
	// the creation time
	Created string `json:"created"`
}

// create a localRepo management structure
func NewFileRegistry() *FileRegistry {
	r := &FileRegistry{
		Artefacts: []*artefact{},
	}
	// load local registry
	r.load()
	return r
}

// the local Path to the local FileRegistry
func (r *FileRegistry) Path() string {
	return core.RegistryPath()
}

// return the FileRegistry full file name
func (r *FileRegistry) file() string {
	return filepath.Join(r.Path(), "repository.json")
}

// save the state of the FileRegistry
func (r *FileRegistry) save() {
	regBytes := core.ToJsonBytes(r)
	err := ioutil.WriteFile(r.file(), regBytes, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

// load the content of the FileRegistry
func (r *FileRegistry) load() {
	// check if localRepo file exist
	_, err := os.Stat(r.file())
	if err != nil {
		// then assume localRepo.json is not there: try and create it
		r.save()
	} else {
		regBytes, err := ioutil.ReadFile(r.file())
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(regBytes, r)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Add the artefact and seal to the FileRegistry
func (r *FileRegistry) Add(filename string, name *core.ArtieName, s *core.Seal) {
	// gets the full base name (with extension)
	basename := filepath.Base(filename)
	// gets the basename directory only
	basenameDir := filepath.Dir(filename)
	// gets the base name extension
	basenameExt := path.Ext(basename)
	// gets the base name without extension
	basenameNoExt := strings.TrimSuffix(basename, path.Ext(basename))
	// if the file to add is not a zip file
	if basenameExt != ".zip" {
		log.Fatal(errors.New(fmt.Sprintf("the localRepo can only accept zip files, the extension provided was %s", basenameExt)))
	}
	// move the zip file to the localRepo folder
	// err := os.Rename(filename, fmt.Sprintf("%s/%s", r.Path(), basename))
	err := RenameFile(filename, fmt.Sprintf("%s/%s", r.Path(), basename), false)
	if err != nil {
		log.Fatal(err)
	}
	// now move the seal file to the localRepo folder
	// err = os.Rename(fmt.Sprintf("%s/%s.json", basenameDir, basenameNoExt), fmt.Sprintf("%s/%s.json", r.Path(), basenameNoExt))
	err = RenameFile(fmt.Sprintf("%s/%s.json", basenameDir, basenameNoExt), fmt.Sprintf("%s/%s.json", r.Path(), basenameNoExt), false)
	if err != nil {
		log.Fatal(err)
	}
	// untag artefact artefact (if any)
	r.unTag(name, name.Tag)
	// creates a new artefact
	artefacts := append(r.Artefacts, &artefact{
		Id:         core.ArtefactId(s),
		Type:       s.Manifest.Type,
		Repository: name.FullyQualifiedName(),
		FileRef:    basenameNoExt,
		Tags:       []string{name.Tag},
		Size:       s.Manifest.Size,
		Created:    s.Manifest.Time,
	})
	r.Artefacts = artefacts
	// persist the changes
	r.save()
}

// removeArtefactByRepository the specified artefact from the slice
func (r *FileRegistry) removeArtefactByRepository(a []*artefact, name *core.ArtieName) []*artefact {
	i := -1
	// find the value to remove
	for ix := 0; ix < len(a); ix++ {
		if a[ix].Repository == name.FullyQualifiedName() {
			i = ix
			break
		}
	}
	if i == -1 {
		return a
	}
	// Remove the element at index i from a.
	a[i] = a[len(a)-1] // Copy last element to index i.
	a[len(a)-1] = nil  // Erase last element (write zero value).
	a = a[:len(a)-1]   // Truncate slice.
	return a
}

func (r *FileRegistry) removeArtefactById(a []*artefact, id string) []*artefact {
	i := -1
	// find the value to remove
	for ix := 0; ix < len(a); ix++ {
		if strings.Contains(a[ix].Id, id) {
			i = ix
			break
		}
	}
	if i == -1 {
		return a
	}
	// Remove the element at index i from a.
	a[i] = a[len(a)-1] // Copy last element to index i.
	a[len(a)-1] = nil  // Erase last element (write zero value).
	a = a[:len(a)-1]   // Truncate slice.
	return a
}

// remove a given tag from an artefact
func (r *FileRegistry) unTag(name *core.ArtieName, tag string) {
	artie := r.GetArtefact(name)
	if artie != nil {
		artie.Tags = core.RemoveElement(artie.Tags, tag)
	}
}

// remove all tags from the specified artefact
func (r *FileRegistry) unTagAll(name *core.ArtieName) {
	if artefs, exists := r.GetArtefactsByName(name); exists {
		// then it has to untag it, leaving a dangling artefact
		for _, artef := range artefs {
			for _, tag := range artef.Tags {
				artef.Tags = core.RemoveElement(artef.Tags, tag)
			}
		}
	}
	// persist changes
	r.save()
}

// List artefacts to stdout
func (r *FileRegistry) List() {
	// get a table writer for the stdout
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 12, ' ', 0)
	// print the header row
	fmt.Fprintln(w, "REPOSITORY\tTAG\tARTEFACT ID\tARTEFACT TYPE\tCREATED\tSIZE")
	// repository, tag, artefact id, created, size
	for _, a := range r.Artefacts {
		// if the artefact is dangling (no tags)
		if len(a.Tags) == 0 {
			fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s",
				a.Repository,
				"<none>",
				a.Id[7:19],
				a.Type,
				toElapsedLabel(a.Created),
				a.Size),
			)
		}
		for _, tag := range a.Tags {
			fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s",
				a.Repository,
				tag,
				a.Id[7:19],
				a.Type,
				toElapsedLabel(a.Created),
				a.Size),
			)
		}
	}
	w.Flush()
}

// list (quiet) artefact IDs only
func (r *FileRegistry) ListQ() {
	// get a table writer for the stdout
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 10, ' ', 0)
	// repository, tag, artefact id, created, size
	for _, a := range r.Artefacts {
		fmt.Fprintln(w, fmt.Sprintf("%s", a.Id[7:19]))
	}
	w.Flush()
}

func (r *FileRegistry) Push(name *core.ArtieName, remote Remote, credentials string) {
	// fetch the artefact info from the local registry
	artie := r.GetArtefact(name)
	if artie == nil {
		log.Fatal(errors.New(fmt.Sprintf("artefact %s not found in the local registry", name)))
	}
	// set up an http client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	// execute the upload
	err := remote.UploadArtefact(client, name, r.Path(), artie.FileRef, credentials)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("pushed %s\n", name.String())
}

func (r *FileRegistry) Remove(names []*core.ArtieName) {
	for _, name := range names {
		// try and get the artefact by complete URI or id ref
		artie := r.GetArtefact(name)
		if artie == nil {
			fmt.Printf("name %s not found\n", name.Name)
			continue
		}
		// try to remove it using full name
		// remove the specified tag
		length := len(artie.Tags)
		artie.Tags = core.RemoveElement(artie.Tags, name.Tag)
		// if the tag was successfully deleted
		if len(artie.Tags) < length {
			// if there are no tags left at the end then remove the artefact
			if len(artie.Tags) == 0 {
				r.Artefacts = r.removeArtefactByRepository(r.Artefacts, name)
				r.removeFiles(artie)
			}
			// persist changes
			r.save()
			log.Print(artie.Id)
		} else {
			// attempt to remove by Id (stored in the Name)
			r.Artefacts = r.removeArtefactById(r.Artefacts, name.Name)
			r.removeFiles(artie)
			r.save()
			log.Print(artie.Id)
		}
	}
}

// remove the files associated with an artefact
func (r *FileRegistry) removeFiles(artie *artefact) {
	// remove the zip file
	err := os.Remove(fmt.Sprintf("%s/%s.zip", r.Path(), artie.FileRef))
	if err != nil {
		log.Fatal(err)
	}
	// remove the json file
	err = os.Remove(fmt.Sprintf("%s/%s.json", r.Path(), artie.FileRef))
	if err != nil {
		log.Fatal(err)
	}
}

func (r *FileRegistry) Pull(name *core.ArtieName, remote Remote) {
}

// returns the elapsed time until now in human friendly format
func toElapsedLabel(rfc850time string) string {
	created, err := time.Parse(time.RFC850, rfc850time)
	if err != nil {
		log.Fatal(err)
	}
	elapsed := time.Since(created)
	seconds := elapsed.Seconds()
	minutes := elapsed.Minutes()
	hours := elapsed.Hours()
	days := hours / 24
	weeks := days / 7
	months := weeks / 4
	years := months / 12

	if math.Trunc(years) > 0 {
		return fmt.Sprintf("%d %s ago", int64(years), plural(int64(years), "year"))
	} else if math.Trunc(months) > 0 {
		return fmt.Sprintf("%d %s ago", int64(months), plural(int64(months), "month"))
	} else if math.Trunc(weeks) > 0 {
		return fmt.Sprintf("%d %s ago", int64(weeks), plural(int64(weeks), "week"))
	} else if math.Trunc(days) > 0 {
		return fmt.Sprintf("%d %s ago", int64(days), plural(int64(days), "day"))
	} else if math.Trunc(hours) > 0 {
		return fmt.Sprintf("%d %s ago", int64(hours), plural(int64(hours), "hour"))
	} else if math.Trunc(minutes) > 0 {
		return fmt.Sprintf("%d %s ago", int64(minutes), plural(int64(minutes), "minute"))
	}
	return fmt.Sprintf("%d %s ago", int64(seconds), plural(int64(seconds), "second"))
}

// turn label into plural if value is greater than one
func plural(value int64, label string) string {
	if value > 1 {
		return fmt.Sprintf("%ss", label)
	}
	return label
}

// the fully qualified name of the json Seal file in the local localReg
func (r *FileRegistry) regDirJsonFilename(uniqueIdName string) string {
	return fmt.Sprintf("%s/%s.json", r.Path(), uniqueIdName)
}

// the fully qualified name of the zip file in the local localReg
func (r *FileRegistry) regDirZipFilename(uniqueIdName string) string {
	return fmt.Sprintf("%s/%s.zip", r.Path(), uniqueIdName)
}
