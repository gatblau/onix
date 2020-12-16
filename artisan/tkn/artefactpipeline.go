/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package tkn

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/asaskevich/govalidator"
	"github.com/gatblau/onix/artisan/build"
	"github.com/gatblau/onix/artisan/core"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"
)

// a tekton-based Artie's CI pipeline
type ArtefactPipelineConfig struct {
	// PIPE_ART_APP_NAME
	AppName string
	// PIPE_ART_GIT_URI
	GitURI string
	// PIPE_ART_BUILDER_IMG
	BuilderImage string
	// PIPE_ART_BUILD_PROFILE
	BuildProfile string
	// PIPE_ART_NAME
	ArtefactName string
	// PIPE_ART_REG_USER
	ArtefactRegistryUser string
	// PIPE_ART_REG_PWD
	ArtefactRegistryPwd string
	// PIPE_ART_APP_ICON
	AppIcon string
	// PIPE_ART_SONAR_URI
	SonarURI string
	// PIPE_ART_SONAR_TOKEN
	SonarToken string
	// PIPE_ART_SONAR_IMAGE
	SonarImage string
	// PIPE_ART_SONAR_PROJ_KEY
	SonarProjectKey string
	// PIPE_ART_SONAR_SOURCES
	SonarSources string
	// PIPE_ART_SONAR_BINARIES
	SonarBinaries string
}

// create a new pipeline
func NewArtPipelineConfig(buildFilePath, profileName string, sonar bool) *ArtefactPipelineConfig {
	var profile *build.Profile
	// load the build file
	buildFile := loadBuildFile(buildFilePath)
	// if no build profile is specified
	if len(profileName) == 0 {
		// try and get the default profile
		if buildFile.DefaultProfile() != nil {
			profile = buildFile.DefaultProfile()
		} else {
			// uses the first profile
			profile = buildFile.Profiles[0]
		}
	} else {
		// pick the specified profile
		profile = buildFile.Profile(profileName)
		// if the profile is not in the build file
		if profile == nil {
			core.RaiseErr("profile '%s' not found", profileName)
			os.Exit(1)
		}
	}
	// create an instance of the pipeline
	p := new(ArtefactPipelineConfig)
	// resolve the builder image using the appType
	p.BuilderImage = builderImage(buildFile.Type)
	// set the build profile
	p.BuildProfile = profile.Name
	// set the application name
	p.AppName = buildFile.Application
	p.AppIcon = buildFile.Icon
	p.ArtefactName = buildFile.Artefact
	// if sonar step is required and there is a Sonar configuration section in buildfile
	if sonar && buildFile.Sonar != nil {
		p.SonarURI = buildFile.Sonar.URI
		p.SonarProjectKey = buildFile.Sonar.ProjectKey
		p.SonarSources = buildFile.Sonar.Sources
		p.SonarBinaries = buildFile.Sonar.Binaries
	}
	// attempt to load the pipeline configuration from the environment
	// NOTE: environment vars can override builder image and/or build profile used (if defined)
	p.loadFromEnv(sonar)
	// survey any variables in the pipeline that has been left undefined
	p.survey(sonar)
	// finally survey any missing variables in the build profile that are not defined
	profile.Survey(buildFile)
	// return the configured pipeline
	return p
}

// try and set ciPipeline variables from the environment
func (p *ArtefactPipelineConfig) loadFromEnv(sonar bool) {
	p.AppName = p.LoadVar("PIPE_ART_APP_NAME", p.AppName)
	p.AppIcon = p.LoadVar("PIPE_ART_PIPE_APP_ICON", p.AppIcon)
	p.GitURI = p.LoadVar("PIPE_ART_GIT_URI", p.GitURI)
	p.BuilderImage = p.LoadVar("PIPE_ART_BUILDER_IMG", p.BuilderImage)
	p.BuildProfile = p.LoadVar("PIPE_ART_BUILD_PROFILE", p.BuildProfile)
	p.ArtefactName = p.LoadVar("PIPE_ART_NAME", p.ArtefactName)
	p.ArtefactRegistryUser = p.LoadVar("PIPE_ART_REG_USER", p.ArtefactRegistryUser)
	p.ArtefactRegistryPwd = p.LoadVar("PIPE_ART_REG_PWD", p.ArtefactRegistryPwd)
	if sonar {
		p.SonarURI = p.LoadVar("PIPE_ART_SONAR_URI", p.SonarURI)
		p.SonarToken = p.LoadVar("PIPE_ART_SONAR_TOKEN", p.SonarToken)
		p.SonarImage = p.LoadVar("PIPE_ART_SONAR_IMAGE", p.SonarImage)
		p.SonarSources = p.LoadVar("PIPE_ART_SONAR_SOURCES", p.SonarSources)
		p.SonarBinaries = p.LoadVar("PIPE_ART_SONAR_BINARIES", p.SonarBinaries)
	}
}

func (p *ArtefactPipelineConfig) LoadVar(name string, value string) string {
	// try and retrieve value from environment variable
	envVarValue := os.Getenv(name)
	// if there is a value use it
	if len(envVarValue) > 0 {
		return envVarValue
	}
	// if not then return the original value
	return value
}

