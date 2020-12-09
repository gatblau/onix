package tkn

type PipelineRun struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

type ResourceRef struct {
	Name string `yaml:"name"`
}
