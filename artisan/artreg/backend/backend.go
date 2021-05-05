package backend

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/registry"
	"mime/multipart"
	"os"
)

// Backend the interface implemented by a backend
type Backend interface {
	// UploadPackage upload an package to the remote repository
	UploadPackage(name *core.PackageName, packageRef string, zipfile multipart.File, jsonFile multipart.File, repo multipart.File, user string, pwd string) error
	// GetRepositoryInfo get repository information
	GetRepositoryInfo(group, name, user, pwd string) (*registry.Repository, error)
	// GetPackageInfo get package information
	GetPackageInfo(group, name, id, user, pwd string) (*registry.Package, error)
	// UpdatePackageInfo update package information
	UpdatePackageInfo(name *core.PackageName, packageInfo *registry.Package, user string, pwd string) error
	// Download open a file for download
	Download(repoGroup, repoName, fileName, user, pwd string) (*os.File, error)
	// Name print usage info
	Name() string
}
