/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package registry

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/crypto"
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/i18n"
	"github.com/gatblau/onix/oxlib/resx"
)

// LocalRegistry the default local registry implemented as a file system
type LocalRegistry struct {
	Repositories []*Repository `json:"repositories"`
	ArtHome      string
}

func (r *LocalRegistry) api(domain, artHome string) *Api {
	return newGenericAPI(domain, artHome)
}

// NewLocalRegistry create a localRepo management structure
func NewLocalRegistry(artHome string) *LocalRegistry {
	r := &LocalRegistry{
		Repositories: []*Repository{},
		ArtHome:      artHome,
	}
	// load local registry
	r.Load()
	return r
}

// GetPackagesByName return all the packages within the same repository
func (r *LocalRegistry) GetPackagesByName(name *core.PackageName) ([]*Package, bool) {
	var packages = make([]*Package, 0)
	for _, repository := range r.Repositories {
		if repository.Repository == name.FullyQualifiedName() {
			for _, packag := range repository.Packages {
				packages = append(packages, packag)
			}
			break
		}
	}
	if len(packages) > 0 {
		return packages, true
	}
	return nil, false
}

func (r *LocalRegistry) Prune() error {
	danglingRepo := r.findDanglingRepo()
	if len(danglingRepo.Packages) > 0 {
		for _, p := range danglingRepo.Packages {
			err := r.removeFiles(p, r.ArtHome)
			if err != nil {
				return err
			}
		}
	}
	danglingRepo.Packages = nil
	r.save()
	// clears the content of the tmp folder
	if pathExist(core.TmpPath(r.ArtHome)) {
		err := cleanFolder(core.TmpPath(r.ArtHome))
		if err != nil {
			return fmt.Errorf("cannot clean tmp folder: %s", err)
		}
	}
	// clears the content of the build folder
	buildPath := path.Join(core.RegistryPath(r.ArtHome), "build")
	if pathExist(buildPath) {
		err := cleanFolder(buildPath)
		if err != nil {
			return fmt.Errorf("cannot clean build folder: %s", err)
		}
	}
	return nil
}

func pathExist(path string) bool {
	// get the absolute path
	abs, _ := filepath.Abs(path)
	// stats the path
	_, err := os.Stat(abs)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			core.WarningLogger.Printf("cannot stat path '%s': %s\n", abs, err)
			return false
		}
	}
	return true
}

// FindPackageByName return the package that matches the specified:
// - domain/group/name:tag
// nil if not found in the LocalRegistry
func (r *LocalRegistry) FindPackageByName(name *core.PackageName) *Package {
	// first gets the repository the package is in
	for _, repository := range r.Repositories {
		if repository.Repository == name.FullyQualifiedName() {
			// try and get it by id first
			for _, packag := range repository.Packages {
				for _, tag := range packag.Tags {
					// try and match against the full URI
					if tag == name.Tag {
						return packag
					}
				}
			}
			break
		}
	}
	return nil
}

// FindPackageNamesById return the packages that matches the specified:
// - package id substring
func (r *LocalRegistry) FindPackageNamesById(id string) []*core.PackageName {
	// go through the packages in the repository and check for Id matches
	names := make([]*core.PackageName, 0)
	// first gets the repository the package is in
	for _, repository := range r.Repositories {
		for _, packag := range repository.Packages {
			// try and match against the package ID substring
			if strings.HasPrefix(packag.Id, id) {
				for _, tag := range packag.Tags {
					// if the package is in the repository for dangling artefacts cannot get a name so returns nil
					if strings.Contains(repository.Repository, "<none>") {
						return nil
					}
					name, err := core.ParseName(fmt.Sprintf("%s:%s", repository.Repository, tag))
					if err != nil {
						log.Fatalf(err.Error())
					}
					names = append(names, name)
				}
			}
		}
	}
	return names
}

func (r *LocalRegistry) FindPackageById(id string) *Package {
	// first gets the repository the package is in
	for _, repository := range r.Repositories {
		for _, packag := range repository.Packages {
			// try and match against the package ID substring
			if strings.HasPrefix(packag.Id, id) {
				return packag
			}
		}
	}
	return nil
}

// Add the package and seal to the LocalRegistry
func (r *LocalRegistry) Add(filename string, name *core.PackageName, s *data.Seal) error {
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
		return errors.New(fmt.Sprintf("the localRepo can only accept zip files, the extension provided was %s", basenameExt))
	}
	// the fully qualified name of the zip package file in the local registry
	registryZipFilename := filepath.Join(core.RegistryPath(r.ArtHome), basename)
	// the fully qualified name of the json seal file in the local registry
	registryJsonFilename := filepath.Join(core.RegistryPath(r.ArtHome), fmt.Sprintf("%s.json", basenameNoExt))
	// if the zip or json files already exist in the local registry
	if fileExists(registryZipFilename) || fileExists(registryJsonFilename) {
		core.InfoLogger.Printf("package '%s' already exists, skipping import", name.FullyQualifiedNameTag())
		return nil
	}
	// move the zip file to the localRepo folder
	if err := MoveFile(filename, filepath.Join(core.RegistryPath(r.ArtHome), basename)); err != nil {
		return fmt.Errorf("failed to move package zip file to the local registry: %s", err)
	}
	// now move the seal file to the localRepo folder
	if err := MoveFile(filepath.Join(basenameDir, fmt.Sprintf("%s.json", basenameNoExt)), filepath.Join(core.RegistryPath(r.ArtHome), fmt.Sprintf("%s.json", basenameNoExt))); err != nil {
		return fmt.Errorf("failed to move package seal file to the local registry: %s", err)
	}
	// check if a package with the same name:tag exists
	old := r.FindPackageByName(name)
	// if a package was found
	if old != nil {
		// moves it to the dangling artefacts repository
		r.moveToDangling(name)
	}
	// find the repository
	repo := r.findRepository(name)
	// if the repo does not exist the creates it
	if repo == nil {
		repo = &Repository{
			Repository: name.FullyQualifiedName(),
			Packages:   make([]*Package, 0),
		}
		r.Repositories = append(r.Repositories, repo)
	}
	pkgId, err := s.PackageId()
	if err != nil {
		return err
	}
	// creates a new package
	packages := append(repo.Packages, &Package{
		Id:      pkgId,
		Type:    s.Manifest.Type,
		FileRef: basenameNoExt,
		Tags:    []string{name.Tag},
		Size:    s.Manifest.Size,
		Created: s.Manifest.Time,
	})
	repo.Packages = packages
	// persist the changes
	r.save()
	return nil
}

