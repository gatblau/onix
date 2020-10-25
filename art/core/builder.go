/*
  Onix Config Manager - Art
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gatblau/onix/pak/sign"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	h "net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const cliName = "art"

type Builder struct {
	zipWriter    *zip.Writer
	workingDir   string
	uniqueIdName string
	repoURI      string
	commit       string
	signer       *sign.Signer
	tagName      string
	buildFile    *BuildFile
}

func NewBuilder() *Builder {
	// create the builder instance
	builder := new(Builder)
	// check the registry directory is there
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
	return builder
}

func (b *Builder) Build(from string, gitToken string, tagName string) {
	// prepare the source ready for the build
	repo := b.prepareSource(from, gitToken, tagName)
	// set the unique identifier name for both the zip file and the seal file
	b.setUniqueIdName(repo)
	// remove any files in the .artignore file
	b.removeIgnored()
	// run commands
	b.run()
	// compress the target(s) defined in the package.yaml
	for _, profile := range b.buildFile.Profiles {
		b.zipPackage(profile.Target, profile.Name)
		b.createSeal(profile)
	}
	// cleanup all relevant folders and move package to target location
	b.cleanUp()
}

// either clone a remote git repo or copy a local one onto the source folder
func (b *Builder) prepareSource(from string, gitToken string, tagName string) *git.Repository {
	var (
		repo *git.Repository
	)
	b.tagName = tagName
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
		err := b.copyFiles(localPath, b.sourceDir())
		if err != nil {
			log.Fatal(err)
		}
		b.repoURI = localPath
		repo = b.openRepo()
	}
	// read package.yaml
	b.buildFile = LoadBuildFile(fmt.Sprintf("%s/package.yaml", b.sourceDir()))
	return repo
}

// compress the target
func (b *Builder) zipPackage(target string, profile string) {
	var targetName = target
	// defines the source for zipping as specified in the package.yaml within the source directory
	source := fmt.Sprintf("%s/%s", b.sourceDir(), targetName)
	// get the target source information
	info, err := os.Stat(source)
	if err != nil {
		log.Fatal(err)
	}
	// if the target is a directory
	if info.IsDir() {
		// then zip it
		err := zipSource(source, b.workDirZipFilename(profile))
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
		contentType, err := b.getFileContentType(file)
		if err != nil {
			log.Fatal(err)
		}
		// if the file is not a zip file
		if contentType != "application/zip" {
			// the zip it
			err := zipSource(source, b.workDirZipFilename(profile))
			if err != nil {
				log.Fatal(err)
			}
			return
		} else {
			// find the file extension
			ext := filepath.Ext(source)
			// if the extension is not zip (e.g. jar files)
			if ext != "zip" {
				// rename the file to .zip
				targetFile := fmt.Sprintf("%s.%s", source[:(len(source)-len(ext))], ext)
				err := os.Rename(source, targetFile)
				if err != nil {
					log.Fatal(err)
				}
				return
			}
			return
		}
	}
}

// clones a remote git repository, it only accepts a token if authentication is required
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

// opens a git repository from the given path
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
	b.removeFromWD("pak")
	// move the package to the user home
	for _, profile := range b.buildFile.Profiles {
		b.moveToHome(profile)
	}
	// remove the working directory
	err := os.RemoveAll(b.workingDir)
	if err != nil {
		log.Fatal(err)
	}
	// set the directory to empty
	b.workingDir = ""
}

// move the specified filename from the working directory to the home directory (~/.pak/)
func (b *Builder) moveToHome(profile Profile) {
	// move the .pak file
	err := os.Rename(b.workDirZipFilename(profile.Name), b.regDirZipFilename(profile.Name))
	if err != nil {
		log.Fatal(err)
	}
	// move the .seal file
	err = os.Rename(b.workDirJsonFilename(profile.Name), b.regDirJsonFilename(profile.Name))
	if err != nil {
		log.Fatal(err)
	}
}

// check the local registry directory exists and if not creates it
func (b *Builder) checkRegistryDir() {
	// check the home directory exists
	_, err := os.Stat(b.registryDirectory())
	// if it does not
	if os.IsNotExist(err) {
		err = os.Mkdir(b.registryDirectory(), os.ModePerm)
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

// zip a file or a folder
func zipSource(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer func() {
		err := zipfile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	archive := zip.NewWriter(zipfile)
	defer func() {
		err := archive.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	info, err := os.Stat(source)
	if err != nil {
		return nil
	}
	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}
	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() {
			err := file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
		_, err = io.Copy(writer, file)
		return err
	})
	return err
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
	b.uniqueIdName = fmt.Sprintf("%s-%s", timeStamp, ref.Hash().String()[:7])
	b.commit = ref.Hash().String()
}

// gets the user home directory
func (b *Builder) homeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

// remove from working directory
func (b *Builder) removeFromWD(path string) {
	err := os.RemoveAll(fmt.Sprintf("%s/%s", b.workingDir, path))
	if err != nil {
		log.Fatal(err)
	}
}

// remove files in the source folder that are specified in the .artignore file
func (b *Builder) removeIgnored() {
	ignoreFilename := ".artignore"
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

// execute all commands in all profiles in sequence
func (b *Builder) run() {
	// construct an environment with the vars at build file level
	env := append(os.Environ(), b.buildFile.getEnv()...)
	// for each build profile
	for _, profile := range b.buildFile.Profiles {
		// for each run statement in the profile
		for _, cmd := range profile.Run {
			// combine the current environment with the profile environment
			profileEnv := append(env, profile.getEnv()...)
			// execute the statement
			execute(cmd, b.sourceDir(), profileEnv)
		}
		// wait for the target to be created in the file system
		b.waitForTargetToBeCreated(b.inSourceDirectory(profile.Target))
	}
}

// executes a single command with arguments
func execute(cmd string, dir string, env []string) {
	strArr := strings.Split(cmd, " ")
	var c *exec.Cmd
	if len(strArr) == 1 {
		c = exec.Command(strArr[0])
	} else {
		c = exec.Command(strArr[0], strArr[1:]...)
	}
	c.Dir = dir
	c.Env = env
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	log.Printf("executing: %s\n", strings.Join(c.Args, " "))
	if err := c.Start(); err != nil {
		log.Fatal(err)
	}
	err := c.Wait()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0
			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				log.Fatal(exitMsg(status.ExitStatus()))
			}
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
		log.Fatal(err)
	}
}

func exitMsg(exitCode int) string {
	switch exitCode {
	case 1:
		return "exit code 1 - general error"
	case 2:
		return "exit code 2 - misuse of shell built-ins"
	case 126:
		return "exit code 126 - command invoked cannot execute"
	case 127:
		return "exit code 127 - command not found"
	case 128:
		return "exit code 128 - invalid argument to exit"
	case 130:
		return "exit code 130 - script terminated by CTRL-C"
	default:
		return fmt.Sprintf("exist code %d", exitCode)
	}
}

// wait a time duration for a file or folder to be created on the path
func (b *Builder) waitForTargetToBeCreated(path string) {
	elapsed := 0
	found := false
	for {
		_, err := os.Stat(path)
		if !os.IsNotExist(err) {
			found = true
			break
		}
		if elapsed > 30 {
			break
		}
		elapsed++
		time.Sleep(500 * time.Millisecond)
	}
	if !found {
		log.Fatal("error: target not found after command execution")
	}
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
	return fmt.Sprintf("%s/%s", b.registryDirectory(), relativePath)
}

// create the package seal
func (b *Builder) createSeal(profile Profile) {
	filename := b.uniqueIdName
	// work out the zip filename
	// if we have a defined name in the profile
	if len(profile.Name) > 0 {
		// append the profile name to the unique Id name
		filename = fmt.Sprintf("%s-%s", b.uniqueIdName, profile.Name)
	}
	// merge the labels in the profile with the ones at the build file level
	labels := b.mergeMaps(b.buildFile.Labels, profile.Labels)
	// prepare the seal info
	info := &manifest{
		Type:    b.buildFile.Type,
		License: b.buildFile.License,
		Name:    fmt.Sprintf("%s.zip", filename),
		Labels:  labels,
		Source:  b.repoURI,
		Commit:  b.commit,
		Branch:  "",
		Tag:     "",
		Target:  profile.Target,
		Time:    time.Now().Format(time.RFC850),
	}
	// take the hash of the zip file and seal info combined
	sum := b.checksum(b.workDirZipFilename(profile.Name), info)
	// create a Base-64 encoded cryptographic signature
	signature, err := b.signer.SignBase64(sum)
	if err != nil {
		log.Fatal(err)
	}
	// construct the seal
	s := &seal{
		// the package
		Manifest: info,
		// the combined checksum of the seal info and the package
		Digest: base64.StdEncoding.EncodeToString(sum),
		// the crypto signature
		Signature: signature,
	}
	// convert the seal to Json
	dest := b.toJsonBytes(s)
	// save the seal
	err = ioutil.WriteFile(b.workDirJsonFilename(profile.Name), dest, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

// convert the passed in parameter to a Json Byte Array
func (b *Builder) toJsonBytes(s interface{}) []byte {
	// serialise the seal to json
	source, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	// indent the json to make it readable
	dest := new(bytes.Buffer)
	err = json.Indent(dest, source, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return dest.Bytes()
}

// takes the combined checksum of the seal information and the compressed file
func (b *Builder) checksum(path string, sealData *manifest) []byte {
	// read the compressed file
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	// serialise the seal info to json
	info := b.toJsonBytes(sealData)
	hash := sha256.New()
	// copy the seal manifest into the hash
	if _, err := io.Copy(hash, bytes.NewReader(info)); err != nil {
		log.Fatal(err)
	}
	// copy the compressed file into the hash
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}
	return hash.Sum(nil)
}

func (b *Builder) sourceDir() string {
	return fmt.Sprintf("%s/%s", b.workingDir, cliName)
}

// the fully qualified name of the json seal in the working directory
func (b *Builder) workDirJsonFilename(profileName string) string {
	if len(profileName) > 0 {
		return fmt.Sprintf("%s/%s-%s.json", b.workingDir, b.uniqueIdName, profileName)
	}
	return fmt.Sprintf("%s/%s.json", b.workingDir, b.uniqueIdName)
}

// the fully qualified name of the zip file in the working directory
// if a profile name is passed in then it is appended to the name
func (b *Builder) workDirZipFilename(profileName string) string {
	if len(profileName) > 0 {
		return fmt.Sprintf("%s/%s-%s.zip", b.workingDir, b.uniqueIdName, profileName)
	}
	return fmt.Sprintf("%s/%s.zip", b.workingDir, b.uniqueIdName)
}

// the fully qualified name of the json seal file in the local registry
func (b *Builder) regDirJsonFilename(profileName string) string {
	if len(profileName) > 0 {
		return fmt.Sprintf("%s/%s-%s.json", b.registryDirectory(), b.uniqueIdName, profileName)
	}
	return fmt.Sprintf("%s/%s.json", b.registryDirectory(), b.uniqueIdName)
}

// the fully qualified name of the zip file in the local registry
func (b *Builder) regDirZipFilename(profileName string) string {
	if len(profileName) > 0 {
		return fmt.Sprintf("%s/%s-%s.zip", b.registryDirectory(), b.uniqueIdName, profileName)
	}
	return fmt.Sprintf("%s/%s.zip", b.registryDirectory(), b.uniqueIdName)
}

// return the art registry directory
func (b *Builder) registryDirectory() string {
	return fmt.Sprintf("%s/.%s", b.homeDir(), cliName)
}

// detect the file content type
func (b *Builder) getFileContentType(f *os.File) (string, error) {
	// get the first 512 bytes to sniff the content type
	buffer := make([]byte, 512)
	_, err := f.Read(buffer)
	if err != nil {
		return "", err
	}
	return h.DetectContentType(buffer), nil
}

// copy a single file
func (b *Builder) copyFile(src, dst string) error {
	var err error
	var srcFd *os.File
	var dstFd *os.File
	var srcInfo os.FileInfo

	if srcFd, err = os.Open(src); err != nil {
		return err
	}
	defer func() {
		err := srcFd.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if dstFd, err = os.Create(dst); err != nil {
		return err
	}
	defer func() {
		err := dstFd.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if _, err = io.Copy(dstFd, srcFd); err != nil {
		return err
	}
	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

// copy the files in a folder recursively
func (b *Builder) copyFiles(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}
	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())
		if fd.IsDir() {
			if err = b.copyFiles(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = b.copyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func (b *Builder) getPathLastSegment(path string) string {
	segments := strings.Split(path, string(os.PathSeparator))
	return segments[len(segments)-1]
}

// merge two or more maps
// the latter map overrides the former if duplicate keys exist across the two maps
func (b *Builder) mergeMaps(ms ...map[string]string) map[string]string {
	res := map[string]string{}
	for _, m := range ms {
		for k, v := range m {
			res[k] = v
		}
	}
	return res
}
