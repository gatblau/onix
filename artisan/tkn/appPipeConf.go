/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
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
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/data"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// a tekton-based Artisan CI pipeline
type AppPipeConf struct {
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
func NewAppPipelineConfig(buildFilePath, profileName string, sonar bool) *AppPipeConf {
	var profile *data.Profile
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
	p := new(AppPipeConf)
	// resolve the builder image using the appType
	p.BuilderImage = builderImage(profile.Type)
	// set the build profile
	p.BuildProfile = profile.Name
	// set the application name
	p.AppName = profile.Application
	p.AppIcon = profile.Icon
	p.ArtefactName = profile.Artefact
	// if sonar step is required and there is a Sonar configuration section in buildfile
	if sonar && profile.Sonar != nil {
		p.SonarURI = profile.Sonar.URI
		p.SonarProjectKey = profile.Sonar.ProjectKey
		p.SonarSources = profile.Sonar.Sources
		p.SonarBinaries = profile.Sonar.Binaries
	}
	// attempt to load the pipeline configuration from the environment
	// NOTE: environment vars can override builder image and/or build profile used (if defined)
	p.loadFromEnv(sonar)
	// survey any variables in the pipeline that has been left undefined
	p.survey(sonar)
	// finally survey any missing variables in the build profile that are not defined
	profile.Survey(buildFile)
	// need this to take the cli cursor to the new line
	fmt.Print("\n")
	// return the configured pipeline
	return p
}

// return the name of the application build task
func (p *AppPipeConf) buildTaskName() string {
	return fmt.Sprintf("%s-app-build-task", p.AppName)
}

// return the name of the code repository resource
func (p *AppPipeConf) codeRepoResourceName() string {
	return fmt.Sprintf("%s-code-repo", p.AppName)
}

// return the name of the code repository resource
func (p *AppPipeConf) pipelineName() string {
	return fmt.Sprintf("%s-app-builder", p.AppName)
}

// return the name of the code repository resource
func (p *AppPipeConf) pipelineRunName() string {
	return fmt.Sprintf("%s-app-pr", p.AppName)
}

// try and set ciPipeline variables from the environment
func (p *AppPipeConf) loadFromEnv(sonar bool) {
	p.AppName = p.LoadVar("ART_PIPE_APP_NAME", p.AppName)
	p.AppIcon = p.LoadVar("ART_PIPE_APP_ICON", p.AppIcon)
	p.GitURI = p.LoadVar("ART_PIPE_APP_GIT_URI", p.GitURI)
	p.BuilderImage = p.LoadVar("ART_PIPE_APP_BUILDER_IMG", p.BuilderImage)
	p.BuildProfile = p.LoadVar("ART_PIPE_APP_BUILD_PROFILE", p.BuildProfile)
	p.ArtefactName = p.LoadVar("ART_PIPE_APP_ART_NAME", p.ArtefactName)
	p.ArtefactRegistryUser = p.LoadVar("ART_PIPE_APP_ART_REG_USER", p.ArtefactRegistryUser)
	p.ArtefactRegistryPwd = p.LoadVar("ART_PIPE_APP_ART_REG_PWD", p.ArtefactRegistryPwd)
	if sonar {
		p.SonarURI = p.LoadVar("ART_PIPE_APP_SONAR_URI", p.SonarURI)
		p.SonarToken = p.LoadVar("ART_PIPE_APP_SONAR_TOKEN", p.SonarToken)
		p.SonarImage = p.LoadVar("ART_PIPE_APP_SONAR_IMAGE", p.SonarImage)
		p.SonarSources = p.LoadVar("ART_PIPE_APP_SONAR_SOURCES", p.SonarSources)
		p.SonarBinaries = p.LoadVar("ART_PIPE_APP_SONAR_BINARIES", p.SonarBinaries)
	}
}

func (p *AppPipeConf) LoadVar(name string, value string) string {
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
func (p *AppPipeConf) survey(sonar bool) {
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
func loadBuildFile(buildFilePath string) *data.BuildFile {
	filePath, err := core.AbsPath(buildFilePath)
	if err != nil {
		core.RaiseErr("invalid build file path: %s", err.Error())
		return nil
	}
	b, err := ioutil.ReadFile(path.Join(filePath, "build.yaml"))
	core.CheckErr(err, "cannot read build file")
	buildFile := new(data.BuildFile)
	err = yaml.Unmarshal(b, buildFile)
	core.CheckErr(err, "cannot unmarshall build file")
	return buildFile
}