// moveToDangling move the specified package to the dangling artefacts repository
// and remove any existing tags
func (r *LocalRegistry) moveToDangling(name *core.PackageName) {
	// get the package repository
	repo := r.findRepository(name)
	// get the package in the repository
	p := r.FindPackageByName(name)
	// get the dangling artefact repository
	dangRepo := r.findDanglingRepo()
	// remove the package from the original repository
	repo.Packages = rmPackage(repo.Packages, p)
	// change the package tag to none
	p.Tags = []string{"<none>"}
	// add the package to the dangling repo
	dangRepo.Packages = append(dangRepo.Packages, p)
}

// findDanglingRepo find the dangling artefacts repository
// if the repository does not exist, it creates one and adds it to the collection of
// repositories of the registry
func (r *LocalRegistry) findDanglingRepo() *Repository {
	for _, r := range r.Repositories {
		if strings.Contains(r.Repository, "none") {
			return r
		}
	}
	// if the dangling repo does not exist, it creates one
	danglingRepo := &Repository{
		Repository: "<none>",
		Packages:   []*Package{},
	}
	// adds it to the collection of repos of the registry
	r.Repositories = append(r.Repositories, danglingRepo)
	// return the repo
	return danglingRepo
}

// Tag remove a given tag from an package
func (r *LocalRegistry) Tag(srcName, tgtName string) error {
	// try the package Id
	var (
		sourceName *core.PackageName
		err        error
	)
	// try to find the source package by its id
	sourcePackage := r.FindPackageById(srcName)
	// if the package was found and is dangling
	if sourcePackage != nil && sourcePackage.IsDangling() {
		// move the package to the target repository
		if err = r.moveDanglingToRepo(sourcePackage, tgtName); err != nil {
			return fmt.Errorf("cannot tag dangling package: %s", err)
		}
		// persist changes
		r.save()
		// return
		return nil
	} else {
		// the package could not be found by Id so try by name
		sourceName, err = core.ParseName(srcName)
		if err != nil {
			return fmt.Errorf("invalid source package name %s; or it does not exist", srcName)
		}
		sourcePackage = r.FindPackageByName(sourceName)
		// if the package is not found by name the exit with error
		if sourcePackage == nil {
			return fmt.Errorf("source package %s does not exist", sourceName)
		}
	}
	targetName, err := core.ParseName(tgtName)
	if err != nil {
		return fmt.Errorf("invalid target package name %s; or it does not exist", tgtName)
	}
	if targetName.IsInTheSameRepositoryAs(sourceName) {
		if !sourcePackage.HasTag(targetName.Tag) {
			// if the source package has the target name tag
			targetPackage := r.FindPackageByName(targetName)
			if targetPackage != nil && targetPackage.HasTag(targetName.Tag) {
				// remove the tag
				targetPackage.Tags = removeItem(targetPackage.Tags, targetName.Tag)
				// if no tags are left, add a default tag equal to the package file reference
				if len(targetPackage.Tags) == 0 {
					targetPackage.Tags = append(targetPackage.Tags, targetPackage.FileRef)
				}
			}
			sourcePackage.Tags = append(sourcePackage.Tags, targetName.Tag)
			r.save()
			return nil
		} else {
			return nil
		}
	} else {
		targetRepository := r.findRepository(targetName)
		newPackage := *sourcePackage
		// if the target package repository does not exist then create it
		if targetRepository == nil {
			newPackage.Tags = []string{targetName.Tag}
			r.Repositories = append(r.Repositories, &Repository{
				Repository: targetName.FullyQualifiedName(),
				Packages: []*Package{
					{
						Id:      sourcePackage.Id,
						Type:    sourcePackage.Type,
						FileRef: sourcePackage.FileRef,
						Tags:    []string{targetName.Tag},
						Size:    sourcePackage.Size,
						Created: sourcePackage.Created,
					},
				},
			})
			r.save()
			return nil
		} else {
			targetPackage := r.FindPackageByName(targetName)
			// if the package exists in the repository
			if targetPackage != nil {
				// check if the tag already exists
				for _, tag := range targetPackage.Tags {
					if tag == targetName.Tag {
					} else {
						// add the tag to the existing package
						targetPackage.Tags = append(targetPackage.Tags, targetName.Tag)
					}
				}
			} else {
				// check that an package with the Id of the source exists
				for _, a := range targetRepository.Packages {
					// if the target repository already contains the package Id
					if a.Id == sourcePackage.Id {
						// add a tag
						a.Tags = append(a.Tags, targetName.Tag)
						r.save()
						return nil
					}
				}
				// add a new package metadata in the existing repository
				targetRepository.Packages = append(targetRepository.Packages,
					&Package{
						Id:      sourcePackage.Id,
						Type:    sourcePackage.Type,
						FileRef: sourcePackage.FileRef,
						Tags:    []string{targetName.Tag},
						Size:    sourcePackage.Size,
						Created: sourcePackage.Created,
					})
				r.save()
				return nil
			}
		}
	}
	return nil
}

func (r *LocalRegistry) AllPackages() []string {
	var packages []string
	for _, repo := range r.Repositories {
		for _, p := range repo.Packages {
			for _, tag := range p.Tags {
				packages = append(packages, fmt.Sprintf("%s:%s", repo.Repository, tag))
			}
		}
	}
	return packages
}

// List packages to stdout
func (r *LocalRegistry) List(artHome string) {
	// get a table writer for the stdout
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	// print the header row
	_, err := fmt.Fprintln(w, i18n.String(artHome, i18n.LBL_LS_HEADER))
	core.CheckErr(err, "failed to write table header")
	// repository, tag, package id, created, size
	for _, repo := range r.Repositories {
		for _, a := range repo.Packages {
			// if the package is dangling (no tags)
			if len(a.Tags) == 0 {
				_, err := fmt.Fprintln(w, fmt.Sprintf("%s\t %s\t %s\t %s\t %s\t %s\t",
					repo.Repository,
					"<none>",
					a.Id[0:12],
					a.Type,
					toElapsedLabel(a.Created),
					a.Size),
				)
				core.CheckErr(err, "failed to write output")
			}
			for _, tag := range a.Tags {
				_, err := fmt.Fprintln(w, fmt.Sprintf("%s\t %s\t %s\t %s\t %s\t %s\t",
					repo.Repository,
					tag,
					a.Id[0:12],
					a.Type,
					toElapsedLabel(a.Created),
					a.Size),
				)
				core.CheckErr(err, "failed to write output")
			}
		}
	}
	err = w.Flush()
	core.CheckErr(err, "failed to flush output")
}

