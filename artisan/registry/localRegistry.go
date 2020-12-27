/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package registry

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/crypto"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"
)

// the default local registry implemented as a file system
type LocalRegistry struct {
	Repositories []*Repository `json:"repositories"`
}

func (r *LocalRegistry) api(domain string, useTLS bool) *Api {
	return NewGenericAPI(domain, useTLS)
}

// find the Repository specified by name
func (r *LocalRegistry) findRepository(name *core.ArtieName) *Repository {
	// find repository using artefact name
	for _, repository := range r.Repositories {
		if repository.Repository == name.FullyQualifiedName() {
			return repository
		}
	}
	// find repository using artefact Id
	for _, repository := range r.Repositories {
		for _, artie := range repository.Artefacts {
			if strings.Contains(artie.Id, name.Name) {
				return repository
			}
		}
	}
	return nil
}

// return all the artefacts within the same repository
func (r *LocalRegistry) GetArtefactsByName(name *core.ArtieName) ([]*Artefact, bool) {
	var artefacts = make([]*Artefact, 0)
	for _, repository := range r.Repositories {
		if repository.Repository == name.FullyQualifiedName() {
			for _, artefact := range repository.Artefacts {
				artefacts = append(artefacts, artefact)
			}
			break
		}
	}
	if len(artefacts) > 0 {
		return artefacts, true
	}
	return nil, false
}

// return the artefact that matches the specified:
// - domain/group/name:tag
// nil if not found in the LocalRegistry
func (r *LocalRegistry) FindArtefact(name *core.ArtieName) *Artefact {
	// first gets the repository the artefact is in
	for _, repository := range r.Repositories {
		if repository.Repository == name.FullyQualifiedName() {
			// try and get it by id first
			for _, artefact := range repository.Artefacts {
				for _, tag := range artefact.Tags {
					// try and match against the full URI
					if tag == name.Tag {
						return artefact
					}
				}
			}
			break
		}
	}
	return nil
}

