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
	zipWriter        *zip.Writer
	workingDir       string
	uniqueIdName     string
	repoURI          string
	commit           string
	signer           *sign.Signer
	repoName         *core.ArtieName
	buildFile        *BuildFile
	localReg         *registry.LocalAPI
	shouldCopySource bool
	loadFrom         string
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
	builder.localReg = registry.NewLocalAPI()
	return builder
}

// build the artefact
// from: the source to build, either http based git repository or local system git repository
// gitToken: if provided it is used to clone a remote repository that has authentication enabled
// artefactName: the full name of the artefact to be built including the tag
// profileName: the name of the profile to be built. If empty then the default profile is built. If no default profile exists, the first profile is built.
// copy: indicates whether a copy should be made of the project files before packaging (only valid for from location in the file system)
func (b *Builder) Build(from, fromPath, gitToken string, name *core.ArtieName, profileName string, copy bool) {
	// prepare the source ready for the build
	repo := b.prepareSource(from, fromPath, gitToken, name, copy)
	// set the unique identifier name for both the zip file and the seal file
	b.setUniqueIdName(repo)
	// run commands
	// set the command execution directory
	execDir := b.loadFrom
	buildProfile := b.run(profileName, execDir)
	// check if a profile target exist, otherwise it cannot package
	if len(buildProfile.Target) == 0 {
		core.RaiseErr("profile '%s' target not specified, cannot continue", buildProfile.Name)
	}
	// wait for the target to be created in the file system
	targetPath := filepath.Join(b.loadFrom, buildProfile.Target)
	waitForTargetToBeCreated(targetPath)
	// compress the target defined in the build.yaml' profile
	b.zipPackage(targetPath)
	// creates a seal
	s := b.createSeal(buildProfile)
	// add the artefact to the local repo
	b.localReg.Add(b.workDirZipFilename(), b.repoName, s)
	// cleanup all relevant folders and move package to target location
	b.cleanUp()
}

// only run the commands in the build file
func (b *Builder) Run(profile string, path string) {
	var localPath = path
	// if a relative path is passed
	if strings.HasPrefix(path, "http") {
		core.RaiseErr("the path must not be an http resource")
	}
	if strings.HasPrefix(path, "./") || (!strings.HasPrefix(path, "/")) {
		// turn it into an absolute path
		absPath, err := filepath.Abs(path)
		if err != nil {
			log.Fatal(err)
		}
		localPath = absPath
	}
	core.Msg("loading build instructions")
	b.buildFile = LoadBuildFile(filepath.Join(localPath, "build.yaml"))
	b.run(profile, localPath)
}

// either clone a remote git repo or copy a local one onto the source folder
func (b *Builder) prepareSource(from string, fromPath string, gitToken string, tagName *core.ArtieName, copy bool) *git.Repository {
	var repo *git.Repository
	b.repoName = tagName
	// if "from" is an http url
	if strings.HasPrefix(strings.ToLower(from), "http") {
		// creates a temporary working directory
		b.newWorkingDir()
		b.loadFrom = b.sourceDir()
		// if a sub-folder was specified
		if len(fromPath) > 0 {
			// add it to the path
			b.loadFrom = filepath.Join(b.loadFrom, fromPath)
		}
		// clone the remote repo
		core.Msg("preparing to clone remote repository '%s'", from)
		repo = b.cloneRepo(from, gitToken)
	} else
	// there is a local repo instead of a downloadable url
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
		// if the user requested a copy of the project before building it
		if copy {
			// creates a temporary working directory to copy
			b.newWorkingDir()
			b.loadFrom = b.sourceDir()
			// if a sub-folder was specified
			if len(fromPath) > 0 {
				// add it to the path
				b.loadFrom = filepath.Join(b.loadFrom, fromPath)
			}
			// copy the folder to the source directory
			core.Msg("preparing to copy local repository '%s'", from)
			err := copyFiles(from, b.sourceDir())
			if err != nil {
				log.Fatal(err)
			}
			b.repoURI = localPath
		} else {
			// the working directory is the current directory
			b.loadFrom = localPath
			// if a sub-folder was specified
			if len(fromPath) > 0 {
				// add it to the path
				b.loadFrom = filepath.Join(b.loadFrom, fromPath)
			}
		}
		repo = b.openRepo(localPath)
	}
	// read build.yaml
	core.Msg("loading build instructions")
	b.buildFile = LoadBuildFile(filepath.Join(b.loadFrom, "build.yaml"))
	return repo
}

