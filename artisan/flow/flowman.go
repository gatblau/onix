package flow

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/merge"
	"github.com/gatblau/onix/artisan/registry"
	"github.com/gatblau/onix/oxlib/httpserver"
	"github.com/gatblau/oxc"
	"gopkg.in/yaml.v2"
)

const GIT_URI_DESC = "the URI of the GIT repository"
const GIT_BRANCH_DESC = "the branch to be used to clone the project from the GIT repository"
const GIT_USER_DESC = "the user name to be used to authenticate with the GIT repository"
const GIT_PASSWORD_DESC = "the password or token to be used to authenticate with the GIT repository"

// Manager manages an Artisan flow
// the pipeline generator requires at least the flow definition
// if a build file is passed then step variables can be inferred from it
type Manager struct {
	Flow         *Flow
	buildFile    *data.BuildFile
	bareFlowPath string
	env          *merge.Envar
	artHome      string
}

func New(bareFlowPath, buildPath, artHome string) (*Manager, error) {
	// check the flow path to see if bare flow is named correctly
	if !strings.HasSuffix(bareFlowPath, "_bare.yaml") {
		core.RaiseErr("a bare flow is required, the naming convention is [flow_name]_bare.yaml")
	}
	m := &Manager{
		bareFlowPath: bareFlowPath,
		artHome:      artHome,
	}
	flow, err := LoadFlow(bareFlowPath, artHome)
	if err != nil {
		return nil, fmt.Errorf("cannot load flow definition from %s: %s", bareFlowPath, err)
	}
	m.Flow = flow
	// if a build file is defined, then load it
	if len(buildPath) > 0 {
		var buildFile *data.BuildFile
		buildPath = core.ToAbs(buildPath)
		buildFile, err = data.LoadBuildFile(path.Join(buildPath, "build.yaml"))
		if err != nil {
			return nil, fmt.Errorf("cannot load build file from %s: %s", buildPath, err)
		}
		m.buildFile = buildFile
	}
	return m, nil
}

func NewWithEnv(bareFlowPath, buildPath string, env *merge.Envar, artHome string) (*Manager, error) {
	m, err := New(bareFlowPath, buildPath, artHome)
	core.CheckErr(err, "cannot load flow")
	m.env = env
	return m, nil
}

func (m *Manager) Merge(interactive bool) error {
	local := registry.NewLocalRegistry(m.artHome)
	if m.Flow.RequiresGitSource() {
		if m.buildFile == nil {
			return fmt.Errorf("a build.yaml file is required to fill the flow")
		}
		m.populateGit(interactive)
		m.Flow.AppIcon = m.buildFile.AppIcon
	}
	for _, step := range m.Flow.Steps {
		// performs a healthcheck of the flow to determine if it can survey inputs
		flowHealthCheck(m.Flow, step)
		if step.surveyManifest() {
			name, err := core.ParseName(step.Package)
			core.CheckErr(err, "invalid step %s package name %s", step.Name, step.Package)
			// get the package manifest
			manifest := local.GetManifest(name)
			step.Input = data.SurveyInputFromManifest(m.Flow.Name, step.Name, step.PackageSource, name.Domain, step.Function, manifest, interactive, false, m.env, m.artHome)
		} else if step.surveyBuildfile(m.Flow.RequiresGitSource()) {
			// add exported inputs to the step
			step.Input = data.SurveyInputFromBuildFile(step.Function, m.buildFile, interactive, false, m.env, m.artHome)
		}
	}
	err := m.Flow.IsValid()
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) YamlString() (string, error) {
	b, err := yaml.Marshal(m.Flow)
	if err != nil {
		return "", fmt.Errorf("cannot marshal execution flow: %s", err)
	}
	return string(b), nil
}

func (m *Manager) JsonString() (string, error) {
	b, err := json.MarshalIndent(m.Flow, "", "   ")
	if err != nil {
		return "", fmt.Errorf("cannot marshal execution flow: %s", err)
	}
	return string(b), nil
}

func (m *Manager) SaveYAML() error {
	y, err := yaml.Marshal(m.Flow)
	if err != nil {
		return fmt.Errorf("cannot marshal bare flow: %s", err)
	}
	err = ioutil.WriteFile(m.path("yaml"), y, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot save merged flow: %s", err)
	}
	return nil
}

