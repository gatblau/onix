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
	"os"

	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/runner/host/parser"
	"github.com/gorilla/mux"
)

type OSpatchingHandler struct {
	accessKey   string
	secretKey   string
	domain      string
	artRegistry string
	artGrp      string
	artUser     string
	artPwd      string
	pkgName     string
	funcName    string
	body        []byte
}

func (h OSpatchingHandler) validate() error {
	msg := ""
	if len(os.Getenv("OBJECT_STORE_ACCESS_KEY")) == 0 {
		msg = fmt.Sprintf("\n%s\n%s", "Error value for environment variable OBJECT_STORE_ACCESS_KEY missing", msg)
	}
	if len(os.Getenv("OBJECT_STORE_SECRET_KEY")) == 0 {
		msg = fmt.Sprintf("\n%s\n%s", "Error value for environment variable OBJECT_STORE_SECRET_KEY missing", msg)
	}
	if len(os.Getenv("OBJECT_STORE_DOMAIN")) == 0 {
		msg = fmt.Sprintf("\n%s\n%s", "Error value for environment variable OBJECT_STORE_DOMAIN missing", msg)
	}
	if len(os.Getenv("ART_REG")) == 0 {
		msg = fmt.Sprintf("\n%s\n%s", "Error value for environment variable ART_REG missing", msg)
	}
	if len(os.Getenv("ART_GROUP")) == 0 {
		msg = fmt.Sprintf("\n%s\n%s", "Error value for environment variable ART_GROUP missing", msg)
	}
	if len(os.Getenv("ART_REG_USER")) == 0 {
		msg = fmt.Sprintf("\n%s\n%s", "Error value for environment variable ART_REG_USER missing", msg)
	}
	if len(os.Getenv("ART_REG_PWD")) == 0 {
		msg = fmt.Sprintf("\n%s\n%s", "Error value for environment variable ART_REG_PWD missing", msg)
	}

	if len(msg) != 0 {
		return errors.New(msg)
	} else {
		return nil
	}

}

func (h *OSpatchingHandler) initialize(r *http.Request) error {

	vars := mux.Vars(r)
	h.pkgName = vars["package"]
	h.funcName = vars["function"]
	h.accessKey = os.Getenv("OBJECT_STORE_ACCESS_KEY")
	h.artGrp = os.Getenv("ART_GROUP")
	h.artPwd = os.Getenv("ART_REG_PWD")
	h.artRegistry = os.Getenv("ART_REG")
	h.artUser = os.Getenv("ART_REG_USER")
	h.domain = os.Getenv("OBJECT_STORE_DOMAIN")
	h.secretKey = os.Getenv("OBJECT_STORE_SECRET_KEY")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := fmt.Sprintf("%s: %s\n", "cannot read request payload:", err)
		fmt.Printf(msg)
		return errors.New(msg)
	}
	h.body = body
	return nil
}

func (h OSpatchingHandler) HandleEvent(w http.ResponseWriter, r *http.Request) {

	err := h.validate()
	if checkErr(w, "validation for prerequired environment variables failed :", err) {
		return
	}

	err = h.initialize(r)
	if checkErr(w, "initialization of OSpatchingHandler failed :", err) {
		return
	}

	err = h.downloadAndCopyFile()
	if checkErr(w, "error while downloading and copying object from s3 bucket :", err) {
		return
	}

	err = h.setArtisanRegistryGPGKey(w)
	if checkErr(w, "error while setting private and public gpg key for artisan registry in the host machine :", err) {
		return
	}

	err = h.executeCommand(w)
	if checkErr(w, "error while executing artisan function in the host machine :", err) {
		return
	}
}

func (h OSpatchingHandler) setArtisanRegistryGPGKey(w http.ResponseWriter) error {

	err := h.setPrivateKey(w)
	if err != nil {
		return err
	}

	err = h.setPublicKey(w)
	if err != nil {
		return err
	}

	return nil
}

