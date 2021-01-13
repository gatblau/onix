package tkn

import (
	"bytes"
	"fmt"
	"github.com/gatblau/onix/artisan/crypto"
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/flow"
	"strings"
)

const (
	ApiVersion              = "v1"
	ApiVersionTekton        = "tekton.dev/v1alpha1"
	ApiVersionTektonTrigger = "triggers.tekton.dev/v1alpha1"
	ServiceAccountName      = "pipeline"
)

// tekton builder
type Builder struct {
	flow *flow.Flow
	pk   *crypto.PGP
}

func NewBuilder(flow *flow.Flow, privateKey *crypto.PGP) *Builder {
	return &Builder{flow: flow, pk: privateKey}
}

func (b *Builder) Create() bytes.Buffer {
	buf := bytes.Buffer{}
	task := b.newTask()
	buf.Write(ToYaml(task, "Task"))
	buf.WriteString("\n---\n")
	secrets := b.newSecrets()
	for _, secret := range secrets {
		buf.Write(ToYaml(secret, "Secret"))
		buf.WriteString("\n---\n")
	}
	pipeline := b.newPipeline()
	buf.Write(ToYaml(pipeline, "Pipeline"))
	buf.WriteString("\n---\n")
	return buf
}

// task
func (b *Builder) newTask() *Task {
	t := new(Task)
	t.APIVersion = ApiVersionTekton
	t.Kind = "Task"
	t.Metadata = &Metadata{
		Name: b.buildTaskName(),
	}
	t.Spec = &Spec{
		Inputs:  b.newInputs(),
		Steps:   b.newSteps(),
		Volumes: b.newVolumes(),
	}
	return t
}

func (b *Builder) newSteps() []*Steps {
	var steps = make([]*Steps, 0)
	for _, step := range b.flow.Steps {
		s := &Steps{
			Name:       step.Name,
			Image:      step.Runtime,
			WorkingDir: "/workspace/source",
		}
		// if the step requires keys
		if step.Input != nil {
			if len(step.Input.Key) > 0 {
				// add a volume mount for the keys
				s.VolumeMounts = []*VolumeMounts{
					{
						Name:      "keys-volume",
						MountPath: "/keys",
					},
				}
			}
			// if the step has vars or secrets
			if len(step.Input.Var)+len(step.Input.Secret) > 0 {
				// add to env
				s.Env = b.getEnv(step)
			}
		}
		steps = append(steps, s)
	}
	return steps
}

func (b *Builder) getEnv(step *flow.Step) []*Env {
	env := make([]*Env, 0)
	// add variables
	for _, variable := range step.Input.Var {
		env = append(env, &Env{
			Name:  variable.Name,
			Value: variable.Value,
		})
	}
	// add secrets
	for _, secret := range step.Input.Secret {
		env = append(env, &Env{
			Name: secret.Name,
			ValueFrom: &ValueFrom{
				SecretKeyRef: &SecretKeyRef{
					Name: b.secretName(secret),
					Key:  strings.ToLower(secret.Name),
				}},
		})
	}
	return env
}

func (b *Builder) newInputs() *Inputs {
	if b.flow.RequiresSource() {
		return &Inputs{
			Resources: []*Resources{
				{
					Name: "source",
					Type: "git",
				},
			},
		}
	}
	return nil
}

func (b *Builder) newVolumes() []*Volumes {
	if b.flow.RequiresKey() {
		return []*Volumes{
			{
				Name: "keys-volume",
				Secret: &Secret{
					SecretName: fmt.Sprintf("%s-key-cm", encode(b.flow.Name)),
				},
			},
		}
	}
	return nil
}

// secret
func (b *Builder) newSecret(secret *data.Secret) *Secret {
	s := new(Secret)
	s.APIVersion = ApiVersion
	s.Kind = "Secret"
	s.Type = "Opaque"
	s.Metadata = &Metadata{
		Name: b.secretName(secret),
	}
	secret.Decrypt(b.pk)
	s.StringData = &map[string]string{
		strings.ToLower(secret.Name): secret.Value,
	}
	return s
}

func (b *Builder) newSecrets() []*Secret {
	var secs []*Secret
	for _, step := range b.flow.Steps {
		if step.Input != nil {
			for _, secret := range step.Input.Secret {
				secs = append(secs, b.newSecret(secret))
			}
		}
	}
	return secs
}

// pipeline
func (b *Builder) newPipeline() *Pipeline {
	p := new(Pipeline)
	p.Kind = "Pipeline"
	p.APIVersion = ApiVersionTekton
	p.Metadata = &Metadata{
		Name: b.pipelineName(),
	}
	p.Spec = &Spec{
		Resources: []*Resources{
			{
				Name: b.codeRepoResourceName(),
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
				Name: b.buildTaskName(),
				TaskRef: &TaskRef{
					Name: b.buildTaskName(),
				},
				Resources: &Resources{
					Inputs: []*Inputs{
						{
							Name:     "source",
							Resource: b.codeRepoResourceName(),
						},
					},
				},
			},
		},
	}
	return p
}

// return the name of the application build task
func (b *Builder) buildTaskName() string {
	return fmt.Sprintf("%s-app-build-task", encode(b.flow.Name))
}

// return the name of the code repository resource
func (b *Builder) codeRepoResourceName() string {
	return fmt.Sprintf("%s-code-repo", encode(b.flow.Name))
}

// return the name of the code repository resource
func (b *Builder) pipelineName() string {
	return fmt.Sprintf("%s-app-builder", encode(b.flow.Name))
}

// return the name of the code repository resource
func (b *Builder) pipelineRunName() string {
	return fmt.Sprintf("%s-app-pr", encode(b.flow.Name))
}

// return the name of the code repository resource
func (b *Builder) secretName(secret *data.Secret) string {
	return fmt.Sprintf("%s-%s-secret", encode(secret.Name), encode(b.flow.Name))
}