// ListQ list (quiet) package IDs only
func (r *LocalRegistry) ListQ() {
	// get a table writer for the stdout
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', 0)
	// repository, tag, package id, created, size
	for _, repo := range r.Repositories {
		for _, a := range repo.Packages {
			_, err := fmt.Fprintln(w, fmt.Sprintf("%s", a.Id[0:12]))
			core.CheckErr(err, "failed to write package Id")
		}
	}
	err := w.Flush()
	core.CheckErr(err, "failed to flush output")
}

func (r *LocalRegistry) Push(name *core.PackageName, credentials string, showWarnings bool) error {
	// get a reference to the remote registry
	api := r.api(name.Domain, r.ArtHome)
	// get registry credentials
	uname, pwd := core.RegUserPwd(credentials)
	// fetch the package info from the local registry
	localPackage := r.FindPackageByName(name)
	if localPackage == nil {
		return fmt.Errorf("package '%s' not found in the local registry\n", name)
	}
	// assume tls enabled
	tls := true
	// check the status of the package in the remote registry
	remotePackage, err := api.GetPackageInfo(name.Group, name.Name, localPackage.Id, uname, pwd, tls)
	if err != nil {
		// try without tls
		var err2 error
		remotePackage, err2 = api.GetPackageInfo(name.Group, name.Name, localPackage.Id, uname, pwd, false)
		if err2 == nil {
			tls = false
			if showWarnings {
				core.WarningLogger.Printf("the connection to the registry is not secure, consider connecting to a TLS enabled registry\n")
			}
		} else {
			if err2 != nil {
				return fmt.Errorf("art push '%s' cannot retrieve remote package information: %s", name.String(), err2)
			}
		}
	}
	// if the package exists in the remote registry
	if remotePackage != nil {
		// if the tag is the same then nothing to do
		if remotePackage.HasTag(name.Tag) {
			// nothing to do, returns
			i18n.Printf(r.ArtHome, i18n.INFO_NOTHING_TO_PUSH)
			return nil
		} else {
			// if the package has a different tag then the metadata has to be updated to include the new tag
			remotePackage.Tags = append(remotePackage.Tags, name.Tag)
			err = api.UpsertPackageInfo(name, remotePackage, uname, pwd, tls)
			if err != nil {
				return fmt.Errorf("cannot update remote package tags: %s", err)
			}
			i18n.Printf(r.ArtHome, i18n.INFO_TAGGED, name.String())
			// once the new tag has been added to the remote repository, exits
			return nil
		}
	}
	// if the package does not exist in the remote registry, it could be that the name:tag is already used by another package
	// so, it checks if the tag has been applied to another package in the remote repository
	repo, err, _ := api.GetRepositoryInfo(name.Group, name.Name, uname, pwd, tls)
	if err != nil {
		return fmt.Errorf("art push '%s' cannot retrieve repository information from registry", name.String())
	}
	// if the tag is in use
	var ok bool
	if remotePackage, ok = repo.GetTag(name.Tag); ok {
		// ==========================
		// removes overridden package
		// ==========================
		// first deletes the old package files
		err = api.DeletePackage(name.Group, name.Name, name.Tag, uname, pwd, tls)
		if err != nil {
			return fmt.Errorf("art push '%s' cannot remove old package from remote registry: %s", name.String(), err)
		}
		// then can remove the old package metadata form the remote repository
		// if not done in this order delete package would fail with 404 not found
		err = api.DeletePackageInfo(name.Group, name.Name, remotePackage.Id, uname, pwd, tls)
		if err != nil {
			return fmt.Errorf("art push '%s' cannot remove old package metadata from remote registry: %s", name.String(), err)
		}
	}
	// ==========================
	// adds new package
	// ==========================
	zipfile := openFile(fmt.Sprintf("%s/%s.zip", core.RegistryPath(r.ArtHome), localPackage.FileRef))
	jsonfile := openFile(fmt.Sprintf("%s/%s.json", core.RegistryPath(r.ArtHome), localPackage.FileRef))
	// prepare the package to upload
	pack := localPackage
	pack.Tags = []string{name.Tag}
	// execute the upload
	return api.UploadPackage(name, localPackage.FileRef, zipfile, jsonfile, pack, uname, pwd, tls, r.ArtHome)
}

func (r *LocalRegistry) Pull(name *core.PackageName, credentials string, showWarnings bool) *Package {
	// get a reference to the remote registry
	api := r.api(name.Domain, r.ArtHome)
	// get registry credentials
	uname, pwd := core.RegUserPwd(credentials)
	// assume tls enabled
	tls := true
	// get remote repository information
	repoInfo := &repositoryInfo{
		name:  *name,
		uname: uname,
		pwd:   pwd,
		tls:   tls,
		api:   *api,
	}
	err := getRepositoryInfoRetry(repoInfo)
	repo := repoInfo.repo
	if err != nil {
		var err2 error
		// attempt not to use tls
		repoInfo.tls = false
		err2 = getRepositoryInfoRetry(repoInfo)
		repo = repoInfo.repo
		// if successful means remote endpoint in not tls enabled
		if err2 == nil {
			// switches tls off
			tls = false
			// issue warning
			if showWarnings {
				core.WarningLogger.Printf("the connection to the registry is not secure, consider connecting to a TLS enabled registry\n")
			}
		} else {
			core.CheckErr(err2, "art pull '%s' cannot retrieve repository information from registry", name.String())
		}
	}
	// find the package to pull in the remote repository
	remoteArt, exists := repo.GetTag(name.Tag)
	if !exists {
		// if it does not exist return
		core.RaiseErr("package '%s', does not exist", name)
	}
	// check the package is not in the local registry
	localPackage := r.findPackageById(remoteArt.Id)
	// if the local registry does not have the package then download it
	if localPackage == nil {
		attempts := 5
		// download package seal file
		sealDownloadInfo := &downloadInfo{
			name:     *name,
			filename: fmt.Sprintf("%s.json", remoteArt.FileRef),
			uname:    uname,
			pwd:      pwd,
			tls:      tls,
			api:      *api,
		}
		downErr := downloadFileRetry(sealDownloadInfo, attempts)
		core.CheckErr(downErr, "failed to download package seal file")
		sealFilename := sealDownloadInfo.downloadedFilename

		// download package zip file
		packageDownloadInfo := &downloadInfo{
			name:     *name,
			filename: fmt.Sprintf("%s.zip", remoteArt.FileRef),
			uname:    uname,
			pwd:      pwd,
			tls:      tls,
			api:      *api,
		}
		downErr = downloadFileRetry(packageDownloadInfo, attempts)
		core.CheckErr(downErr, "failed to download package zip file")
		packageFilename := packageDownloadInfo.downloadedFilename

		var (
			seal  *data.Seal
			valid bool
		)
		seal, err = r.loadSeal(sealFilename)
		core.CheckErr(err, "cannot load package seal")

		// if the downloaded package digest does not match the one stored in the seal manifest
		if valid, err = seal.Valid(packageFilename); !valid {
			core.InfoLogger.Printf("package files corruption detected after download: %s, retrying %d times, stand by...\n", err, attempts)

			// retry the download of the package seal file
			downErr = downloadFileRetry(sealDownloadInfo, attempts)
			core.CheckErr(downErr, "retry failed to download the package seal file")
			sealFilename = sealDownloadInfo.downloadedFilename

			// retry the download of the package zip file
			downErr = downloadFileRetry(packageDownloadInfo, attempts)
			core.CheckErr(downErr, "retry failed to download the package zip file")
			packageFilename = packageDownloadInfo.downloadedFilename

			if valid, err = seal.Valid(packageFilename); !valid {
				core.RaiseErr("package files corruption detected after retry: %s", err)
			}
		}

		// add the package to the local registry
		err2 := r.Add(packageFilename, name, seal)
		core.CheckErr(err2, "cannot add package to local registry")
	} else {
		// if the remote package repository exists locally
		if r.findPackageByRepoAndId(name, remoteArt.Id) != nil {
			// check if the local package does not have the remote tag
			if !localPackage.HasTag(name.Tag) {
				// find the local package coordinates
				repoIx, packageIx := r.artCoords(name, localPackage)
				// add the tag locally
				r.Repositories[repoIx].Packages[packageIx].Tags = append(r.Repositories[repoIx].Packages[packageIx].Tags, name.Tag)
				// persist the changes
				r.save()
				fmt.Printf("tagged '%s' with '%s'\n", name.FullyQualifiedName(), name.Tag)
			}
		} else { // at this point the package exists locally but in a different repository or repositories
			// it needs to create the repository metadata and link it to the package
			r.Repositories = append(r.Repositories, &Repository{
				Repository: name.FullyQualifiedName(), // the local registry needs the fully qualified name because is multi repository
				Packages:   []*Package{remoteArt},
			})
			// persist the changes
			r.save()
			fmt.Printf("added package '%s' to repository '%s'\n", localPackage.Id, name.FullyQualifiedName())
		}
	}
	return r.FindPackageByName(name)
}

