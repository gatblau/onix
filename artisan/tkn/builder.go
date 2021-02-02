package tkn

import (
	"bytes"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
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
}

func NewBuilder(flow *flow.Flow) *Builder {
	return &Builder{flow: flow}
}

func (b *Builder) Create() bytes.Buffer {
	buf := bytes.Buffer{}
	task := b.newTask()
	buf.Write(ToYaml(task, "Task"))
	buf.WriteString("\n---\n")
	secrets := b.newCredentialsSecret()
	if secrets != nil {
		buf.Write(ToYaml(secrets, "Secret"))
		buf.WriteString("\n---\n")
	}
	keysSecret := b.newKeySecrets()
	if keysSecret != nil {
		buf.Write(ToYaml(keysSecret, "Keys Secret"))
		buf.WriteString("\n---\n")
	}
	pipeline := b.newPipeline()
	buf.Write(ToYaml(pipeline, "Pipeline"))
	buf.WriteString("\n---\n")
	// if source code repository is required by the pipeline
	if b.flow.RequiresSource() {
		// add the following resources
		pipelineResource := b.newPipelineResource()
		buf.Write(ToYaml(pipelineResource, "PipelineResource"))
		buf.WriteString("\n---\n")
		eventListener := b.newEventListener()
		buf.Write(ToYaml(eventListener, "EventListener"))
		buf.WriteString("\n---\n")
		route := b.newRoute()
		buf.Write(ToYaml(route, "Route"))
		buf.WriteString("\n---\n")
		triggerBinding := b.newTriggerBinding()
		buf.Write(ToYaml(triggerBinding, "TriggerBinding"))
		buf.WriteString("\n---\n")
		triggerTemplate := b.newTriggerTemplate()
		buf.Write(ToYaml(triggerTemplate, "TriggerTemplate"))
		buf.WriteString("\n---\n")
	}
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
		// add the environment information required by the Artisan runtime to work
		// see here: https://github.com/gatblau/artisan/tree/master/runtime
		s.Env = b.addRuntimeInterfaceVars(step, s.Env)
		steps = append(steps, s)
	}
	return steps
}

func (b *Builder) getEnv(step *flow.Step) []*Env {
	env := make([]*Env, 0)
	// if there is an input defined
	if step.Input != nil {
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
	}
	return env
}

