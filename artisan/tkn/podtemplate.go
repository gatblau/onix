package tkn

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

type PodTemplate struct {
	SecurityContext *PipelineSecurityContext `yaml:"securityContext,omitempty"`
}

type PipelineSecurityContext struct {
	RunAsNonRoot bool `yaml:"runAsNonRoot"`
	FsGroup      int  `yaml:"fsGroup"`
	RunAsUser    int  `yaml:"runAsUser"`
}

type StepSecurityContext struct {
	Privileged   bool `yaml:"privileged,omitempty"`
	RunAsNonRoot bool `yaml:"runAsNonRoot"`
	RunAsUser    int  `yaml:"runAsUser"`
}
