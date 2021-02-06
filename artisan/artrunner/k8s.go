/*
  Onix Config Manager - Artisan Runner
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

type K8S struct {
	cfg       *rest.Config
	decoder   runtime.Serializer
	inCluster bool
}

func NewK8S() (*K8S, error) {
	config, inCluster, err := getKubeConfig()
	if err != nil {
		return nil, err
	}
	return &K8S{
		cfg:       config,
		decoder:   yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme),
		inCluster: inCluster,
	}, nil
}

// Kubernetes Server-Side Apply
func (k *K8S) Apply(yamlResource string, ctx context.Context) error {
	// 1. Prepare a RESTMapper to find GVR
	dc, err := discovery.NewDiscoveryClientForConfig(k.cfg)
	if err != nil {
		return err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	// 2. Prepare the dynamic client
	dyn, err := dynamic.NewForConfig(k.cfg)
	if err != nil {
		return err
	}

	// 3. Decode YAML manifest into unstructured.Unstructured
	obj := &unstructured.Unstructured{}
	_, gvk, err := k.decoder.Decode([]byte(yamlResource), nil, obj)
	if err != nil {
		return err
	}

	// 4. Find GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	// 5. Obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(mapping.Resource)
	}

	// 6. Marshal object into JSON
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	// 7. Create or Update the object with SSA
	//     types.ApplyPatchType indicates SSA.
	//     FieldManager specifies the field owner ID.
	_, err = dr.Patch(ctx, obj.GetName(), types.ApplyPatchType, data, metav1.PatchOptions{
		FieldManager: "artisan-runner",
	})
	return err
}

// gets the K8S client configuration either inside or outside of the cluster depending on
// whether the kube config file could be found
func getKubeConfig() (*rest.Config, bool, error) {
	// k8s client configuration
	var config *rest.Config

	inCluster := false

	// gets the path to the kube config file
	kubeConfigFile := fmt.Sprintf("%s/.kube/config", os.Getenv("HOME"))

	// if a .kube/config file is found, then not running in K8S
	if _, err := os.Stat(kubeConfigFile); err == nil {
		log.Printf("%s file found: attempting connection from outside of the cluster", kubeConfigFile)
		// if the kube config file exists then do an outside of cluster configuration
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigFile)
		if err != nil {
			log.Printf("could not create out of cluster configuration: %v", err)
			return nil, inCluster, err
		}
	} else if os.IsNotExist(err) {
		log.Printf("%s file not found: attempting connection from within the cluster", kubeConfigFile)
		// the kube config file was not found then do in cluster configuration
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Printf("could not find the K8S client configuration. "+
				"you cannot run the Build Manager in a container that has not been deployed in Kubernetes?\n "+
				"the error message was: %v.", err)
			return nil, inCluster, err
		}
		inCluster = true
	} else {
		// kube config might be there or not but it failed anyway :(
		if err != nil {
			log.Printf("could not figure out the Kube client configuration: %v", err)
			return nil, inCluster, err
		}
	}
	return config, inCluster, nil
}
