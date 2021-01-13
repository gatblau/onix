package tkn2

type PipelineRun struct {
	APIVersion string    `yaml:"apiVersion",omitempty`
	Kind       string    `yaml:"kind",omitempty`
	Metadata   *Metadata `yaml:"metadata",omitempty`
	Spec       *Spec     `yaml:"spec",omitempty`
}

type ResourceRef struct {
	Name string `yaml:"name",omitempty`
}