func (r *LocalRegistry) loadSeal(sealFilename string) (*data.Seal, error) {
	// unmarshal the seal
	sealFile, err := os.Open(sealFilename)
	if err != nil {
		return nil, fmt.Errorf("cannot read package seal file: %s", err)
	}
	seal := new(data.Seal)
	sealBytes, err := ioutil.ReadAll(sealFile)
	// exit if it failed to read the seal
	if err != nil {
		return nil, fmt.Errorf("cannot read package seal file: %s", err)
	}
	// release the handle on the seal
	err = sealFile.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close seal file stream: %s", err)
	}
	// unmarshal the seal
	err = json.Unmarshal(sealBytes, seal)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal package seal file: %s", err)
	}
	if seal.Manifest.Labels == nil {
		seal.Manifest.Labels = map[string]string{}
	}
	if seal.Manifest.Functions == nil {
		seal.Manifest.Functions = []*data.FxInfo{}
	}
	return seal, nil
}

func (r *LocalRegistry) Open(name *core.PackageName, credentials string, targetPath string, certPath string, ignoreSignature bool, v Verifier) {
	var (
		pubKeyPath = certPath
		err        error
	)
	if len(targetPath) == 0 {
		targetPath = core.WorkDir()
	} else {
		if !filepath.IsAbs(targetPath) {
			targetPath, err = filepath.Abs(targetPath)
			core.CheckErr(err, "cannot convert open path to absolute path")
		}
	}
	// fetch from local registry
	pkg := r.FindPackageByName(name)
	// if not found locally
	if pkg == nil {
		// pull it
		pkg = r.Pull(name, credentials, true)
	}
	// get the package seal
	seal, err := r.GetSeal(pkg)
	core.CheckErr(err, "cannot read package seal")
	if !ignoreSignature && v != nil {
		// get the location of the package
		zipFilename := filepath.Join(core.RegistryPath(r.ArtHome), fmt.Sprintf("%s.zip", pkg.FileRef))
		core.CheckErr(v.Verify(name, pubKeyPath, seal, zipFilename, r.ArtHome), "invalid signature")
	}
	// now we are ready to open it
	// if the target was already compressed (e.g. jar file, etc) then it should not unzip it but rename it
	// to ist original file extension
	if seal.Manifest.Zip {
		_, filename := filepath.Split(seal.Manifest.Target)
		if _, err = os.Stat(targetPath); os.IsNotExist(err) {
			err = os.MkdirAll(targetPath, os.ModePerm)
			core.CheckErr(err, "cannot create path to open package: %s", targetPath)
		}
		src := path.Join(core.RegistryPath(r.ArtHome), fmt.Sprintf("%s.zip", pkg.FileRef))
		dst := path.Join(targetPath, filename)
		err = CopyFile(src, dst)
		core.CheckErr(err, "cannot rename package %s", fmt.Sprintf("%s.zip", pkg.FileRef))
	} else {
		// otherwise, unzip the target
		err = unzip(path.Join(core.RegistryPath(r.ArtHome), fmt.Sprintf("%s.zip", pkg.FileRef)), targetPath)
		core.CheckErr(err, "cannot unzip package %s", fmt.Sprintf("%s.zip", pkg.FileRef))
		// check if the target path is a folder
		var info os.FileInfo
		info, err = os.Stat(targetPath)
		core.CheckErr(err, "cannot stat target path %s", targetPath)
		// only get rid of the target folder if there is one
		if info.IsDir() {
			srcPath := path.Join(targetPath, seal.Manifest.Target)
			info, err = os.Stat(srcPath)
			core.CheckErr(err, "cannot stat source path %s", srcPath)
			// if the source path is a folder
			if info.IsDir() {
				// unwrap the folder
				err = MoveFolderContent(srcPath, targetPath)
				core.CheckErr(err, "cannot move target folder content")
			}
		}
	}
}

