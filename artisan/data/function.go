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
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/crypto"
	"net/url"
	"path/filepath"
	"reflect"
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
	// defines any bindings to inputs required to run this function
	Input *InputBinding `yaml:"input,omitempty"`
	// the runtime to run this function
	Runtime string `yaml:"runtime,omitempty"`
}

// list the names of the inputs required by a function
type InputBinding struct {
	Var    []string `yaml:"var"`
	Secret []string `yaml:"secret"`
	Key    []string `yaml:"key"`
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
	PackageGroup string `yaml:"package_group,omitempty" json:"package_group,omitempty"`
	// the artisan package name used to select the key
	PackageName string `yaml:"package_name,omitempty" json:"package_name,omitempty"`
	// indicates if this key should be aggregated with other keys
	Aggregate bool `yaml:"aggregate,omitempty" json:"aggregate,omitempty"`
	// the key content
	Value string `yaml:"value,omitempty" json:"value,omitempty"`
	// the path to the key in the Artisan registry
	Path string `yaml:"path,omitempty" json:"path,omitempty"`
}

func (k *Key) Encrypt(pubKey *crypto.PGP) error {
	encValue, err := pubKey.Encrypt([]byte(k.Value))
	if err != nil {
		return fmt.Errorf("cannot encrypt PGP key %s: %s", k.Name, err)
	}
	k.Value = base64.StdEncoding.EncodeToString(encValue)
	return nil
}

func (k *Key) Decrypt(privateKey *crypto.PGP) error {
	decoded, err := base64.StdEncoding.DecodeString(k.Value)
	core.CheckErr(err, "cannot base64 decode key '%s'", k.Name)
	decValue, err := privateKey.Decrypt(decoded)
	if err != nil {
		return fmt.Errorf("cannot decrypt PGP key %s: %s", k.Name, err)
	}
	k.Value = string(decValue)
	return nil
}

// describes the secrets required by functions
type Secret struct {
	// the unique reference for the secret
	Name string `yaml:"name" json:"name"`
	// a description of the intended use or meaning of this secret
	Description string `yaml:"description" json:"description"`
	// the value of the secret
	Value string `yaml:"value,omitempty" json:"value,omitempty"`
}

func (s *Secret) Encrypt(pubKey *crypto.PGP) error {
	encValue, err := pubKey.Encrypt([]byte(s.Value))
	if err != nil {
		return fmt.Errorf("cannot encrypt secret %s: %s", s.Name, err)
	}
	s.Value = base64.StdEncoding.EncodeToString(encValue)
	return nil
}

func (s *Secret) Decrypt(pk *crypto.PGP) {
	if !pk.HasPrivate() {
		core.RaiseErr("provided key is not private")
	}
	// decode encrypted value
	decoded, err := base64.StdEncoding.DecodeString(s.Value)
	core.CheckErr(err, "cannot decode encrypted value using base64")
	decValueBytes, err := pk.Decrypt([]byte(decoded))
	core.CheckErr(err, "cannot decrypt secret")
	s.Value = string(decValueBytes)
}

type Var struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
	Required    bool   `yaml:"required" json:"required"`
	Type        string `yaml:"type" json:"type"`
	Value       string `yaml:"value,omitempty" json:"value,omitempty"`
	Default     string `yaml:"default,omitempty" json:"default,omitempty"`
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
