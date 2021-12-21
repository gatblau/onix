/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package tkn

import (
	"github.com/gatblau/onix/artisan/core"
	"gopkg.in/yaml.v2"
)

type Task struct {
	APIVersion string    `yaml:"apiVersion,omitempty"`
	Kind       string    `yaml:"kind,omitempty"`
	Metadata   *Metadata `yaml:"metadata,omitempty"`
	Spec       *TaskSpec `yaml:"spec,omitempty"`
}

func ToYaml(o interface{}, ref string) []byte {
	b, err := yaml.Marshal(o)
	core.CheckErr(err, "cannot marshal %s", ref)
	return b
}

type SecretKeyRef struct {
	Key  string `yaml:"key,omitempty"`
	Name string `yaml:"name,omitempty"`
}

type ValueFrom struct {
	SecretKeyRef *SecretKeyRef `yaml:"secretKeyRef,omitempty"`
}

type Env struct {
	Name      string     `yaml:"name,omitempty"`
	Value     string     `yaml:"value,omitempty"`
	ValueFrom *ValueFrom `yaml:"valueFrom,omitempty"`
}

type VolumeMounts struct {
	Name      string `yaml:"name,omitempty"`
	MountPath string `yaml:"mountPath,omitempty"`
}

type Steps struct {
	Name            string               `yaml:"name,omitempty"`
	Image           string               `yaml:"image,omitempty"`
	Env             []*Env               `yaml:"env,omitempty"`
	Command         []string             `yaml:"command,omitempty,flow"`
	WorkingDir      string               `yaml:"workingDir,omitempty"`
	VolumeMounts    []*VolumeMounts      `yaml:"volumeMounts,omitempty"`
	Args            string               `yaml:"args,omitempty"`
	SecurityContext *StepSecurityContext `yaml:"securityContext,omitempty"`
}

type ConfigMap struct {
	Name string `yaml:"name,omitempty"`
}

type Volumes struct {
	Name      string     `yaml:"name,omitempty"`
	Secret    *Secret    `yaml:"secret,omitempty"`
	ConfigMap *ConfigMap `yaml:"configMap,omitempty"`
}
