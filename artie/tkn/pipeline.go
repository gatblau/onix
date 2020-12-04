/*
  Onix Config Manager - Artie
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package tkn

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/asaskevich/govalidator"
	"github.com/gatblau/onix/artie/build"
	"github.com/gatblau/onix/artie/core"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

// a tekton-based Artie's pipeline
type Pipeline struct {
	// ART_APP_NAME
	AppName string
	// ART_GIT_URI
	GitURI string
	// ART_BUILDER_IMG
	BuilderImage string
	// ART_BUILD_PROFILE
	BuildProfile string
	// ART_NAME
	ArtefactName string
	// ART_REG_USER
	ArtefactRegistryUser string
	// ART_REG_PWD
	ArtefactRegistryPwd string
}

// create a new pipeline
func NewPipeline(buildFilePath, buildProfile string) *Pipeline {
	var profile = buildProfile
	buildFile := loadBuildFile(buildFilePath)
	// if no build profile is specified
	if len(buildProfile) == 0 {
		// try and get the default profile
		if buildFile.DefaultProfile() != nil {
			profile = buildFile.DefaultProfile().Name
		} else {
			// uses the first profile
			profile = buildFile.Profiles[0].Name
		}
	}
	// create an instance of the pipeline
	p := new(Pipeline)
	// resolve the builder image using the appType
	p.BuilderImage = builderImage(buildFile.Type)
	// set the build profile
	p.BuildProfile = profile
	// attempt to load the pipeline configuration from the environment
	// NOTE: environment vars can override builder image and/or build profile used (if defined)
	p.loadFromEnv()
	// survey whatever is left undefined
	p.survey()
	// return the configured pipeline
	return p
}

// try and set pipeline variables from the environment
func (p *Pipeline) loadFromEnv() {
	p.AppName = os.Getenv("ART_APP_NAME")
	if len(p.AppName) > 0 {
		fmt.Printf("using ART_APP_NAME=%s\n", p.AppName)
	}
	p.GitURI = os.Getenv("ART_GIT_URI")
	p.BuilderImage = os.Getenv("ART_BUILDER_IMG")
	p.BuildProfile = os.Getenv("ART_BUILD_PROFILE")
	p.ArtefactName = os.Getenv("ART_NAME")
	p.ArtefactRegistryUser = os.Getenv("ART_REG_USER")
	p.ArtefactRegistryPwd = os.Getenv("ART_REG_PWD")
}

// collect missing variables on the command line
func (p *Pipeline) survey() {
	// if the application name is not defined prompt for it
	if len(p.AppName) == 0 {
		prompt := &survey.Input{
			Message: "application name:",
		}
		survey.AskOne(prompt, &p.AppName, survey.WithValidator(survey.Required))
	}
	// if the GIT URI is not defined, prompt for it
	if len(p.GitURI) == 0 {
		prompt := &survey.Input{
			Message: "git repo url:",
		}
		survey.AskOne(prompt, &p.AppName, survey.WithValidator(validURL))
	}
	// if the artefact name is not defined prompt for it
	if len(p.ArtefactName) == 0 {
		prompt := &survey.Input{
			Message: "artefact name:",
		}
		survey.AskOne(prompt, &p.AppName, survey.WithValidator(survey.Required))
	}
	// if the artefact registry user is not defined prompt for it
	if len(p.ArtefactRegistryUser) == 0 {
		prompt := &survey.Input{
			Message: "artefact registry username:",
		}
		survey.AskOne(prompt, &p.ArtefactRegistryUser, survey.WithValidator(survey.Required))
	}
	// if the artefact registry pwd is not defined prompt for it
	if len(p.ArtefactRegistryPwd) == 0 {
		prompt := &survey.Password{
			Message: "artefact registry password:",
		}
		survey.AskOne(prompt, &p.ArtefactRegistryPwd, survey.WithValidator(survey.Required))
	}
}

// merges the template and its values into the passed in writer
func (p *Pipeline) Merge(w io.Writer) error {
	t, err := template.New("pipeline").Parse(pipeline)
	if err != nil {
		return err
	}
	return t.Execute(w, p)
}

// validates url is valid
func validURL(url interface{}) error {
	if str, ok := url.(string); !ok || !govalidator.IsURL(str) {
		return errors.New("a valid URL must be provided")
	}
	return nil
}

// resolve the builder image to use
func builderImage(appType string) string {
	switch strings.ToUpper(appType) {
	case "GOLANG":
		return "quay.io/gatblau/art-go"
	case "JAVA":
		return "quay.io/gatblau/art-java"
	case "NODEJS":
		return "quay.io/gatblau/art-node"
	case "PYTHON":
		return "quay.io/gatblau/art-python"
	default:
		core.RaiseErr("a pipeline for an application of type '%s' is not supported by Artie")
	}
	return ""
}

// load the build file from a file system location
// the passed in path can be either relative or absolute
func loadBuildFile(buildFilePath string) *build.BuildFile {
	filePath, err := core.AbsPath(buildFilePath)
	if err != nil {
		core.RaiseErr("invalid build file path: %s", err.Error())
		return nil
	}
	b, err := ioutil.ReadFile(filePath)
	core.CheckErr(err, "cannot read build file")
	buildFile := new(build.BuildFile)
	err = json.Unmarshal(b, buildFile)
	core.CheckErr(err, "cannot unmarshall build file")
	return buildFile
}

// the pipeline template containing parameterised resource definitions
const pipeline = `
apiVersion: tekton.dev/v1alpha1
kind: Task
metadata:
  name: {{.AppName}}-build-artefacts
spec:
  inputs:
    resources:
      - {type: git, name: source}
  steps:
    - name: apply
      image: {{.BuilderImage} 
      env:
        - name: ARTEFACT_NAME
          value: {{.ArtefactName}
        - name: BUILD_PROFILE
          value: {{.BuildProfile}
        - name: ARTEFACT_UNAME
          value: {{.ArtefactUser}
        - name: ARTEFACT_PWD
          value: {{.ArtefactPwd}
      workingDir: /workspace/source
      volumeMounts:
        - name: config-volume
          mountPath: /keys
  volumes:
    - name: config-volume
      configMap:
        name: sap-config-map
---
apiVersion: tekton.dev/v1alpha1
kind: Pipeline
metadata:
  name: {{.AppName}}-build-and-deploy
spec:
  resources:
  - name: {{.AppName}}-git-repo
    type: git
  params:
  - name: deployment-name
    type: string
    description: name of the deployment to be patched
  tasks:
  - name: build-artefacts
    taskRef:
      name: {{.AppName}}-build-artefacts
    resources:
      inputs:
      - name: source
        resource: {{.AppName}}-git-repo
---
apiVersion: tekton.dev/v1alpha1
kind: PipelineResource
metadata:
  name: {{.AppName}}-git-repo
spec:
  type: git
  params:
  - name: url
    value: {{.GitURI}}
---
apiVersion: tekton.dev/v1alpha1
kind: PipelineRun
metadata:
  name: build-deploy-{{.AppName}}-pipelinerun
spec:
  serviceAccountName: pipeline
  pipelineRef:
    name: {{.AppName}}-build-and-deploy
  resources:
  - name: {{.AppName}}-git-repo
    resourceRef:
      name: {{.AppName}}-git-repo
  params:
  - name: deployment-name
    value: {{.AppName}}
`
