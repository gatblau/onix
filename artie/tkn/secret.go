package tkn

type Secret struct {
	APIVersion string     `yaml:"apiVersion",omitempty`
	Kind       string     `yaml:"kind",omitempty`
	Metadata   Metadata   `yaml:"metadata",omitempty`
	Type       string     `yaml:"type",omitempty`
	StringData StringData `yaml:"stringData",omitempty`
}

type StringData struct {
	Pwd string `yaml:"pwd",omitempty`
}
