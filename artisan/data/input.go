/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package data

import (
	"encoding/base64"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/crypto"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
)

// describes exported input information required by functions or runtimes
type Input struct {
	// required PGP keys
	Key []*Key `yaml:"key,omitempty" json:"key,omitempty"`
	// required string value secrets
	Secret []*Secret `yaml:"secret,omitempty" json:"secret,omitempty"`
	// required variables
	Var []*Var `yaml:"var,omitempty" json:"var,omitempty"`
}

func (i *Input) ContainsVar(binding string) bool {
	for _, variable := range i.Var {
		if variable.Name == binding {
			return true
		}
	}
	return false
}

func (i *Input) ContainsSecret(binding string) bool {
	for _, secret := range i.Secret {
		if secret.Name == binding {
			return true
		}
	}
	return false
}

func (i *Input) ContainsKey(binding string) bool {
	for _, key := range i.Key {
		if key.Name == binding {
			return true
		}
	}
	return false
}

// extracts the build file Input that is relevant to a function (using its bindings)
func ExportInput(fxName string, buildFile *BuildFile, encPubKey *crypto.PGP, interactive bool) *Input {
	result := &Input{
		Key:    make([]*Key, 0),
		Secret: make([]*Secret, 0),
		Var:    make([]*Var, 0),
	}
	if buildFile == nil {
		core.RaiseErr("build file is required")
	}
	// get the build file function to inspect
	fx := buildFile.Fx(fxName)
	if fx == nil {
		core.RaiseErr("function %s cannot be found in build file", fxName)
	}
	// if the function has bindings
	if fx.Input != nil {
		// collects exported vars
		for _, varBinding := range fx.Input.Var {
			for _, variable := range buildFile.Input.Var {
				if variable.Name == varBinding {
					result.Var = append(result.Var, variable)
					// if interactive mode is enabled then prompt the user to enter the variable value
					if interactive {
						surveyVar(variable)
					}
				}
			}
		}
		// collect exported secrets
		for _, secretBinding := range fx.Input.Secret {
			for _, secret := range buildFile.Input.Secret {
				if secret.Name == secretBinding {
					result.Secret = append(result.Secret, secret)
					// if interactive mode is enabled then prompt the user to enter the variable value
					if interactive {
						surveySecret(secret, encPubKey)
					}
				}
			}
		}
		// collect exported keys
		for _, keyBinding := range fx.Input.Key {
			for _, key := range buildFile.Input.Key {
				if key.Name == keyBinding {
					result.Key = append(result.Key, key)
					surveyKey(key, encPubKey)
				}
			}
		}
	}
	return result
}

func surveyVar(variable *Var) {
	var validator survey.Validator
	desc := ""
	// if a description is available use it
	if len(variable.Description) > 0 {
		desc = variable.Description
	}
	// prompt for the value
	prompt := &survey.Input{
		Message: fmt.Sprintf("var => %s (%s):", variable.Name, desc),
	}
	// if required then add required validator
	if variable.Required {
		validator = survey.ComposeValidators(survey.Required)
	}
	// add type validators
	switch strings.ToLower(variable.Type) {
	case "path":
		validator = survey.ComposeValidators(validator, isPath)
	case "uri":
		validator = survey.ComposeValidators(validator, isURI)
	case "name":
		validator = survey.ComposeValidators(validator, isPackageName)
	}
	core.HandleCtrlC(survey.AskOne(prompt, &variable.Value, survey.WithValidator(validator)))
}

func surveySecret(secret *Secret, encPubKey *crypto.PGP) {
	desc := ""
	// if a description is available use it
	if len(secret.Description) > 0 {
		desc = secret.Description
	}
	// prompt for the value
	prompt := &survey.Password{
		Message: fmt.Sprintf("secret => %s (%s):", secret.Name, desc),
	}
	core.HandleCtrlC(survey.AskOne(prompt, &secret.Value, survey.WithValidator(survey.Required)))
	// and encrypts the secret value
	err := secret.Encrypt(encPubKey)
	core.CheckErr(err, "cannot encrypt secret")
}

func surveyKey(key *Key, encPubKey *crypto.PGP) {
	desc := ""
	// if a description is available use it
	if len(key.Description) > 0 {
		desc = key.Description
	}
	// prompt for the value
	prompt := &survey.Input{
		Message: fmt.Sprintf("PGP key => %s PATH (%s):", key.Name, desc),
		Default: "/",
		Help:    "/ indicates root keys; /group-name indicates group level keys; /group-name/package-name indicates package level keys",
	}
	var (
		keyPath, pk, pub string
		keyBytes         []byte
		err              error
	)
	core.HandleCtrlC(survey.AskOne(prompt, &keyPath, survey.WithValidator(keyPathExist)))
	// load the keys
	parts := strings.Split(keyPath, "/")
	switch len(parts) {
	case 2:
		// root level keys
		if len(parts[1]) == 0 {
			pk, pub = crypto.KeyNames(core.KeysPath(), "root", "pgp")
			key.PackageGroup = ""
			key.PackageName = ""
		} else {
			// group level keys
			pk, pub = crypto.KeyNames(path.Join(core.KeysPath(), parts[1]), parts[1], "pgp")
			key.PackageGroup = parts[1]
			key.PackageName = ""
		}
	// package level keys
	case 3:
		pk, pub = crypto.KeyNames(path.Join(core.KeysPath(), parts[1], parts[2]), fmt.Sprintf("%s_%s", parts[1], parts[2]), "pgp")
		key.PackageGroup = parts[1]
		key.PackageName = parts[2]
	// error
	default:
		core.RaiseErr("the provided path %s is invalid", keyPath)
	}
	if key.Private {
		keyBytes, err = ioutil.ReadFile(pk)
		core.CheckErr(err, "cannot read private key from registry")
	} else {
		keyBytes, err = ioutil.ReadFile(pub)
		core.CheckErr(err, "cannot read public key from registry")
	}
	// and encrypts the key value
	encValue, err := encPubKey.Encrypt(keyBytes)
	core.CheckErr(err, "cannot encrypt PGP key %s: %s", key.Name, err)
	key.Value = base64.StdEncoding.EncodeToString(encValue)
}

func keyPathExist(val interface{}) error {
	// the reflect value of the result
	value := reflect.ValueOf(val)

	// if the value passed in is a string
	if value.Kind() == reflect.String {
		if len(value.String()) > 0 {
			if !strings.HasPrefix(value.String(), "/") {
				// it is not a valid package name
				return fmt.Errorf("key path '%s' must start with a forward slash", value.String())
			}
			_, err := os.Stat(filepath.Join(core.KeysPath(), value.String()))
			// if the path to the group does not exist
			if os.IsNotExist(err) {
				// it is not a valid package name
				return fmt.Errorf("key path '%s' does not exist", value.String())
			}
		}
	} else {
		// if the value is not of a string type it cannot be a path
		return fmt.Errorf("key group must be a string")
	}
	return nil
}
