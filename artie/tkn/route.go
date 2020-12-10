package tkn

type Route struct {
	APIVersion string   `yaml:"apiVersion",omitempty`
	Kind       string   `yaml:"kind",omitempty`
	Metadata   Metadata `yaml:"metadata",omitempty`
	Spec       Spec     `yaml:"spec",omitempty`
}

type Annotations struct {
	Description string `yaml:"description",omitempty`
}

type Port struct {
	TargetPort string `yaml:"targetPort",omitempty`
}

type TLS struct {
	InsecureEdgeTerminationPolicy string `yaml:"insecureEdgeTerminationPolicy",omitempty`
	Termination                   string `yaml:"termination",omitempty`
}

type To struct {
	Kind string `yaml:"kind",omitempty`
	Name string `yaml:"name",omitempty`
}
