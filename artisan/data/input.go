/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package data

import (
	"bytes"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/crypto"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
)

// describes exported input information required by functions or runtimes
type Input struct {
	// reguired by configuration files
	File Files `yaml:"file,omitempty" json:"file,omitempty"`
	// required PGP keys
	Key Keys `yaml:"key,omitempty" json:"key,omitempty"`
	// required string value secrets
	Secret Secrets `yaml:"secret,omitempty" json:"secret,omitempty"`
	// required variables
	Var Vars `yaml:"var,omitempty" json:"var,omitempty"`
}

func (i *Input) HasVarBinding(binding string) bool {
	for _, variable := range i.Var {
		if variable.Name == binding {
			return true
		}
	}
	return false
}

func (i *Input) HasSecretBinding(binding string) bool {
	for _, secret := range i.Secret {
		if secret.Name == binding {
			return true
		}
	}
	return false
}

func (i *Input) HasKeyBinding(binding string) bool {
	for _, key := range i.Key {
		if key.Name == binding {
			return true
		}
	}
	return false
}

func (i *Input) HasVar(name string) bool {
	if i.Var != nil {
		for _, v := range i.Var {
			if v.Name == name {
				return true
			}
		}
	}
	return false
}

func (i *Input) HasSecret(name string) bool {
	if i.Secret != nil {
		for _, s := range i.Secret {
			if s.Name == name {
				return true
			}
		}
	}
	return false
}

func (i *Input) HasKey(name string) bool {
	if i.Key != nil {
		for _, k := range i.Key {
			if k.Name == name {
				return true
			}
		}
	}
	return false
}

func (i *Input) Encrypt(pub *crypto.PGP) {
	encryptInput(i, pub)
}

func (i *Input) SurveyRegistryCreds(flowName, stepName, packageSource, domain string, prompt, defOnly bool, env *core.Envar) {
	if packageSource != "read" {
		// check for art_reg_user
		userName := fmt.Sprintf("%s_%s_OXART_REG_USER", NormInputName(flowName), NormInputName(stepName))
		if !i.HasSecret(userName) {
			userSecret := &Secret{
				Name:        userName,
				Description: fmt.Sprintf("the username to authenticate with the registry at %s'", domain),
			}
			if !defOnly {
				EvalSecret(userSecret, prompt, env)
			}
			i.Secret = append(i.Secret, userSecret)
		}
		// check for art_reg_pwd
		pwd := fmt.Sprintf("%s_%s_OXART_REG_PWD", NormInputName(flowName), NormInputName(stepName))
		if !i.HasSecret(pwd) {
			pwdSecret := &Secret{
				Name:        pwd,
				Description: fmt.Sprintf("the password to authenticate with the registry at '%s'", domain),
			}
			if !defOnly {
				EvalSecret(pwdSecret, prompt, env)
			}
			i.Secret = append(i.Secret, pwdSecret)
		}
		// as we need to open this package a verification (public PGP) key is needed
		keyName := fmt.Sprintf("%s_%s_OXART_VERIFICATION_KEY", NormInputName(flowName), NormInputName(stepName))
		if !i.HasKey(keyName) {
			key := &Key{
				Name:        keyName,
				Description: fmt.Sprintf("the public PGP key required to open the package %s", domain),
				Private:     false,
			}
			if !defOnly {
				EvalKey(key, prompt, env)
			}
			i.Key = append(i.Key, key)
		}
	}
}

func (i *Input) Env(includeKeys bool) *core.Envar {
	env := make(map[string]string)
	for _, v := range i.Var {
		env[v.Name] = v.Value
	}
	for _, s := range i.Secret {
		env[s.Name] = s.Value
	}
	if includeKeys {
		for _, k := range i.Key {
			env[k.Name] = k.Value
		}
	}
	return core.NewEnVarFromMap(env)
}

// merges the passed in input with the current input
func (i *Input) Merge(in *Input) {
	if in == nil {
		// nothing to merge
		return
	}
	for _, v := range in.Var {
		// dedup
		found := false
		for _, iV := range i.Var {
			// if the variable to be merged is already in the source
			if iV.Name == v.Name {
				found = true
				break
			}
		}
		if !found {
			i.Var = append(i.Var, v)
		}
	}
	sort.Sort(i.Var)
	for _, s := range in.Secret {
		// dedup
		found := false
		for _, iS := range i.Secret {
			// if the secret to be merged is already in the source
			if iS.Name == s.Name {
				found = true
				break
			}
		}
		if !found {
			i.Secret = append(i.Secret, s)
		}
	}
	sort.Sort(i.Secret)
	for _, k := range in.Key {
		// dedup
		found := false
		for _, kV := range i.Key {
			// if the key to be merged is already in the source
			if kV.Name == k.Name {
				found = true
				break
			}
		}
		if !found {
			i.Key = append(i.Key, k)
		}
	}
	sort.Sort(i.Key)
}

