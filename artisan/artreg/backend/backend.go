package backend

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/registry"
	"mime/multipart"
	"os"
)

// Backend the interface implemented by a backend
// the interface
type Backend interface {
	// UploadPackage upload a package to the remote repository
	UploadPackage(group, name, packageRef string, zipfile multipart.File, jsonFile multipart.File, repo multipart.File, user string, pwd string) error
	// DeletePackage a package
	DeletePackage(group, name, packageRef, user, pwd string) error
	// GetPackageInfo get package information
	GetPackageInfo(group, name, id, user, pwd string) (*registry.Package, error)
	// UpdatePackageInfo update package information
	UpdatePackageInfo(group, name string, packageInfo *registry.Package, user string, pwd string) error
	// GetPackageManifest get the package manifest
	GetPackageManifest(group, name, tag, user, pwd string) (*data.Manifest, error)

	// GetAllRepositoryInfo get information for all repositories in the remote repository
	GetAllRepositoryInfo(user, pwd string) ([]*registry.Repository, error)
	// GetRepositoryInfo get information for a specific repository in the remote repository
	GetRepositoryInfo(group, name, user, pwd string) (*registry.Repository, error)

	// Download open a file for download
	Download(repoGroup, repoName, fileName, user, pwd string) (*os.File, error)
	// Name print usage info
	Name() string
}
