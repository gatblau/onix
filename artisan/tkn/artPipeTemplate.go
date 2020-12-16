/*
  Onix Config Manager - Artisan
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
func MergeArtPipe(c *ArtefactPipelineConfig, sonar bool) bytes.Buffer {
	buf := bytes.Buffer{}
	task := newArtPipeTask(c, sonar)
	buf.Write(ToYaml(task, "Task"))
	buf.WriteString("\n---\n")
	regSecret := newArtRegSecret(c)
	buf.Write(ToYaml(regSecret, "Secret"))
	buf.WriteString("\n---\n")
	sonarSecret := newSonarSecret(c)
	buf.Write(ToYaml(sonarSecret, "Secret"))
	buf.WriteString("\n---\n")
	pipe := newArtPipe(c)
	buf.Write(ToYaml(pipe, "Pipeline"))
	buf.WriteString("\n---\n")
	pipeResx := newArtPipeResource(c)
	buf.Write(ToYaml(pipeResx, "PipelineResource"))
	buf.WriteString("\n---\n")
	pipeRun := newArtPipeRun(c)
	buf.Write(ToYaml(pipeRun, "PipelineRun"))
	buf.WriteString("\n---\n")
	el := newArtPipeEventListener(c)
	buf.Write(ToYaml(el, "EventListener"))
	buf.WriteString("\n---\n")
	route := newArtPipeRoute(c)
	buf.Write(ToYaml(route, "Route"))
	buf.WriteString("\n---\n")
	tb := newArtPipeTriggerBinding(c)
	buf.Write(ToYaml(tb, "TriggerBinding"))
	buf.WriteString("\n---\n")
	tt := newArtPipeTriggerTemplate(c)
	buf.Write(ToYaml(tt, "TriggerTemplate"))
	buf.WriteString("\n---\n")
	return buf
}

func newArtPipeTask(c *ArtefactPipelineConfig, sonar bool) *Task {
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
		Steps: getSteps(c, sonar),
		Volumes: []*Volumes{
			{
				Name: "keys-volume",
				Secret: &Secret{
					SecretName: fmt.Sprintf("%s-key-cm", c.AppName),
				},
			},
		},
	}
	return t
}

func getSteps(c *ArtefactPipelineConfig, sonar bool) []*Steps {
	stepCount := 2
	if sonar {
		stepCount = 3
	}
	var (
		ix    = 0
		steps = make([]*Steps, stepCount)
	)
	steps[ix] = &Steps{
		Name:       "build-app",
		Image:      c.BuilderImage,
		Command:    []string{"artie", "run", "build-app"},
		WorkingDir: "/workspace/source",
		VolumeMounts: []*VolumeMounts{
			{
				Name:      "keys-volume",
				MountPath: "/keys",
			},
		},
	}
	if sonar {
		ix++
		steps[ix] = &Steps{
			Name:  "scan-app",
			Image: c.SonarImage,
			Env: []*Env{
				{
					Name:  "SONAR_PROJECT_KEY",
					Value: c.SonarProjectKey,
				},
				{
					Name:  "SONAR_URI",
					Value: c.SonarURI,
				},
				{
					Name:  "SONAR_SOURCES",
					Value: c.SonarSources,
				},
				{
					Name:  "SONAR_BINARIES",
					Value: c.SonarBinaries,
				},
				{
					Name: "SONAR_TOKEN",
					ValueFrom: &ValueFrom{
						SecretKeyRef: &SecretKeyRef{
							Name: fmt.Sprintf("%s-sonar-token", c.AppName),
							Key:  "token",
						}},
				},
			},
			WorkingDir: "/workspace/source",
		}
	}
	ix++
	steps[ix] = &Steps{
		Name:  "package-app",
		Image: c.BuilderImage,
		Env: []*Env{
			{
				Name:  "ARTEFACT_NAME",
				Value: c.ArtefactName,
			},
			{
				Name:  "BUILD_PROFILE",
				Value: c.BuildProfile,
			},
			{
				Name: "ARTEFACT_REG_USER",
				ValueFrom: &ValueFrom{
					SecretKeyRef: &SecretKeyRef{
						Name: fmt.Sprintf("%s-art-registry-creds", c.AppName),
						Key:  "user",
					}},
			},
			{
				Name: "ARTEFACT_REG_PWD",
				ValueFrom: &ValueFrom{
					SecretKeyRef: &SecretKeyRef{
						Name: fmt.Sprintf("%s-art-registry-creds", c.AppName),
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
	}
	return steps
}

func newArtPipe(c *ArtefactPipelineConfig) *Pipeline {
	p := new(Pipeline)
	p.Kind = "Pipeline"
	p.APIVersion = ApiVersionTekton
	p.Metadata = &Metadata{
		Name: fmt.Sprintf("%s-artefact-pipeline", c.AppName),
	}
	p.Spec = &Spec{
		Resources: []*Resources{
			{
				Name: fmt.Sprintf("%s-code-repository", c.AppName),
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
					Name: fmt.Sprintf("%s-build-artefacts", c.AppName),
				},
				Resources: &Resources{
					Inputs: []*Inputs{
						{
							Name:     "source",
							Resource: fmt.Sprintf("%s-code-repository", c.AppName),
						},
					},
				},
			},
		},
	}
	return p
}

func newArtPipeResource(c *ArtefactPipelineConfig) *PipelineResource {
	r := new(PipelineResource)
	r.APIVersion = ApiVersionTekton
	r.Kind = "PipelineResource"
	r.Metadata = &Metadata{
		Name: fmt.Sprintf("%s-code-repository", c.AppName),
	}
	r.Spec = &Spec{
		Type: "git",
		Params: []*Params{
			{
				Name:  "url",
				Value: c.GitURI,
			},
		},
	}
	return r
}

func newArtPipeRun(c *ArtefactPipelineConfig) *PipelineRun {
	r := new(PipelineRun)
	r.Kind = "PipelineRun"
	r.APIVersion = ApiVersionTekton
	r.Metadata = &Metadata{
		Name: fmt.Sprintf("build-deploy-%s-pipelinerun", c.AppName),
	}
	r.Spec = &Spec{
		Resources: []*Resources{
			{
				Name: fmt.Sprintf("%s-code-repository", c.AppName),
				ResourceRef: &ResourceRef{
					Name: fmt.Sprintf("%s-code-repository", c.AppName),
				},
			},
		},
		Params: []*Params{
			{
				Name:  "deployment-name",
				Value: c.AppName,
			},
		},
		ServiceAccountName: ServiceAccountName,
		PipelineRef: &PipelineRef{
			Name: fmt.Sprintf("%s-artefact-builder", c.AppName),
		},
	}
	return r
}

func newArtPipeEventListener(c *ArtefactPipelineConfig) *EventListener {
	e := new(EventListener)
	e.APIVersion = ApiVersionTektonTrigger
	e.Kind = "EventListener"
	e.Metadata = &Metadata{
		Name: c.AppName,
		Labels: &Labels{
			AppOpenshiftIoRuntime: c.AppIcon,
		},
	}
	e.Spec = &Spec{
		ServiceAccountName: ServiceAccountName,
		Triggers: []*Triggers{
			{
				Bindings: []*Bindings{
					{
						Name: c.AppName,
					},
				},
				Template: &Template{
					Name: c.AppName,
				},
			},
		},
	}
	return e
}

func newArtPipeRoute(c *ArtefactPipelineConfig) *Route {
	r := new(Route)
	r.APIVersion = ApiVersion
	r.Kind = "Route"
	r.Metadata = &Metadata{
		Name: fmt.Sprintf("el-%s", c.AppName),
		Labels: &Labels{
			Application: fmt.Sprintf("%s-https", c.AppName),
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
			Name: fmt.Sprintf("el-%s", c.AppName),
		},
	}
	return r
}

func newArtPipeTriggerBinding(c *ArtefactPipelineConfig) *TriggerBinding {
	t := new(TriggerBinding)
	t.APIVersion = ApiVersionTektonTrigger
	t.Kind = "TriggerBinding"
	t.Metadata = &Metadata{
		Name: c.AppName,
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

func newArtPipeTriggerTemplate(c *ArtefactPipelineConfig) *PipelineRun {
	pipeResx := newArtPipeResourceTriggerTemplate(c)
	pipeRun := newArtPipeRunTriggerTemplate(c)

	t := new(PipelineRun)
	t.APIVersion = ApiVersionTektonTrigger
	t.Kind = "TriggerTemplate"
	t.Metadata = &Metadata{
		Name: c.AppName,
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

func newArtPipeResourceTriggerTemplate(c *ArtefactPipelineConfig) *PipelineResource {
	r := new(PipelineResource)
	r.APIVersion = ApiVersionTekton
	r.Kind = "PipelineResource"
	r.Metadata = &Metadata{
		Name: "$(params.git-repo-name)-git-repo-$(uid)",
	}
	r.Spec = &Spec{
		ServiceAccountName: ServiceAccountName,
		PipelineRef: &PipelineRef{
			Name: fmt.Sprintf("%s-artefact-builder", c.AppName),
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
				Value: c.GitURI,
			},
		},
	}
	return r
}

func newArtPipeRunTriggerTemplate(c *ArtefactPipelineConfig) *PipelineRun {
	r := new(PipelineRun)
	r.Kind = "PipelineRun"
	r.APIVersion = ApiVersionTekton
	r.Metadata = &Metadata{
		Name: "build-deploy-$(params.git-repo-name)-$(uid)",
	}
	r.Spec = &Spec{
		Resources: []*Resources{
			{
				Name: fmt.Sprintf("%s-code-repository", c.AppName),
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
			Name: fmt.Sprintf("%s-artefact-builder", c.AppName),
		},
	}
	return r
}

func newArtRegSecret(c *ArtefactPipelineConfig) *Secret {
	s := new(Secret)
	s.APIVersion = ApiVersion
	s.Kind = "Secret"
	s.Type = "Opaque"
	s.Metadata = &Metadata{
		Name: fmt.Sprintf("%s-art-registry-creds", c.AppName),
	}
	s.StringData = &StringData{
		Pwd:  c.ArtefactRegistryUser,
		User: c.ArtefactRegistryPwd,
	}
	return s
}

func newSonarSecret(c *ArtefactPipelineConfig) *Secret {
	s := new(Secret)
	s.APIVersion = ApiVersion
	s.Kind = "Secret"
	s.Type = "Opaque"
	s.Metadata = &Metadata{
		Name: fmt.Sprintf("%s-sonar-token", c.AppName),
	}
	s.StringData = &StringData{
		Token: c.SonarToken,
	}
	return s
}