func (i *Input) VarExist(name string) bool {
	for _, v := range i.Var {
		if v.Name == name {
			return true
		}
	}
	return false
}

func (i *Input) SecretExist(name string) bool {
	for _, s := range i.Secret {
		if s.Name == name {
			return true
		}
	}
	return false
}

// extracts the build file Input that is relevant to a function (using its bindings)
func SurveyInputFromBuildFile(fxName string, buildFile *BuildFile, prompt, defOnly bool, env *core.Envar) *Input {
	if buildFile == nil {
		core.RaiseErr("build file is required")
	}
	// get the build file function to inspect
	fx := buildFile.Fx(fxName)
	if fx == nil {
		core.RaiseErr("function '%s' cannot be found in build file", fxName)
	}
	return getBoundInput(fx.Input, buildFile.Input, prompt, defOnly, env)
}

// extracts the package manifest Input in an exported function
func SurveyInputFromManifest(flowName, stepName, packageSource, domain string, fxName string, manifest *Manifest, prompt, defOnly bool, env *core.Envar) *Input {
	// get the function in the manifest
	fx := manifest.Fx(fxName)
	if fx == nil {
		core.RaiseErr("function '%s' does not exist in or has not been exported", fxName)
	}
	input := fx.Input
	if input == nil {
		input = &Input{
			Key:    make([]*Key, 0),
			Secret: make([]*Secret, 0),
			Var:    make([]*Var, 0),
		}
	}
	// first evaluates the existing inputs
	input = evalInput(input, prompt, defOnly, env)
	// then add registry credential inputs
	input.SurveyRegistryCreds(flowName, stepName, packageSource, domain, prompt, defOnly, env)
	return input
}

// ensure the passed in name is formatted as a valid environment variable name
func NormInputName(name string) string {
	result := strings.Replace(strings.ToUpper(name), "-", "_", -1)
	result = strings.Replace(result, ".", "_", -1)
	return result
}

func SurveyInputFromURI(uri string, prompt, defOnly bool, env *core.Envar) *Input {
	response, err := core.Get(uri, "", "")
	core.CheckErr(err, "cannot fetch runtime manifest")
	body, err := ioutil.ReadAll(response.Body)
	core.CheckErr(err, "cannot read runtime manifest http response")
	// need a wrapper object for the input for the unmarshaller to work so using buildfile
	var buildFile = new(BuildFile)
	err = yaml.Unmarshal(body, buildFile)
	return evalInput(buildFile.Input, prompt, defOnly, env)
}

func evalInput(input *Input, interactive, defOnly bool, env *core.Envar) *Input {
	// makes a shallow copy of the input
	result := *input
	// collect values from command line interface
	for _, v := range result.Var {
		if !defOnly {
			EvalVar(v, interactive, env)
		}
	}
	for _, secret := range result.Secret {
		if !defOnly {
			EvalSecret(secret, interactive, env)
		}
	}
	for _, key := range result.Key {
		if !defOnly {
			EvalKey(key, interactive, env)
		}
	}
	// return pointer to new object
	return &result
}

func EvalVar(inputVar *Var, prompt bool, env *core.Envar) {
	// do not evaluate it if there is already a value
	if len(inputVar.Value) > 0 {
		return
	}
	// check if there is an env variable
	varValue, ok := env.Vars[inputVar.Name]
	// if so
	if ok {
		// set the var value to the env variable's
		inputVar.Value = varValue
	} else if prompt {
		// survey the var value
		surveyVar(inputVar)
	} else {
		// otherwise error
		core.RaiseErr("%s is required", inputVar.Name)
	}
}

func EvalSecret(inputSecret *Secret, prompt bool, env *core.Envar) {
	// do not evaluate it if there is already a value
	if len(inputSecret.Value) > 0 {
		return
	}
	// check if there is an env variable
	secretValue, ok := env.Vars[inputSecret.Name]
	// if so
	if ok {
		// set the secret value to the env variable's
		inputSecret.Value = secretValue
	} else if prompt {
		// survey the secret value
		surveySecret(inputSecret)
	} else {
		// otherwise error
		core.RaiseErr("%s is required", inputSecret.Name)
	}
}

