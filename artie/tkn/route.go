package tkn

type Route struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

type Annotations struct {
	Description string `yaml:"description"`
}

type Port struct {
	TargetPort string `yaml:"targetPort"`
}

type TLS struct {
	InsecureEdgeTerminationPolicy string `yaml:"insecureEdgeTerminationPolicy"`
	Termination                   string `yaml:"termination"`
}

type To struct {
	Kind string `yaml:"kind"`
	Name string `yaml:"name"`
}
