/*
  Onix Config Manager - Pak
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
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Builder struct {
	zipWriter   *zip.Writer
	workingDir  string
	pakFilename string
	cmds        []string
	labels      map[string]string
	target      string
	pType       string
	repoURI     string
	commit      string
	signer      *sign.Signer
}

func NewBuilder() *Builder {
	// create the builder instance
	builder := &Builder{
		cmds:   []string{},
		labels: make(map[string]string),
	}
	// check the registry directory is there
	builder.checkRegistryDir()
	// retrieve the signing key
	bytes, err := ioutil.ReadFile(builder.inRegistryDirectory("keys/private.pem"))
	if err != nil {
		log.Fatal(err)
	}
	// creates a signer
	signer, err := sign.NewSigner(bytes)
	if err != nil {
		log.Fatal(err)
	}
	builder.signer = signer
	return builder
}

func (b *Builder) Build(repoUrl string, gitToken string) {
	// creates a temporary working directory
	b.newWorkingDir()
	repo := b.clone(repoUrl, gitToken)
	// set the package name
	b.pakName(repo)
	// remove any files in the .pakignore file
	b.removeIgnored()
	// load the pakfile
	b.loadPakfile()
	// run commands
	b.run()
	// compress the target defined in the Pakfile
	b.zipPackage()
	// create seal
	b.seal()
	// cleanup all relevant folders and move package to target location
	b.cleanUp()
}

// compress the target
func (b *Builder) zipPackage() {
	// defines the source for zipping as specified in the Pakfile within the source directory
	source := fmt.Sprintf("%s/%s", b.sourceDir(), b.target)
	// get the target source information
	info, err := os.Stat(source)
	if err != nil {
		log.Fatal(err)
	}
	// if the target is a directory
	if info.IsDir() {
		// then zip it
		zipSource(source, b.pakWDirFullFilename())
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
			zipSource(source, b.pakWDirFullFilename())
			return
		} else {
			// find the file extension
			ext := filepath.Ext(source)
			// if the extension is not zip (e.g. jar files)
			if ext != "zip" {
				// rename the file to .zip
				targetFile := fmt.Sprintf("%s.%s", source[:(len(source)-len(ext))], ext)
				os.Rename(source, targetFile)
				return
			}
			return
		}
	}
}

// clone a remote git repository, it only accepts a token if authentication is required
// if the token is not provided (empty string) then no authentication is used
func (b *Builder) clone(repoUrl string, gitToken string) *git.Repository {
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

// cleanup all relevant folders and move package to target location
func (b *Builder) cleanUp() {
	// remove the zip folder
	b.removeFromWD("pak")
	// move the package to the user home
	b.moveToHome()
	// remove the working directory
	err := os.RemoveAll(b.workingDir)
	if err != nil {
		log.Fatal(err)
	}
	// set the directory to empty
	b.workingDir = ""
}

// move the specified filename from the working directory to the home directory (~/.pak/)
func (b *Builder) moveToHome() {
	// move the .pak file
	err := os.Rename(b.pakWDirFullFilename(), b.pakRegDirFullFilename())
	if err != nil {
		log.Fatal(err)
	}
	// move the .seal file
	err = os.Rename(b.sealWDirFullFilename(), b.sealRegDirFullFilename())
	if err != nil {
		log.Fatal(err)
	}
}

// check the local registry directory exists and if not creates it
func (b *Builder) checkRegistryDir() {
	// check the home directory exists
	_, err := os.Stat(fmt.Sprintf("%s/.pak", b.homeDir()))
	// if it does not
	if os.IsNotExist(err) {
		err = os.Mkdir(fmt.Sprintf("%s/.pak", b.homeDir()), os.ModePerm)
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
	defer zipfile.Close()
	archive := zip.NewWriter(zipfile)
	defer archive.Close()
	info, err := os.Stat(source)
	if err != nil {
		return nil
	}
	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}
	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
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
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})
	return err
}

// construct a unique name for the package using the short HEAD commit hash and current time
func (b *Builder) pakName(repo *git.Repository) {
	ref, err := repo.Head()
	if err != nil {
		log.Fatal(err)
	}
	// get the current time
	t := time.Now()
	timeStamp := fmt.Sprintf("%d%d%s%d%d%d%s", t.Day(), t.Month(), strconv.Itoa(t.Year())[:2], t.Hour(), t.Minute(), t.Second(), strconv.Itoa(t.Nanosecond())[:3])
	b.pakFilename = fmt.Sprintf("%s-%s", timeStamp, ref.Hash().String()[:7])
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

// remove files in the source folder that are specified in the .pakignore file
func (b *Builder) removeIgnored() {
	// retrieve .pakignore
	bytes, err := ioutil.ReadFile(b.inSourceDirectory(".pakignore"))
	if err != nil {
		// assume no .pakignore exists, do nothing
		log.Printf(".pakignore not found")
		return
	}
	// get the lines in the ignore file
	lines := strings.Split(string(bytes), "\n")
	// adds the .pakignore
	lines = append(lines, ".pakignore")
	// loop and remove the included files or folders
	for _, line := range lines {
		path := b.inSourceDirectory(line)
		err := os.RemoveAll(path)
		if err != nil {
			log.Printf("failed to ignore file %s", path)
		}
	}
}

// parse the Pakfile build instructions
func (b *Builder) loadPakfile() {
	// retrieve Pakfile
	bytes, err := ioutil.ReadFile(fmt.Sprintf("%s/Pakfile", b.sourceDir()))
	if err != nil {
		log.Fatal(err)
	}
	// get the lines in the Pakfile
	lines := strings.Split(string(bytes), "\n")
	// loop load info
	for _, line := range lines {
		// add labels
		if strings.HasPrefix(line, "LABEL ") {
			value := line[6:]
			parts := strings.Split(value, "=")
			b.labels[strings.Trim(parts[0], " ")] = strings.Trim(strings.Trim(parts[1], " "), "\"")
		}
		// add commands
		if strings.HasPrefix(line, "RUN ") {
			value := line[4:]
			b.cmds = append(b.cmds, value)
		}
		// add the name of the file or folder being zipped
		if strings.HasPrefix(line, "TARGET ") {
			b.target = line[7:]
		}
		// add the name of the file or folder being zipped
		if strings.HasPrefix(line, "TYPE ") {
			b.pType = line[5:]
		}
	}
}

func (b *Builder) run() {
	for _, cmd := range b.cmds {
		execute(cmd, b.sourceDir())
	}
	b.waitForFileExist(b.inSourceDirectory(b.target), 5*time.Second)
}

// executes a command
func execute(cmd string, dir string) {
	strArr := strings.Split(cmd, " ")
	var c *exec.Cmd
	if len(strArr) == 1 {
		//nolint:gosec
		c = exec.Command(strArr[0])
	} else {
		//nolint:gosec
		c = exec.Command(strArr[0], strArr[1:]...)
	}
	c.Dir = dir
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	log.Printf("executing: %s\n", strings.Join(c.Args, " "))
	if err := c.Start(); err != nil {
		log.Fatal(err)
	}
	err := c.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

// wait a time duration for a file to be created on the path
func (b *Builder) waitForFileExist(path string, d time.Duration) {
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
	return fmt.Sprintf("%s/.pak/%s", b.homeDir(), relativePath)
}

// create the package seal
func (b *Builder) seal() {
	// prepare the seal info
	info := &manifest{
		Type:   b.pType,
		Name:   fmt.Sprintf("%s.zip", b.pakFilename),
		Labels: b.labels,
		Source: b.repoURI,
		Commit: b.commit,
		Branch: "",
		Tag:    "",
		Target: b.target,
		Time:   time.Now().Format(time.RFC850),
	}
	// take the hash of the zip file and seal info combined
	sum := b.checksum(b.pakWDirFullFilename(), info)
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
	dest, err := b.toJsonBytes(s)
	// save the seal
	err = ioutil.WriteFile(b.sealWDirFullFilename(), dest, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

// convert the passed in parameter to a Json Byte Array
func (b *Builder) toJsonBytes(s interface{}) ([]byte, error) {
	// serialise the seal to json
	source, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	// indent the json to make it readable
	dest := new(bytes.Buffer)
	json.Indent(dest, source, "", "  ")
	return dest.Bytes(), err
}

// takes the combined checksum of the seal information and the compressed file
func (b *Builder) checksum(path string, sealData *manifest) []byte {
	// read the compressed file
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// serialise the seal info to json
	info, err := b.toJsonBytes(sealData)
	if err != nil {
		log.Fatal(err)
	}
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
	return fmt.Sprintf("%s/pak", b.workingDir)
}

func (b *Builder) sealWDirFullFilename() string {
	return fmt.Sprintf("%s/%s_seal.json", b.workingDir, b.pakFilename)
}

func (b *Builder) pakWDirFullFilename() string {
	return fmt.Sprintf("%s/%s.zip", b.workingDir, b.pakFilename)
}

func (b *Builder) sealRegDirFullFilename() string {
	return fmt.Sprintf("%s/.pak/%s.json", b.homeDir(), b.pakFilename)
}

func (b *Builder) pakRegDirFullFilename() string {
	return fmt.Sprintf("%s/.pak/%s.zip", b.homeDir(), b.pakFilename)
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
