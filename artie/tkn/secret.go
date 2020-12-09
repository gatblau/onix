package tkn

type Secret struct {
	APIVersion string     `yaml:"apiVersion"`
	Kind       string     `yaml:"kind"`
	Metadata   Metadata   `yaml:"metadata"`
	Type       string     `yaml:"type"`
	StringData StringData `yaml:"stringData"`
}

type StringData struct {
	Pwd string `yaml:"pwd"`
}