func EvalKey(inputKey *Key, prompt bool, env *core.Envar) {
	// do not evaluate it if there is already a value
	if len(inputKey.Value) > 0 {
		return
	}
	// check if there is an env variable
	keyPath, ok := env.Vars[inputKey.Name]
	// if so
	if ok {
		// load the correct key using the provided path
		loadKeyFromPath(inputKey, keyPath)
	} else if prompt {
		surveyKey(inputKey)
	} else {
		core.RaiseErr("%s is required", inputKey.Name)
	}
}

func EvalFile(inputFile *File, prompt bool, env *core.Envar) {
	// do not evaluate it if there is already a value
	if len(inputFile.Content) > 0 {
		return
	}
	// check if there is an env variable
	keyPath, ok := env.Vars[inputFile.Name]
	// if so
	if ok {
		// load the correct key using the provided path
		loadFileFromPath(inputFile, keyPath)
	} else if prompt {
		surveyFile(inputFile)
	} else {
		core.RaiseErr("%s is required", inputFile.Name)
	}
}

func (i *Input) ToEnvFile() []byte {
	buf := &bytes.Buffer{}
	buf.WriteString("# ===================================================\n")
	buf.WriteString("# VARIABLES\n")
	buf.WriteString("# ===================================================\n")
	for _, v := range i.Var {
		buf.WriteString(toEnvComments(v.Description))
		if len(v.Default) > 0 {
			buf.WriteString(fmt.Sprintf("%s=%s\n", v.Name, v.Default))
		} else {
			buf.WriteString(fmt.Sprintf("%s=\n", v.Name))
		}
	}
	buf.WriteString("\n# ===================================================\n")
	buf.WriteString("# SECRETS\n")
	buf.WriteString("# ===================================================\n")
	for _, s := range i.Secret {
		buf.WriteString(toEnvComments(s.Description))
		buf.WriteString(fmt.Sprintf("%s=\n", s.Name))
	}
	buf.WriteString("\n# ===================================================\n")
	buf.WriteString("# KEY PATHS\n")
	buf.WriteString("# ===================================================\n")
	for _, k := range i.Key {
		buf.WriteString(fmt.Sprint("# the path of the key in the artisan registry as described below:\n"))
		buf.WriteString(toEnvComments(k.Description))
		buf.WriteString(fmt.Sprintf("%s=\n", k.Name))
	}
	return buf.Bytes()
}

func toEnvComments(value string) string {
	out := new(bytes.Buffer)
	values := strings.Split(value, "\n")
	for _, v := range values {
		out.WriteString(fmt.Sprintf("# %s\n", v))
	}
	return out.String()
}

// extract any Input data from the source that have a binding
func getBoundInput(fxInput *InputBinding, sourceInput *Input, prompt, defOnly bool, env *core.Envar) *Input {
	result := &Input{
		Key:    make([]*Key, 0),
		Secret: make([]*Secret, 0),
		Var:    make([]*Var, 0),
		File:   make([]*File, 0),
	}
	// if no bindings then return an empty Input
	if fxInput == nil {
		return result
	}
	// collects exported vars
	for _, varBinding := range fxInput.Var {
		for _, variable := range sourceInput.Var {
			if variable.Name == varBinding {
				result.Var = append(result.Var, variable)
				// if not definition only it should evaluate the variable
				if !defOnly {
					EvalVar(variable, prompt, env)
				}
			}
		}
	}
	// collect exported secrets
	for _, secretBinding := range fxInput.Secret {
		for _, secret := range sourceInput.Secret {
			if secret.Name == secretBinding {
				result.Secret = append(result.Secret, secret)
				// if not definition only it should evaluate the secret
				if !defOnly {
					EvalSecret(secret, prompt, env)
				}
			}
		}
	}
	// collect exported keys
	for _, keyBinding := range fxInput.Key {
		for _, key := range sourceInput.Key {
			if key.Name == keyBinding {
				result.Key = append(result.Key, key)
				// if not definition only it should evaluate the key
				if !defOnly {
					EvalKey(key, prompt, env)
				}
			}
		}
	}
	for _, fileBinding := range fxInput.File {
		for _, file := range sourceInput.File {
			if file.Name == fileBinding {
				result.File = append(result.File, file)
				// if not definition only it should evaluate the file
				if !defOnly {
					EvalFile(file, prompt, env)
				}
			}
		}
	}
	return result
}