func (h OSpatchingHandler) setPrivateKey(w http.ResponseWriter) error {
	pkg := fmt.Sprintf("%s/%s/keys/pk-registry:latest", h.artRegistry, h.artGrp)
	cred := fmt.Sprintf("%s:%s", h.artUser, h.artPwd)
	argsPull := []string{"pull", pkg, "-u", cred}
	err := execute(w, "art", argsPull)
	if err != nil {
		msg := fmt.Sprintf("%s: %s\n", "Error while pulling Artisan registry private key  :", err)
		fmt.Printf(msg)
		return errors.New(msg)
	}

	argsOpen := []string{"open", pkg, "-u", cred, "-s"}
	err = execute(w, "art", argsOpen)
	if err != nil {
		msg := fmt.Sprintf("%s: %s\n", "Error while opening Artisan registry private key  :", err)
		fmt.Printf(msg)
		return errors.New(msg)
	}

	key := "ecp_reg_rsa_key.pgp"
	argsimp := []string{"pgp", "import", key, "-k"}
	err = execute(w, "art", argsimp)
	if err != nil {
		msg := fmt.Sprintf("%s: %s\n", "Error while importing Artisan registry private key  :", err)
		fmt.Printf(msg)
		return errors.New(msg)
	}
	return nil
}

func (h OSpatchingHandler) setPublicKey(w http.ResponseWriter) error {
	pkg := fmt.Sprintf("%s/%s/registry-deploy-publickey", h.artRegistry, h.artGrp)
	cred := fmt.Sprintf("%s:%s", h.artUser, h.artPwd)
	argsPull := []string{"pull", pkg, "-u", cred}
	err := execute(w, "art", argsPull)
	if err != nil {
		msg := fmt.Sprintf("%s: %s\n", "Error while pulling Artisan registry public key  :", err)
		fmt.Printf(msg)
		return errors.New(msg)
	}

	funct := "import"
	argsexe := []string{"exe", pkg, funct}
	err = execute(w, "art", argsexe)
	if err != nil {
		msg := fmt.Sprintf("%s: %s\n", "Error while executing import on Artisan registry public key  :", err)
		fmt.Printf(msg)
		return errors.New(msg)
	}

	return nil
}

func (h OSpatchingHandler) executeCommand(w http.ResponseWriter) error {
	//art exe aps-edge-registry.amosonline.io/aps/patching-package-builder build-package
	pkg := fmt.Sprintf("%s/%s/%s", h.artRegistry, h.artGrp, h.pkgName)
	cred := fmt.Sprintf("%s:%s", h.artUser, h.artPwd)
	argsPull := []string{"pull", pkg, "-u", cred}
	err := execute(w, "art", argsPull)
	if err != nil {
		msg := fmt.Sprintf("%s [ %s %s ] : %s\n", "Error while pulling artisan package using command :", pkg, cred, err)
		fmt.Printf(msg)
		return errors.New(msg)
	}

	argsexe := []string{"exe", pkg, h.funcName}
	err = execute(w, "art", argsexe)
	if err != nil {
		msg := fmt.Sprintf("%s [%s %s ] : %s\n", "Error while executing artisan package function using command", pkg, h.funcName, err)
		fmt.Printf(msg)
		return errors.New(msg)
	}

	return nil
}

func (h OSpatchingHandler) downloadAndCopyFile() error {

	fmt.Printf("\n========================================================\n")
	fmt.Printf("\n $$$ body $$$ \n", string(h.body))
	fmt.Printf("\n========================================================\n")
	ev, err := parser.NewS3Event(h.body)
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

	cred := fmt.Sprintf("%s:%s", h.accessKey, h.secretKey)
	f, err := core.ReadFile(link, cred)
	if err != nil {
		msg := fmt.Sprintf("%s %s: %s\n", "Error while downloading object from s3 bucket using url", link, err)
		fmt.Printf(msg)
		return errors.New(msg)
	}
	//TODO find the value of XLXS_FILE_PATH_WITH_NAME from build.yaml which gives where csv to be copied with name of file
	filePath := "/home/ubuntu/env/patch.csv"
	err = core.WriteFile(f, filePath, "")
	if err != nil {
		msg := fmt.Sprintf("%s %s: %s\n", "Error while copying file downloaded from S3 bucket ", filePath, err)
		fmt.Printf(msg)
		return errors.New(msg)
	}

	return nil
}