func (r *LocalRegistry) Verify(name *core.PackageName, pubKeyPath string, seal *data.Seal, zipFilename string, artHome string) error {
	var (
		primaryKey, backupKey *crypto.PGP
		err                   error
	)
	if len(pubKeyPath) > 0 {
		// retrieve the verification key from the specified location
		primaryKey, err = crypto.LoadPGP(pubKeyPath, "")
		core.CheckErr(err, "cannot load public key, cannot verify signature")
	} else {
		// otherwise, loads it from the registry store
		primaryKey, backupKey, err = crypto.LoadKeys(*name, false, artHome)
		core.CheckErr(err, "cannot load public key, cannot verify signature")
	}
	// get a slice to have the unencrypted signature
	sum, _ := seal.Checksum(zipFilename)
	// if in debug mode prints out signature
	core.Debug("seal stored base64 encoded signature:\n>> start on next line\n%s\n>> ended on previous line\n", seal.Signature)
	// decode the signature in the seal
	sig, err := base64.StdEncoding.DecodeString(seal.Signature)
	core.CheckErr(err, "cannot decode signature in the seal")
	// if in debug mode prints out base64 decoded signature
	core.Debug("seal stored signature:\n>> start on next line\n%s\n>> ended on previous line\n", string(sig))
	// verify the signature using the primary key
	err = primaryKey.Verify(sum, sig)
	// if the verification failed
	if err != nil {
		// if a backup key exists
		if backupKey != nil {
			core.InfoLogger.Printf("invalid digital signature using primary key, attempting verification using backup key")
			// verify the signature using the backup key
			err = backupKey.Verify(sum, sig)
			core.CheckErr(err, "invalid digital signature (used both, primary and backup keys)")
		} else {
			// raise the error as no backup key exists
			core.CheckErr(err, "invalid digital signature (used primary key)")
		}
	}
	return err
}

func (r *LocalRegistry) removePkg(pkg *Package, artHome string) error {
	repoIxList := r.findRepositoryIxByPackageId(pkg.Id)
	for _, repoIx := range repoIxList {
		// if the repository contains the package
		if r.Repositories[repoIx].FindPackage(pkg.Id) != nil {
			// remove the package
			r.Repositories[repoIx].Packages = r.removePackageById(r.Repositories[repoIx].Packages, pkg.Id)
		}
		// if the repo does not have more packages
		if len(r.Repositories[repoIx].Packages) == 0 {
			// removes the repo
			r.Repositories = r.removeRepo(r.Repositories, *r.Repositories[repoIx])
		}
	}
	err := r.removeFiles(pkg, artHome)
	if err != nil {
		return err
	}
	r.save()
	return nil
}

func (r *LocalRegistry) removeByName(name *core.PackageName) error {
	pkg := r.FindPackageByName(name)
	if pkg == nil {
		return fmt.Errorf("package %s does not exist", name.FullyQualifiedNameTag())
	}
	// get repository by name
	repoIx := -1
	for ix, repository := range r.Repositories {
		if repository.Repository == name.FullyQualifiedName() {
			repoIx = ix
			break
		}
	}
	if repoIx == -1 {
		return fmt.Errorf("package %s does not exist", name.FullyQualifiedNameTag())
	}
	for pix, p := range r.Repositories[repoIx].Packages {
		for _, tag := range p.Tags {
			if tag == name.Tag {
				// remove the tag
				r.Repositories[repoIx].Packages[pix].Tags = removeItem(r.Repositories[repoIx].Packages[pix].Tags, tag)
				// if there are no more tags
				if len(r.Repositories[repoIx].Packages[pix].Tags) == 0 {
					// remove the package
					r.Repositories[repoIx].Packages = removePackage(r.Repositories[repoIx].Packages, p)
				}
				break
			}
		}
	}
	// if there are no more packages in the repo
	if len(r.Repositories[repoIx].Packages) == 0 {
		// remove the repo
		r.Repositories = r.removeRepo(r.Repositories, *r.Repositories[repoIx])
	}
	// check if there is any repos left for the package
	rIx := r.findRepositoryIxByPackageId(pkg.Id)
	// if not, then
	if len(rIx) == 0 {
		// remove the files
		err := r.removeFiles(pkg, r.ArtHome)
		if err != nil {
			return err
		}
	}
	r.save()
	return nil
}

