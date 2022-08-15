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
    "strings"
)

type Repository struct {
    // the package repository (name without tag)
    Repository string `json:"repository"`
    // the reference name of the package corresponding to different builds
    Packages []*Package `json:"artefacts"`
}

// GetFileRef get the file reference for a specific tag
func (r *Repository) GetFileRef(tag string) (string, error) {
    for _, pack := range r.Packages {
        for _, t := range pack.Tags {
            if t == tag {
                return pack.FileRef, nil
            }
        }
    }
    return "", fmt.Errorf("package not found using tag %s\n", tag)
}

func (r *Repository) ToJsonBytes() ([]byte, error) {
    return json.Marshal(r)
}

func (r *Repository) IsDangling() bool {
    return strings.Contains(r.Repository, "<none>")
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

// UpsertPackage updates the specified package information if exists, otherwise adds it to the repository
// returns true if the package info was updated or false if the package info was added
func (r *Repository) UpsertPackage(a *Package) (updated bool) {
    position := -1
    // try and find the package using its unique Id
    for ix, pack := range r.Packages {
        if pack.Id == a.Id {
            position = ix
            break
        }
    }
    // if the package was found
    if position != -1 {
        // replaces the package info
        r.Packages[position] = a
        return true
    } else {
        // adds the package to the list of packages
        r.Packages = append(r.Packages, a)
    }
    return false
}

// GetTag determines if the repository contains a package with the specified tag
func (r *Repository) GetTag(tag string) (*Package, bool) {
    for _, pack := range r.Packages {
        if pack.HasTag(tag) {
            return pack, true
        }
    }
    return nil, false
}

func (r *Repository) RemovePackage(id string) error {
    for i := 0; i < len(r.Packages); i++ {
        if r.Packages[i].Id == id {
            r.Packages = removePackage(r.Packages, r.Packages[i])
            return nil
        }
    }
    return fmt.Errorf("cannot find package with id=%s to remove", id)
}

func (r *Repository) FindPackageByRef(ref string) *Package {
    for _, p := range r.Packages {
        if p.FileRef == ref {
            return p
        }
    }
    return nil
}

func (r *Repository) Diff(repo *Repository) (*RepoDiff, error) {
    diff := new(RepoDiff)
    // if the original repository is empty
    if len(repo.Repository) == 0 && len(repo.Packages) == 0 && len(r.Repository) > 0 {
        // then the new repository packages should be added to the diff result
        for _, p := range r.Packages {
            diff.Added = append(diff.Added, p)
        }
        return diff, nil
    }
    if !strings.EqualFold(r.Repository, repo.Repository) {
        return nil, fmt.Errorf("cannot diff two different repositories %s and %s", r.Repository, repo.Repository)
    }
    // work out added
    // loops through the source repo packages
    for _, source := range r.Packages {
        var found bool
        // for each source package, loops through the target repo packages
        for _, target := range repo.Packages {
            // if the target repo contains the source package
            if source.Id == target.Id {
                // track it was found
                found = true
                // check for tags difference
                addedTags := difference(source.Tags, target.Tags)
                removedTags := difference(target.Tags, source.Tags)
                if len(addedTags) > 0 || len(removedTags) > 0 {
                    // the package has been updated
                    diff.Updated = append(diff.Updated, &UpdatedPackage{
                        Package:     source,
                        TagsAdded:   addedTags,
                        TagsRemoved: removedTags,
                    })
                }
                break
            }
        }
        // if the source package is not in the target repo, it means it is "added" in the source
        if !found {
            diff.Added = append(diff.Added, source)
        }
    }
    // work out removed
    // loops through the target repo packages
    for _, target := range repo.Packages {
        var found bool
        // for each target package, loops through the source repo packages
        for _, source := range r.Packages {
            // if the target repo contains the source package
            if target.Id == source.Id {
                // track it was found
                found = true
                break
            }
        }
        // if the target package is not in the source repo, it means it is "removed" in the source
        if !found {
            diff.Removed = append(diff.Removed, target)
        }
    }
    return diff, nil
}

// DeepCopy creates a copy that is totally unrelated to the original entity
func (r *Repository) DeepCopy() *Repository {
    repo := new(Repository)
    repo.Repository = r.Repository
    for _, p := range r.Packages {
        repo.Packages = append(repo.Packages, &Package{
            Id:      p.Id,
            Type:    p.Type,
            FileRef: p.FileRef,
            Tags:    p.Tags,
            Size:    p.Size,
            Created: p.Created,
        })
    }
    return repo
}

func (r *Repository) Group() string {
    name, _ := core.ParseName(r.Repository)
    return name.Group
}

func (r *Repository) Name() string {
    name, _ := core.ParseName(r.Repository)
    return name.Name
}

type RepoDiff struct {
    Added   []*Package
    Removed []*Package
    Updated []*UpdatedPackage
}

type UpdatedPackage struct {
    Package     *Package
    TagsAdded   []string
    TagsRemoved []string
}

// difference returns the elements in `a` that aren't in `b`.
func difference(a, b []string) []string {
    mb := make(map[string]struct{}, len(b))
    for _, x := range b {
        mb[x] = struct{}{}
    }
    var diff []string
    for _, x := range a {
        if _, found := mb[x]; !found {
            diff = append(diff, x)
        }
    }
    return diff
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

func (a *Package) IsDangling() bool {
    return len(a.Tags) == 1 && strings.Contains(a.Tags[0], "<none>")
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

type DigestInfo struct {
    Date  string `json:"date"`
    Value string `json:"value"`
}