// return the artefacts that matches the specified:
// - artefact id substring
func (r *LocalRegistry) FindArtefactsById(id string) []*core.ArtieName {
	// go through the artefacts in the repository and check for Id matches
	names := make([]*core.ArtieName, 0)
	// first gets the repository the artefact is in
	for _, repository := range r.Repositories {
		for _, artefact := range repository.Artefacts {
			// try and match against the artefact ID substring
			if strings.Contains(artefact.Id, id) {
				for _, tag := range artefact.Tags {
					names = append(names, core.ParseName(fmt.Sprintf("%s:%s", repository.Repository, tag)))
				}
			}
		}
	}
	return names
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

// check the local localReg directory exists and if not creates it
func (r *LocalRegistry) checkRegistryDir() {
	// check the home directory exists
	_, err := os.Stat(r.Path())
	// if it does not
	if os.IsNotExist(err) {
		err = os.Mkdir(r.Path(), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
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
}

// the local Path to the local LocalRegistry
func (r *LocalRegistry) Path() string {
	return core.RegistryPath()
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

// Add the artefact and seal to the LocalRegistry
func (r *LocalRegistry) Add(filename string, name *core.ArtieName, s *core.Seal) {
	core.Msg("adding artefact to local registry: %s", name)
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
	core.CheckErr(RenameFile(filename, filepath.Join(r.Path(), basename), false), "failed to move artefact zip file to the local registry")
	// now move the seal file to the localRepo folder
	core.CheckErr(RenameFile(filepath.Join(basenameDir, fmt.Sprintf("%s.json", basenameNoExt)), filepath.Join(r.Path(), fmt.Sprintf("%s.json", basenameNoExt)), false), "failed to move artefact seal file to the local registry")
	// untag artefact artefact (if any)
	r.unTag(name, name.Tag)
	// remove any dangling artefacts
	r.removeDangling(name)
	// find the repository
	repo := r.findRepository(name)
	// if the repo does not exist the creates it
	if repo == nil {
		repo = &Repository{
			Repository: name.FullyQualifiedName(),
			Artefacts:  make([]*Artefact, 0),
		}
		r.Repositories = append(r.Repositories, repo)
	}
	// creates a new artefact
	artefacts := append(repo.Artefacts, &Artefact{
		Id:      core.ArtefactId(s),
		Type:    s.Manifest.Type,
		FileRef: basenameNoExt,
		Tags:    []string{name.Tag},
		Size:    s.Manifest.Size,
		Created: s.Manifest.Time,
	})
	repo.Artefacts = artefacts
	// persist the changes
	r.save()
}

func (r *LocalRegistry) removeArtefactById(a []*Artefact, id string) []*Artefact {
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

func (r *LocalRegistry) removeRepoByName(a []*Repository, name *core.ArtieName) []*Repository {
	i := -1
	// find an artefact with the specified tag
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

// remove a given tag from an Artefact
func (r *LocalRegistry) unTag(name *core.ArtieName, tag string) {
	artie := r.FindArtefact(name)
	if artie != nil {
		artie.Tags = core.RemoveElement(artie.Tags, tag)
	}
}

// remove a given tag from an artefact
func (r *LocalRegistry) Tag(sourceName *core.ArtieName, targetName *core.ArtieName) {
	sourceArtie := r.FindArtefact(sourceName)
	if sourceArtie == nil {
		core.RaiseErr("source artefact %s does not exit", sourceName)
	}
	if targetName.IsInTheSameRepositoryAs(sourceName) {
		if !sourceArtie.HasTag(targetName.Tag) {
			core.Msg("tagging %s", sourceName)
			sourceArtie.Tags = append(sourceArtie.Tags, targetName.Tag)
			r.save()
			return
		} else {
			core.Msg("already tagged")
			return
		}
	} else {
		targetRepository := r.findRepository(targetName)
		newArtie := *sourceArtie
		// if the target artefact repository does not exist then create it
		if targetRepository == nil {
			core.Msg("tagging %s", sourceName)
			newArtie.Tags = []string{targetName.Tag}
			r.Repositories = append(r.Repositories, &Repository{
				Repository: targetName.FullyQualifiedName(),
				Artefacts: []*Artefact{
					{
						Id:      sourceArtie.Id,
						Type:    sourceArtie.Type,
						FileRef: sourceArtie.FileRef,
						Tags:    []string{targetName.Tag},
						Size:    sourceArtie.Size,
						Created: sourceArtie.Created,
					},
				},
			})
			r.save()
			return
		} else {
			targetArtie := r.FindArtefact(targetName)
			// if the artefact exists in the repository
			if targetArtie != nil {
				// check if the tag already exists
				for _, tag := range targetArtie.Tags {
					if tag == targetName.Tag {
						core.Msg("already tagged")
					} else {
						// add the tag to the existing artefact
						targetArtie.Tags = append(targetArtie.Tags, targetName.Tag)
					}
				}
			} else {
				// check that an artefact with the Id of the source exists
				for _, a := range targetRepository.Artefacts {
					// if the target repository already contains the artefact Id
					if a.Id == sourceArtie.Id {
						// add a tag
						a.Tags = append(a.Tags, targetName.Tag)
						r.save()
						return
					}
				}
				// add a new artefact metadata in the existing repository
				targetRepository.Artefacts = append(targetRepository.Artefacts,
					&Artefact{
						Id:      sourceArtie.Id,
						Type:    sourceArtie.Type,
						FileRef: sourceArtie.FileRef,
						Tags:    []string{targetName.Tag},
						Size:    sourceArtie.Size,
						Created: sourceArtie.Created,
					})
				r.save()
				return
			}
		}
	}
}

// remove all tags from the specified Artefact
func (r *LocalRegistry) unTagAll(name *core.ArtieName) {
	if artefs, exists := r.GetArtefactsByName(name); exists {
		// then it has to untag it, leaving a dangling artefact
		for _, artef := range artefs {
			for _, tag := range artef.Tags {
				artef.Tags = core.RemoveElement(artef.Tags, tag)
			}
		}
	}
	// persist changes
	r.save()
}

// List artefacts to stdout
func (r *LocalRegistry) List() {
	// get a table writer for the stdout
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 12, ' ', 0)
	// print the header row
	_, err := fmt.Fprintln(w, "REPOSITORY\tTAG\tARTEFACT ID\tARTEFACT TYPE\tCREATED\tSIZE")
	core.CheckErr(err, "failed to write table header")
	// repository, tag, artefact id, created, size
	for _, repo := range r.Repositories {
		for _, a := range repo.Artefacts {
			// if the artefact is dangling (no tags)
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

// list (quiet) artefact IDs only
func (r *LocalRegistry) ListQ() {
	// get a table writer for the stdout
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 10, ' ', 0)
	// repository, tag, artefact id, created, size
	for _, repo := range r.Repositories {
		for _, a := range repo.Artefacts {
			_, err := fmt.Fprintln(w, fmt.Sprintf("%s", a.Id[7:19]))
			core.CheckErr(err, "failed to write artefact Id")
		}
	}
	err := w.Flush()
	core.CheckErr(err, "failed to flush output")
}

func (r *LocalRegistry) Push(name *core.ArtieName, credentials string, useTLS bool) {
	// get a reference to the remote registry
	api := r.api(name.Domain, useTLS)
	// get registry credentials
	uname, pwd := core.UserPwd(credentials)
	// fetch the artefact info from the local registry
	artie := r.FindArtefact(name)
	if artie == nil {
		fmt.Printf("artefact %s not found in the local registry\n", name)
		return
	}
	// check the status of the artefact in the remote registry
	remoteArt, err := api.GetArtefactInfo(name.Group, name.Name, artie.Id, uname, pwd)
	core.CheckErr(err, "cannot retrieve remote artefact information")
	// if the artefact exists in the remote registry
	if remoteArt != nil {
		// check if the tag already exist in the remote repository
		if remoteArt.HasTag(name.Tag) {
			// nothing to do, returns
			fmt.Printf("tag already exists, nothing to push\n")
			return
		} else {
			// the metadata has to be updated to include the new tag
			remoteArt.Tags = append(remoteArt.Tags, name.Tag)
			err = api.UpdateArtefactInfo(name, remoteArt, uname, pwd)
			core.CheckErr(err, "cannot update remote artefact tags")
			fmt.Printf("tag %s pushed\n", name.Tag)
			return
		}
	}
	// if the artefact does not exist in the remote registry
	// check if the tag has been applied to another artefact in the repository
	repo, err := api.GetRepositoryInfo(name.Group, name.Name, uname, pwd)
	core.CheckErr(err, "cannot retrieve repository information from backend")
	// if so
	if a, ok := repo.GetTag(name.Tag); ok {
		// remove the tag from the artefact as it will be applied to the new artefact
		a.RemoveTag(name.Tag)
		// if the artefact has no tags left
		if len(a.Tags) == 0 {
			// adds a default tag matching the artefact file reference
			a.Tags = append(a.Tags, a.FileRef)
			// updates the metadata in the remote repo
			core.CheckErr(api.UpdateArtefactInfo(name, a, uname, pwd), "cannot update artefact info")
		}
	}
	zipfile := openFile(fmt.Sprintf("%s/%s.zip", r.Path(), artie.FileRef))
	jsonfile := openFile(fmt.Sprintf("%s/%s.json", r.Path(), artie.FileRef))
	// prepare the artefact to upload
	artefact := artie
	artefact.Tags = []string{name.Tag}
	// execute the upload
	err = api.UploadArtefact(name, artie.FileRef, zipfile, jsonfile, artefact, uname, pwd)
	core.CheckErr(err, "cannot push artefact")
	fmt.Printf("pushed %s\n", name.String())
}

func (r *LocalRegistry) Pull(name *core.ArtieName, credentials string, useTLS bool) *Artefact {
	// get a reference to the remote registry
	api := r.api(name.Domain, useTLS)
	// get registry credentials
	uname, pwd := core.UserPwd(credentials)
	// get remote repository information
	repo, err := api.GetRepositoryInfo(name.Group, name.Name, uname, pwd)
	core.CheckErr(err, "cannot retrieve repository information from backend")
	// find the artefact to pull in the remote repository
	remoteArt, exists := repo.GetTag(name.Tag)
	if !exists {
		// if it does not exist return
		core.RaiseErr("artefact %s, does not exist", name)
	}
	// check the artefact is not in the local registry
	localArt := r.findArtefactByRepoAndId(name, remoteArt.Id)
	// if the local registry does not have the artefact then download it
	if localArt == nil {
		// download artefact seal file from registry
		sealFilename, err := api.Download(name.Group, name.Name, fmt.Sprintf("%s.json", remoteArt.FileRef), uname, pwd)
		core.CheckErr(err, "failed to download artefact seal file")

		// download artefact file from registry
		artieFilename, err := api.Download(name.Group, name.Name, fmt.Sprintf("%s.zip", remoteArt.FileRef), uname, pwd)
		core.CheckErr(err, "failed to download artefact file")

		// unmarshal the seal
		sealFile, err := os.Open(sealFilename)
		core.CheckErr(err, "cannot read artefact seal file")
		seal := new(core.Seal)
		sealBytes, err := ioutil.ReadAll(sealFile)
		core.CheckErr(err, "cannot read artefact seal file")
		err = json.Unmarshal(sealBytes, seal)
		core.CheckErr(err, "cannot unmarshal artefact seal file")

		// add the artefact to the local registry
		r.Add(artieFilename, name, seal)
	} else {
		// the local registry has the artefact
		// if the local artefact does not have the tag
		if !localArt.HasTag(name.Tag) {
			// find the local artefact coordinates
			rIx, aIx := r.artCoords(name, localArt)
			// add the tag locally
			r.Repositories[rIx].Artefacts[aIx].Tags = append(r.Repositories[rIx].Artefacts[aIx].Tags, name.Tag)
			// persist the changes
			r.save()
			fmt.Printf("artefact already exist, tag '%s' has been added\n", name.Tag)
		} else {
			// the artefact exists and has the requested tag
			fmt.Printf("artefact already exist, tag '%s' already exist, nothing to do\n", name.Tag)
		}
	}
	return r.FindArtefact(name)
}

func (r *LocalRegistry) Open(name *core.ArtieName, credentials string, useTLS bool, targetPath string, certPath string, verify bool) {
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
	artie := r.FindArtefact(name)
	// if not found locally
	if artie == nil {
		// pull it
		artie = r.Pull(name, credentials, useTLS)
	}
	// get the path to the public key
	if len(pubKeyPath) > 0 {
		if !path.IsAbs(pubKeyPath) {
			pubKeyPath, err = filepath.Abs(pubKeyPath)
			core.CheckErr(err, "cannot retrieve absolute path for public key")
		}
	}
	// get the artefact seal
	seal, err := r.getSeal(artie)
	core.CheckErr(err, "cannot read artefact seal")
	if verify {
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
		// get the location of the artefact
		zipFilename := filepath.Join(core.RegistryPath(), fmt.Sprintf("%s.zip", artie.FileRef))
		// get a slice to have the unencrypted signature
		sum := core.SealChecksum(zipFilename, seal.Manifest)
		// decode the signature in the seal
		sig, err := base64.StdEncoding.DecodeString(seal.Signature)
		core.CheckErr(err, "cannot decode signature in the seal")
		// verify the signature
		err = pgp.Verify(sum, sig)
		core.CheckErr(err, "invalid digital signature")
	}
	// now we are ready to open it
	// if the target was already compressed (e.g. jar file, etc) then it should not unzip it but rename it
	// to ist original file extension
	if seal.Manifest.Zip {
		_, filename := filepath.Split(seal.Manifest.Target)
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			err = os.MkdirAll(targetPath, os.ModePerm)
			core.CheckErr(err, "cannot create path to open package: %s", targetPath)
		}
		src := path.Join(r.Path(), fmt.Sprintf("%s.zip", artie.FileRef))
		dst := path.Join(targetPath, filename)
		err := CopyFile(src, dst)
		core.CheckErr(err, "cannot rename package %s", fmt.Sprintf("%s.zip", artie.FileRef))
	} else {
		// otherwise unzip the target
		err = unzip(path.Join(r.Path(), fmt.Sprintf("%s.zip", artie.FileRef)), targetPath)
		core.CheckErr(err, "cannot unzip package %s", fmt.Sprintf("%s.zip", artie.FileRef))
	}
}

func (r *LocalRegistry) Remove(names []*core.ArtieName) {
	for _, name := range names {
		// try and get the artefact by complete URI or id ref
		artie := r.FindArtefact(name)
		// if the artefact is not found by name:tag
		if artie == nil {
			// try finding it by Id (passed in the name part of the artefact name)
			list := r.FindArtefactsById(name.Name)
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
			length := len(artie.Tags)
			r.unTag(name, name.Tag)
			// if the tag was successfully deleted
			if len(artie.Tags) < length {
				// if there are no tags left at the end then remove the repository
				if len(artie.Tags) == 0 {
					r.Repositories = r.removeRepoByName(r.Repositories, name)
					// only remove the files if there are no other repositories containing the same artefact!
					found := false
				Loop:
					for _, repo := range r.Repositories {
						for _, art := range repo.Artefacts {
							if art.Id == artie.Id {
								found = true
								break Loop
							}
						}
					}
					// no other repo contains the artefact so safe to remove the files
					if !found {
						r.removeFiles(artie)
					}
				}
				// persist changes
				r.save()
				log.Print(name)
			} else {
				// attempt to remove by Id (stored in the Name)
				repo := r.findRepository(name)
				repo.Artefacts = r.removeArtefactById(repo.Artefacts, name.Name)
				r.removeFiles(artie)
				r.save()
				log.Print(artie.Id)
			}
		}
	}
}

// remove the files associated with an Artefact
func (r *LocalRegistry) removeFiles(artie *Artefact) {
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

// find the artefact specified by ist id
func (r *LocalRegistry) findArtefactByRepoAndId(name *core.ArtieName, id string) *Artefact {
	for _, repository := range r.Repositories {
		rep := fmt.Sprintf("%s/%s/%s", name.Domain, name.Group, name.Name)
		if rep == repository.Repository {
			for _, artefact := range repository.Artefacts {
				if artefact.Id == id {
					return artefact
				}
			}
		}
	}
	return nil
}

// returns the artefact coordinates in the repository as (repo index, artefact index)
func (r *LocalRegistry) artCoords(name *core.ArtieName, art *Artefact) (int, int) {
	for rIx, repository := range r.Repositories {
		rep := fmt.Sprintf("%s/%s/%s", name.Domain, name.Group, name.Name)
		if rep == repository.Repository {
			for aIx, artefact := range repository.Artefacts {
				if artefact.Id == art.Id {
					return rIx, aIx
				}
			}
		}
	}
	return -1, -1
}

func (r *LocalRegistry) getSeal(name *Artefact) (*core.Seal, error) {
	sealFilename := path.Join(r.Path(), fmt.Sprintf("%s.json", name.FileRef))
	sealFile, err := os.Open(sealFilename)
	if err != nil {
		return nil, fmt.Errorf("cannot open seal file %s: %s", sealFilename, err)
	}
	sealBytes, err := ioutil.ReadAll(sealFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read seal file %s: %s", sealFilename, err)
	}
	seal := new(core.Seal)
	err = json.Unmarshal(sealBytes, seal)
	return seal, err
}

func (r *LocalRegistry) ImportKey(keyPath string, isPrivate bool, repoGroup string, repoName string) {
	if !filepath.IsAbs(keyPath) {
		keyPath, err := filepath.Abs(keyPath)
		core.CheckErr(err, "cannot get an absolute representation of path '%s'", keyPath)
	}
	destPath, prefix := r.keyDestinationFolder(repoName, repoGroup)
	key, err := crypto.LoadPGP(keyPath)
	core.CheckErr(err, "cannot read pgp key '%s'", keyPath)
	if isPrivate {
		privateKeyFilename := path.Join(destPath, crypto.PrivateKeyName(prefix, "pgp"))
		key.SavePrivateKey(privateKeyFilename)
	} else {
		publicKeyFilename := path.Join(destPath, crypto.PublicKeyName(prefix, "pgp"))
		key.SavePublicKey(publicKeyFilename)
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

// removes any artefacts with no tags
func (r *LocalRegistry) removeDangling(name *core.ArtieName) {
	repo := r.findRepository(name)
	if repo != nil {
		for _, artefact := range repo.Artefacts {
			// if the artefact has no tags then remove it
			if len(artefact.Tags) == 0 {
				// remove the artefact metadata using its Id
				repo.Artefacts = r.removeArtefactById(repo.Artefacts, artefact.Id)
				// remove the artefact files
				r.removeFiles(artefact)
			}
		}
	}
}