func (r *LocalRegistry) Remove(names []string) error {
	for _, name := range names {
		// try and find the package using its unique ID
		pkg := r.FindPackageById(name)
		// if the package was found
		if pkg != nil {
			// remove it completely including files and references in repositories
			return r.removePkg(pkg, r.ArtHome)
		}
		// the package was not found by its ID, so try by name
		pkgName, err := core.ParseName(name)
		if err != nil {
			return fmt.Errorf("invalid package name: %s", err)
		}
		// try and find the package by name
		pkg = r.FindPackageByName(pkgName)
		// if a package with the name was found
		if pkg != nil {
			// remove the package name (if there is not more associated names then removes the package files)
			err = r.removeByName(pkgName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *LocalRegistry) GetSeal(name *Package) (*data.Seal, error) {
	sealFilename := path.Join(core.RegistryPath(r.ArtHome), fmt.Sprintf("%s.json", name.FileRef))
	sealFile, err := os.Open(sealFilename)
	if err != nil {
		return nil, fmt.Errorf("cannot open seal file %s: %s", sealFilename, err)
	}
	sealBytes, err := ioutil.ReadAll(sealFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read seal file %s: %s", sealFilename, err)
	}
	seal := new(data.Seal)
	err = json.Unmarshal(sealBytes, seal)
	return seal, err
}

func (r *LocalRegistry) ImportKey(keyPath string, isPrivate, isBackup bool, repoGroup string, repoName string, artHome string) error {
	var err error
	if !filepath.IsAbs(keyPath) {
		keyPath, err = filepath.Abs(keyPath)
		core.CheckErr(err, "cannot get an absolute representation of path '%s'", keyPath)
	}
	// only check it can read the key
	_, err = crypto.LoadPGP(keyPath, "")
	core.CheckErr(err, "cannot read pgp key '%s'", keyPath)
	destFile := crypto.KeyPath(repoGroup, repoName, isPrivate, isBackup, artHome)
	destFolder := path.Dir(destFile)
	// check if the target directory exists and if not creates it
	if _, err = os.Stat(destFolder); os.IsNotExist(err) {
		err = os.MkdirAll(destFolder, os.ModePerm)
		if err != nil {
			return fmt.Errorf("cannot create key directory '%s': %s", destFolder, err)
		}
	}
	// if so, then move the key to the correct location to preserve PEM block data
	return CopyFile(keyPath, destFile)
}

func (r *LocalRegistry) ExportKey(keyPath string, isPrivate, isBackup bool, repoGroup string, repoName string, artHome string) error {
	var err error
	if !filepath.IsAbs(keyPath) {
		keyPath, err = filepath.Abs(keyPath)
		core.CheckErr(err, "cannot get an absolute representation of path '%s'", keyPath)
	}
	return CopyFile(crypto.KeyPath(repoGroup, repoName, isPrivate, isBackup, artHome), keyPath)
}

func (r *LocalRegistry) GetManifest(name *core.PackageName) *data.Manifest {
	// find the package in the local registry
	a := r.FindPackageByName(name)
	if a == nil {
		core.RaiseErr("package '%s' not found in the local registry, pull it from remote first", name)
	}
	seal, err := r.GetSeal(a)
	core.CheckErr(err, "cannot get package seal")
	return seal.Manifest
}

// ExportPackage exports one or more packages as a tar archive to the target URI
// names: the slice of packages to save
// sourceCreds: the artisan registry credentials to pull the packages to save (in the format user:password)
// targetUri: the URI where the tar archive should be saved (could be S3 or file system)
// targetCreds: the credentials to connect to the targetUri (if it is authenticated S3 in the format user:password)
func (r *LocalRegistry) ExportPackage(names []core.PackageName, sourceCreds, targetUri, targetCreds string) error {
	var (
		pack  *Package
		repo  *Repository
		files []core.TarFile
		reg   = LocalRegistry{}
	)
	for _, name := range names {
		// find the package metadata
		repo = r.findRepository(&name)
		// if not found locally, pull the package from remote (needs credentials)
		if repo == nil {
			pack = r.Pull(&name, sourceCreds, true)
		} else {
			pack = r.FindPackageByName(&name)
		}
		// check if the package exists
		if pack == nil {
			return fmt.Errorf("package %s does not exist", name)
		}
		// works out the path to the package files in the local registry
		zipFile := filepath.Join(core.RegistryPath(r.ArtHome), fmt.Sprintf("%s.zip", pack.FileRef))
		jsonFile := filepath.Join(core.RegistryPath(r.ArtHome), fmt.Sprintf("%s.json", pack.FileRef))
		// append the package index data
		reg.Repositories = append(reg.Repositories, &Repository{
			Repository: repo.Repository,
			Packages: []*Package{ // only the exported package
				{
					Id:      pack.Id,
					Type:    pack.Type,
					FileRef: pack.FileRef,
					Size:    pack.Size,
					Created: pack.Created,
					Tags:    []string{name.Tag}, // only the exported tag
				},
			},
		})
		// add the package files to the archive list
		files = append(files, []core.TarFile{
			// add package seal
			{Path: jsonFile},
			// add package content
			{Path: zipFile},
		}...)
	}
	// add repository metadata to the archive list
	files = append(files, core.TarFile{
		Bytes: core.ToJsonBytes(reg),
		Name:  "repository.json",
	})
	// creates a bytes buffer to record content of tar
	tar := &bytes.Buffer{}
	// tar the package files without preserving directory structure
	err := core.Tar(files, tar, false, r.ArtHome)
	if err != nil {
		return err
	}

	content := tar.Bytes()

	// if no output has been specified
	if len(targetUri) == 0 {
		// prints to the stdout
		fmt.Print(string(content[:]))
	} else {
		// otherwise, if the path does not implement an URI scheme (i.e. is a file path)
		if !strings.Contains(targetUri, "://") {
			targetUri, err = filepath.Abs(targetUri)
			core.CheckErr(err, "cannot obtain the absolute output path")
			ext := filepath.Ext(targetUri)
			if len(ext) == 0 || ext != ".tar" {
				core.RaiseErr("output path must contain a filename with .tar extension")
			}
			// creates target directory
			err = os.MkdirAll(filepath.Dir(targetUri), 0755)
			if err != nil {
				return err
			}
		}
		core.InfoLogger.Printf("writing package tarball to %s", targetUri)
		err = resx.WriteFile(content, targetUri, targetCreds)
		if err != nil {
			return err
		}
	}
	return nil
}

// Import a package tar archive into the local registry
// uri: the uri of the package to import (can be file path or S3 bucket uri)
// creds: the credentials to connect to the endpoint if it is authenticated S3 in the format user:password
// localPath: if specified, it downloads the remote files to a target folder
func (r *LocalRegistry) Import(uri []string, creds, pubKeyPath string, v Verifier) error {
	for _, path := range uri {
		if err := r.importTar(path, creds, pubKeyPath, v); err != nil {
			return err
		}
	}
	return nil
}

func (r *LocalRegistry) importTar(uri, creds, pubKeyPath string, v Verifier) error {
	core.InfoLogger.Printf("reading => %s\n", uri)
	tarBytes, err := resx.ReadFile(uri, creds)
	if err != nil {
		return err
	}
	tmp, err := core.NewTempDir(r.ArtHome)
	if err != nil {
		return err
	}
	core.InfoLogger.Printf("untarring => %s\n", uri)
	err = core.Untar(bytes.NewReader(tarBytes), tmp)
	if err != nil {
		return err
	}
	// loop through extracted packages
	entries, err := os.ReadDir(tmp)
	if err != nil {
		return err
	}
	// load the repository index
	repoIndex, err := loadIndexFromPath(tmp)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		// if the entry is a package seal
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") && !strings.Contains(entry.Name(), "repository.json") {
			// Name() returns filename without path
			sealFilename := entry.Name()
			// load the package seal
			seal, err2 := r.loadSeal(filepath.Join(tmp, sealFilename))
			if err2 != nil {
				return fmt.Errorf("cannot load package seal: %s", err2)
			}
			packageName, err2 := getPackageName(*repoIndex, seal)
			if err2 != nil {
				return fmt.Errorf("cannot parse package name: %s", err2)
			}
			// works out the path to the package zip file
			packageFilename := filepath.Join(tmp, fmt.Sprintf("%s.zip", seal.Manifest.Ref))
			// if a verifier has been provided
			if v != nil {
				// use it to check the package digital signature
				err = v.Verify(packageName, pubKeyPath, seal, packageFilename, r.ArtHome)
				if err != nil {
					return err
				}
			}
			core.InfoLogger.Printf("importing => %s\n", packageName.FullyQualifiedNameTag())
			if err2 = r.Add(packageFilename, packageName, seal); err2 != nil {
				// cleanup tmp folder
				os.RemoveAll(tmp)
				// return error
				return err2
			}
		}
	}
	// cleanup tmp folder
	os.RemoveAll(tmp)
	return nil
}

// -----------------
// utility functions
// -----------------

// load the repository.json index file from a path
func loadIndexFromPath(path string) (*LocalRegistry, error) {
	repos := new(LocalRegistry)
	repoBytes, err := ioutil.ReadFile(filepath.Join(path, "repository.json"))
	if err != nil {
		return nil, fmt.Errorf("cannot read repository index in tar archive: %s", err)
	}
	err = json.Unmarshal(repoBytes, repos)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal repository index in tar archive: %s", err)
	}
	return repos, nil
}

