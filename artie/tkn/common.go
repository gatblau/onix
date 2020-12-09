package tkn

type Metadata struct {
	Name        string      `yaml:"name"`
	Labels      Labels      `yaml:"labels"`
	Annotations Annotations `yaml:"annotations"`
}

type Resources struct {
	Name        string      `yaml:"name"`
	Type        string      `yaml:"type"`
	Inputs      []Inputs    `yaml:"inputs"`
	ResourceRef ResourceRef `yaml:"resourceRef"`
}

type Params struct {
	Name        string `yaml:"name"`
	Value       string `yaml:"value"`
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	Default     string `yaml:"default,omitempty"`
}

type Inputs struct {
	Name      string      `yaml:"name"`
	Resource  string      `yaml:"resource"`
	Resources []Resources `yaml:"resources"`
}

type Spec struct {
	Resources          []Resources   `yaml:"resources",omitempty`
	ResourceTemplates  []interface{} `yaml:"resourcetemplates",omitempty`
	Params             []Params      `yaml:"params",omitempty`
	Tasks              []Tasks       `yaml:"tasks",omitempty`
	Inputs             Inputs        `yaml:"inputs",omitempty`
	Steps              []Steps       `yaml:"steps",omitempty`
	Volumes            []Volumes     `yaml:"volumes",omitempty`
	Type               string        `yaml:"type",omitempty`
	ServiceAccountName string        `yaml:"serviceAccountName",omitempty`
	PipelineRef        PipelineRef   `yaml:"pipelineRef",omitempty`
	Triggers           []Triggers    `yaml:"triggers",omitempty`
	Port               Port          `yaml:"port",omitempty`
	TLS                TLS           `yaml:"tls",omitempty`
	To                 To            `yaml:"to",omitempty`
}

type Labels struct {
	Application           string `yaml:"application"`
	AppOpenshiftIoRuntime string `yaml:"app.openshift.io/runtime"`
}

type PipelineRef struct {
	Name string `yaml:"name"`
}
