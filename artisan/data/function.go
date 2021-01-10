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
	"net/url"
	"path/filepath"
	"reflect"
	"strings"
)

type Function struct {
	// the name of the function
	Name string `yaml:"name"`
	// the description for the function
	Description string `yaml:"description,omitempty"`
	// a set of environment variables required by the run commands
	Env map[string]string `yaml:"env,omitempty"`
	// the commands to be executed by the function
	Run []string `yaml:"run,omitempty"`
	// is this function to be available in the manifest
	Export *bool `yaml:"export,omitempty"`
	// defines any input variables required to run this function
	Input *Input `yaml:"input,omitempty"`
}

// describes external input information required by functions or runtimes
type Input struct {
	// required PGP keys
	Key []*Key `yaml:"key,omitempty"`
	// required string value secrets
	Secret []*Secret `yaml:"secret,omitempty"`
	// required variables
	Var []*Var `yaml:"var,omitempty"`
}

// describes PGP keys required by functions
type Key struct {
	// the unique reference for the PGP key
	Name string `yaml:"name"`
	// a description of the intended use of this key
	Description string `yaml:"description"`
	// indicates if the referred key is private or public
	Private bool `yaml:"private"`
	// the artisan package group used to select the key
	PackageGroup string `yaml:"package_group"`
	// the artisan package name used to select the key
	PackageName string `yaml:"package_name"`
	// indicates if this key should be aggregated with other keys
	Aggregate bool `yaml:"aggregate,omitempty"`
	// the key content
	Value string `yaml:"value,omitempty"`
}

// describes the secrets required by functions
type Secret struct {
	// the unique reference for the secret
	Name string `yaml:"name"`
	// a description of the intended use or meaning of this secret
	Description string `yaml:"description"`
	// indicates if this secret should be aggregated with other secrets
	Aggregate bool `yaml:"aggregate,omitempty"`
	// the value of the secret
	Value string `yaml:"value,omitempty"`
}

func (s *Secret) Encrypt(pubKey *crypto.PGP) error {
	encValue, err := pubKey.Encrypt([]byte(s.Value))
	if err != nil {
		return fmt.Errorf("cannot encrypt secret %s: %s", s.Name, err)
	}
	s.Value = base64.StdEncoding.EncodeToString(encValue)
	return nil
}

type Var struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
	Type        string `yaml:"type"`
	Value       string `yaml:"value,omitempty"`
}

// gets a slice of string with each element containing key=value
func (f *Function) GetEnv() map[string]string {
	return f.Env
}

// survey all missing variables in the function
// pass in any available environment variables so that they are not surveyed
func (f *Function) Survey(env map[string]string) map[string]string {
	// merges the function environment with the passed in environment
	for k, v := range f.Env {
		env[k] = v
	}
	// attempt to merge any environment variable in the run commands
	// run the merge in interactive mode so that any variables not available in the build file environment are surveyed
	_, updatedEnvironment := core.MergeEnvironmentVars(f.Run, env, true)
	return updatedEnvironment
}

// go through any defined inputs with no value and prompts the user to complete them
func (f *Function) SurveyInputs() {
	if f.Input != nil {
		if f.Input.Var != nil {
			// survey the variables
			f.surveyInputVar()
		}
		if f.Input.Secret != nil {
			f.surveyInputSecret()
		}
	}
}

// survey function Input.Var section
func (f *Function) surveyInputVar() {
	var validator survey.Validator
	for _, variable := range f.Input.Var {
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
}

func (f *Function) surveyInputSecret() {
	for _, secret := range f.Input.Secret {
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
	}
}

// requires the value conforms to a path
func isPath(val interface{}) error {
	// the reflect value of the result
	value := reflect.ValueOf(val)

	// if the value passed in is a string
	if value.Kind() == reflect.String {
		// try and convert the value to an absolute path
		_, err := filepath.Abs(value.String())
		// if the value cannot be converted to an absolute path
		if err != nil {
			// assumes it is not a valid path
			return fmt.Errorf("value is not a valid path: %s", err)
		}
	} else {
		// if the value is not of a string type it cannot be a path
		return fmt.Errorf("value must be a string")
	}
	return nil
}

// requires the value conforms to a URI
func isURI(val interface{}) error {
	// the reflect value of the result
	value := reflect.ValueOf(val)

	// if the value passed in is a string
	if value.Kind() == reflect.String {
		// try and parse the URI
		_, err := url.ParseRequestURI(value.String())

		// if the value cannot be converted to an absolute path
		if err != nil {
			// assumes it is not a valid path
			return fmt.Errorf("value is not a valid URI: %s", err)
		}
	} else {
		// if the value is not of a string type it cannot be a path
		return fmt.Errorf("value must be a string")
	}
	return nil
}

// requires the value conforms to an Artisan package name
func isPackageName(val interface{}) error {
	// the reflect value of the result
	value := reflect.ValueOf(val)

	// if the value passed in is a string
	if value.Kind() == reflect.String {
		// try and parse the package name
		_, err := core.ParseName(value.String())
		// if the value cannot be parsed
		if err != nil {
			// it is not a valid package name
			return fmt.Errorf("value is not a valid package name: %s", err)
		}
	} else {
		// if the value is not of a string type it cannot be a path
		return fmt.Errorf("value must be a string")
	}
	return nil
}