// given the ID of a package, returns the package name with the last available tag (avoid latest)
func getPackageName(repoIx LocalRegistry, seal *data.Seal) (*core.PackageName, error) {
	// compute the package Id as a hex encoded hash of the seal
	pkgId, err := seal.PackageId()
	if err != nil {
		return nil, err
	}
	for _, repo := range repoIx.Repositories {
		for _, pack := range repo.Packages {
			if pack.Id == pkgId {
				tag := ""
				if len(pack.Tags) > 0 {
					// pick the last tag
					tag = pack.Tags[len(pack.Tags)-1]
				}
				return core.ParseName(fmt.Sprintf("%s:%s", repo.Repository, tag))
			}
		}
	}
	return nil, fmt.Errorf("either the seal or repository index content are corrupted, " +
		"the seal checksum does not match any entry held in the repository index")
}

// works out the destination folder and prefix for the key
func (r *LocalRegistry) keyDestinationFolder(repoName string, repoGroup string, artHome string) (destPath string, prefix string) {
	if len(repoName) > 0 {
		// use the repo name location
		destPath = path.Join(core.RegistryPath(artHome), "keys", repoGroup, repoName)
		prefix = fmt.Sprintf("%s_%s", repoGroup, repoName)
	} else if len(repoGroup) > 0 {
		// use the repo group location
		destPath = path.Join(core.RegistryPath(artHome), "keys", repoGroup)
		prefix = repoGroup
	} else {
		// use the registry root location
		destPath = path.Join(core.RegistryPath(artHome), "keys")
		prefix = "root"
	}
	_, err := os.Stat(destPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(destPath, os.ModePerm)
		core.CheckErr(err, "cannot create private key destination '%s'", destPath)
	}
	return destPath, prefix
}