// compress the target
func (b *Builder) zipPackage(targetPath string) {
	ignored := b.getIgnored()
	core.Msg("compressing target '%s'", targetPath)
	// get the target source information
	info, err := os.Stat(targetPath)
	core.CheckErr(err, "failed to retrieve target to compress: '%s'", targetPath)
	// if the target is a directory
	if info.IsDir() {
		// then zip it
		core.Msg("compressing folder")
		core.CheckErr(zipSource(targetPath, b.workDirZipFilename(), ignored), "failed to compress folder")
	} else {
		// if it is a file open it to check its type
		core.Msg("checking type of file target: '%s'", targetPath)
		file, err := os.Open(targetPath)
		core.CheckErr(err, "failed to open target: %s", targetPath)
		// find the content type
		contentType, err := findContentType(file)
		core.CheckErr(err, "failed to find target content type")
		// if the file is not a zip file
		if contentType != "application/zip" {
			core.Msg("target is not a zip file, proceeding to compress it")
			// the zip it
			core.CheckErr(zipSource(targetPath, b.workDirZipFilename(), ignored), "failed to compress file target")
			return
		} else {
			core.Msg("cannot compress file target, already compressed. checking target file extension")
			// find the file extension
			ext := filepath.Ext(targetPath)
			// if the extension is not zip (e.g. jar files)
			if ext != ".zip" {
				core.Msg("renaming file target to .zip extension")
				// rename the file to .zip
				core.CheckErr(os.Rename(targetPath, b.workDirZipFilename()), "failed to rename file target to .zip extension")
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
func (b *Builder) openRepo(path string) *git.Repository {
	repo, err := git.PlainOpen(path)
	if err != nil {
		log.Fatal(err)
	}
	return repo
}

// cleanup all relevant folders and move package to target location
func (b *Builder) cleanUp() {
	core.Msg("cleaning up temporary build directory")
	// remove the working directory
	core.CheckErr(os.RemoveAll(b.workingDir), "failed to remove temporary build directory")
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
	// the working directory will be a build folder within the registry directory
	basePath := filepath.Join(core.RegistryPath(), "build")
	uid := uuid.New()
	folder := strings.Replace(uid.String(), "-", "", -1)[:12]
	workingDirPath := filepath.Join(basePath, folder)
	// creates a temporary working directory
	err := os.MkdirAll(workingDirPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	b.workingDir = workingDirPath
	// create a sub-folder to zip
	err = os.MkdirAll(b.sourceDir(), os.ModePerm)
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
	core.Msg("creating artefact filename reference '%s'", b.uniqueIdName)
}

// remove files in the source folder that are specified in the .buildignore file
func (b *Builder) getIgnored() []string {
	ignoreFilename := ".buildignore"
	// retrieve the ignore file
	ignoreFileBytes, err := ioutil.ReadFile(filepath.Join(b.loadFrom, ".buildignore"))
	if err != nil {
		// assume no ignore file exists, do nothing
		core.Msg("nothing to remove, %s file not found in project", ignoreFilename)
		return []string{}
	}
	// get the lines in the ignore file
	lines := strings.Split(string(ignoreFileBytes), "\n")
	// adds the .ignore file
	lines = append(lines, ignoreFilename)
	// turns relative paths into absolute paths
	var output []string
	for _, line := range lines {
		if !filepath.IsAbs(line) {
			line, err = filepath.Abs(line)
			if err != nil {
				core.RaiseErr("cannot convert relation path to absolute path: %s", err)
			}
		}
		output = append(output, line)
	}
	return output
}

// execute all commands in the specified profile
// if not profile is specified, it uses the default profile
// if a default profile has not been defined, then uses the first profile in the build file
// returns the profile used
func (b *Builder) run(profileName string, execDir string) *Profile {
	core.Msg("preparing to execute build commands")
	// construct an environment with the vars at build file level
	env := append(os.Environ(), b.buildFile.getEnv()...)
	// for each build profile
	for _, profile := range b.buildFile.Profiles {
		// if a profile name has been provided then build it
		if len(profileName) > 0 && profile.Name == profileName {
			core.Msg("building profile '%s'", profileName)
			// for each run statement in the profile
			for _, cmd := range profile.Run {
				// combine the current environment with the profile environment
				profileEnv := append(env, profile.getEnv()...)
				// execute the statement
				err := execute(cmd, execDir, profileEnv)
				core.CheckErr(err, "cannot execute command: %s", cmd)
			}
			return &profile
		}
		// if the profile has not been provided
		if len(profileName) == 0 {
			// check if a default profile has been set
			defaultProfile := b.buildFile.defaultProfile()
			// use the default profile
			if defaultProfile != nil {
				core.Msg("building the default profile '%s'", defaultProfile.Name)
				return b.run(defaultProfile.Name, execDir)
			} else {
				core.Msg("building the first profile in the build file: '%s'", b.buildFile.Profiles[0].Name)
				// there is no default profile defined so use the first profile
				return b.run(b.buildFile.Profiles[0].Name, execDir)
			}
		}
	}
	// if we got to this point then a specific profile was requested but not defined
	// so cannot continue
	core.RaiseErr("the requested profile '%s' is not defined in artie's build configuration", profileName)
	return nil
}

// return an absolute path using the working directory as base
func (b *Builder) inWorkingDirectory(relativePath string) string {
	return filepath.Join(b.workingDir, relativePath)
}

// return an absolute path using the source directory as base
func (b *Builder) inSourceDirectory(relativePath string) string {
	return filepath.Join(b.sourceDir(), relativePath)
}

// return an absolute path using the home directory as base
func (b *Builder) inRegistryDirectory(relativePath string) string {
	return filepath.Join(b.localReg.Path(), relativePath)
}

// create the package Seal
func (b *Builder) createSeal(profile *Profile) *core.Seal {
	core.Msg("creating artefact seal")
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
	core.Msg("creating artefact cryptographic signature")
	// take the hash of the zip file and seal info combined
	sum := core.SealChecksum(b.workDirZipFilename(), info)
	// create a Base-64 encoded cryptographic signature
	signature, err := b.signer.SignBase64(sum)
	core.CheckErr(err, "failed to create cryptographic signature")
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
	core.CheckErr(ioutil.WriteFile(b.workDirJsonFilename(), dest, os.ModePerm), "failed to write artefact seal file")
	return s
}

func (b *Builder) sourceDir() string {
	return filepath.Join(b.workingDir, core.CliName)
}

// the fully qualified name of the json Seal in the working directory
func (b *Builder) workDirJsonFilename() string {
	return filepath.Join(b.workingDir, fmt.Sprintf("%s.json", b.uniqueIdName))
}

// the fully qualified name of the zip file in the working directory
func (b *Builder) workDirZipFilename() string {
	return filepath.Join(b.workingDir, fmt.Sprintf("%s.zip", b.uniqueIdName))
}

// determine if the from location is a file system path
func (b *Builder) copySource(from string, profile *Profile) bool {
	// location is in the file system and no target is specified for the profile
	// should only run commands where the source is
	return !(!strings.HasPrefix(from, "http") && len(profile.Target) == 0)
}
