/*
   Onix CMDB - Copyright (c) 2018-2019 by www.gatblau.org

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
   Unless required by applicable law or agreed to in writing, software distributed under
   the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
   either express or implied.
   See the License for the specific language governing permissions and limitations under the License.

   Contributors to this project, hereby assign copyright in this code to the project,
   to be licensed under the same terms as the rest of the code.
*/
package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	appsV1 "k8s.io/api/apps/v1"
	batchV1 "k8s.io/api/batch/v1"
	apiV1 "k8s.io/api/core/v1"
	extV1beta1 "k8s.io/api/extensions/v1beta1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

// gets an environment variable or the defaultValue if the variable is not set
func getEnv(key string, defaultValue string) string {
	value := ""
	if os.Getenv(key) != "" {
		value = os.Getenv(key)
	} else {
		value = defaultValue
	}
	return value
}

// gets the kube config default path
func getKubeConfigPath() string {
	return getEnv("KUBECONFIG", fmt.Sprintf("%s/.kube/config", os.Getenv("HOME")))
}

// gets the K8S client configuration either inside or outside of the cluster depending on
// whether the kube config file could be found
func getKubeConfig() (*rest.Config, error) {
	// k8s client configuration
	var config *rest.Config

	// gets the path to the kube config file
	kubeConfigFile := getKubeConfigPath()

	if _, err := os.Stat(kubeConfigFile); err == nil {
		logrus.Info("Kube config file found: attempting out of cluster configuration")
		// if the kube config file exists then do an outside of cluster configuration
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigFile)
		if err != nil {
			logrus.Errorf("Could create out of cluster configuration: %v", err)
			return nil, err
		}

	} else if os.IsNotExist(err) {
		logrus.Info("Kube config file not found: attempting in cluster configuration")
		// the kube config file was not found then do in cluster configuration
		config, err = rest.InClusterConfig()
		if err != nil {
			logrus.Errorf("Could create in cluster configuration: %v", err)
			return nil, err
		}
	} else {
		// kube config might be there or not but it failed anyway :(
		if err != nil {
			logrus.Errorf("Could not figure out the Kube client configuration: %v", err)
			return nil, err
		}
	}
	return config, nil
}

// gets an instance of the k8s client
func getKubeClient() (kubernetes.Interface, error) {
	config, err := getKubeConfig()
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Fatalf("Can not create kubernetes client: %v", err)
		return nil, err
	}
	return client, nil
}

// gets the metadata for the persisted resource
func getMetaData(obj interface{}) metaV1.ObjectMeta {
	var objectMeta metaV1.ObjectMeta
	switch object := obj.(type) {
	case *appsV1.Deployment:
		objectMeta = object.ObjectMeta
	case *apiV1.ReplicationController:
		objectMeta = object.ObjectMeta
	case *appsV1.ReplicaSet:
		objectMeta = object.ObjectMeta
	case *appsV1.DaemonSet:
		objectMeta = object.ObjectMeta
	case *apiV1.Service:
		objectMeta = object.ObjectMeta
	case *apiV1.Pod:
		objectMeta = object.ObjectMeta
	case *batchV1.Job:
		objectMeta = object.ObjectMeta
	case *apiV1.PersistentVolume:
		objectMeta = object.ObjectMeta
	case *apiV1.Namespace:
		objectMeta = object.ObjectMeta
	case *apiV1.Secret:
		objectMeta = object.ObjectMeta
	case *extV1beta1.Ingress:
		objectMeta = object.ObjectMeta
	}
	return objectMeta
}