// remove the files associated with an Package
func (r *LocalRegistry) removeFiles(pack *Package, artHome string) error {
	// remove the zip file
	err := os.Remove(fmt.Sprintf("%s/%s.zip", core.RegistryPath(artHome), pack.FileRef))
	if err != nil {
		return err
	}
	// remove the json file
	return os.Remove(fmt.Sprintf("%s/%s.json", core.RegistryPath(artHome), pack.FileRef))
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

// the fully qualified name of the json Seal file in the local localReg
func (r *LocalRegistry) regDirJsonFilename(uniqueIdName string) string {
	return fmt.Sprintf("%s/%s.json", core.RegistryPath(r.ArtHome), uniqueIdName)
}

// the fully qualified name of the zip file in the local localReg
func (r *LocalRegistry) regDirZipFilename(uniqueIdName string) string {
	return fmt.Sprintf("%s/%s.zip", core.RegistryPath(r.ArtHome), uniqueIdName)
}

// find the package specified by its id
func (r *LocalRegistry) findPackageByRepoAndId(name *core.PackageName, id string) *Package {
	for _, repository := range r.Repositories {
		rep := fmt.Sprintf("%s/%s/%s", name.Domain, name.Group, name.Name)
		if rep == repository.Repository {
			for _, pack := range repository.Packages {
				if pack.Id == id {
					return pack
				}
			}
		}
	}
	return nil
}

func (r *LocalRegistry) findPackageById(id string) *Package {
	for _, repository := range r.Repositories {
		for _, pack := range repository.Packages {
			if pack.Id == id {
				return pack
			}
		}
	}
	return nil
}

// returns the package coordinates in the repository as (repo index, package index)
func (r *LocalRegistry) artCoords(name *core.PackageName, art *Package) (int, int) {
	for rIx, repository := range r.Repositories {
		rep := fmt.Sprintf("%s/%s/%s", name.Domain, name.Group, name.Name)
		if rep == repository.Repository {
			for aIx, pack := range repository.Packages {
				if pack.Id == art.Id {
					return rIx, aIx
				}
			}
		}
	}
	return -1, -1
}

// remove all tags from the specified Package
func (r *LocalRegistry) unTagAll(name *core.PackageName) {
	if packages, exists := r.GetPackagesByName(name); exists {
		// then it has to untag it, leaving a dangling package
		for _, pack := range packages {
			for _, tag := range pack.Tags {
				pack.Tags = core.RemoveElement(pack.Tags, tag)
			}
		}
	}
	// persist changes
	r.save()
}

// remove a given tag from an Package
func (r *LocalRegistry) unTag(name *core.PackageName, tag string) {
	pkg := r.FindPackageByName(name)
	if pkg != nil {
		pkg.Tags = core.RemoveElement(pkg.Tags, tag)
	}
}

func (r *LocalRegistry) removePackageById(a []*Package, id string) []*Package {
	i := -1
	// find the value to remove
	for ix := 0; ix < len(a); ix++ {
		if strings.Contains(a[ix].Id, id) {
			i = ix
			break
		}
	}
	if i == -1 {
		return a
	}
	// Remove the element at index i from a.
	a[i] = a[len(a)-1] // Copy last element to index i.
	a[len(a)-1] = nil  // Erase last element (write zero value).
	a = a[:len(a)-1]   // Truncate slice.
	return a
}

func (r *LocalRegistry) removeRepoByName(a []*Repository, name *core.PackageName) []*Repository {
	i := -1
	// find an package with the specified tag
	for ix := 0; ix < len(a); ix++ {
		if a[ix].Repository == name.FullyQualifiedName() {
			i = ix
			break
		}
	}
	if i == -1 {
		return a
	}
	// Remove the element at index i from a.
	a[i] = a[len(a)-1] // Copy last element to index i.
	a[len(a)-1] = nil  // Erase last element (write zero value).
	a = a[:len(a)-1]   // Truncate slice.
	return a
}

func (r *LocalRegistry) removeRepo(a []*Repository, repo Repository) []*Repository {
	i := -1
	// find an package with the specified tag
	for ix := 0; ix < len(a); ix++ {
		if a[ix].Repository == repo.Repository {
			i = ix
			break
		}
	}
	if i == -1 {
		return a
	}
	// Remove the element at index i from a.
	a[i] = a[len(a)-1] // Copy last element to index i.
	a[len(a)-1] = nil  // Erase last element (write zero value).
	a = a[:len(a)-1]   // Truncate slice.
	return a
}

// return the LocalRegistry full file name
func (r *LocalRegistry) file() string {
	return filepath.Join(core.RegistryPath(r.ArtHome), "repository.json")
}

// save the state of the LocalRegistry
func (r *LocalRegistry) save() {
	regBytes := core.ToJsonBytes(r)
	core.CheckErr(ioutil.WriteFile(r.file(), regBytes, os.ModePerm), "fail to update local registry metadata")
}

// Load the content of the LocalRegistry
func (r *LocalRegistry) Load() {
	var (
		regBytes []byte
		err      error
	)
	// check if localRepo file exist
	_, err = os.Stat(r.file())
	if err != nil {
		// then assume localRepo.json is not there: try and create it
		r.save()
	} else {
		regBytes, err = ioutil.ReadFile(r.file())
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(regBytes, r)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// find the Repository specified by name
func (r *LocalRegistry) findRepository(name *core.PackageName) *Repository {
	// find repository using package name
	for _, repository := range r.Repositories {
		if repository.Repository == name.FullyQualifiedName() {
			return repository
		}
	}
	// find repository using package Id
	for _, repository := range r.Repositories {
		for _, artie := range repository.Packages {
			if strings.Contains(artie.Id, name.Name) {
				return repository
			}
		}
	}
	return nil
}

// moveDanglingToRepo move the passed in dangling package to the repo specified by the target name
func (r *LocalRegistry) moveDanglingToRepo(srcPackage *Package, targetName string) error {
	// check the package is dangling
	srcRepo := r.findRepositoryByPackageId(srcPackage.Id)
	if srcRepo == nil {
		return fmt.Errorf("cannot find repository for package Id %6s", srcPackage.Id)
	}
	if srcRepo != nil && len(srcRepo) > 1 || len(srcRepo) == 0 {
		return fmt.Errorf("the package with Id %6s is not in the dangling repository", srcPackage.Id)
	}
	// parse the target name
	tgtName, err := core.ParseName(targetName)
	if err != nil {
		return fmt.Errorf("invalid target package name %s", targetName)
	}
	// check if the target name already exists in the target repo
	tgtPackage := r.FindPackageByName(tgtName)
	if tgtPackage != nil {
		// cannot override this package with dangling
		return fmt.Errorf("package name %s already exists, use a name that is not in use or re-tag or delete the existing package", targetName)
	}
	// check if the target repo exists
	tgtRepo := r.findRepository(tgtName)
	// if it does not
	if tgtRepo == nil {
		// creates a new repo
		tgtRepo = &Repository{
			Repository: tgtName.FullyQualifiedName(),
			Packages:   []*Package{},
		}
		// add the repo to the local registry
		r.Repositories = append(r.Repositories, tgtRepo)
	}
	// remove the package from the dangling repo
	srcRepo[0].Packages = r.removePackageById(srcRepo[0].Packages, srcPackage.Id)
	// update the src package tag
	srcPackage.Tags = []string{tgtName.Tag}
	// add the package  to the target repo
	tgtRepo.Packages = append(tgtRepo.Packages, srcPackage)
	return nil
}

// findRepositoryByPackageId find the repository a package with a specific Id is in
func (r *LocalRegistry) findRepositoryByPackageId(id string) []Repository {
	var repos []Repository
	for _, repository := range r.Repositories {
		for _, p := range repository.Packages {
			if p.Id == id {
				repos = append(repos, *repository)
			}
		}
	}
	return repos
}

func (r *LocalRegistry) findRepositoryIxByPackageId(id string) []int {
	var ix []int
	for rix, repository := range r.Repositories {
		for _, p := range repository.Packages {
			if p.Id == id {
				ix = append(ix, rix)
			}
		}
	}
	return ix
}

func (r *LocalRegistry) Sign(pac, pkPath, pubPath string, v Verifier) error {
	// parses the package name
	packageName, err := core.ParseName(pac)
	if err != nil {
		return err
	}
	// find the package by name
	pkg := r.FindPackageByName(packageName)
	if pkg == nil {
		return fmt.Errorf("package %s not found", pac)
	}
	// works out the seal filename
	sealFilename := r.regDirJsonFilename(pkg.FileRef)
	// works out the zip filename
	zipFilename := r.regDirZipFilename(pkg.FileRef)
	// load the package seal
	s, sealErr := r.loadSeal(sealFilename)
	if sealErr != nil {
		return sealErr
	}
	// if a public key has been provided, use it to verify the package digital signature
	if len(pubPath) > 0 && v != nil {
		err = v.Verify(packageName, pubPath, s, zipFilename, r.ArtHome)
		if err != nil {
			return err
		}
	}
	// gets a timestamp
	t := time.Now()
	timeStamp := fmt.Sprintf("%04s%02d%02d%02d%02d%02d%s", strconv.Itoa(t.Year()), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), strconv.Itoa(t.Nanosecond())[:3])
	// add labels to the manifest to keep an audit trail of the re-signing operation
	s.Manifest.Labels[fmt.Sprintf("source-signature:%s", timeStamp)] = s.Signature
	// gets the combined checksum of the manifest and the package
	sum, _ := s.Checksum(zipFilename)
	// load private key
	var pk *crypto.PGP
	// if no private key path has been provided
	if len(pkPath) == 0 {
		// load the key from the local registry
		pk, _, err = crypto.LoadKeys(*packageName, true, r.ArtHome)
		if err != nil {
			return fmt.Errorf("cannot load signing key: %s", err)
		}
	} else {
		// uses the path provided
		path, absErr := filepath.Abs(pkPath)
		if absErr != nil {
			return absErr
		}
		pk, err = crypto.LoadPGP(path, "")
		if err != nil {
			return fmt.Errorf("cannot load signing key: %s", err)
		}
	}
	// create a PGP cryptographic signature
	signature, err := pk.Sign(sum)
	if err != nil {
		return fmt.Errorf("cannot create cryptographic signature: %s", err)
	}
	// replace the signature
	s.Signature = base64.StdEncoding.EncodeToString(signature)
	// convert the seal to Json
	dest := core.ToJsonBytes(s)
	// save the seal
	err = ioutil.WriteFile(sealFilename, dest, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot update package seal file: %s", err)
	}
	return nil
}

// checks if a file exists
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// rmPackage removes a package from a slice
func rmPackage(a []*Package, value *Package) []*Package {
	i := -1
	// find the value to remove
	for ix := 0; ix < len(a); ix++ {
		if a[ix] == value {
			i = ix
			break
		}
	}
	if i == -1 {
		return a
	}
	// Remove the element at index i from a.
	a[i] = a[len(a)-1] // Copy last element to index i.
	a[len(a)-1] = nil  // Erase last element (write zero value).
	a = a[:len(a)-1]   // Truncate slice.
	return a
}

type Verifier interface {
	Verify(name *core.PackageName, pubKeyPath string, seal *data.Seal, zipFilename, artHome string) error
}

type Signer interface {
	Verify(name *core.PackageName, pubKeyPath string, seal *data.Seal, zipFilename string) error
}
