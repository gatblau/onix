/*
  Onix Config Manager - Artisan Runner
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package handlers

// @title Artisan Host Runner
// @version 0.0.4
// @description Run Artisan flows in host
// @contact.name gatblau
// @contact.url http://onix.gatblau.org/
// @contact.email onix@gatblau.org
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/runner/host/parser"
	"github.com/gorilla/mux"
)

type S3EventHandler struct {
}

func (h S3EventHandler) HandleEvent(w http.ResponseWriter, r *http.Request) {
	err := downloadAndCopyFile(r)
	if checkErr(w, "error while downloading and copying object from s3 bucket :", err) {
		return
	}

	err = setArtisanRegistryGPGKey(w)
	if checkErr(w, "error while setting private and public gpg key for artisan registry in the host machine :", err) {
		return
	}
}

func setArtisanRegistryGPGKey(w http.ResponseWriter) error {
	ART_REG := ""
	ART_GROUP := ""
	ART_REG_USER := ""
	ART_REG_PWD := ""

	err := setPrivateKey(ART_REG_USER, ART_REG_PWD, ART_REG, ART_GROUP, w)
	if err != nil {
		return err
	}

	err = setPublicKey(ART_REG_USER, ART_REG_PWD, ART_REG, ART_GROUP, w)
	if err != nil {
		return err
	}

	return nil
}

func setPrivateKey(name string, pwd string, reg string, grp string, w http.ResponseWriter) error {
	pkgPullCmd := fmt.Sprintf("art pull %s/%s/keys/pk-registry", reg, grp)
	cred := fmt.Sprintf("-u %s:%s", name, pwd)
	args := []string{cred}

	err := execute(pkgPullCmd, w, args)
	if err != nil {
		msg := fmt.Sprintf("%s: %s\n", "Error while pulling Artisan registry private key  :", err)
		fmt.Printf(msg)
		return errors.New(msg)
	}
	pkgOpenCmd := fmt.Sprintf("art open %s/%s/keys/pk-registry", reg, grp)
	args = []string{fmt.Sprintf("%s %s", cred, "-s")}
	err = execute(pkgOpenCmd, w, args)
	if err != nil {
		msg := fmt.Sprintf("%s: %s\n", "Error while opening Artisan registry private key  :", err)
		fmt.Printf(msg)
		return errors.New(msg)
	}

	importPvtKeyCmd := "art pgp import ecp_reg_rsa_key.pgp"
	args = []string{"-k"}
	err = execute(importPvtKeyCmd, w, args)
	if err != nil {
		msg := fmt.Sprintf("%s: %s\n", "Error while importing Artisan registry private key  :", err)
		fmt.Printf(msg)
		return errors.New(msg)
	}
	return nil
}

func setPublicKey(name string, pwd string, reg string, grp string, w http.ResponseWriter) error {
	pkgPullCmd := fmt.Sprintf("art pull %s/%s/registry-deploy-publickey", reg, grp)
	cred := fmt.Sprintf("-u %s:%s", name, pwd)
	args := []string{cred}
	err := execute(pkgPullCmd, w, args)
	if err != nil {
		msg := fmt.Sprintf("%s: %s\n", "Error while pulling Artisan registry private key  :", err)
		fmt.Printf(msg)
		return errors.New(msg)
	}

	pkgExeCmd := fmt.Sprintf("art exe %s/%s/registry-deploy-publickey %s", reg, grp, "import")
	err = execute(pkgExeCmd, w, nil)
	if err != nil {
		msg := fmt.Sprintf("%s: %s\n", "Error while opening Artisan registry private key  :", err)
		fmt.Printf(msg)
		return errors.New(msg)
	}

	return nil
}

func executeCommand(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	// read path variables
	pkg := vars["package"]
	fun := vars["function"]

	//art exe aps-edge-registry.amosonline.io/aps/patching-package-builder build-package
	ART_REG := ""
	ART_GROUP := ""
	ART_REG_USER := ""
	ART_REG_PWD := ""

	pkgPullCmd := fmt.Sprintf("art pull %s/%s/%s %s", ART_REG, ART_GROUP, pkg)
	cred := fmt.Sprintf("-u %s:%s", ART_REG_USER, ART_REG_PWD)
	args := []string{cred}
	err := execute(pkgPullCmd, w, args)
	if err != nil {
		msg := fmt.Sprintf("%s [ %s %s ] : %s\n", "Error while pulling artisan package using command :", pkgPullCmd, cred, err)
		fmt.Printf(msg)
		return errors.New(msg)
	}

	funcCmd := fmt.Sprintf("art exe %s/%s/%s %s", ART_REG, ART_GROUP, pkg, fun)
	err = execute(funcCmd, w, nil)
	if err != nil {
		msg := fmt.Sprintf("%s [ %s ] : %s\n", "Error while executing artisan package function using command", funcCmd, err)
		fmt.Printf(msg)
		return errors.New(msg)
	}

	return nil
}

func downloadAndCopyFile(r *http.Request) error {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := fmt.Sprintf("%s: %s\n", "cannot read request payload:", err)
		fmt.Printf(msg)
		return errors.New(msg)
	}
	fmt.Println("body", body)
	ev, err := parser.NewS3Event(body)
	if err != nil {
		msg := fmt.Sprintf("%s: %s\n", "Error while parsing s3 event :", err)
		fmt.Printf(msg)
		return errors.New(msg)
	}

	link, err := ev.GetObjectDownloadURL()
	if err != nil {
		msg := fmt.Sprintf("%s: %s\n", "Error while getting object download url from s3 event :", err)
		fmt.Printf(msg)
		return errors.New(msg)
	}

	f, err := core.ReadFile(link, "")
	if err != nil {
		msg := fmt.Sprintf("%s %s: %s\n", "Error while downloading object from s3 bucket using url", link, err)
		fmt.Printf(msg)
		return errors.New(msg)
	}
	//TODO find the value of XLXS_FILE_PATH_WITH_NAME from build.yaml which gives where csv to be copied with name of file
	filePath := ""
	err = core.WriteFile(f, filePath, "")
	if err != nil {
		msg := fmt.Sprintf("%s %s: %s\n", "Error while copying object downloaded from s3 bucket at path ", filePath, err)
		fmt.Printf(msg)
		return errors.New(msg)
	}

	return nil
}
