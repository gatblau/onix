/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package flow

import (
	"encoding/base64"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/crypto"
	"github.com/gatblau/onix/artisan/data"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// the pipeline generator requires at least the flow definition
// if a build file is passed then step variables can be inferred from it
type Manager struct {
	flow      *Flow
	buildFile *data.BuildFile
	pub       *crypto.PGP
}

func NewFromPath(flowPath, pubKeyPath, buildPath string) (*Manager, error) {
	m := new(Manager)
	flow, err := loadFlow(flowPath)
	if err != nil {
		return nil, fmt.Errorf("cannot load flow definition from %s: %s", flowPath, err)
	}
	m.flow = flow
	// if a build file is defined, then load it
	if len(buildPath) > 0 {
		buildFile, err := data.LoadBuildFile(buildPath)
		if err != nil {
			return nil, fmt.Errorf("cannot load build file from %s: %s", buildPath, err)
		}
		m.buildFile = buildFile
	}
	m.pub, err = crypto.LoadPGP(pubKeyPath)
	core.CheckErr(err, "cannot load public PGP encryption key")
	if m.pub.HasPrivate() {
		return nil, fmt.Errorf("a private PGP key has been provided but a public PGP key is required")
	}
	err = m.validate()
	if err != nil {
		return nil, fmt.Errorf("invalid generator: %s", err)
	}
	return m, nil
}

func NewFromRemote(remotePath string) *Manager {
	return &Manager{}
}

func (m *Manager) FillIn() {
	for _, step := range m.flow.Steps {
		if len(step.Package) > 0 {
			// todo
		} else {
			if len(step.Function) > 0 {
				m.setStepInput(step)
			} else {
				// do nothing
			}
		}
		// m.surveyInputs(step)
	}
}

func (m *Manager) YamlString() (string, error) {
	b, err := yaml.Marshal(m.flow)
	if err != nil {
		return "", fmt.Errorf("cannot marshal execution flow: %s", err)
	}
	return string(b), nil
}

func loadFlow(path string) (*Flow, error) {
	var err error
	if len(path) == 0 {
		return nil, fmt.Errorf("flow definition is required")
	}
	if !filepath.IsAbs(path) {
		path, err = filepath.Abs(path)
		if err != nil {
			fmt.Errorf("cannot get absolute path for %s: %s", path, err)
		}
	}
	flowBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read flow definition %s: %s", path, err)
	}
	flow := new(Flow)
	err = yaml.Unmarshal(flowBytes, flow)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal flow definition %s: %s", path, err)
	}
	return flow, nil
}

func (m *Manager) validate() error {
	// check that the steps have the required attributes set
	for _, step := range m.flow.Steps {
		if len(step.Runtime) == 0 {
			return fmt.Errorf("invalid step %s, runtime is missing", step.Name)
		}
	}
	return nil
}

func (m *Manager) setStepInput(step *Step) {
	if m.buildFile != nil {
		// get the function in question
		fx := m.buildFile.Fx(step.Function)
		// if the function has inputs
		if fx.Input != nil {
			if step.Input == nil {
				step.Input = &data.Input{
					Key:    make([]*data.Key, 0),
					Secret: make([]*data.Secret, 0),
					Var:    make([]*data.Var, 0),
				}
			}
			// add input vars for the exported function
			for _, varBinding := range fx.Input.Var {
				for _, variable := range m.buildFile.Input.Var {
					if variable.Name == varBinding {
						step.Input.Var = append(step.Input.Var, variable)
						m.surveyVar(variable)
					}
				}
			}
			// add input secrets for the exported function
			for _, secretBinding := range fx.Input.Secret {
				for _, secret := range m.buildFile.Input.Secret {
					if secret.Name == secretBinding {
						step.Input.Secret = append(step.Input.Secret, secret)
						m.surveySecret(secret)
					}
				}
			}
			// encrypt the secrets
			step.Input.Secret = m.encryptSecrets(step.Input.Secret)
			// add input keys for the exported function
			for _, keyBinding := range fx.Input.Key {
				for _, key := range m.buildFile.Input.Key {
					if key.Name == keyBinding {
						step.Input.Key = append(step.Input.Key, key)
						m.surveyKey(key)
					}
				}
			}
		}
	}
}

func (m *Manager) encryptSecrets(secrets []*data.Secret) []*data.Secret {
	for _, secret := range secrets {
		err := secret.Encrypt(m.pub)
		core.CheckErr(err, "cannot encrypt secret")
	}
	return secrets
}

func (m *Manager) surveyVar(variable *data.Var) {
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

func (m *Manager) surveySecret(secret *data.Secret) {
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

func (m *Manager) surveyKey(key *data.Key) {
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
	// root level keys
	case 2:
		pk, pub = crypto.KeyNames(core.KeysPath(), "root", "pgp")
	// group level keys
	case 3:
		pk, pub = crypto.KeyNames(core.KeysPath(), parts[1], "pgp")
	// package level keys
	case 4:
		pk, pub = crypto.KeyNames(core.KeysPath(), fmt.Sprintf("%s_%s", parts[1], parts[2]), "pgp")
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
	key.Value = base64.StdEncoding.EncodeToString(keyBytes)
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
