/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package backend

import (
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/registry"
	"mime/multipart"
	"os"
)

// the interface implemented by a backend
type Backend interface {
	// upload an package to the remote repository
	UploadPackage(name *core.PackageName, packageRef string, zipfile multipart.File, jsonFile multipart.File, repo multipart.File, user string, pwd string) error
	// get repository information
	GetRepositoryInfo(group, name, user, pwd string) (*registry.Repository, error)
	// get package information
	GetPackageInfo(group, name, id, user, pwd string) (*registry.Package, error)
	// update package information
	UpdatePackageInfo(group string, name string, packageInfo *registry.Package, user string, pwd string) error
	// open a file for download
	Download(repoGroup, repoName, fileName, user, pwd string) (*os.File, error)
}
