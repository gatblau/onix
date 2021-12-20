/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
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
	Command []string `yaml:"command"`
}

type ReadinessProbe struct {
	Exec                Exec `yaml:"exec"`
	InitialDelaySeconds int  `yaml:"initialDelaySeconds"`
	PeriodSeconds       int  `yaml:"periodSeconds"`
	TimeoutSeconds      int  `yaml:"timeoutSeconds"`
	FailureThreshold    int  `yaml:"failureThreshold"`
}

type Ports struct {
	ContainerPort int    `yaml:"containerPort"`
	Name          string `yaml:"name"`
	Protocol      string `yaml:"protocol"`
	Port          int    `yaml:"port"`
	TargetPort    int    `yaml:"targetPort"`
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
	ReadinessProbe  ReadinessProbe `yaml:"readinessProbe"`
	ImagePullPolicy string         `yaml:"imagePullPolicy"`
	Name            string         `yaml:"name"`
	Ports           []Ports        `yaml:"ports"`
	Resources       Resources      `yaml:"resources"`
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
	Replicas int      `yaml:"replicas"`
	Selector Selector `yaml:"selector"`
	Template Template `yaml:"template"`
	Ports    []Ports  `yaml:"ports"`
	Type     string   `yaml:"type"`
	TLS      []TLS    `yaml:"tls"`
	Rules    []Rules  `yaml:"rules"`
}
