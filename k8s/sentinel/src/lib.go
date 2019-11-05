/*
   Sentinel - Copyright (c) 2019 by www.gatblau.org

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
	"bytes"
	"encoding/json"
	"fmt"
	appsV1 "k8s.io/api/apps/v1"
	batchV1 "k8s.io/api/batch/v1"
	coreV1 "k8s.io/api/core/v1"
	extV1beta1 "k8s.io/api/extensions/v1beta1"
	netV1 "k8s.io/api/networking/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"math/rand"
)

// gets the metadata for the persisted resource
func getMetaData(obj interface{}) metaV1.ObjectMeta {
	var objectMeta metaV1.ObjectMeta
	switch object := obj.(type) {
	case *appsV1.Deployment:
		objectMeta = object.ObjectMeta
	case *coreV1.ReplicationController:
		objectMeta = object.ObjectMeta
	case *appsV1.ReplicaSet:
		objectMeta = object.ObjectMeta
	case *appsV1.DaemonSet:
		objectMeta = object.ObjectMeta
	case *coreV1.Service:
		objectMeta = object.ObjectMeta
	case *coreV1.Pod:
		objectMeta = object.ObjectMeta
	case *batchV1.Job:
		objectMeta = object.ObjectMeta
	case *coreV1.PersistentVolume:
		objectMeta = object.ObjectMeta
	case *coreV1.PersistentVolumeClaim:
		objectMeta = object.ObjectMeta
	case *coreV1.Namespace:
		objectMeta = object.ObjectMeta
	case *coreV1.ConfigMap:
		objectMeta = object.ObjectMeta
	case *coreV1.Secret:
		objectMeta = object.ObjectMeta
	case *extV1beta1.Ingress:
		objectMeta = object.ObjectMeta
	case *coreV1.ServiceAccount:
		objectMeta = object.ObjectMeta
	case *rbacV1.ClusterRole:
		objectMeta = object.ObjectMeta
	case *coreV1.ResourceQuota:
		objectMeta = object.ObjectMeta
	case *netV1.NetworkPolicy:
		objectMeta = object.ObjectMeta
	}
	return objectMeta
}

// returns a watch interface for the specified resource type (e.g. pod)
func newWatch(client kubernetes.Interface, options metaV1.ListOptions, resourceType string) (watch.Interface, error) {
	switch resourceType {
	case "configmap":
		return client.CoreV1().ConfigMaps(metaV1.NamespaceAll).Watch(options)
	case "daemonset":
		return client.ExtensionsV1beta1().DaemonSets(metaV1.NamespaceAll).Watch(options)
	case "deployment":
		return client.AppsV1().Deployments(metaV1.NamespaceAll).Watch(options)
	case "namespace":
		return client.CoreV1().Namespaces().Watch(options)
	case "ingress":
		return client.ExtensionsV1beta1().Ingresses(metaV1.NamespaceAll).Watch(options)
	case "job":
		return client.BatchV1().Jobs(metaV1.NamespaceAll).Watch(options)
	case "persistent_volume":
		return client.CoreV1().PersistentVolumes().Watch(options)
	case "persistent_volume_claim":
		return client.CoreV1().PersistentVolumeClaims(metaV1.NamespaceAll).Watch(options)
	case "pod":
		return client.CoreV1().Pods(metaV1.NamespaceAll).Watch(options)
	case "replicaset":
		return client.ExtensionsV1beta1().ReplicaSets(metaV1.NamespaceAll).Watch(options)
	case "replication_controller":
		return client.CoreV1().ReplicationControllers(metaV1.NamespaceAll).Watch(options)
	case "secret":
		return client.CoreV1().Secrets(metaV1.NamespaceAll).Watch(options)
	case "service":
		return client.CoreV1().Services(metaV1.NamespaceAll).Watch(options)
	case "service_account":
		return client.CoreV1().ServiceAccounts(metaV1.NamespaceAll).Watch(options)
	case "cluster_role":
		return client.RbacV1().ClusterRoles().Watch(options)
	case "resource_quota":
		return client.CoreV1().ResourceQuotas(metaV1.NamespaceAll).Watch(options)
	case "network_policy":
		return client.NetworkingV1().NetworkPolicies(metaV1.NamespaceAll).Watch(options)
	default:
		return nil, nil
	}
}

// returns a list (runtime object) for the specified resource type (e.g. pod)
func newList(client kubernetes.Interface, options metaV1.ListOptions, resourceType string) (runtime.Object, error) {
	switch resourceType {
	case "configmap":
		return client.CoreV1().ConfigMaps(metaV1.NamespaceAll).List(options)
	case "daemonset":
		return client.ExtensionsV1beta1().DaemonSets(metaV1.NamespaceAll).List(options)
	case "deployment":
		return client.AppsV1beta1().Deployments(metaV1.NamespaceAll).List(options)
	case "namespace":
		return client.CoreV1().Namespaces().List(options)
	case "ingress":
		return client.ExtensionsV1beta1().Ingresses(metaV1.NamespaceAll).List(options)
	case "job":
		return client.BatchV1().Jobs(metaV1.NamespaceAll).List(options)
	case "persistent_volume":
		return client.CoreV1().PersistentVolumes().List(options)
	case "persistent_volume_claim":
		return client.CoreV1().PersistentVolumeClaims(metaV1.NamespaceAll).List(options)
	case "pod":
		return client.CoreV1().Pods(metaV1.NamespaceAll).List(options)
	case "replicaset":
		return client.ExtensionsV1beta1().ReplicaSets(metaV1.NamespaceAll).List(options)
	case "replication_controller":
		return client.CoreV1().ReplicationControllers(metaV1.NamespaceAll).List(options)
	case "secret":
		return client.CoreV1().Secrets(metaV1.NamespaceAll).List(options)
	case "service":
		return client.CoreV1().Services(metaV1.NamespaceAll).List(options)
	case "service_account":
		return client.CoreV1().ServiceAccounts(metaV1.NamespaceAll).List(options)
	case "cluster_role":
		return client.RbacV1().ClusterRoles().List(options)
	case "resource_quota":
		return client.CoreV1().ResourceQuotas(metaV1.NamespaceAll).List(options)
	case "network_policy":
		return client.NetworkingV1().NetworkPolicies(metaV1.NamespaceAll).List(options)
	default:
		return nil, nil
	}
}

// creates a random string of the specified length
func randString(strLen int) string {
	bytes := make([]byte, strLen)
	for i := 0; i < strLen; i++ {
		// 65 - 90: range of characters to use in the random string
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

// gets a random integer
func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

// converts the specified object into a pretty looking JSON string
func toJSON(obj interface{}) ([]byte, error) {
	// serialises the object into JSON applying indentation to format the output
	jsonBytes, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("Failed to marshall object: %s", err)
	}
	return jsonBytes, nil
}

// gets a byte reader from a specified object
func getJSONBytesReader(data interface{}) (*bytes.Reader, error) {
	jsonBytes, err := json.Marshal(data)
	return bytes.NewReader(jsonBytes), err
}

// checks if the specified string value is contained in the passed-in slice
func contains(slic []string, value string) bool {
	for _, element := range slic {
		if element == value {
			return true
		}
	}
	return false
}
