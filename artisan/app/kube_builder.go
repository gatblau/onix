/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package app

import (
	"encoding/base64"
	"fmt"
	"github.com/gatblau/onix/artisan/app/k8s"
	"gopkg.in/yaml.v2"
	"regexp"
	"strconv"
	"strings"
)

type KubeBuilder struct {
	Manifest Manifest
}

// newKubeBuilder called internally by NewBuilder()
func newKubeBuilder(appMan Manifest) Builder {
	return &KubeBuilder{Manifest: appMan}
}

func (b *KubeBuilder) Build() ([]DeploymentRsx, error) {
	rsx := make([]DeploymentRsx, 0)
	for _, s := range b.Manifest.Services {
		deployment, err := b.buildDeployment(s)
		if err != nil {
			return nil, err
		}
		rsx = append(rsx, *deployment)
		secrets, err := b.buildSecrets(s)
		if err != nil {
			return nil, err
		}
		rsx = append(rsx, secrets...)
	}
	return rsx, nil
}

func (b *KubeBuilder) buildSecrets(svc SvcRef) ([]DeploymentRsx, error) {
	rsx := make([]DeploymentRsx, 0)
	for _, v := range svc.Info.Var {
		if v.Secret {
			value, err := b.getVarValue(v.Value)
			if err != nil {
				return rsx, err
			}
			secret := k8s.Secret{
				APIVersion: k8s.SecretsVersion,
				Kind:       "Secret",
				Metadata: &k8s.Metadata{
					Name:        secretName(v),
					Annotations: k8s.Annotations{Description: v.Description},
				},
				Type: "Opaque",
				Data: &map[string]string{
					v.Name: base64.StdEncoding.EncodeToString([]byte(value)),
				},
			}
			content, err := yaml.Marshal(secret)
			if err != nil {
				return rsx, nil
			}
			rsx = append(rsx, DeploymentRsx{
				Name:    fmt.Sprintf("%s-%s-secret.yaml", svc.Name, strings.Replace(strings.ToLower(v.Name), "_", "-", -1)),
				Content: content,
				Type:    K8SResource,
			})
		}
	}
	return rsx, nil
}

func (b *KubeBuilder) getVarValue(name string) (string, error) {
	// check if the name contains other variables
	vs := nestedVars(name)
	var vName string
	var merged bool
	for _, vv := range vs {
		if strings.HasPrefix(vv, "${") && strings.HasSuffix(vv, "}") {
			vName = vv[2 : len(vv)-1]
		}
		for _, v := range b.Manifest.Var.Items {
			if v.Name == vName {
				name = strings.Replace(name, vv, v.Value, -1)
				merged = true
			}
		}
	}
	if merged {
		return name, nil
	}
	return "", fmt.Errorf("variable %s not found", name)
}

func (b *KubeBuilder) buildDeployment(svc SvcRef) (*DeploymentRsx, error) {
	containers, err := getContainers(svc)
	if err != nil {
		return nil, err
	}
	d := &k8s.Deployment{
		APIVersion: k8s.AppsVersion,
		Kind:       "Deployment",
		Metadata: k8s.Metadata{
			Labels: k8s.Labels{App: svc.Name},
			Name:   svc.Name,
			Annotations: k8s.Annotations{
				Description: fmt.Sprintf("Deployment configuration for %s", svc.Description),
			},
		},
		Spec: k8s.Spec{
			Replicas: 1,
			Selector: k8s.Selector{
				MatchLabels: k8s.MatchLabels{
					App: svc.Name,
				},
			},
			Template: k8s.Template{
				Metadata: k8s.Metadata{
					Labels: k8s.Labels{
						App: svc.Name,
					},
					Annotations: k8s.Annotations{
						Description: svc.Description,
					},
					Name: fmt.Sprintf("%s-deployment", strings.ToLower(svc.Name)),
				},
				Spec: k8s.TemplateSpec{
					Containers:                    containers,
					TerminationGracePeriodSeconds: 30,
				},
			},
		},
	}
	content, err := yaml.Marshal(d)
	if err != nil {
		return nil, err
	}
	return &DeploymentRsx{
		Name:    fmt.Sprintf("%s-deployment.yaml", svc.Name),
		Content: content,
		Type:    K8SResource,
	}, nil
}

func getContainers(svc SvcRef) ([]k8s.Containers, error) {
	ports, err := getPorts(svc)
	if err != nil {
		return nil, err
	}
	return []k8s.Containers{
		{
			Env:             getK8SEnv(svc),
			Image:           svc.Image,
			ImagePullPolicy: getImagePullPolicy(svc.Image),
			Name:            fmt.Sprintf("%s Application", svc.Name),
			Ports:           ports,
		},
	}, nil
}

func getPorts(svc SvcRef) ([]k8s.Ports, error) {
	p, err := strconv.Atoi(svc.Info.Port)
	if err != nil {
		return nil, err
	}
	return []k8s.Ports{
		{
			ContainerPort: p,
		},
	}, nil
}

// getImagePullPolicy work out what pull policy to use based on the  image tag
func getImagePullPolicy(image string) string {
	// if the image has a tag and the tag is not "latest", pull if not present
	r, _ := regexp.Compile(":[\\w][\\w.-]{0,127}")
	tag := r.FindString(image)
	if len(tag) > 0 && tag[1:] != "latest" {
		return "IfNotPresent"
	}
	// if image does not have a tag or tag is "latest" then pull always
	return "Always"
}

func getK8SEnv(svc SvcRef) []k8s.Env {
	env := make([]k8s.Env, 0)
	for _, v := range svc.Info.Var {
		if v.Secret {
			env = append(env, k8s.Env{
				Name: v.Name,
				ValueFrom: k8s.ValueFrom{
					SecretKeyRef: k8s.SecretKeyRef{
						Key:  secretName(v),
						Name: v.Name,
					}},
			})
		} else {
			env = append(env, k8s.Env{
				Name:  v.Name,
				Value: v.Value,
			})
		}
	}
	return env
}

func secretName(v Var) string {
	return fmt.Sprintf("%s-secret", strings.ToLower(strings.Replace(v.Name, "_", "-", -1)))
}

func nestedVars(value string) []string {
	r, _ := regexp.Compile("\\${(?P<NAME>[^}]+)}")
	return r.FindAllString(value, -1)
}
