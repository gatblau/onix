/*
  Onix Config Manager - Artie
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

// a location in the local machine where to cache build artefacts
type LocalRegistry struct {
	// the reference name of the artefact corresponding to different builds
	Artefacts []*artefact `json:"artefacts"`
}

// return the artefact that matches the specified name tag or nil if not found in the LocalRegistry
func (r *LocalRegistry) GetArtefactsByRepo(repoName string) ([]*artefact, bool) {
	var artefacts = make([]*artefact, 0)
	for _, artefact := range r.Artefacts {
		if artefact.Repository == repoName {
			artefacts = append(artefacts, artefact)
		}
	}
	if len(artefacts) > 0 {
		return artefacts, true
	}
	return nil, false
}

type artefact struct {
	// a unique identifier for the artefact calculated as the checksum of the complete seal
	Id string `json:"id"`
	// the artefact repository (name without without tag)
	Repository string `json:"repository"`
	// the artefact actual file name
	File string `json:"file"`
	// the list of Tags associated with the artefact
	Tags []string `json:"tags"`
	// the size
	Size string `json:"size"`
	// the creation time
	Created string `json:"created"`
}

// create a localRepo management structure
func NewRepository() *LocalRegistry {
	r := &LocalRegistry{
		Artefacts: []*artefact{},
	}
	// load localRepo
	r.load()
	return r
}

// the local path to the local LocalRegistry
func (r *LocalRegistry) path() string {
	return fmt.Sprintf("%s/.%s", homeDir(), cliName)
}

// return the LocalRegistry full file name
func (r *LocalRegistry) file() string {
	return fmt.Sprintf("%s/repository.json", r.path())
}

// save the state of the LocalRegistry
func (r *LocalRegistry) save() {
	regBytes := toJsonBytes(r)
	err := ioutil.WriteFile(r.file(), regBytes, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

// load the content of the LocalRegistry
func (r *LocalRegistry) load() {
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

// add the artefact and seal to the LocalRegistry
func (r *LocalRegistry) add(filename string, artefactName Named, s *seal) {
	var repo namedRepository
	if namedRepo, ok := artefactName.(namedRepository); ok {
		repo = namedRepo
	} else {
		log.Fatal(errors.New("artefact name not supported"))
	}
	// tag by default is latest
	tag := "latest"
	// if a tag has been provided then us it
	if tagged, ok := artefactName.(Tagged); ok {
		tag = tagged.Tag()
	}
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
	err := os.Rename(filename, fmt.Sprintf("%s/%s", r.path(), basename))
	if err != nil {
		log.Fatal(err)
	}
	// now move the seal file to the localRepo folder
	err = os.Rename(fmt.Sprintf("%s/%s.json", basenameDir, basenameNoExt), fmt.Sprintf("%s/%s.json", r.path(), basenameNoExt))
	if err != nil {
		log.Fatal(err)
	}
	// if the artefact already exists in the repository
	if artefs, exists := r.GetArtefactsByRepo(artefactName.Name()); exists {
		// then it has to untag it, leaving a dangling artefact
		for _, artef := range artefs {
			artef.Tags = removeElement(artef.Tags, tag)
		}
	}
	// creates a new artefact
	artefacts := append(r.Artefacts, &artefact{
		Id:         artefactId(s),
		Repository: repo.Name(),
		File:       fmt.Sprintf("%s.zip", basenameNoExt),
		Tags:       []string{tag},
		Size:       s.Manifest.Size,
		Created:    s.Manifest.Time,
	})
	r.Artefacts = artefacts
	// persist the changes
	r.save()
}

// List packages to stdout
func (r *LocalRegistry) List() {
	// get a table writer for the stdout
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 10, ' ', 0)
	// print the header row
	fmt.Fprintln(w, "REPOSITORY\tTAG\tARTEFACT ID\tCREATED\tSIZE")
	// repository, tag, artefact id, created, size
	for _, a := range r.Artefacts {
		// if the artefact is dangling (no tags)
		if len(a.Tags) == 0 {
			fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%s\t%s",
				a.Repository,
				"<none>",
				a.Id[:12],
				toElapsedLabel(a.Created),
				a.Size),
			)
		}
		for _, tag := range a.Tags {
			fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%s\t%s",
				a.Repository,
				tag,
				a.Id[:12],
				toElapsedLabel(a.Created),
				a.Size),
			)
		}
	}
	w.Flush()
}
