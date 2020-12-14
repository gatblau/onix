/*
  Onix Config Manager - Artie
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package tkn

import (
	"bytes"
	"fmt"
)

const (
	ApiVersion              = "v1"
	ApiVersionTekton        = "tekton.dev/v1alpha1"
	ApiVersionTektonTrigger = "triggers.tekton.dev/v1alpha1"
	ServiceAccountName      = "pipeline"
)

// return the full configuration for an Artefact Tekton Pipeline
func MergeArtPipe(applicationName, builderImage, artefactName, buildProfile, signingKeyName, gitURI, applicationIcon string) string {
	buf := bytes.Buffer{}
	task := newArtPipeTask(applicationName, builderImage, "", artefactName, buildProfile, signingKeyName, "", "", "", "")
	buf.Write(ToYaml(task, "Task"))
	buf.WriteString("\n---\n")
	pipe := newArtPipe(applicationName)
	buf.Write(ToYaml(pipe, "Pipeline"))
	buf.WriteString("\n---\n")
	pipeResx := newArtPipeResource(applicationName, gitURI)
	buf.Write(ToYaml(pipeResx, "PipelineResource"))
	buf.WriteString("\n---\n")
	pipeRun := newArtPipeRun(applicationName)
	buf.Write(ToYaml(pipeRun, "PipelineRun"))
	buf.WriteString("\n---\n")
	el := newArtPipeEventListener(applicationName, applicationIcon)
	buf.Write(ToYaml(el, "EventListener"))
	buf.WriteString("\n---\n")
	route := newArtPipeRoute(applicationName)
	buf.Write(ToYaml(route, "Route"))
	buf.WriteString("\n---\n")
	tb := newArtPipeTriggerBinding(applicationName)
	buf.Write(ToYaml(tb, "TriggerBinding"))
	buf.WriteString("\n---\n")
	tt := newArtPipeTriggerTemplate(applicationName, gitURI)
	buf.Write(ToYaml(tt, "TriggerTemplate"))
	buf.WriteString("\n---\n")
	return buf.String()
}

func newArtPipeTask(applicationName, builderImage, sonarImage, artefactName, buildProfile, signingKeyName, sonarURI, sonarProjectKey, sonarSources, sonarBinaries string) *Task {
	t := new(Task)
	t.APIVersion = ApiVersionTekton
	t.Kind = "Task"
	t.Spec = &Spec{
		Inputs: &Inputs{
			Resources: []*Resources{
				{
					Name: "source",
					Type: "git",
				},
			},
		},
		Steps: []*Steps{
			{
				Name:       "build-app",
				Image:      builderImage,
				Command:    []string{"artie", "run", "build-app"},
				WorkingDir: "/workspace/source",
			},
			{
				Name:    "scan-app",
				Image:   sonarImage,
				Command: []string{"artie", "run", "build-app"},
				Env: []*Env{
					{
						Name:  "SONAR_PROJECT_KEY",
						Value: sonarProjectKey,
					},
					{
						Name:  "SONAR_URI",
						Value: sonarURI,
					},
					{
						Name:  "SONAR_SOURCES",
						Value: sonarSources,
					},
					{
						Name:  "SONAR_BINARIES",
						Value: sonarBinaries,
					},
					{
						Name: "SONAR_TOKEN",
						ValueFrom: &ValueFrom{
							SecretKeyRef: &SecretKeyRef{
								Name: fmt.Sprintf("%s-sonar-token", applicationName),
								Key:  "token",
							}},
					},
				},
				WorkingDir: "/workspace/source",
			},
			{
				Name:    "package-app",
				Image:   builderImage,
				Command: []string{"artie", "run", "build-app"},
				Env: []*Env{
					{
						Name:  "ARTEFACT_NAME",
						Value: artefactName,
					},
					{
						Name:  "BUILD_PROFILE",
						Value: buildProfile,
					},
					{
						Name: "ARTEFACT_REG_USER",
						ValueFrom: &ValueFrom{
							SecretKeyRef: &SecretKeyRef{
								Name: fmt.Sprintf("%s-art-registry-creds", applicationName),
								Key:  "user",
							}},
					},
					{
						Name: "ARTEFACT_REG_PWD",
						ValueFrom: &ValueFrom{
							SecretKeyRef: &SecretKeyRef{
								Name: fmt.Sprintf("%s-art-registry-creds", applicationName),
								Key:  "pwd",
							}},
					},
				},
				WorkingDir: "/workspace/source",
				VolumeMounts: []*VolumeMounts{
					{
						Name:      "keys-volume",
						MountPath: "/keys",
					},
				},
			},
		},
		Volumes: []*Volumes{
			{
				Name: "keys-volume",
				Secret: &Secret{
					SecretName: fmt.Sprintf("%s-key-cm", applicationName),
				},
			},
		},
	}
	return t
}

func newArtPipe(applicationName string) *Pipeline {
	p := new(Pipeline)
	p.Kind = "Pipeline"
	p.APIVersion = ApiVersionTekton
	p.Metadata = &Metadata{
		Name: fmt.Sprintf("%s-artefact-pipeline", applicationName),
	}
	p.Spec = &Spec{
		Resources: []*Resources{
			{
				Name: fmt.Sprintf("%s-code-repository", applicationName),
				Type: "git",
			},
		},
		Params: []*Params{
			{
				Name:        "deployment-name",
				Type:        "string",
				Description: "the unique name for this deployment",
			},
		},
		Tasks: []*Tasks{
			{
				Name: "build-artefacts",
				TaskRef: &TaskRef{
					Name: fmt.Sprintf("%s-build-artefacts", applicationName),
				},
				Resources: &Resources{
					Inputs: []*Inputs{
						{
							Name:     "source",
							Resource: fmt.Sprintf("%s-code-repository", applicationName),
						},
					},
				},
			},
		},
	}
	return p
}

func newArtPipeResource(applicationName, gitURI string) *PipelineResource {
	r := new(PipelineResource)
	r.APIVersion = ApiVersionTekton
	r.Kind = "PipelineResource"
	r.Metadata = &Metadata{
		Name: fmt.Sprintf("%s-code-repository", applicationName),
	}
	r.Spec = &Spec{
		Type: "git",
		Params: []*Params{
			{
				Name:  "url",
				Value: gitURI,
			},
		},
	}
	return r
}

func newArtPipeRun(applicationName string) *PipelineRun {
	r := new(PipelineRun)
	r.Kind = "PipelineRun"
	r.APIVersion = ApiVersionTekton
	r.Metadata = &Metadata{
		Name: fmt.Sprintf("build-deploy-%s-pipelinerun", applicationName),
	}
	r.Spec = &Spec{
		Resources: []*Resources{
			{
				Name: fmt.Sprintf("%s-code-repository", applicationName),
				ResourceRef: &ResourceRef{
					Name: fmt.Sprintf("%s-code-repository", applicationName),
				},
			},
		},
		Params: []*Params{
			{
				Name:  "deployment-name",
				Value: applicationName,
			},
		},
		ServiceAccountName: ServiceAccountName,
		PipelineRef: &PipelineRef{
			Name: fmt.Sprintf("%s-artefact-builder", applicationName),
		},
	}
	return r
}

func newArtPipeEventListener(applicationName, appIcon string) *EventListener {
	e := new(EventListener)
	e.APIVersion = ApiVersionTektonTrigger
	e.Kind = "EventListener"
	e.Metadata = &Metadata{
		Name: applicationName,
		Labels: &Labels{
			AppOpenshiftIoRuntime: appIcon,
		},
	}
	e.Spec = &Spec{
		ServiceAccountName: ServiceAccountName,
		Triggers: []*Triggers{
			{
				Bindings: []*Bindings{
					{
						Name: applicationName,
					},
				},
				Template: &Template{
					Name: applicationName,
				},
			},
		},
	}
	return e
}

func newArtPipeRoute(applicationName string) *Route {
	r := new(Route)
	r.APIVersion = ApiVersion
	r.Kind = "Route"
	r.Metadata = &Metadata{
		Name: fmt.Sprintf("el-%s", applicationName),
		Labels: &Labels{
			Application: fmt.Sprintf("%s-https", applicationName),
		},
		Annotations: &Annotations{
			Description: "Route for the Artefact Pipeline Event Listener.",
		},
	}
	r.Spec = &Spec{
		Port: &Port{
			TargetPort: "8080",
		},
		TLS: &TLS{
			InsecureEdgeTerminationPolicy: "Redirect",
			Termination:                   "edge",
		},
		To: &To{
			Kind: "Service",
			Name: fmt.Sprintf("el-%s", applicationName),
		},
	}
	return r
}

func newArtPipeTriggerBinding(applicationName string) *TriggerBinding {
	t := new(TriggerBinding)
	t.APIVersion = ApiVersionTektonTrigger
	t.Kind = "TriggerBinding"
	t.Metadata = &Metadata{
		Name: applicationName,
	}
	t.Spec = &Spec{
		Params: []*Params{
			{
				Name:  "git-repo-url",
				Value: "$(body.project.web_url)",
			},
			{
				Name:  "git-repo-name",
				Value: "$(body.repository.name)",
			},
			{
				Name:  "git-revision",
				Value: "$(body.commits[0].id)",
			},
		},
	}
	return t
}

func newArtPipeTriggerTemplate(applicationName, gitURI string) *PipelineRun {
	pipeResx := newArtPipeResourceTriggerTemplate(applicationName, gitURI)
	pipeRun := newArtPipeRunTriggerTemplate(applicationName)

	t := new(PipelineRun)
	t.APIVersion = ApiVersionTektonTrigger
	t.Kind = "TriggerTemplate"
	t.Metadata = &Metadata{
		Name: applicationName,
	}
	t.Spec = &Spec{
		Params: []*Params{
			{
				Name:  "git-repo-url",
				Value: "The git repository url",
			},
			{
				Name:  "git-repo-name",
				Value: "The git repository name",
			},
			{
				Name:  "git-revision",
				Value: "The git revision",
			},
		},
		ResourceTemplates: []interface{}{pipeResx, pipeRun},
	}
	return t
}

func newArtPipeResourceTriggerTemplate(applicationName, gitURI string) *PipelineResource {
	r := new(PipelineResource)
	r.APIVersion = ApiVersionTekton
	r.Kind = "PipelineResource"
	r.Metadata = &Metadata{
		Name: "$(params.git-repo-name)-git-repo-$(uid)",
	}
	r.Spec = &Spec{
		ServiceAccountName: ServiceAccountName,
		PipelineRef: &PipelineRef{
			Name: fmt.Sprintf("%s-artefact-builder", applicationName),
		},
		Resources: []*Resources{
			{
				Name: "",
			},
		},
		Params: []*Params{
			{
				Name:  "revision",
				Value: "$(params.git-revision)",
			},
			{
				Name:  "url",
				Value: gitURI,
			},
		},
	}
	return r
}

func newArtPipeRunTriggerTemplate(applicationName string) *PipelineRun {
	r := new(PipelineRun)
	r.Kind = "PipelineRun"
	r.APIVersion = ApiVersionTekton
	r.Metadata = &Metadata{
		Name: "build-deploy-$(params.git-repo-name)-$(uid)",
	}
	r.Spec = &Spec{
		Resources: []*Resources{
			{
				Name: fmt.Sprintf("%s-code-repository", applicationName),
				ResourceRef: &ResourceRef{
					Name: "$(params.git-repo-name)-git-repo-$(uid)",
				},
			},
		},
		Params: []*Params{
			{
				Name:  "deployment-name",
				Value: "$(params.git-repo-name)",
			},
		},
		ServiceAccountName: ServiceAccountName,
		PipelineRef: &PipelineRef{
			Name: fmt.Sprintf("%s-artefact-builder", applicationName),
		},
	}
	return r
}

func newArtSecret() *Secret {
	return nil
}
