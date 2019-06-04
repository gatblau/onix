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
	"errors"
	apiV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const maxRetries = 5

var startTime time.Time

// watch for changes to one or more kubernetes resources and calls registered
// events handlers
type Sentinel struct {
	client kubernetes.Interface
	config *Config
}

// starts observing for K8S resource changes
func (s *Sentinel) Start() error {
	if s.config == nil {
		return errors.New("Config must me provided.")
	}

	// gets a k8s client
	client, err := getKubeClient()

	if err != nil {
		return err
	}

	// creates controllers to listen to events
	if s.config.Observe.Pod {
		s.startWatcher(client, &apiV1.Pod{}, "pod")
	}
	if s.config.Observe.Namespace {
		s.startWatcher(client, &apiV1.Namespace{}, "namespace")
	}
	if s.config.Observe.PersistentVolume {
		s.startWatcher(client, &apiV1.PersistentVolume{}, "persistent_volume")
	}
	if s.config.Observe.Service {
		s.startWatcher(client, &apiV1.Service{}, "service")
	}

	// creates a channel to pass kernel signals to terminate the main process
	terminateCh := make(chan os.Signal, 1)

	// sends any SIGINT or SIGTERM signals to the channel
	signal.Notify(terminateCh, syscall.SIGINT)  // interrupt ctrl+c
	signal.Notify(terminateCh, syscall.SIGTERM) // terminate

	// waits until any termination signals are raised
	<-terminateCh

	return nil
}

// starts a new watcher K8S controller to listen for status change events and trigger a handling function
func (s *Sentinel) startWatcher(client kubernetes.Interface, objType runtime.Object, resourceType string) {
	// creates an informer to receive notifications of state changes for a given collection of objects.
	// objects are identified by its API group, kind/resource, namespace, and name.
	informer := cache.NewSharedIndexInformer(
		// the controller wants to list and watch all pods in all namespaces
		&cache.ListWatch{
			ListFunc: func(options metaV1.ListOptions) (runtime.Object, error) {
				return s.newList(client, options, resourceType)
			},
			WatchFunc: func(options metaV1.ListOptions) (watch.Interface, error) {
				return s.newWatch(client, options, resourceType)
			},
		},
		objType,
		0, // skip re-sync
		cache.Indexers{},
	)

	// creates a new controller to handle object status changes
	watcher := newWatcher(informer, resourceType, *s.config)

	// run the controller
	go watcher.run()
}

// returns a watch interface for the specified resource type (e.g. pod)
func (s *Sentinel) newWatch(client kubernetes.Interface, options metaV1.ListOptions, resourceType string) (watch.Interface, error) {
	switch resourceType {
	case "configmap":
		return client.CoreV1().ConfigMaps(metaV1.NamespaceAll).Watch(options)
	case "daemonset":
		return client.ExtensionsV1beta1().DaemonSets(metaV1.NamespaceAll).Watch(options)
	case "deployment":
		return client.AppsV1beta1().Deployments(metaV1.NamespaceAll).Watch(options)
	case "namespace":
		return client.CoreV1().Namespaces().Watch(options)
	case "ingress":
		return client.ExtensionsV1beta1().Ingresses(metaV1.NamespaceAll).Watch(options)
	case "job":
		return client.BatchV1().Jobs(metaV1.NamespaceAll).Watch(options)
	case "persistent_volume":
		return client.CoreV1().PersistentVolumes().Watch(options)
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
	default:
		return nil, nil
	}
}

// returns a list (runtime object) for the specified resource type (e.g. pod)
func (s *Sentinel) newList(client kubernetes.Interface, options metaV1.ListOptions, resourceType string) (runtime.Object, error) {
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
	default:
		return nil, nil
	}
}
