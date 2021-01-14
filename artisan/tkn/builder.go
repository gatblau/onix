package tkn

import (
	"bytes"
	"fmt"
	"github.com/gatblau/onix/artisan/crypto"
	"github.com/gatblau/onix/artisan/flow"
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
	secrets := b.newCredentialsSecret()
	buf.Write(ToYaml(secrets, "Secret"))
	buf.WriteString("\n---\n")
	keysSecret := b.newKeySecrets()
	buf.Write(ToYaml(keysSecret, "Keys Secret"))
	buf.WriteString("\n---\n")
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
					Name: b.secretName(),
					Key:  secret.Name,
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
					SecretName: b.keysSecretName(),
				},
			},
		}
	}
	return nil
}

func (b *Builder) newCredentialsSecret() *Secret {
	for _, step := range b.flow.Steps {
		if step.Input != nil && step.Input.Key != nil {
			secrets := step.Input.Secret
			s := new(Secret)
			s.APIVersion = ApiVersion
			s.Kind = "Secret"
			s.Type = "Opaque"
			s.Metadata = &Metadata{
				Name: b.secretName(),
			}
			credentials := make(map[string]string)
			for _, secret := range secrets {
				name := secret.Name
				secret.Decrypt(b.pk)
				credentials[name] = secret.Value
			}
			s.StringData = &credentials
			return s
		}
	}
	return nil
}

func (b *Builder) newKeySecrets() *Secret {
	for _, step := range b.flow.Steps {
		if step.Input != nil && step.Input.Key != nil {
			keys := step.Input.Key
			s := new(Secret)
			s.APIVersion = ApiVersion
			s.Kind = "Secret"
			s.Type = "Opaque"
			s.Metadata = &Metadata{
				Name: b.keysSecretName(),
			}
			keysDict := make(map[string]string)
			var name string
			for _, key := range keys {
				prefix := crypto.KeyNamePrefix(key.PackageGroup, key.PackageName)
				if key.Private {
					name = crypto.PrivateKeyName(prefix, "pgp")
				} else {
					name = crypto.PublicKeyName(prefix, "pgp")
				}
				key.Decrypt(b.pk)
				keysDict[name] = key.Value
			}
			s.StringData = &keysDict
			return s
		}
	}
	return nil
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
	return fmt.Sprintf("%s-build-task", encode(b.flow.Name))
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
func (b *Builder) secretName() string {
	return fmt.Sprintf("%s-creds-secret", encode(b.flow.Name))
}

func (b *Builder) keysSecretName() string {
	return fmt.Sprintf("%s-keys-secret", encode(b.flow.Name))
}
