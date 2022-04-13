/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package backend

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
	"path/filepath"
	"regexp"
)

// FsBackend file system backend
type FsBackend struct {
	path string
}

func NewFsBackend() *FsBackend {
	fs := &FsBackend{
		path: "data",
	}
	fs.checkPath()
	return fs
}

// UpsertPackageInfo insert or update the package information in the repository index
func (fs *FsBackend) UpsertPackageInfo(group, name string, packageInfo *registry.Package, user string, pwd string) error {
	repo, err := fs.GetRepositoryInfo(group, name, user, pwd)
	if err != nil {
		return err
	}
	repo.UpsertPackage(packageInfo)
	return fs.saveIndex(group, name, repo)
}

// DeletePackageInfo delete the package information from the repository index
func (fs *FsBackend) DeletePackageInfo(group, name string, packageId string, user string, pwd string) error {
	repo, err := fs.GetRepositoryInfo(group, name, user, pwd)
	if err != nil {
		return err
	}
	repo.RemovePackage(packageId)
	return fs.saveIndex(group, name, repo)
}

// DeletePackage delete a specific package
func (fs *FsBackend) DeletePackage(group, name, packageRef, user, pwd string) error {
	root := fs.packagePath(group, name)
	sealFile := fmt.Sprintf("%s.json", packageRef)
	err := os.Remove(path.Join(root, sealFile))
	if err != nil {
		return fmt.Errorf("cannot remove package seal %s: %s", sealFile, err)
	}
	pacFile := fmt.Sprintf("%s.zip", packageRef)
	err = os.Remove(path.Join(root, sealFile))
	if err != nil {
		return fmt.Errorf("cannot remove package file %s: %s", pacFile, err)
	}
	return nil
}

// GetPackageManifest get the manifest for a specific package
func (fs *FsBackend) GetPackageManifest(group, name, tag, user, pwd string) (*data.Manifest, error) {
	repo, err := fs.GetRepositoryInfo(group, name, user, pwd)
	if err != nil {
		return nil, err
	}
	if len(tag) == 0 {
		tag = "latest"
	}
	ref, err := repo.GetFileRef(tag)
	if err != nil {
		return nil, err
	}
	manifestFile, err := fs.Download(group, name, fmt.Sprintf("%s.json", ref), user, pwd)
	if err != nil {
		return nil, err
	}
	defer manifestFile.Close()
	bytes, err := ioutil.ReadAll(manifestFile)
	if err != nil {
		return nil, err
	}
	var seal data.Seal
	err = json.Unmarshal(bytes, &seal)
	if err != nil {
		return nil, err
	}
	return seal.Manifest, nil
}

func (fs *FsBackend) GetAllRepositoryInfo(user, pwd string) ([]*registry.Repository, error) {
	infos := make([]*registry.Repository, 0)
	libRegEx, err := regexp.Compile("^.*repository.json$")
	if err != nil {
		return nil, err
	}
	err = filepath.Walk(fs.dataPath(), func(path string, info os.FileInfo, err error) error {
		if err == nil && libRegEx.MatchString(info.Name()) {
			bytes, err2 := os.ReadFile(path)
			if err2 != nil {
				return fmt.Errorf("cannot read file: %s", err)
			}
			repo := new(registry.Repository)
			err2 = json.Unmarshal(bytes, repo)
			if err2 != nil {
				return fmt.Errorf("cannot unmarshal file: %s", err)
			}
			infos = append(infos, repo)
		}
		return nil
	})
	return infos, err
}

// Name of the backend
func (fs *FsBackend) Name() string {
	return "FILE_SYSTEM"
}