func (m *Manager) SaveJSON() error {
	y, err := json.MarshalIndent(m.Flow, "", "   ")
	if err != nil {
		return fmt.Errorf("cannot marshal bare flow: %s", err)
	}
	err = ioutil.WriteFile(m.path("json"), y, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot save merged flow: %s", err)
	}
	return nil
}

func (m *Manager) SaveOnixJSON() error {
	meta, err := m.Flow.Map()
	if err != nil {
		return fmt.Errorf("cannot convert flow to map: %s", err)
	}
	item := &oxc.Item{
		Key:         fmt.Sprintf("ART_FLOW_%s", m.Flow.Name),
		Name:        m.Flow.Name,
		Description: fmt.Sprintf("defines the execution flow for %s", m.Flow.Name),
		Type:        "ART_FLOW",
		Meta:        meta,
	}
	y, err := json.MarshalIndent(item, "", "   ")
	if err != nil {
		return fmt.Errorf("cannot marshal onix item flow: %s", err)
	}
	err = ioutil.WriteFile(m.path("json"), y, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot save onix item flow: %s", err)
	}
	return nil
}

// Run merge and send a flow to a runner
func (m *Manager) Run(runnerName, creds string, interactive bool) error {
	err := m.Merge(interactive)
	if err != nil {
		return err
	}
	body, err := m.Flow.JsonBytes()
	if err != nil {
		return err
	}
	token := httpserver.BasicToken(core.UserPwd(creds))
	// assume tls enabled
	response, err := m.postFlow(runnerName, err, true, body, token)
	if err != nil {
		// try without using tls
		var err2 error
		response, err2 = m.postFlow(runnerName, err, false, body, token)
		// if succeeded warn the registry is not secured
		if err2 == nil {
			core.WarningLogger.Printf("remote registry does not use TLS - this is a security risk\n")
		} else {
			// if failed not using tls then return the original error
			return err
		}
	}
	if response.StatusCode > 300 {
		bodyBytes, _ := ioutil.ReadAll(response.Body)
		return fmt.Errorf("%s, %s", response.Status, string(bodyBytes))
	}
	// copy the response body to the stdout
	_, err = io.Copy(os.Stdout, response.Body)
	return err
}

func (m *Manager) postFlow(runnerName string, err error, tls bool, body []byte, token string) (*http.Response, error) {
	request, err := http.NewRequest("POST", m.flowURI(runnerName, tls), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", token)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (m *Manager) flowURI(runnerName string, https bool) string {
	scheme := "http"
	if https {
		scheme = fmt.Sprintf("%ss", scheme)
	}
	// {scheme}://{runner-name}/flow
	return fmt.Sprintf("%s://%s/flow", scheme, runnerName)
}

// get the merged Flow path
func (m *Manager) path(extension string) string {
	dir, file := filepath.Split(m.bareFlowPath)
	filename := core.FilenameWithoutExtension(file)
	return filepath.Join(dir, fmt.Sprintf("%s.%s", filename[0:len(filename)-len("_bare")], extension))
}

func (m *Manager) AddLabels(labels []string) {
	for _, label := range labels {
		parts := strings.Split(label, "=")
		if len(parts) != 2 {
			core.RaiseErr("invalid labels")
		}
		m.Flow.Labels[parts[0]] = parts[1]
	}
}

func (m *Manager) populateGit(interactive bool) {
	gitUri := &data.Var{
		Name:        "GIT_URI",
		Description: GIT_URI_DESC,
		Required:    true,
		Type:        "uri",
	}
	data.EvalVar(gitUri, interactive, m.env)
	branch := &data.Var{
		Name:        "GIT_BRANCH",
		Description: GIT_BRANCH_DESC,
		Required:    false,
		Type:        "string",
	}
	data.EvalVar(branch, interactive, m.env)
	gitLogin := &data.Var{
		Name:        "GIT_USER",
		Description: GIT_USER_DESC,
		Required:    false,
		Type:        "string",
	}
	data.EvalVar(gitLogin, interactive, m.env)
	pwd := &data.Var{
		Name:        "GIT_PASSWORD",
		Description: GIT_PASSWORD_DESC,
		Required:    false,
		Type:        "string",
	}
	data.EvalVar(pwd, interactive, m.env)
	m.Flow.Git = &Git{
		Uri:      gitUri.Value,
		Branch:   branch.Value,
		Login:    gitLogin.Value,
		Password: pwd.Value,
	}
}
