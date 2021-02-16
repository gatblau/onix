/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package tkn

import (
	"io"
	"text/template"
)

// event listener merging information
type EL struct {
	AppName string
	AppIcon string
	GitURI  string
}

// merges the template and its values into the passed in writer
func (p *EL) Merge(w io.Writer) error {
	t, err := template.New("event-listener").Parse(evListener)
	if err != nil {
		return err
	}
	return t.Execute(w, p)
}

// event listener template definition
const evListener = `
apiVersion: triggers.tekton.dev/v1alpha1
kind: EventListener
metadata:
  labels:
      app.openshift.io/runtime: {{.AppIcon}}
  name: {{.AppName}}
spec:
  serviceAccountName: pipeline
  triggers:
  - bindings:
    - name: {{.AppName}}
    template:
      name: {{.AppName}}
---
apiVersion: v1
kind: Route
metadata:
  annotations:
    description: Route for application's https service.
  labels:
    application: {{.AppName}}-https
  name: el-{{.AppName}}
spec:
  port:
    targetPort: "8080"
  tls:
    insecureEdgeTerminationPolicy: Redirect
    termination: edge
  to:
    kind: Service
    name: el-{{.AppName}}
---
apiVersion: triggers.tekton.dev/v1alpha1
kind: TriggerBinding
metadata:
  name: {{.AppName}}
spec:
  params:
  - name: git-repo-url
    value: $(body.project.web_url)
  - name: git-repo-name
    value: $(body.repository.name)
  - name: git-revision
    value: $(body.commits[0].id)
---
apiVersion: triggers.tekton.dev/v1alpha1
kind: TriggerTemplate
metadata:
  name: {{.AppName}}
spec:
  params:
  - name: git-repo-url
    description: The git repository url
  - name: git-revision
    description: The git revision
    default: master
  - name: git-repo-name
    description: The name of the deployment to be created / patched

  resourcetemplates:
  - apiVersion: tekton.dev/v1alpha1
    kind: PipelineResource
    metadata:
      name: $(params.git-repo-name)-git-repo-$(uid)
    spec:
      type: git
      params:
      - name: revision
        value: $(params.git-revision)
      - name: url
        value: {{.GitURI}}

  - apiVersion: tekton.dev/v1beta1
    kind: PipelineRun
    metadata:
      name: build-deploy-$(params.git-repo-name)-$(uid)
    spec:
      serviceAccountName: pipeline
      pipelineRef:
        name: {{.AppName}}-build-and-deploy
      resources:
      - name: {{.AppName}}-git-repo
        resourceRef:
          name: $(params.git-repo-name)-git-repo-$(uid)
      params:
      - name: deployment-name
        value: $(params.git-repo-name)
`
