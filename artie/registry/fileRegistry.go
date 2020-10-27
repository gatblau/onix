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
func (r *FileRegistry) GetArtefactsByRepo(repoName string) ([]*artefact, bool) {
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

// return the artefact that matches the specified name:tag or nil if not found in the FileRegistry
func (r *FileRegistry) GetArtefactByName(artefactName string) *artefact {
	for _, artefact := range r.Artefacts {
		for _, tag := range artefact.Tags {
			if (tag != "latest" && fmt.Sprintf("%s:%s", artefact.Repository, tag) == artefactName) ||
				(tag == "latest" && artefact.Repository == artefactName) {
				return artefact
			}
		}
	}
	return nil
}

type artefact struct {
	// a unique identifier for the artefact calculated as the checksum of the complete seal
	Id string `json:"id"`
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
	return fmt.Sprintf("%s/.%s", core.HomeDir(), core.CliName)
}

// return the FileRegistry full file name
func (r *FileRegistry) file() string {
	return fmt.Sprintf("%s/repository.json", r.Path())
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
func (r *FileRegistry) Add(filename string, artefactName core.Named, s *core.Seal) {
	var repo core.NamedRepository
	if namedRepo, ok := artefactName.(core.NamedRepository); ok {
		repo = namedRepo
	} else {
		log.Fatal(errors.New("artefact name not supported"))
	}
	// tag by default is latest
	tag := "latest"
	// if a tag has been provided then us it
	if tagged, ok := artefactName.(core.Tagged); ok {
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
	err := os.Rename(filename, fmt.Sprintf("%s/%s", r.Path(), basename))
	if err != nil {
		log.Fatal(err)
	}
	// now move the seal file to the localRepo folder
	err = os.Rename(fmt.Sprintf("%s/%s.json", basenameDir, basenameNoExt), fmt.Sprintf("%s/%s.json", r.Path(), basenameNoExt))
	if err != nil {
		log.Fatal(err)
	}
	// if the artefact already exists in the repository
	if artefs, exists := r.GetArtefactsByRepo(artefactName.Name()); exists {
		// then it has to untag it, leaving a dangling artefact
		for _, artef := range artefs {
			artef.Tags = core.RemoveElement(artef.Tags, tag)
		}
	}
	// creates a new artefact
	artefacts := append(r.Artefacts, &artefact{
		Id:         core.ArtefactId(s),
		Repository: repo.Name(),
		FileRef:    basenameNoExt,
		Tags:       []string{tag},
		Size:       s.Manifest.Size,
		Created:    s.Manifest.Time,
	})
	r.Artefacts = artefacts
	// persist the changes
	r.save()
}

// List packages to stdout
func (r *FileRegistry) List() {
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

func (r *FileRegistry) Push(name core.Named, remote Remote, credentials string) {
	// fetch the artefact info from the local registry
	artie := r.GetArtefactByName(name.String())
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
	remote.UploadArtefact(client, name, r.Path(), artie.FileRef, credentials)
}

func (r *FileRegistry) Pull(name core.Named, remote Remote) {
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