func (b *Builder) addRuntimeInterfaceVars(step *flow.Step, env []*Env) []*Env {
	if len(step.Function) > 0 {
		env = append(env, &Env{
			Name:  "FX_NAME",
			Value: step.Function,
		})
	}
	if len(step.Package) > 0 {
		env = append(env, &Env{
			Name:  "PACKAGE_NAME",
			Value: step.Package,
		})
		name, _ := core.ParseName(step.Package)
		env = append(env, &Env{
			Name: "ART_REG_USER",
			ValueFrom: &ValueFrom{
				SecretKeyRef: &SecretKeyRef{
					Name: b.secretName(),
					Key:  fmt.Sprintf("ART_REG_USER_%s", name.Domain),
				},
			},
		})
		env = append(env, &Env{
			Name: "ART_REG_PWD",
			ValueFrom: &ValueFrom{
				SecretKeyRef: &SecretKeyRef{
					Name: b.secretName(),
					Key:  fmt.Sprintf("ART_REG_PWD_%s", name.Domain),
				},
			},
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
	if b.flow.RequiresSecrets() {
		s := new(Secret)
		s.APIVersion = ApiVersion
		s.Kind = "Secret"
		s.Type = "Opaque"
		s.Metadata = &Metadata{
			Name: b.secretName(),
		}
		credentials := make(map[string]string)
		for _, step := range b.flow.Steps {
			if step.Input != nil && step.Input.Secret != nil {
				for _, secret := range step.Input.Secret {
					name := secret.Name
					credentials[name] = secret.Value
				}
			}
		}
		// add flow level secrets
		for _, cred := range b.flow.Credential {
			credentials[fmt.Sprintf("ART_REG_USER_%s", cred.Domain)] = cred.User
			credentials[fmt.Sprintf("ART_REG_PWD_%s", cred.Domain)] = cred.Password
		}
		s.StringData = &credentials
		return s
	}
	return nil
}

func (b *Builder) newKeySecrets() *Secret {
	if b.flow.RequiresKey() {
		s := new(Secret)
		s.APIVersion = ApiVersion
		s.Kind = "Secret"
		s.Type = "Opaque"
		s.Metadata = &Metadata{
			Name: b.keysSecretName(),
		}
		keysDict := make(map[string]string)
		var name string
		for _, step := range b.flow.Steps {
			if step.Input != nil {
				keys := step.Input.Key
				for _, key := range keys {
					prefix := crypto.KeyNamePrefix(key.PackageGroup, key.PackageName)
					if key.Private {
						name = crypto.PrivateKeyName(prefix, "pgp")
					} else {
						name = crypto.PublicKeyName(prefix, "pgp")
					}
					keysDict[name] = key.Value
				}
			}
		}
		s.StringData = &keysDict
		return s
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
	var (
		inputs    []*Inputs
		resources []*Resources
	)
	if b.flow.RequiresSource() {
		inputs = []*Inputs{
			{
				Name:     "source",
				Resource: b.codeRepoResourceName(),
			},
		}
		resources = []*Resources{
			{
				Name: b.codeRepoResourceName(),
				Type: "git",
			},
		}
	}
	p.Spec = &Spec{
		Resources: resources,
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
					Inputs: inputs,
				},
			},
		},
	}
	return p
}

// pipeline resource
func (b *Builder) newPipelineResource() *PipelineResource {
	r := new(PipelineResource)
	r.APIVersion = ApiVersionTekton
	r.Kind = "PipelineResource"
	r.Metadata = &Metadata{
		Name: b.codeRepoResourceName(),
	}
	r.Spec = &Spec{
		Type: "git",
		Params: []*Params{
			{
				Name:  "url",
				Value: b.flow.GitURI,
			},
		},
	}
	return r
}

// event listener
func (b *Builder) newEventListener() *EventListener {
	e := new(EventListener)
	e.APIVersion = ApiVersionTektonTrigger
	e.Kind = "EventListener"
	e.Metadata = &Metadata{
		Name: encode(b.flow.Name),
		Labels: &Labels{
			AppOpenshiftIoRuntime: b.flow.AppIcon,
		},
	}
	e.Spec = &Spec{
		ServiceAccountName: ServiceAccountName,
		Triggers: []*Triggers{
			{
				Bindings: []*Bindings{
					{
						Name: encode(b.flow.Name),
					},
				},
				Template: &Template{
					Name: encode(b.flow.Name),
				},
			},
		},
	}
	return e
}

// route
func (b *Builder) newRoute() *Route {
	r := new(Route)
	r.APIVersion = ApiVersion
	r.Kind = "Route"
	r.Metadata = &Metadata{
		Name: fmt.Sprintf("el-%s", encode(b.flow.Name)),
		Labels: &Labels{
			Application: fmt.Sprintf("%s-https", encode(b.flow.Name)),
		},
		Annotations: &Annotations{
			Description: "Route for the Pipeline Event Listener.",
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
			Name: fmt.Sprintf("el-%s", encode(b.flow.Name)),
		},
	}
	return r
}

// trigger binding
func (b *Builder) newTriggerBinding() *TriggerBinding {
	t := new(TriggerBinding)
	t.APIVersion = ApiVersionTektonTrigger
	t.Kind = "TriggerBinding"
	t.Metadata = &Metadata{
		Name: encode(b.flow.Name),
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

// trigger template
func (b *Builder) newTriggerTemplate() *PipelineRun {
	pipeResx := b.newPipelineResourceTriggerTemplate()
	pipeRun := b.newPipelineRunTriggerTemplate()

	t := new(PipelineRun)
	t.APIVersion = ApiVersionTektonTrigger
	t.Kind = "TriggerTemplate"
	t.Metadata = &Metadata{
		Name: encode(b.flow.Name),
	}
	t.Spec = &Spec{
		Params: []*Params{
			{
				Name:        "git-repo-url",
				Description: "The git repository url",
			},
			{
				Name:        "git-repo-name",
				Description: "The git repository name",
			},
			{
				Name:        "git-revision",
				Description: "The git revision",
				Default:     "master",
			},
		},
		ResourceTemplates: []interface{}{pipeResx, pipeRun},
	}
	return t
}

func (b *Builder) newPipelineResourceTriggerTemplate() *PipelineResource {
	r := new(PipelineResource)
	r.APIVersion = ApiVersionTekton
	r.Kind = "PipelineResource"
	r.Metadata = &Metadata{
		Name: "$(params.git-repo-name)-git-repo-$(uid)",
	}
	r.Spec = &Spec{
		Type: "git",
		Params: []*Params{
			{
				Name:  "revision",
				Value: "$(params.git-revision)",
			},
			{
				Name:  "url",
				Value: b.flow.GitURI,
			},
		},
	}
	return r
}

func (b *Builder) newPipelineRunTriggerTemplate() *PipelineRun {
	r := new(PipelineRun)
	r.Kind = "PipelineRun"
	r.APIVersion = ApiVersionTekton
	r.Metadata = &Metadata{
		Name: "$(params.git-repo-name)-app-pr-$(uid)",
	}
	r.Spec = &Spec{
		ServiceAccountName: ServiceAccountName,
		PipelineRef: &PipelineRef{
			Name: b.pipelineName(),
		},
		Resources: []*Resources{
			{
				Name: b.codeRepoResourceName(),
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
	}
	return r
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
	return fmt.Sprintf("%s-builder", encode(b.flow.Name))
}

// return the name of the code repository resource
func (b *Builder) pipelineRunName() string {
	return fmt.Sprintf("%s-pr", encode(b.flow.Name))
}

// return the name of the code repository resource
func (b *Builder) secretName() string {
	return fmt.Sprintf("%s-creds-secret", encode(b.flow.Name))
}

func (b *Builder) keysSecretName() string {
	return fmt.Sprintf("%s-keys-secret", encode(b.flow.Name))
}
