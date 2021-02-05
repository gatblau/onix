/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package flow

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/registry"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// the pipeline generator requires at least the flow definition
// if a build file is passed then step variables can be inferred from it
type Manager struct {
	Flow         *Flow
	buildFile    *data.BuildFile
	bareFlowPath string
	envFile      string
}

func New(bareFlowPath, buildPath string) (*Manager, error) {
	// check the flow path to see if bare flow is named correctly
	if !strings.HasSuffix(bareFlowPath, "_bare.yaml") {
		core.RaiseErr("a bare flow is required, the naming convention is [flow_name]_bare.yaml")
	}
	m := &Manager{
		bareFlowPath: bareFlowPath,
	}
	flow, err := LoadFlow(bareFlowPath)
	if err != nil {
		return nil, fmt.Errorf("cannot load flow definition from %s: %s", bareFlowPath, err)
	}
	m.Flow = flow
	// if a build file is defined, then load it
	if len(buildPath) > 0 {
		buildPath = core.ToAbs(buildPath)
		buildFile, err := data.LoadBuildFile(path.Join(buildPath, "build.yaml"))
		if err != nil {
			return nil, fmt.Errorf("cannot load build file from %s: %s", buildPath, err)
		}
		m.buildFile = buildFile
	}
	err = m.validate()
	if err != nil {
		return nil, fmt.Errorf("invalid generator: %s", err)
	}
	return m, nil
}

func NewWithEnv(bareFlowPath, buildPath, envFile string) (*Manager, error) {
	m, err := New(bareFlowPath, buildPath)
	core.CheckErr(err, "cannot load flow")
	m.envFile = core.ToAbs(envFile)
	return m, nil
}

func (m *Manager) Merge(interactive bool) error {
	// load environment variables from file, if file not specified then try loading .env
	core.LoadEnvFromFile(m.envFile)
	local := registry.NewLocalRegistry()
	if m.Flow.RequiresSource() {
		if m.buildFile == nil {
			return fmt.Errorf("a build.yaml file is required to fill the flow")
		}
		// if git uri is not defined
		if len(m.buildFile.GitURI) == 0 {
			// survey its value
			gitUri := &data.Var{
				Name:        "GIT_URI",
				Description: "the URI of the git repository for the project",
				Required:    true,
				Type:        "uri",
			}
			if interactive {
				data.SurveyVar(gitUri)
				// set the git uri in the build file
				m.buildFile.GitURI = gitUri.Value
			} else {
				core.RaiseErr("GIT_URI is required")
			}
		}
		m.Flow.GitURI = m.buildFile.GitURI
		m.Flow.AppIcon = m.buildFile.AppIcon
	}
	for _, step := range m.Flow.Steps {
		step.Runtime = core.QualifyRuntime(step.Runtime)
		if len(step.Package) > 0 {
			name, err := core.ParseName(step.Package)
			core.CheckErr(err, "invalid step %s package name %s", step.Name, step.Package)
			// get the package manifest
			manifest := local.GetManifest(name)
			step.Input = data.SurveyInputFromManifest(name, step.Function, manifest, interactive)
			// collects credentials to retrieve package from registry
			m.surveyRegistryCreds(step.Package, interactive)
		} else {
			// if the step has a function
			if len(step.Function) > 0 {
				// add exported inputs to the step
				step.Input = data.SurveyInputFromBuildFile(step.Function, m.buildFile, interactive)
			} else {
				// read input from from runtime_uri
				step.Input = data.SurveyInputFromURI(step.RuntimeManifest, interactive)
			}
		}
	}
	return nil
}

func (m *Manager) surveyRegistryCreds(packageName string, prompt bool) {
	if m.Flow.Input == nil {
		m.Flow.Input = &data.Input{
			Key:    make([]*data.Key, 0),
			Secret: make([]*data.Secret, 0),
			Var:    make([]*data.Var, 0),
		}
	}
	name, _ := core.ParseName(packageName)
	// check for art_reg_user
	userName := fmt.Sprintf("ART_REG_USER_%s", name.Domain)
	if !m.Flow.HasSecret(userName) {
		userSecret := &data.Secret{
			Name:        userName,
			Description: fmt.Sprintf("the username to authenticate with the registry at '%s'", name.Domain),
		}
		if prompt {
			data.SurveySecret(userSecret)
		} else {
			core.RaiseErr("%s is required", userName)
		}
		m.Flow.Input.Secret = append(m.Flow.Input.Secret, userSecret)
	}
	// check for art_reg_pwd
	pwd := fmt.Sprintf("ART_REG_PWD_%s", name.Domain)
	if !m.Flow.HasSecret(pwd) {
		pwdSecret := &data.Secret{
			Name:        pwd,
			Description: fmt.Sprintf("the password to authenticate with the registry at '%s'", name.Domain),
		}
		if prompt {
			data.SurveySecret(pwdSecret)
		} else {
			core.RaiseErr("%s is required", pwd)
		}
		m.Flow.Input.Secret = append(m.Flow.Input.Secret, pwdSecret)
	}
}

func (m *Manager) YamlString() (string, error) {
	b, err := yaml.Marshal(m.Flow)
	if err != nil {
		return "", fmt.Errorf("cannot marshal execution flow: %s", err)
	}
	return string(b), nil
}

func LoadFlow(path string) (*Flow, error) {
	var err error
	if len(path) == 0 {
		return nil, fmt.Errorf("flow definition is required")
	}
	if !filepath.IsAbs(path) {
		path, err = filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("cannot get absolute path for %s: %s", path, err)
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
	for _, step := range m.Flow.Steps {
		if len(step.Runtime) == 0 {
			return fmt.Errorf("invalid step %s, runtime is missing", step.Name)
		}
	}
	return nil
}

func (m *Manager) Save() error {
	y, err := yaml.Marshal(m.Flow)
	if err != nil {
		return fmt.Errorf("cannot marshal bare flow: %s", err)
	}
	err = ioutil.WriteFile(m.path(), y, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot save merged flow: %s", err)
	}
	return nil
}

// get the merged Flow path
func (m *Manager) path() string {
	dir, file := filepath.Split(m.bareFlowPath)
	filename := core.FilenameWithoutExtension(file)
	return filepath.Join(dir, fmt.Sprintf("%s.yaml", filename[0:len(filename)-len("_bare")]))
}
