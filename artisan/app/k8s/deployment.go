/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package k8s

// Deployment Kubernetes deployment resource definition
type Deployment struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

type MatchLabels struct {
	App string `yaml:"app"`
}

type SecretKeyRef struct {
	Key  string `yaml:"key"`
	Name string `yaml:"name"`
}

type ValueFrom struct {
	SecretKeyRef SecretKeyRef `yaml:"secretKeyRef"`
}

type Env struct {
	Name      string    `yaml:"name"`
	Value     string    `yaml:"value,omitempty"`
	ValueFrom ValueFrom `yaml:"valueFrom,omitempty"`
}

type Exec struct {
	Command []string `yaml:"command,omitempty"`
}

type ReadinessProbe struct {
	Exec                Exec `yaml:"exec,omitempty"`
	InitialDelaySeconds int  `yaml:"initialDelaySeconds,omitempty"`
	PeriodSeconds       int  `yaml:"periodSeconds,omitempty"`
	TimeoutSeconds      int  `yaml:"timeoutSeconds,omitempty"`
	FailureThreshold    int  `yaml:"failureThreshold,omitempty"`
}

type Ports struct {
	ContainerPort int    `yaml:"containerPort"`
	Name          string `yaml:"name,omitempty"`
	Protocol      string `yaml:"protocol,omitempty"`
	Port          int    `yaml:"port,omitempty"`
	TargetPort    int    `yaml:"targetPort,omitempty"`
}

type Requests struct {
	Memory string `yaml:"memory"`
	CPU    string `yaml:"cpu"`
}

type Limits struct {
	Memory string `yaml:"memory"`
	CPU    string `yaml:"cpu"`
}

type Resources struct {
	Requests Requests `yaml:"requests"`
	Limits   Limits   `yaml:"limits"`
}

type Containers struct {
	Env             []Env          `yaml:"env"`
	Image           string         `yaml:"image"`
	ReadinessProbe  ReadinessProbe `yaml:"readinessProbe,omitempty"`
	ImagePullPolicy string         `yaml:"imagePullPolicy"`
	Name            string         `yaml:"name"`
	Ports           []Ports        `yaml:"ports"`
	Resources       Resources      `yaml:"resources,omitempty"`
}

type TemplateSpec struct {
	Containers                    []Containers `yaml:"containers"`
	TerminationGracePeriodSeconds int          `yaml:"terminationGracePeriodSeconds"`
}

type Template struct {
	Metadata Metadata     `yaml:"metadata"`
	Spec     TemplateSpec `yaml:"spec"`
}

type Spec struct {
	Replicas int      `yaml:"replicas,omitempty"`
	Selector Selector `yaml:"selector,omitempty"`
	Template Template `yaml:"template,omitempty"`
	Ports    []Ports  `yaml:"ports,omitempty"`
	Type     string   `yaml:"type,omitempty"`
	TLS      []TLS    `yaml:"tls,omitempty"`
	Rules    []Rules  `yaml:"rules,omitempty"`
}