// collect missing variables on the command line
func (p *ArtefactPipelineConfig) survey(sonar bool) {
	// the sonar scanner image to use
	p.SonarImage = "quay.io/gatblau/art-sonar"

	// if the application name is not defined prompt for it
	if len(p.AppName) == 0 {
		prompt := &survey.Input{
			Message: "application name:",
		}
		core.HandleCtrlC(survey.AskOne(prompt, &p.AppName, survey.WithValidator(survey.Required)))
	} else {
		fmt.Printf("application name: %s\n", p.AppName)
	}
	// if the GIT URI is not defined, prompt for it
	if len(p.GitURI) == 0 {
		prompt := &survey.Input{
			Message: "git repo url:",
		}
		core.HandleCtrlC(survey.AskOne(prompt, &p.GitURI, survey.WithValidator(validURL)))
	} else {
		fmt.Printf("git repo url: %s", p.GitURI)
	}
	// if the artefact name is not defined prompt for it
	if len(p.ArtefactName) == 0 {
		prompt := &survey.Input{
			Message: "artefact name:",
		}
		core.HandleCtrlC(survey.AskOne(prompt, &p.ArtefactName, survey.WithValidator(survey.Required)))
	} else {
		fmt.Printf("artefact name: %s", p.ArtefactName)
	}
	// if the artefact registry user is not defined prompt for it
	if len(p.ArtefactRegistryUser) == 0 {
		prompt := &survey.Input{
			Message: "artefact registry username:",
		}
		core.HandleCtrlC(survey.AskOne(prompt, &p.ArtefactRegistryUser, survey.WithValidator(survey.Required)))
	} else {
		fmt.Printf("artefact registry username: %s", p.ArtefactRegistryUser)
	}
	// if the artefact registry pwd is not defined prompt for it
	if len(p.ArtefactRegistryPwd) == 0 {
		prompt := &survey.Password{
			Message: "artefact registry password:",
		}
		core.HandleCtrlC(survey.AskOne(prompt, &p.ArtefactRegistryPwd, survey.WithValidator(survey.Required)))
	}
	if sonar {
		// if the Sonar URI is not defined, prompt for it
		if len(p.SonarURI) == 0 {
			prompt := &survey.Input{
				Message: "Sonar URI:",
			}
			core.HandleCtrlC(survey.AskOne(prompt, &p.SonarURI, survey.WithValidator(validURL)))
		} else {
			fmt.Printf("Sonar URI: %s", p.SonarURI)
		}
		// if the Sonar token is not defined prompt for it
		if len(p.SonarToken) == 0 {
			prompt := &survey.Password{
				Message: "Sonar Token:",
			}
			core.HandleCtrlC(survey.AskOne(prompt, &p.SonarToken, survey.WithValidator(survey.Required)))
		}
		// if the Sonar Project Key is not defined prompt for it
		if len(p.SonarProjectKey) == 0 {
			prompt := &survey.Input{
				Message: "Sonar project key:",
			}
			core.HandleCtrlC(survey.AskOne(prompt, &p.SonarProjectKey, survey.WithValidator(survey.Required)))
		} else {
			fmt.Printf("Sonar project key: %s\n", p.SonarProjectKey)
		}
		// if the Sonar Project Key is not defined prompt for it
		if len(p.SonarSources) == 0 {
			prompt := &survey.Input{
				Message: "Sonar sources:",
			}
			core.HandleCtrlC(survey.AskOne(prompt, &p.SonarSources, survey.WithValidator(survey.Required)))
		} else {
			fmt.Printf("Sonar sources: %s\n", p.SonarSources)
		}
		// if the Sonar Project Key is not defined prompt for it
		if len(p.SonarBinaries) == 0 {
			prompt := &survey.Input{
				Message: "Sonar binaries:",
			}
			core.HandleCtrlC(survey.AskOne(prompt, &p.SonarBinaries, survey.WithValidator(survey.Required)))
		} else {
			fmt.Printf("Sonar binaries: %s\n", p.SonarBinaries)
		}
	}
}

// merges the template and its values into the passed in writer
func (p *ArtefactPipelineConfig) Merge(w io.Writer) error {
	t, err := template.New("pipeline").Parse(ciPipeline)
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
		core.RaiseErr("a pipeline for an application of type '%s' is not supported by Artie", appType)
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
	b, err := ioutil.ReadFile(path.Join(filePath, "build.yaml"))
	core.CheckErr(err, "cannot read build file")
	buildFile := new(build.BuildFile)
	err = yaml.Unmarshal(b, buildFile)
	core.CheckErr(err, "cannot unmarshall build file")
	return buildFile
}

// the ciPipeline template containing parameterised resource definitions
const ciPipeline = `
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
      image: {{.BuilderImage}} 
      env:
        - name: ARTEFACT_NAME
          value: {{.ArtefactName}}
        - name: BUILD_PROFILE
          value: {{.BuildProfile}}
        - name: ARTEFACT_UNAME
          value: {{.ArtefactRegistryUser}}
        - name: ARTEFACT_PWD
          value: {{.ArtefactRegistryPwd}}
      workingDir: /workspace/source
      volumeMounts:
        - name: config-volume
          mountPath: /keys
  volumes:
    - name: config-volume
      configMap:
        name: signing-key-config-map
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
