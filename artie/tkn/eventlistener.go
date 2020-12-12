package tkn

type EventListener struct {
	APIVersion string    `yaml:"apiVersion"`
	Kind       string    `yaml:"kind"`
	Metadata   *Metadata `yaml:"metadata"`
	Spec       *Spec     `yaml:"spec"`
}

type Bindings struct {
	Name string `yaml:"name"`
}

type Template struct {
	Name string `yaml:"name"`
}

type Triggers struct {
	Bindings []*Bindings `yaml:"bindings"`
	Template *Template   `yaml:"template"`
}
