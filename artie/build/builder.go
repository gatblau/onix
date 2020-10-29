/*
  Onix Config Manager - Artie
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package build

import (
	"archive/zip"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gatblau/onix/artie/core"
	"github.com/gatblau/onix/artie/registry"
	"github.com/gatblau/onix/artie/sign"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Builder struct {
	zipWriter    *zip.Writer
	workingDir   string
	uniqueIdName string
	repoURI      string
	commit       string
	signer       *sign.Signer
	repoName     *core.ArtieName
	buildFile    *BuildFile
	localReg     *registry.FileRegistry
}

func NewBuilder() *Builder {
	// create the builder instance
	builder := new(Builder)
	// check the localRepo directory is there
	builder.checkRegistryDir()
	// retrieve the signing key
	privateKeyBytes, err := ioutil.ReadFile(builder.inRegistryDirectory("keys/private.pem"))
	if err != nil {
		log.Fatal(err)
	}
	// creates a signer
	signer, err := sign.NewSigner(privateKeyBytes)
	if err != nil {
		log.Fatal(err)
	}
	builder.signer = signer
	builder.localReg = registry.NewFileRegistry()
	return builder
}

// build the artefact
// from: the source to build, either http based git repository or local system git repository
// gitToken: if provided it is used to clone a remote repository that has authentication enabled
// artefactName: the full name of the artefact to be built including the tag
// profileName: the name of the profile to be built. If empty then the default profile is built. If no default profile exists, the first profile is built.
func (b *Builder) Build(from, fromPath, gitToken string, name *core.ArtieName, profileName string) {
	// prepare the source ready for the build
	repo := b.prepareSource(from, fromPath, gitToken, name)
	// set the unique identifier name for both the zip file and the seal file
	b.setUniqueIdName(repo)
	// run commands
	buildProfile := b.run(profileName, fromPath)
	// remove any files in the .buildignore file
	b.removeIgnored()
	// compress the target defined in the build.yaml' profile
	b.zipPackage(buildProfile.Target, fromPath)
	// creates a seal
	s := b.createSeal(buildProfile)
	// add the artefact to the local repo
	b.localReg.Add(b.workDirZipFilename(), b.repoName, s)
	// cleanup all relevant folders and move package to target location
	b.cleanUp()
}

// either clone a remote git repo or copy a local one onto the source folder
func (b *Builder) prepareSource(from string, fromPath string, gitToken string, tagName *core.ArtieName) *git.Repository {
	var (
		repo *git.Repository
	)
	b.repoName = tagName
	// creates a temporary working directory
	b.newWorkingDir()
	// if "from" is an http url
	if strings.HasPrefix(strings.ToLower(from), "http") {
		// clone the remote repo
		repo = b.cloneRepo(from, gitToken)
	} else
	// there is a local repo so copy it to the source folder and then open it
	{
		var localPath = from
		// if a relative path is passed
		if strings.HasPrefix(from, "./") || (!strings.HasPrefix(from, "/")) {
			// turn it into an absolute path
			absPath, err := filepath.Abs(from)
			if err != nil {
				log.Fatal(err)
			}
			localPath = absPath
		}
		// copy the folder to the source directory
		err := copyFiles(localPath, b.sourceDir())
		if err != nil {
			log.Fatal(err)
		}
		b.repoURI = localPath
		repo = b.openRepo()
	}
	loadFrom := b.sourceDir()
	// if a sub-folder was specified
	if len(fromPath) > 0 {
		// add it to the path
		loadFrom = filepath.Join(loadFrom, fromPath)
	}
	// read build.yaml
	b.buildFile = LoadBuildFile(fmt.Sprintf("%s/build.yaml", loadFrom))
	return repo
}

// compress the target
func (b *Builder) zipPackage(target, fromPath string) {
	var targetName = target
	// defines the source for zipping as specified in the build.yaml within the source directory
	source := filepath.Join(b.sourceDir(), targetName)
	// if a sub-project folder exist then add it to the path
	if len(fromPath) > 0 {
		source = filepath.Join(source, fromPath)
	}
	// get the target source information
	info, err := os.Stat(source)
	if err != nil {
		log.Fatal(err)
	}
	// if the target is a directory
	if info.IsDir() {
		// then zip it
		err := zipSource(source, b.workDirZipFilename())
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// if it is a file open it to check its type
		file, err := os.Open(source)
		if err != nil {
			log.Fatal(err)
		}
		// find the content type
		contentType, err := findContentType(file)
		if err != nil {
			log.Fatal(err)
		}
		// if the file is not a zip file
		if contentType != "application/zip" {
			// the zip it
			err := zipSource(source, b.workDirZipFilename())
			if err != nil {
				log.Fatal(err)
			}
			return
		} else {
			// find the file extension
			ext := filepath.Ext(source)
			// if the extension is not zip (e.g. jar files)
			if ext != ".zip" {
				// rename the file to .zip
				err := os.Rename(source, b.workDirZipFilename())
				if err != nil {
					log.Fatal(err)
				}
				return
			}
			return
		}
	}
}

// clones a remote git LocalRegistry, it only accepts a token if authentication is required
// if the token is not provided (empty string) then no authentication is used
func (b *Builder) cloneRepo(repoUrl string, gitToken string) *git.Repository {
	b.repoURI = repoUrl
	// clone the remote repository
	opts := &git.CloneOptions{
		URL:      repoUrl,
		Progress: os.Stdout,
	}
	// if authentication token has been provided
	if len(gitToken) > 0 {
		// The intended use of a GitHub personal access token is in replace of your password
		// because access tokens can easily be revoked.
		// https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
		opts.Auth = &http.BasicAuth{
			Username: "abc123", // yes, this can be anything except an empty string
			Password: gitToken,
		}
	}
	repo, err := git.PlainClone(b.sourceDir(), false, opts)
	if err != nil {
		_ = os.RemoveAll(b.workingDir)
		log.Fatal(err)
	}
	return repo
}

// opens a git LocalRegistry from the given path
func (b *Builder) openRepo() *git.Repository {
	repo, err := git.PlainOpen(b.sourceDir())
	if err != nil {
		log.Fatal(err)
	}
	return repo
}

// cleanup all relevant folders and move package to target location
func (b *Builder) cleanUp() {
	// remove the zip folder
	b.removeFromWD("art")
	// remove the working directory
	err := os.RemoveAll(b.workingDir)
	if err != nil {
		log.Fatal(err)
	}
	// set the directory to empty
	b.workingDir = ""
}

// check the local localReg directory exists and if not creates it
func (b *Builder) checkRegistryDir() {
	// check the home directory exists
	_, err := os.Stat(b.localReg.Path())
	// if it does not
	if os.IsNotExist(err) {
		err = os.Mkdir(b.localReg.Path(), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	// check the keys directory exists
	_, err = os.Stat(b.inRegistryDirectory("keys"))
	// if it does not
	if os.IsNotExist(err) {
		// create a key pair
		err = os.Mkdir(b.inRegistryDirectory("keys"), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		key := sign.NewKeyPair()
		sign.SavePrivateKey(b.inRegistryDirectory("keys/private.pem"), key)
		sign.SavePublicKey(b.inRegistryDirectory("keys/public.pem"), key.PublicKey)
	}
}

// create a new working directory and return its path
func (b *Builder) newWorkingDir() {
	basePath, _ := os.Getwd()
	uid := uuid.New()
	folder := strings.Replace(uid.String(), "-", "", -1)
	workingDirPath := fmt.Sprintf("%s/.%s", basePath, folder)
	// creates a temporary working directory
	err := os.Mkdir(workingDirPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	b.workingDir = workingDirPath
	// create a sub-folder to zip
	err = os.Mkdir(b.sourceDir(), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

// construct a unique name for the package using the short HEAD commit hash and current time
func (b *Builder) setUniqueIdName(repo *git.Repository) {
	ref, err := repo.Head()
	if err != nil {
		log.Fatal(err)
	}
	// get the current time
	t := time.Now()
	timeStamp := fmt.Sprintf("%02d%02d%02s%02d%02d%02d%s", t.Day(), t.Month(), strconv.Itoa(t.Year())[:2], t.Hour(), t.Minute(), t.Second(), strconv.Itoa(t.Nanosecond())[:3])
	b.uniqueIdName = fmt.Sprintf("%s-%s", timeStamp, ref.Hash().String()[:10])
	b.commit = ref.Hash().String()
}

// remove from working directory
func (b *Builder) removeFromWD(path string) {
	err := os.RemoveAll(fmt.Sprintf("%s/%s", b.workingDir, path))
	if err != nil {
		log.Fatal(err)
	}
}

// remove files in the source folder that are specified in the .buildignore file
func (b *Builder) removeIgnored() {
	ignoreFilename := ".buildignore"
	// retrieve the ignore file
	ignoreFileBytes, err := ioutil.ReadFile(b.inSourceDirectory(ignoreFilename))
	if err != nil {
		// assume no ignore file exists, do nothing
		log.Printf("%s not found", ignoreFilename)
		return
	}
	// get the lines in the ignore file
	lines := strings.Split(string(ignoreFileBytes), "\n")
	// adds the .ignore file
	lines = append(lines, ignoreFilename)
	// loop and remove the included files or folders
	for _, line := range lines {
		sourcePath := b.inSourceDirectory(line)
		err := os.RemoveAll(sourcePath)
		if err != nil {
			log.Printf("failed to ignore file %s", sourcePath)
		}
	}
}

// execute all commands in the specified profile
// if not profile is specified, it uses the default profile
// if a default profile has not been defined, then uses the first profile in the build file
// returns the profile used
func (b *Builder) run(profileName string, fromPath string) *Profile {
	// set the command execution directory
	execDir := b.sourceDir()
	// check if there is a sub-folder specified
	if len(fromPath) > 0 {
		execDir = filepath.Join(b.sourceDir(), fromPath)
	}
	// construct an environment with the vars at build file level
	env := append(os.Environ(), b.buildFile.getEnv()...)
	// for each build profile
	for _, profile := range b.buildFile.Profiles {
		// if a profile name has been provided then build it
		if len(profileName) > 0 && profile.Name == profileName {
			// for each run statement in the profile
			for _, cmd := range profile.Run {
				// combine the current environment with the profile environment
				profileEnv := append(env, profile.getEnv()...)
				// execute the statement
				execute(cmd, execDir, profileEnv)
			}
			// wait for the target to be created in the file system
			waitForTargetToBeCreated(b.inSourceDirectory(profile.Target))
			return &profile
		}
		// if the profile has not been provided
		if len(profileName) == 0 {
			// check if a default profile has been set
			defaultProfile := b.buildFile.defaultProfile()
			// use the default profile
			if defaultProfile != nil {
				return b.run(defaultProfile.Name, fromPath)
			} else {
				// there is no default profile defined so use the first profile
				return b.run(b.buildFile.Profiles[0].Name, fromPath)
			}
		}
	}
	log.Fatal(errors.New("no profile was found suitable"))
	return nil
}

// return an absolute path using the working directory as base
func (b *Builder) inWorkingDirectory(relativePath string) string {
	return fmt.Sprintf("%s/%s", b.workingDir, relativePath)
}

// return an absolute path using the source directory as base
func (b *Builder) inSourceDirectory(relativePath string) string {
	return fmt.Sprintf("%s/%s", b.sourceDir(), relativePath)
}

// return an absolute path using the home directory as base
func (b *Builder) inRegistryDirectory(relativePath string) string {
	return fmt.Sprintf("%s/%s", b.localReg.Path(), relativePath)
}

// create the package Seal
func (b *Builder) createSeal(profile *Profile) *core.Seal {
	filename := b.uniqueIdName
	// merge the labels in the profile with the ones at the build file level
	labels := mergeMaps(b.buildFile.Labels, profile.Labels)
	// gets the size of the artefact
	zipInfo, err := os.Stat(b.workDirZipFilename())
	if err != nil {
		log.Fatal(err)
	}
	// prepare the seal info
	info := &core.Manifest{
		Type:    b.buildFile.Type,
		License: b.buildFile.License,
		Ref:     filename,
		Profile: profile.Name,
		Labels:  labels,
		Source:  b.repoURI,
		Commit:  b.commit,
		Branch:  "",
		Tag:     "",
		Target:  profile.Target,
		Time:    time.Now().Format(time.RFC850),
		Size:    bytesToLabel(zipInfo.Size()),
	}
	// take the hash of the zip file and seal info combined
	sum := checksum(b.workDirZipFilename(), info)
	// create a Base-64 encoded cryptographic signature
	signature, err := b.signer.SignBase64(sum)
	if err != nil {
		log.Fatal(err)
	}
	// construct the seal
	s := &core.Seal{
		// the package
		Manifest: info,
		// the combined checksum of the seal info and the package
		Digest: fmt.Sprintf("sha256:%s", base64.StdEncoding.EncodeToString(sum)),
		// the crypto signature
		Signature: signature,
	}
	// convert the seal to Json
	dest := core.ToJsonBytes(s)
	// save the seal
	err = ioutil.WriteFile(b.workDirJsonFilename(), dest, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func (b *Builder) sourceDir() string {
	return fmt.Sprintf("%s/%s", b.workingDir, core.CliName)
}

// the fully qualified name of the json Seal in the working directory
func (b *Builder) workDirJsonFilename() string {
	return fmt.Sprintf("%s/%s.json", b.workingDir, b.uniqueIdName)
}

// the fully qualified name of the zip file in the working directory
func (b *Builder) workDirZipFilename() string {
	return fmt.Sprintf("%s/%s.zip", b.workingDir, b.uniqueIdName)
}
