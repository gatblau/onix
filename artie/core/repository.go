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
)

// a repository in the localRepo
type repository struct {
	// the reference name of the artefact corresponding to different builds
	Artefacts []*artefact `json:"artefacts"`
}

// return the artefact that matches the specified name tag or nil if not found in the repository
func (r *repository) artefact(artefactName string) (*artefact, bool) {
	for _, artefact := range r.Artefacts {
		for _, tag := range artefact.Tags {
			if tag == artefactName {
				return artefact, true
			}
		}
	}
	return nil, false
}

type artefact struct {
	// a unique identifier for the artefact calculated as the checksum of the complete seal
	Id string `json:"id"`
	// the artefact actual file name
	File string `json:"file"`
	// the list of Tags associated with the artefact
	Tags []string `json:"tags"`
}

// create a localRepo management structure
func NewRepository() *repository {
	r := &repository{
		Artefacts: []*artefact{},
	}
	// load localRepo
	r.load()
	return r
}

// the local path to the local repository
func (r *repository) path() string {
	return fmt.Sprintf("%s/.%s", homeDir(), cliName)
}

// return the repository full file name
func (r *repository) file() string {
	return fmt.Sprintf("%s/repository.json", r.path())
}

// save the state of the repository
func (r *repository) save() {
	regBytes := toJsonBytes(r)
	err := ioutil.WriteFile(r.file(), regBytes, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

// load the content of the repository
func (r *repository) load() {
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

// add the artefact and seal to the repository
func (r *repository) add(filename, artefactName string, s *seal) {
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
	// work out the tag to use
	tagMarks := strings.Count(artefactName, ":")
	switch tagMarks {
	// the artefact name does not have a tag
	case 0:
		// add the "latest" tag
		artefactName = fmt.Sprintf("%s:latest", artefactName)
	// the artefact name does contain a tag
	case 1:
		// all good do nothing
	default:
		// any other number is an error
		log.Fatal(errors.New(fmt.Sprintf("the package name-tag cannot contain more than 1 colon. found %d colons", tagMarks)))
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
	if artef, exists := r.artefact(artefactName); exists {
		// then it has to untag it, leaving a dangling artefact
		artef.Tags = removeElement(artef.Tags, artefactName)
	}
	// creates a new artefact
	artefacts := append(r.Artefacts, &artefact{
		Id:   artefactId(s),
		File: fmt.Sprintf("%s.zip", basenameNoExt),
		Tags: []string{artefactName},
	})
	r.Artefacts = artefacts
	// persist the changes
	r.save()
}
