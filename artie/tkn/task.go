package tkn

import (
	"github.com/gatblau/onix/artie/core"
	"gopkg.in/yaml.v2"
)

type Task struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata",omitempty`
	Spec       Spec     `yaml:"spec",omitempty`
}

func (t *Task) ToYaml() string {
	b, err := yaml.Marshal(t)
	core.CheckErr(err, "cannot marshall pipeline Type")
	return string(b)
}

type SecretKeyRef struct {
	Key  string `yaml:"key"`
	Name string `yaml:"name"`
}

type ValueFrom struct {
	SecretKeyRef SecretKeyRef `yaml:"secretKeyRef",omitempty`
}

type Env struct {
	Name      string    `yaml:"name"`
	Value     string    `yaml:"value,omitempty"`
	ValueFrom ValueFrom `yaml:"valueFrom,omitempty"`
}

type VolumeMounts struct {
	Name      string `yaml:"name"`
	MountPath string `yaml:"mountPath"`
}

type Steps struct {
	Name         string         `yaml:"name",omitempty`
	Image        string         `yaml:"image",omitempty`
	Env          []Env          `yaml:"env",omitempty`
	WorkingDir   string         `yaml:"workingDir",omitempty`
	VolumeMounts []VolumeMounts `yaml:"volumeMounts",omitempty`
}

type ConfigMap struct {
	Name string `yaml:"name"`
}

type Volumes struct {
	Name      string    `yaml:"name"`
	ConfigMap ConfigMap `yaml:"configMap"`
}
