package tkn

type PipelineResource struct {
	APIVersion string   `yaml:"apiVersion",omitempty`
	Kind       string   `yaml:"kind",omitempty`
	Metadata   Metadata `yaml:"metadata",omitempty`
	Spec       Spec     `yaml:"spec",omitempty`
}
