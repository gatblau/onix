/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package registry

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
)

type Repository struct {
	// the package repository (name without tag)
	Repository string `json:"repository"`
	// the reference name of the package corresponding to different builds
	Packages []*Package `json:"artefacts"`
}

func (r *Repository) ToJsonBytes() ([]byte, error) {
	return json.Marshal(r)
}

// FindPackage find a package using its unique id
func (r *Repository) FindPackage(id string) *Package {
	for _, pack := range r.Packages {
		if pack.Id == id {
			return pack
		}
	}
	return nil
}

// UpdatePackage updates the specified package
func (r *Repository) UpdatePackage(a *Package) bool {
	position := -1
	for ix, pack := range r.Packages {
		if pack.Id == a.Id {
			position = ix
			break
		}
	}
	if position != -1 {
		r.Packages[position] = a
		return true
	}
	return false
}

// GetTag determines if the repository contains an package with the specified tag
func (r *Repository) GetTag(tag string) (*Package, bool) {
	for _, pack := range r.Packages {
		if pack.HasTag(tag) {
			return pack, true
		}
	}
	return nil, false
}

// Package metadata for an Artisan package
type Package struct {
	// a unique identifier for the package calculated as the checksum of the complete seal
	Id string `json:"id"`
	// the type of application in the package
	Type string `json:"type"`
	// the package actual file name
	FileRef string `json:"file_ref"`
	// the list of Tags associated with the package
	Tags []string `json:"tags"`
	// the size
	Size string `json:"size"`
	// the creation time
	Created string `json:"created"`
}

func (a *Package) String() string {
	return fmt.Sprintf("%s-%s", a.Id[0:12], a.FileRef)
}

// HasTag determines if the package has the specified tag
func (a *Package) HasTag(tag string) bool {
	for _, t := range a.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// RemoveTag removes a specified tag
// returns true if the tag was found and removed, otherwise false
func (a *Package) RemoveTag(tag string) bool {
	before := len(a.Tags)
	a.Tags = core.RemoveElement(a.Tags, tag)
	after := len(a.Tags)
	return before > after
}

func (a *Package) ToJson() (string, error) {
	bs, err := json.Marshal(a)
	return base64.StdEncoding.EncodeToString(bs), err
}
