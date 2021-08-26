package registry

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/crypto"
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/i18n"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"
)

// LocalRegistry the default local registry implemented as a file system
type LocalRegistry struct {
	Repositories []*Repository `json:"repositories"`
}

func (r *LocalRegistry) api(domain string) *Api {
	return NewGenericAPI(domain)
}

// create a localRepo management structure
func NewLocalRegistry() *LocalRegistry {
	r := &LocalRegistry{
		Repositories: []*Repository{},
	}
	// check the registry directory is in place
	r.checkRegistryDir()
	// load local registry
	r.load()
	return r
}

// return all the packages within the same repository
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

// return the package that matches the specified:
// - domain/group/name:tag
// nil if not found in the LocalRegistry
func (r *LocalRegistry) FindPackage(name *core.PackageName) *Package {
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

// return the packages that matches the specified:
// - package id substring
func (r *LocalRegistry) FindPackagesById(id string) []*core.PackageName {
	// go through the packages in the repository and check for Id matches
	names := make([]*core.PackageName, 0)
	// first gets the repository the package is in
	for _, repository := range r.Repositories {
		for _, packag := range repository.Packages {
			// try and match against the package ID substring
			if strings.Contains(packag.Id, id) {
				for _, tag := range packag.Tags {
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

// the local Path to the local LocalRegistry
func (r *LocalRegistry) Path() string {
	return core.RegistryPath()
}

// Add the package and seal to the LocalRegistry
func (r *LocalRegistry) Add(filename string, name *core.PackageName, s *data.Seal) {
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
	core.CheckErr(MoveFile(filename, filepath.Join(r.Path(), basename)), "failed to move package zip file to the local registry")
	// now move the seal file to the localRepo folder
	core.CheckErr(MoveFile(filepath.Join(basenameDir, fmt.Sprintf("%s.json", basenameNoExt)), filepath.Join(r.Path(), fmt.Sprintf("%s.json", basenameNoExt))), "failed to move package seal file to the local registry")
	// untag package package (if any)
	r.unTag(name, name.Tag)
	// remove any dangling packages
	r.removeDangling(name)
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
	// creates a new package
	packages := append(repo.Packages, &Package{
		Id:      s.PackageId(),
		Type:    s.Manifest.Type,
		FileRef: basenameNoExt,
		Tags:    []string{name.Tag},
		Size:    s.Manifest.Size,
		Created: s.Manifest.Time,
	})
	repo.Packages = packages
	// persist the changes
	r.save()
}

// remove a given tag from an package
func (r *LocalRegistry) Tag(sourceName *core.PackageName, targetName *core.PackageName) {
	sourcePackage := r.FindPackage(sourceName)
	if sourcePackage == nil {
		core.RaiseErr("source package %s does not exit", sourceName)
	}
	if targetName.IsInTheSameRepositoryAs(sourceName) {
		if !sourcePackage.HasTag(targetName.Tag) {
			sourcePackage.Tags = append(sourcePackage.Tags, targetName.Tag)
			r.save()
			return
		} else {
			return
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
			return
		} else {
			targetPackage := r.FindPackage(targetName)
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
						return
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
				return
			}
		}
	}
}

// List packages to stdout
func (r *LocalRegistry) List() {
	// get a table writer for the stdout
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 12, ' ', 0)
	// print the header row
	_, err := fmt.Fprintln(w, i18n.String(i18n.LBL_LS_HEADER))
	core.CheckErr(err, "failed to write table header")
	// repository, tag, package id, created, size
	for _, repo := range r.Repositories {
		for _, a := range repo.Packages {
			// if the package is dangling (no tags)
			if len(a.Tags) == 0 {
				_, err := fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s",
					repo.Repository,
					"<none>",
					a.Id[7:19],
					a.Type,
					toElapsedLabel(a.Created),
					a.Size),
				)
				core.CheckErr(err, "failed to write output")
			}
			for _, tag := range a.Tags {
				_, err := fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s",
					repo.Repository,
					tag,
					a.Id[7:19],
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

// list (quiet) package IDs only
func (r *LocalRegistry) ListQ() {
	// get a table writer for the stdout
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 10, ' ', 0)
	// repository, tag, package id, created, size
	for _, repo := range r.Repositories {
		for _, a := range repo.Packages {
			_, err := fmt.Fprintln(w, fmt.Sprintf("%s", a.Id[7:19]))
			core.CheckErr(err, "failed to write package Id")
		}
	}
	err := w.Flush()
	core.CheckErr(err, "failed to flush output")
}

func (r *LocalRegistry) Push(name *core.PackageName, credentials string) {
	// get a reference to the remote registry
	api := r.api(name.Domain)
	// get registry credentials
	uname, pwd := core.UserPwd(credentials)
	// fetch the package info from the local registry
	localPackage := r.FindPackage(name)
	if localPackage == nil {
		fmt.Printf("package '%s' not found in the local registry\n", name)
		return
	}
	// assume tls enabled
	tls := true
	// check the status of the package in the remote registry
	remoteArt, err := api.GetPackageInfo(name.Group, name.Name, localPackage.Id, uname, pwd, tls)
	if err != nil {
		// try without tls
		var err2 error
		remoteArt, err2 = api.GetPackageInfo(name.Group, name.Name, localPackage.Id, uname, pwd, false)
		if err2 == nil {
			tls = false
			core.WarningLogger.Printf("artisan registry does not use TLS: the connection to the registry is not secure\n")
		} else {
			core.CheckErr(err, "art push '%s' cannot retrieve remote package information", name.String())
		}
	}
	// if the package exists in the remote registry
	if remoteArt != nil {
		// check if the tag already exist in the remote repository
		if remoteArt.HasTag(name.Tag) {
			// nothing to do, returns
			i18n.Printf(i18n.INFO_NOTHING_TO_PUSH)
			return
		} else {
			// the metadata has to be updated to include the new tag
			remoteArt.Tags = append(remoteArt.Tags, name.Tag)
			err = api.UpdatePackageInfo(name, remoteArt, uname, pwd, tls)
			core.CheckErr(err, "cannot update remote package tags")
			return
		}
	}
	// if the package does not exist in the remote registry
	// check if the tag has been applied to another package in the repository
	repo, err := api.GetRepositoryInfo(name.Group, name.Name, uname, pwd, tls)
	core.CheckErr(err, "art push '%s' cannot retrieve repository information from registry", name.String())
	// if so
	if a, ok := repo.GetTag(name.Tag); ok {
		// remove the tag from the package as it will be applied to the new package
		a.RemoveTag(name.Tag)
		// if the package has no tags left
		if len(a.Tags) == 0 {
			// adds a default tag matching the package file reference
			a.Tags = append(a.Tags, a.FileRef)
			// updates the metadata in the remote repo
			core.CheckErr(api.UpdatePackageInfo(name, a, uname, pwd, tls), "cannot update package info")
		}
	}
	zipfile := openFile(fmt.Sprintf("%s/%s.zip", r.Path(), localPackage.FileRef))
	jsonfile := openFile(fmt.Sprintf("%s/%s.json", r.Path(), localPackage.FileRef))
	// prepare the package to upload
	pack := localPackage
	pack.Tags = []string{name.Tag}
	// execute the upload
	err = api.UploadPackage(name, localPackage.FileRef, zipfile, jsonfile, pack, uname, pwd, tls)
	i18n.Err(err, i18n.ERR_CANT_PUSH_PACKAGE)
}

func (r *LocalRegistry) Pull(name *core.PackageName, credentials string) *Package {
	// get a reference to the remote registry
	api := r.api(name.Domain)
	// get registry credentials
	uname, pwd := core.UserPwd(credentials)
	// assume tls enabled
	tls := true
	// get remote repository information
	repo, err := api.GetRepositoryInfo(name.Group, name.Name, uname, pwd, tls)
	if err != nil {
		var err2 error
		// attempt not to use tls
		repo, err2 = api.GetRepositoryInfo(name.Group, name.Name, uname, pwd, false)
		// if successful means remote endpoint in not tls enabled
		if err2 == nil {
			// switches tls off
			tls = false
			// issue warning
			core.WarningLogger.Printf("artisan registry does not use TLS: the connection to the registry is not secure\n")
		} else {
			core.CheckErr(err, "art pull '%s' cannot retrieve repository information from registry", name.String())
		}
	}
	// find the package to pull in the remote repository
	remoteArt, exists := repo.GetTag(name.Tag)
	if !exists {
		// if it does not exist return
		core.RaiseErr("package '%s', does not exist", name)
	}
	// check the package is not in the local registry
	localPackage := r.findPackageByRepoAndId(name, remoteArt.Id)
	// if the local registry does not have the package then download it
	if localPackage == nil {
		// download package seal file from registry
		sealFilename, err := api.Download(name.Group, name.Name, fmt.Sprintf("%s.json", remoteArt.FileRef), uname, pwd, tls)
		core.CheckErr(err, "failed to download package seal file")

		// download package file from registry
		packageFilename, err := api.Download(name.Group, name.Name, fmt.Sprintf("%s.zip", remoteArt.FileRef), uname, pwd, tls)
		core.CheckErr(err, "failed to download package file")

		// unmarshal the seal
		sealFile, err := os.Open(sealFilename)
		core.CheckErr(err, "cannot read package seal file")
		seal := new(data.Seal)
		sealBytes, err := ioutil.ReadAll(sealFile)
		// exit if it failed to read the seal
		core.CheckErr(err, "cannot read package seal file")
		// release the handle on the seal
		core.CheckErr(sealFile.Close(), "failed to close seal file stream")
		// unmarshal the seal
		err = json.Unmarshal(sealBytes, seal)
		core.CheckErr(err, "cannot unmarshal package seal file")
		// add the package to the local registry
		r.Add(packageFilename, name, seal)
	} else {
		// the local registry has the package
		// if the local package does not have the tag
		if !localPackage.HasTag(name.Tag) {
			// find the local package coordinates
			rIx, aIx := r.artCoords(name, localPackage)
			// add the tag locally
			r.Repositories[rIx].Packages[aIx].Tags = append(r.Repositories[rIx].Packages[aIx].Tags, name.Tag)
			// persist the changes
			r.save()
			fmt.Printf("package already exist, tag '%s' has been added\n", name.Tag)
		} else {
			// the package exists and has the requested tag
			fmt.Printf("package already exist, tag '%s' already exist, nothing to pull\n", name.Tag)
		}
	}
	return r.FindPackage(name)
}

func (r *LocalRegistry) Open(name *core.PackageName, credentials string, noTLS bool, targetPath string, certPath string, ignoreSignature bool) {
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
	artie := r.FindPackage(name)
	// if not found locally
	if artie == nil {
		// pull it
		artie = r.Pull(name, credentials)
	}
	// get the path to the public key
	if len(pubKeyPath) > 0 {
		if !path.IsAbs(pubKeyPath) {
			pubKeyPath, err = filepath.Abs(pubKeyPath)
			core.CheckErr(err, "cannot retrieve absolute path for public key")
		}
	}
	// get the package seal
	seal, err := r.GetSeal(artie)
	core.CheckErr(err, "cannot read package seal")
	if !ignoreSignature {
		// var pubKey *rsa.PublicKey
		var pgp *crypto.PGP
		if len(pubKeyPath) > 0 {
			// retrieve the verification key from the specified location
			pgp, err = crypto.LoadPGP(pubKeyPath)
			core.CheckErr(err, "cannot load public key, cannot verify signature")
		} else {
			// otherwise load it from the registry store
			pgp, err = crypto.LoadPGPPublicKey(name.Group, name.Name)
			core.CheckErr(err, "cannot load public key, cannot verify signature")
		}
		// get the location of the package
		zipFilename := filepath.Join(core.RegistryPath(), fmt.Sprintf("%s.zip", artie.FileRef))
		// get a slice to have the unencrypted signature
		sum := seal.Checksum(zipFilename)
		// if in debug mode prints out signature
		core.Debug("seal stored base64 encoded signature:\n>> start on next line\n%s\n>> ended on previous line\n", seal.Signature)
		// decode the signature in the seal
		sig, err := base64.StdEncoding.DecodeString(seal.Signature)
		core.CheckErr(err, "cannot decode signature in the seal")
		// if in debug mode prints out base64 decoded signature
		core.Debug("seal stored signature:\n>> start on next line\n%s\n>> ended on previous line\n", string(sig))
		// verify the signature
		err = pgp.Verify(sum, sig)
		core.CheckErr(err, "invalid digital signature")
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
		src := path.Join(r.Path(), fmt.Sprintf("%s.zip", artie.FileRef))
		dst := path.Join(targetPath, filename)
		err = CopyFile(src, dst)
		core.CheckErr(err, "cannot rename package %s", fmt.Sprintf("%s.zip", artie.FileRef))
	} else {
		// otherwise, unzip the target
		err = unzip(path.Join(r.Path(), fmt.Sprintf("%s.zip", artie.FileRef)), targetPath)
		core.CheckErr(err, "cannot unzip package %s", fmt.Sprintf("%s.zip", artie.FileRef))
		// check if the target path is a folder
		info, err := os.Stat(targetPath)
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

func (r *LocalRegistry) Remove(names []*core.PackageName) {
	for _, name := range names {
		// try and get the package by complete URI or id ref
		localPackage := r.FindPackage(name)
		// if the package is not found by name:tag
		if localPackage == nil {
			// try finding it by Id (passed in the name part of the package name)
			list := r.FindPackagesById(name.Name)
			if len(list) == 0 {
				fmt.Printf("name %s not found\n", name.Name)
				continue
			} else {
				// call the remove with the new names
				r.Remove(list)
			}
		} else {
			// try to remove it using full name
			// remove the specified tag
			length := len(localPackage.Tags)
			r.unTag(name, name.Tag)
			// if the tag was successfully deleted
			if len(localPackage.Tags) < length {
				// if there are no tags left at the end then remove the repository
				if len(localPackage.Tags) == 0 {
					r.Repositories = r.removeRepoByName(r.Repositories, name)
					// only remove the files if there are no other repositories containing the same package!
					found := false
				Loop:
					for _, repo := range r.Repositories {
						for _, pack := range repo.Packages {
							if pack.Id == localPackage.Id {
								found = true
								break Loop
							}
						}
					}
					// no other repo contains the package so safe to remove the files
					if !found {
						r.removeFiles(localPackage)
					}
				}
				// persist changes
				r.save()
				log.Print(name)
			} else {
				// attempt to remove by Id (stored in the Name)
				repo := r.findRepository(name)
				repo.Packages = r.removePackageById(repo.Packages, name.Name)
				r.removeFiles(localPackage)
				r.save()
				log.Print(localPackage.Id)
			}
		}
	}
}

func (r *LocalRegistry) GetSeal(name *Package) (*data.Seal, error) {
	sealFilename := path.Join(r.Path(), fmt.Sprintf("%s.json", name.FileRef))
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

func (r *LocalRegistry) ImportKey(keyPath string, isPrivate bool, repoGroup string, repoName string) {
	if !filepath.IsAbs(keyPath) {
		keyPath, err := filepath.Abs(keyPath)
		core.CheckErr(err, "cannot get an absolute representation of path '%s'", keyPath)
	}
	destPath, prefix := r.keyDestinationFolder(repoName, repoGroup)
	// only check it can read the key
	_, err := crypto.LoadPGP(keyPath)
	core.CheckErr(err, "cannot read pgp key '%s'", keyPath)
	// if so, then move the key to the correct location to preserve PEM block data
	if isPrivate {
		CopyFile(keyPath, path.Join(destPath, crypto.PrivateKeyName(prefix, "pgp")))
	} else {
		CopyFile(keyPath, path.Join(destPath, crypto.PublicKeyName(prefix, "pgp")))
	}
}

func (r *LocalRegistry) ExportKey(keyPath string, isPrivate bool, repoGroup string, repoName string) {
	if !filepath.IsAbs(keyPath) {
		keyPath, err := filepath.Abs(keyPath)
		core.CheckErr(err, "cannot get an absolute representation of path '%s'", keyPath)
	}
	destPath, prefix := r.keyDestinationFolder(repoName, repoGroup)
	if isPrivate {
		keyName := crypto.PrivateKeyName(prefix, "pgp")
		err := CopyFile(path.Join(destPath, keyName), path.Join(keyPath, keyName))
		core.CheckErr(err, "cannot export private key")
	} else {
		keyName := crypto.PublicKeyName(prefix, "pgp")
		err := CopyFile(path.Join(destPath, keyName), path.Join(keyPath, keyName))
		core.CheckErr(err, "cannot export public key")
	}
}

func (r *LocalRegistry) GetManifest(name *core.PackageName) *data.Manifest {
	// find the package in the local registry
	a := r.FindPackage(name)
	if a == nil {
		core.RaiseErr("package '%s' not found in the local registry, pull it from remote first", name)
	}
	seal, err := r.GetSeal(a)
	core.CheckErr(err, "cannot get package seal")
	return seal.Manifest
}

// -----------------
// utility functions
// -----------------

// works out the destination folder and prefix for the key
func (r *LocalRegistry) keyDestinationFolder(repoName string, repoGroup string) (destPath string, prefix string) {
	if len(repoName) > 0 {
		// use the repo name location
		destPath = path.Join(r.Path(), "keys", repoGroup, repoName)
		prefix = fmt.Sprintf("%s_%s", repoGroup, repoName)
	} else if len(repoGroup) > 0 {
		// use the repo group location
		destPath = path.Join(r.Path(), "keys", repoGroup)
		prefix = repoGroup
	} else {
		// use the registry root location
		destPath = path.Join(r.Path(), "keys")
		prefix = "root"
	}
	_, err := os.Stat(destPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(destPath, os.ModePerm)
		core.CheckErr(err, "cannot create private key destination '%s'", destPath)
	}
	return destPath, prefix
}

// removes any packages with no tags
func (r *LocalRegistry) removeDangling(name *core.PackageName) {
	repo := r.findRepository(name)
	if repo != nil {
		for _, pack := range repo.Packages {
			// if the package has no tags then remove it
			if len(pack.Tags) == 0 {
				// remove the package metadata using its Id
				repo.Packages = r.removePackageById(repo.Packages, pack.Id)
				// remove the package files
				r.removeFiles(pack)
			}
		}
	}
}

// remove the files associated with an Package
func (r *LocalRegistry) removeFiles(artie *Package) {
	// remove the zip file
	err := os.Remove(fmt.Sprintf("%s/%s.zip", r.Path(), artie.FileRef))
	if err != nil {
		log.Fatal(err)
	}
	// remove the json file
	err = os.Remove(fmt.Sprintf("%s/%s.json", r.Path(), artie.FileRef))
	if err != nil {
		log.Fatal(err)
	}
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
	return fmt.Sprintf("%s/%s.json", r.Path(), uniqueIdName)
}

// the fully qualified name of the zip file in the local localReg
func (r *LocalRegistry) regDirZipFilename(uniqueIdName string) string {
	return fmt.Sprintf("%s/%s.zip", r.Path(), uniqueIdName)
}

// find the package specified by ist id
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
	artie := r.FindPackage(name)
	if artie != nil {
		artie.Tags = core.RemoveElement(artie.Tags, tag)
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

// return the LocalRegistry full file name
func (r *LocalRegistry) file() string {
	return filepath.Join(r.Path(), "repository.json")
}

// save the state of the LocalRegistry
func (r *LocalRegistry) save() {
	regBytes := core.ToJsonBytes(r)
	core.CheckErr(ioutil.WriteFile(r.file(), regBytes, os.ModePerm), "fail to update local registry metadata")
}

// load the content of the LocalRegistry
func (r *LocalRegistry) load() {
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

// check the local localReg directory exists and if not creates it
func (r *LocalRegistry) checkRegistryDir() {
	// check the home directory exists
	_, err := os.Stat(r.Path())
	// if it does not
	if os.IsNotExist(err) {
		if runtime.GOOS == "linux" && os.Geteuid() == 0 {
			core.WarningLogger.Printf("if the root user creates the local registry then runc commands will fail\n" +
				"as the runtime user will not be able to access its content when it is bind mounted\n" +
				"ensure the local registry path is not owned by the root user\n")
		}
		err = os.Mkdir(r.Path(), os.ModePerm)
		i18n.Err(err, i18n.ERR_CANT_CREATE_REGISTRY_FOLDER, r.Path(), core.HomeDir())
	}
	keysPath := path.Join(r.Path(), "keys")
	// check the keys directory exists
	_, err = os.Stat(keysPath)
	// if it does not
	if os.IsNotExist(err) {
		// create a key pair
		err = os.Mkdir(keysPath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		host, _ := os.Hostname()
		crypto.GeneratePGPKeys(keysPath, "root", fmt.Sprintf("root-%s", host), "", "", 2048)
	}
	filesPath := core.FilesPath()
	// check the files directory exists
	_, err = os.Stat(filesPath)
	// if it does not
	if os.IsNotExist(err) {
		// create a key pair
		err = os.Mkdir(filesPath, os.ModePerm)
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