// UploadPackage upload a package to the remote repository
func (fs *FsBackend) UploadPackage(group, name string, packageRef string, zipfile multipart.File, jsonFile multipart.File, repo multipart.File, user string, pwd string) error {
	// ensure files are properly closed
	defer zipfile.Close()
	defer jsonFile.Close()
	defer repo.Close()

	fs.checkPackagePath(group, name)

	// seal file
	seal := new(data.Seal)
	sealBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("cannot read package seal file: %s", err)
	}
	err = json.Unmarshal(sealBytes, seal)
	if err != nil {
		return fmt.Errorf("cannot unmarshal package seal file: %s", err)
	}
	err = os.WriteFile(fs.sealFilename(group, name, seal), sealBytes, 0666)
	if err != nil {
		return fmt.Errorf("cannot write package seal file to the backend file system: %s", err)
	}
	// zip file
	packageBytes, err := ioutil.ReadAll(zipfile)
	if err != nil {
		return fmt.Errorf("cannot read package file: %s", err)
	}
	err = os.WriteFile(fs.packFilename(group, name, seal), packageBytes, 0666)
	if err != nil {
		return fmt.Errorf("cannot write package file to the backend file system: %s", err)
	}
	// repository.json
	repoBytes, err := ioutil.ReadAll(repo)
	if err != nil {
		return fmt.Errorf("cannot read repository.json file: %s", err)
	}
	err = os.WriteFile(fs.indexFilename(group, name), repoBytes, 0666)
	if err != nil {
		return fmt.Errorf("cannot write repository.json file to the backend file system: %s", err)
	}
	return nil
}

// GetRepositoryInfo get repository information
func (fs *FsBackend) GetRepositoryInfo(group, name, user, pwd string) (*registry.Repository, error) {
	repoFile := fs.indexFilename(group, name)
	if _, err := os.Stat(repoFile); os.IsNotExist(err) {
		// return an empty repository
		return &registry.Repository{
			Repository: fmt.Sprintf("%s/%s", group, name),
			Packages:   make([]*registry.Package, 0),
		}, nil
	}
	repoBytes, err := os.ReadFile(repoFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read repository index: %s", err)
	}
	repository := new(registry.Repository)
	err = json.Unmarshal(repoBytes, repository)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal repository index: %s", err)
	}
	return repository, nil
}

// GetPackageInfo get package information
func (fs *FsBackend) GetPackageInfo(group, name, id, user, pwd string) (*registry.Package, error) {
	repo, err := fs.GetRepositoryInfo(group, name, user, pwd)
	if err != nil {
		return nil, err
	}
	if repo != nil {
		return repo.FindPackage(id), nil
	}
	return nil, nil
}

// Download open a file for download
func (fs *FsBackend) Download(repoGroup, repoName, fileName, user, pwd string) (*os.File, error) {
	fqn := filepath.Join(fs.packagePath(repoGroup, repoName), fileName)
	f, err := os.Open(fqn)
	if err != nil {
		return nil, fmt.Errorf("cannot open package file %s: file", fileName)
	}
	return f, nil
}

func (fs *FsBackend) dataPath() string {
	return path.Join(core.RegistryPath(), fs.path)
}

func (fs *FsBackend) checkPath() {
	_, err := os.Stat(fs.dataPath())
	if os.IsNotExist(err) {
		err = os.MkdirAll(fs.dataPath(), os.ModePerm)
		core.CheckErr(err, "cannot create Artisan registry file system backend path")
	}
}

func (fs *FsBackend) indexFilename(group, name string) string {
	return path.Join(fs.dataPath(), group, name, "repository.json")
}

func (fs *FsBackend) packagePath(group, name string) string {
	return path.Join(fs.dataPath(), group, name)
}

func (fs *FsBackend) checkPackagePath(group, name string) {
	packagePath := fs.packagePath(group, name)
	if _, err := os.Stat(packagePath); os.IsNotExist(err) {
		_ = os.MkdirAll(packagePath, os.ModePerm)
	}
}

func (fs *FsBackend) sealFilename(group, name string, seal *data.Seal) string {
	return path.Join(fs.packagePath(group, name), fmt.Sprintf("%s.json", seal.Manifest.Ref))
}

func (fs *FsBackend) packFilename(group, name string, seal *data.Seal) string {
	return path.Join(fs.packagePath(group, name), fmt.Sprintf("%s.zip", seal.Manifest.Ref))
}

func (fs *FsBackend) saveIndex(group string, name string, repo *registry.Repository) error {
	repoBytes, err := json.Marshal(repo)
	if err != nil {
		return fmt.Errorf("cannot marshal repository.json file: %s", err)
	}
	err = os.WriteFile(fs.indexFilename(group, name), repoBytes, 0666)
	if err != nil {
		return fmt.Errorf("cannot write repository.json file to the backend file system: %s", err)
	}
	return nil
}
