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
	"github.com/gatblau/onix/artisan/crypto"
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
		service, err := b.buildService(s)
		if err != nil {
			return nil, err
		}
		rsx = append(rsx, *service)
		ingress, err := b.buildIngress(s)
		if err != nil {
			return nil, err
		}
		if ingress != nil {
			rsx = append(rsx, *ingress)
		}
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
				APIVersion: k8s.CoreVersion,
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
	if value, exists := svc.Attributes["tls"]; exists {
		switch strings.ToLower(value) {
		// auto generate tls certificate secret
		case "auto":
			cert, key, err := crypto.SelfSignedBase64()
			if err != nil {
				return nil, err
			}
			secret := k8s.Secret{
				APIVersion: k8s.CoreVersion,
				Kind:       "Secret",
				Metadata: &k8s.Metadata{
					Name: tlsSecretName(svc),
					Annotations: k8s.Annotations{
						Description: fmt.Sprintf("certificate for TLS encryption of ingress endpoint"),
					},
				},
				Type: "kubernetes.io/tls",
				Data: &map[string]string{
					"tls.crt": cert,
					"tls.key": key,
				},
			}
			content, err := yaml.Marshal(secret)
			if err != nil {
				return rsx, nil
			}
			rsx = append(rsx, DeploymentRsx{
				Name:    fmt.Sprintf("%s.yaml", tlsSecretName(svc)),
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
			Labels: k8s.Labels{
				App:     svc.Name,
				Version: b.Manifest.Version,
			},
			Name: svc.Name,
			Annotations: k8s.Annotations{
				Description: fmt.Sprintf("Deployment configuration for %s", svc.Description),
			},
		},
		Spec: k8s.Spec{
			Replicas: getReplicas(svc),
			Selector: k8s.Selector{
				MatchLabels: k8s.MatchLabels{
					App:     svc.Name,
					Version: b.Manifest.Version,
				},
			},
			Template: k8s.Template{
				Metadata: k8s.Metadata{
					Labels: k8s.Labels{
						App:     svc.Name,
						Version: b.Manifest.Version,
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

// getReplicas get the number of pod replicas based on the highly_available attribute
func getReplicas(svc SvcRef) int {
	if value, exists := svc.Attributes["highly_available"]; exists {
		replicas, err := strconv.Atoi(value)
		if err != nil {
			return 1
		}
		return replicas
	}
	return 1
}

func (b *KubeBuilder) buildService(svc SvcRef) (*DeploymentRsx, error) {
	ports, err := getServicePorts(svc)
	if err != nil {
		return nil, err
	}
	s := k8s.Service{
		APIVersion: k8s.CoreVersion,
		Kind:       "Service",
		Metadata: k8s.Metadata{
			Annotations: k8s.Annotations{Description: svc.Description},
			Labels: k8s.Labels{
				App:     svc.Name,
				Version: b.Manifest.Version,
			},
			Name: svcName(svc),
		},
		Spec: k8s.Spec{
			Selector: k8s.Selector{
				App:     svc.Name,
				Version: b.Manifest.Version,
			},
			Ports: ports,
		},
	}
	content, err := yaml.Marshal(s)
	if err != nil {
		return nil, err
	}
	return &DeploymentRsx{
		Name:    fmt.Sprintf("%s-service.yaml", svc.Name),
		Content: content,
		Type:    K8SResource,
	}, nil
}

func svcName(svc SvcRef) string {
	return fmt.Sprintf("%s-service", normalisedName(svc.Name))
}

func (b *KubeBuilder) buildIngress(svc SvcRef) (*DeploymentRsx, error) {
	if host, exists := svc.Attributes["publish"]; exists {
		port, err := strconv.Atoi(svc.Port)
		if err != nil {
			return nil, err
		}
		tls, err := getTLS(svc, host)
		i := &k8s.Ingress{
			APIVersion: k8s.NetVersion,
			Kind:       "Ingress",
			Metadata: k8s.Metadata{
				Annotations: k8s.Annotations{Description: fmt.Sprintf("publishes the %s service", svc.Name)},
				Labels: k8s.Labels{
					App:     svc.Name,
					Version: b.Manifest.Version,
				},
				Name: fmt.Sprintf("%s-ingress", strings.Replace(strings.ToLower(svc.Name), "_", "-", -1)),
			},
			Spec: k8s.Spec{
				TLS: tls,
				Rules: []k8s.Rules{
					{
						Host: host,
						HTTP: k8s.HTTP{Paths: []k8s.Paths{
							{Path: "/",
								Backend: k8s.Backend{
									ServiceName: svcName(svc),
									ServicePort: port,
								},
							},
						}},
					},
				},
			},
		}
		content, err := yaml.Marshal(i)
		if err != nil {
			return nil, err
		}
		return &DeploymentRsx{
			Name:    fmt.Sprintf("%s-ingress.yaml", svcName(svc)),
			Content: content,
			Type:    K8SResource,
		}, nil
	}
	return nil, nil
}

func getTLS(svc SvcRef, host string) ([]k8s.TLS, error) {
	if value, exists := svc.Attributes["tls"]; exists {
		switch strings.ToLower(value) {
		case "auto":
			return []k8s.TLS{
				{
					Hosts:      []string{host},
					SecretName: ingressTlsSecretName(svc),
				},
			}, nil
		default:
			return nil, nil
		}
	}
	return nil, nil
}

func ingressTlsSecretName(svc SvcRef) string {
	return fmt.Sprintf("%s-ingress-tls-secret", normalisedName(svc.Name))
}

func getContainers(svc SvcRef) ([]k8s.Containers, error) {
	ports, err := getDeploymentPorts(svc)
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

func getDeploymentPorts(svc SvcRef) ([]k8s.Ports, error) {
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

func getServicePorts(svc SvcRef) ([]k8s.Ports, error) {
	target, err := strconv.Atoi(svc.Info.Port)
	if err != nil {
		return nil, err
	}
	published, err := strconv.Atoi(svc.Port)
	if err != nil {
		return nil, err
	}
	return []k8s.Ports{
		{
			Port:       published,
			TargetPort: target,
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
	return fmt.Sprintf("%s-secret", normalisedName(v.Name))
}

func nestedVars(value string) []string {
	r, _ := regexp.Compile("\\${(?P<NAME>[^}]+)}")
	return r.FindAllString(value, -1)
}

func normalisedName(name string) string {
	return strings.Replace(strings.Replace(strings.ToLower(name), "_", "-", -1), " ", "-", -1)
}

func tlsSecretName(svc SvcRef) string {
	return fmt.Sprintf("%s-ingress-tls-cert-secret", normalisedName(svc.Name))
}
