package tkn

type Metadata struct {
	Name        string       `yaml:"name,omitempty"`
	Labels      *Labels      `yaml:"labels,omitempty"`
	Annotations *Annotations `yaml:"annotations,omitempty"`
}

type Resources struct {
	Name        string       `yaml:"name,omitempty"`
	Type        string       `yaml:"type,omitempty"`
	Inputs      []*Inputs    `yaml:"inputs,omitempty"`
	ResourceRef *ResourceRef `yaml:"resourceRef,omitempty"`
}

type Params struct {
	Name        string `yaml:"name,omitempty"`
	Value       string `yaml:"value,omitempty"`
	Type        string `yaml:"type,omitempty"`
	Description string `yaml:"description,omitempty"`
	Default     string `yaml:"default,omitempty"`
}

type Inputs struct {
	Name      string       `yaml:"name,omitempty"`
	Resource  string       `yaml:"resource,omitempty"`
	Resources []*Resources `yaml:"resources,omitempty"`
}

type Spec struct {
	Resources          []*Resources  `yaml:"resources,omitempty"`
	ResourceTemplates  []interface{} `yaml:"resourcetemplates,omitempty"`
	Params             []*Params     `yaml:"params,omitempty"`
	Tasks              []*Tasks      `yaml:"tasks,omitempty"`
	Inputs             *Inputs       `yaml:"inputs,omitempty"`
	Steps              []*Steps      `yaml:"steps,omitempty"`
	Volumes            []*Volumes    `yaml:"volumes,omitempty"`
	Type               string        `yaml:"type,omitempty"`
	ServiceAccountName string        `yaml:"serviceAccountName,omitempty"`
	PipelineRef        *PipelineRef  `yaml:"pipelineRef,omitempty"`
	Triggers           []*Triggers   `yaml:"triggers,omitempty"`
	Port               *Port         `yaml:"port,omitempty"`
	TLS                *TLS          `yaml:"tls,omitempty"`
	To                 *To           `yaml:"to,omitempty"`
}

type Labels struct {
	Application           string `yaml:"application,omitempty"`
	AppOpenshiftIoRuntime string `yaml:"app.openshift.io/runtime,omitempty"`
}

type PipelineRef struct {
	Name string `yaml:"name,omitempty"`
}