// encrypts secret and key values
func encryptInput(input *Input, encPubKey *crypto.PGP) {
	if input == nil {
		return
	}
	for _, secret := range input.Secret {
		// and encrypts the secret value
		err := secret.Encrypt(encPubKey)
		core.CheckErr(err, "cannot encrypt secret")
	}
	for _, key := range input.Key {
		// and encrypts the key value
		err := key.Encrypt(encPubKey)
		core.CheckErr(err, "cannot encrypt PGP key %s: %s", key.Name, err)
	}
}

func surveyVar(variable *Var) {
	// check if an env var has been set
	envVal := os.Getenv(variable.Name)
	// if so, skip survey
	if len(envVal) > 0 {
		return
	}
	// otherwise prompts the user to enter it
	var validator survey.Validator
	desc := ""
	// if a description is available use it
	if len(variable.Description) > 0 {
		desc = variable.Description
	}
	// prompt for the value
	prompt := &survey.Input{
		Message: fmt.Sprintf("var => %s (%s):", variable.Name, desc),
		Default: variable.Default,
	}
	// if required then add required validator
	if variable.Required {
		validator = survey.ComposeValidators(survey.Required)
	}
	// add type validators
	switch strings.ToLower(variable.Type) {
	case "path":
		validator = survey.ComposeValidators(validator, core.IsPath)
	case "uri":
		validator = survey.ComposeValidators(validator, core.IsURI)
	case "name":
		validator = survey.ComposeValidators(validator, core.IsPackageName)
	}
	core.HandleCtrlC(survey.AskOne(prompt, &variable.Value, survey.WithValidator(validator)))
}

func surveySecret(secret *Secret) {
	// check if an env var has been set
	envVal := os.Getenv(secret.Name)
	// if so, skip survey
	if len(envVal) > 0 {
		return
	}
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

func surveyKey(key *Key) {
	// check if an env var has been set
	envVal := os.Getenv(key.Name)
	// if so, skip survey
	if len(envVal) > 0 {
		// load the key using the env var path value specified
		loadKeyFromPath(key, envVal)
		return
	}
	desc := ""
	// if a description is available use it
	if len(key.Description) > 0 {
		desc = key.Description
	}
	// takes default path from input
	defaultPath := key.Path
	// if not defined in input
	if len(defaultPath) == 0 {
		// defaults to root path
		defaultPath = "/"
	}
	// prompt for the value
	prompt := &survey.Input{
		Message: fmt.Sprintf("PGP key => path to %s (%s):", key.Name, desc),
		Default: defaultPath,
		Help:    "/ indicates root keys; /group-name indicates group level keys; /group-name/package-name indicates package level keys",
	}
	var keyPath string
	// survey the key path
	core.HandleCtrlC(survey.AskOne(prompt, &keyPath, survey.WithValidator(keyPathExist)))
	// load the keys
	loadKeyFromPath(key, keyPath)
}

// load the PGP in the key object using the passed in key path
func loadKeyFromPath(key *Key, keyPath string) {
	var (
		pk, pub  string
		keyBytes []byte
		err      error
	)
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
	key.Value = string(keyBytes)
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

func surveyFile(file *File) {
	// check if an env var has been set
	envVal := os.Getenv(file.Name)
	// if so, skip survey
	if len(envVal) > 0 {
		// load the file using the env var path value specified
		loadFileFromPath(file, envVal)
		return
	}
	desc := ""
	// if a description is available use it
	if len(file.Description) > 0 {
		desc = file.Description
	}
	// takes default path from input
	defaultPath := file.Path
	// prompt for the value
	prompt := &survey.Input{
		Message: fmt.Sprintf("File => path to %s (%s):", file.Name, desc),
		Default: defaultPath,
		Help:    "the path to the file to load from the Artisan registry",
	}
	var keyPath string
	// survey the key path
	core.HandleCtrlC(survey.AskOne(prompt, &keyPath, survey.WithValidator(keyPathExist)))
	// load the keys
	loadFileFromPath(file, keyPath)
}

// load the file content in the file object using the passed in file path
func loadFileFromPath(file *File, filePath string) {
	var (
		contentBytes []byte
		err          error
	)
	contentBytes, err = ioutil.ReadFile(path.Join(core.FilesPath(), filePath))
	core.CheckErr(err, "cannot load file from registry")
	file.Content = string(contentBytes)
}
