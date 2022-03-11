package backend

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"encoding/json"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/registry"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
)

type FsBackend struct {
	path string
}

func (fs *FsBackend) DeletePackage(group, name, packageRef, user, pwd string) error {
	panic("implement me")
}

func (fs *FsBackend) GetPackageManifest(group, name, tag, user, pwd string) (*data.Manifest, error) {
	panic("implement me")
}

func NewFsBackend() *FsBackend {
	fs := &FsBackend{
		path: "data",
	}
	fs.checkPath()
	return fs
}

func (fs *FsBackend) GetManifest(group, name, tag, user, pwd string) (*data.Manifest, error) {
	panic("implement me")
}

func (fs *FsBackend) GetAllRepositoryInfo(user, pwd string) ([]*registry.Repository, error) {
	panic("implement me")
}

func (fs *FsBackend) Name() string {
	return "FILE_SYSTEM"
}

// UploadPackage upload a package to the remote repository
func (fs *FsBackend) UploadPackage(group, name string, packageRef string, zipfile multipart.File, jsonFile multipart.File, repo multipart.File, user string, pwd string) error {
	// ensure files are properly closed
	defer zipfile.Close()
	defer jsonFile.Close()
	defer repo.Close()

	seal := new(data.Seal)
	sealBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("cannot read package seal file: %s", err)
	}
	err = json.Unmarshal(sealBytes, seal)
	if err != nil {
		return fmt.Errorf("cannot unmarshal package seal file: %s", err)
	}
	err = ioutil.WriteFile(fs.sealFilename(seal), sealBytes, 0666)
	if err != nil {
		return fmt.Errorf("cannot write package seal file to the backend file system: %s", err)
	}
	packageBytes, err := ioutil.ReadAll(zipfile)
	if err != nil {
		return fmt.Errorf("cannot read package file: %s", err)
	}
	err = ioutil.WriteFile(fs.packFilename(seal), packageBytes, 0666)
	if err != nil {
		return fmt.Errorf("cannot write package file to the backend file system: %s", err)
	}
	return nil
}

// GetRepositoryInfo get repository information
func (fs *FsBackend) GetRepositoryInfo(group, name, user, pwd string) (*registry.Repository, error) {
	// return an empty repository
	return &registry.Repository{
		Repository: fmt.Sprintf("%s/%s", group, name),
		Packages:   make([]*registry.Package, 0),
	}, nil
}

// GetPackageInfo get package information
func (fs *FsBackend) GetPackageInfo(group, name, id, user, pwd string) (*registry.Package, error) {
	return nil, nil
}

// UpdatePackageInfo update package information
func (fs *FsBackend) UpdatePackageInfo(group, name string, packageInfo *registry.Package, user string, pwd string) error {
	return nil
}

// Download open a file for download
func (fs *FsBackend) Download(repoGroup, repoName, fileName, user, pwd string) (*os.File, error) {
	return nil, nil
}

func (fs *FsBackend) dataPath() string {
	return path.Join(core.RegistryPath(), fs.path)
}

func (fs *FsBackend) checkPath() {
	_, err := os.Stat(fs.dataPath())
	if os.IsNotExist(err) {
		err := os.MkdirAll(fs.dataPath(), os.ModePerm)
		core.CheckErr(err, "cannot create Artisan registry file system backend path")
	}
}

func (fs *FsBackend) sealFilename(seal *data.Seal) string {
	return path.Join(fs.dataPath(), fmt.Sprintf("%s.json", seal.Manifest.Ref))
}

func (fs *FsBackend) packFilename(seal *data.Seal) string {
	return path.Join(fs.dataPath(), fmt.Sprintf("%s.zip", seal.Manifest.Ref))
}
