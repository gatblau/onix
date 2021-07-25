package flow

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/merge"
	"github.com/gatblau/onix/artisan/registry"
	"github.com/gatblau/oxc"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Manager manages an Artisan flow
// the pipeline generator requires at least the flow definition
// if a build file is passed then step variables can be inferred from it
type Manager struct {
	Flow         *Flow
	buildFile    *data.BuildFile
	bareFlowPath string
	env          *merge.Envar
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

func NewWithEnv(bareFlowPath, buildPath string, env *merge.Envar) (*Manager, error) {
	m, err := New(bareFlowPath, buildPath)
	core.CheckErr(err, "cannot load flow")
	m.env = env
	return m, nil
}

func (m *Manager) Merge(interactive bool) error {
	local := registry.NewLocalRegistry()
	if m.Flow.RequiresGitSource() {
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
			data.EvalVar(gitUri, interactive, m.env)
			m.buildFile.GitURI = gitUri.Value
		}
		m.Flow.GitURI = m.buildFile.GitURI
		m.Flow.AppIcon = m.buildFile.AppIcon
	}
	for _, step := range m.Flow.Steps {
		step.Runtime = core.QualifyRuntime(step.Runtime)
		// performs a healthcheck of the flow to determine if it can survey inputs
		flowHealthCheck(m.Flow, step)
		if step.surveyManifest() {
			name, err := core.ParseName(step.Package)
			core.CheckErr(err, "invalid step %s package name %s", step.Name, step.Package)
			// get the package manifest
			manifest := local.GetManifest(name)
			step.Input = data.SurveyInputFromManifest(m.Flow.Name, step.Name, step.PackageSource, name.Domain, step.Function, manifest, interactive, false, m.env)
		} else if step.surveyBuildfile(m.Flow.RequiresGitSource()) {
			// add exported inputs to the step
			step.Input = data.SurveyInputFromBuildFile(step.Function, m.buildFile, interactive, false, m.env)
		} else if step.surveyRuntime() {
			// read input from from runtime_uri
			step.Input = data.SurveyInputFromURI(step.RuntimeManifest, interactive, false, m.env)
		}
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
	if flow.Labels == nil {
		flow.Labels = make(map[string]string)
	}
	return flow, nil
}

func NewFlow(flowJSONBytes []byte) (*Flow, error) {
	flow := new(Flow)
	err := json.Unmarshal(flowJSONBytes, flow)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal flow definition %s", err)
	}
	return flow, nil
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
	token := core.BasicToken(core.UserPwd(creds))
	// assume tls enabled
	response, err := m.postFlow(runnerName, err, true, body, token)
	if err != nil {
		// try without using tls
		var err2 error
		response, err2 = m.postFlow(runnerName, err, false, body, token)
		// if succeeded warn the registry is not secured
		if err2 == nil {
			core.Msg("WARNING: remote registry does not use TLS - this is a security risk")
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

func (m *Manager) validate() error {
	// check that the steps have the required attributes set
	for _, step := range m.Flow.Steps {
		if len(step.Runtime) == 0 {
			return fmt.Errorf("invalid step %s, runtime is missing", step.Name)
		}
	}
	return nil
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
