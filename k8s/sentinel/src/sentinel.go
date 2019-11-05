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
	"fmt"
	"github.com/sirupsen/logrus"
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
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const maxRetries = 5

var startTime time.Time

// watch for changes to one or more kubernetes resources and calls the configured publisher
type Sentinel struct {
	client    kubernetes.Interface
	config    *Config
	publisher Publisher
	log       *logrus.Entry
}

// starts observing for K8S resource changes
func (s *Sentinel) Start() error {
	// reset the start time for Sentinel
	startTime = time.Now()

	// loads the configuration
	c, err := NewConfig()
	if err == nil {
		s.config = &c
	} else {
		return err
	}

	// adds the platform field to the logger
	s.log = logrus.WithFields(logrus.Fields{
		"platform": s.config.Platform,
	})

	// try and parse the logging level in the configuration
	level, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		// if the value was not recognised then return the error
		s.log.Errorf("Failed to recognise value LogLevel entry in the configuration: %s.", err)
		return err
	}
	// otherwise sets the logging level for the entire system
	logrus.SetLevel(level)
	s.log.Infof("%s has been set as the logger level.", strings.ToUpper(c.LogLevel))

	// registers the publisher used by Sentinel
	s.publisher, err = s.getPublisher()

	if err != nil {
		// if the resolution of the publisher failed then return the error
		return err
	} else {
		s.log.Infof("%s publisher has been registered.", strings.ToUpper(s.config.Publishers.Publisher))
	}

	// gets a k8s client
	client, err := s.getKubeClient()

	if err != nil {
		// if the creation of the k8s client failed then return the error
		return err
	} else {
		s.log.Infof("Kubernetes client created.")
	}

	// launch the required controllers to listen to object state changes
	s.startWatchers(client)

	s.log.Infof("Sentinel is ready and looking out for changes.")

	// waits until a SIGINT or SIGTERM signal is raised
	// creates a channel to pass kernel signals to terminate the main process
	terminateCh := make(chan os.Signal, 1)
	// sends any SIGINT or SIGTERM signals to the channel
	signal.Notify(terminateCh, syscall.SIGINT)
	// interrupt ctrl+c
	signal.Notify(terminateCh, syscall.SIGTERM)
	// terminate
	// waits until any termination signals are raised
	<-terminateCh

	return nil
}

// launch k8s controllers to listen to object state changes
func (s *Sentinel) startWatchers(client kubernetes.Interface) {
	if s.config.Observe.Service {
		s.startWatcher(client, &coreV1.Service{}, "service")
	}
	if s.config.Observe.Pod {
		s.startWatcher(client, &coreV1.Pod{}, "pod")
	}
	if s.config.Observe.PersistentVolume {
		s.startWatcher(client, &coreV1.PersistentVolume{}, "persistent_volume")
	}
	if s.config.Observe.PersistentVolumeClaim {
		s.startWatcher(client, &coreV1.PersistentVolumeClaim{}, "persistent_volume_claim")
	}
	if s.config.Observe.Namespace {
		s.startWatcher(client, &coreV1.Namespace{}, "namespace")
	}
	if s.config.Observe.Deployment {
		s.startWatcher(client, &appsV1.Deployment{}, "deployment")
	}
	if s.config.Observe.ReplicationController {
		s.startWatcher(client, &coreV1.ReplicationController{}, "replication_controller")
	}
	if s.config.Observe.ReplicaSet {
		s.startWatcher(client, &appsV1.ReplicaSet{}, "replicaset")
	}
	if s.config.Observe.DaemonSet {
		s.startWatcher(client, &extV1beta1.DaemonSet{}, "daemonset")
	}
	if s.config.Observe.Job {
		s.startWatcher(client, &batchV1.Job{}, "job")
	}
	if s.config.Observe.Secret {
		s.startWatcher(client, &coreV1.Secret{}, "secret")
	}
	if s.config.Observe.ConfigMap {
		s.startWatcher(client, &coreV1.ConfigMap{}, "configmap")
	}
	if s.config.Observe.Ingress {
		s.startWatcher(client, &extV1beta1.Ingress{}, "ingress")
	}
	if s.config.Observe.ServiceAccount {
		s.startWatcher(client, &coreV1.ServiceAccount{}, "service_account")
	}
	if s.config.Observe.ClusterRole {
		s.startWatcher(client, &rbacV1.ClusterRole{}, "cluster_role")
	}
	if s.config.Observe.ResourceQuota {
		s.startWatcher(client, &coreV1.ResourceQuota{}, "resource_quota")
	}
	if s.config.Observe.NetworkPolicy {
		s.startWatcher(client, &netV1.NetworkPolicy{}, "network_policy")
	}
}

// starts a new watcher K8S controller to listen for status change events and trigger a handling function
func (s *Sentinel) startWatcher(client kubernetes.Interface, objType runtime.Object, resourceType string) {
	// creates an informer to receive notifications of state changes for a given collection of objects.
	// objects are identified by its API group, kind/resource, namespace, and name.
	informer := cache.NewSharedIndexInformer(
		// the controller wants to list and watch all pods in all namespaces
		&cache.ListWatch{
			ListFunc: func(options metaV1.ListOptions) (runtime.Object, error) {
				return newList(client, options, resourceType)
			},
			WatchFunc: func(options metaV1.ListOptions) (watch.Interface, error) {
				return newWatch(client, options, resourceType)
			},
		},
		objType,
		0, // skip re-sync
		cache.Indexers{},
	)

	// creates a new controller to handle object status changes
	watcher := newWatcher(informer, resourceType, *s)

	// run the controller
	go watcher.run()
}

// gets the publisher specified in the configuration
func (s *Sentinel) getPublisher() (Publisher, error) {
	var pub Publisher
	switch s.config.Publishers.Publisher {
	case "webhook":
		pub = new(WebhookPub)
	case "broker":
		pub = new(BrokerPub)
	case "logger":
		pub = new(LoggerPub)
	default:
		return nil, fmt.Errorf(
			"Failed to register a publisher: the value '%s' in the configuration could not be recognised.",
			s.config.Publishers.Publisher)
	}
	pub.Init(s.config, s.log)
	return pub, nil
}

// gets an instance of the k8s client
func (s *Sentinel) getKubeClient() (kubernetes.Interface, error) {
	config, err := s.getKubeConfig()
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		s.log.Fatalf("Can not create kubernetes client: %v.", err)
		return nil, err
	}
	return client, nil
}

// gets the K8S client configuration either inside or outside of the cluster depending on
// whether the kube config file could be found
func (s *Sentinel) getKubeConfig() (*rest.Config, error) {
	// k8s client configuration
	var config *rest.Config

	// gets the path to the kube config file
	kubeConfigFile := fmt.Sprintf("%s/%s", os.Getenv("HOME"), s.config.KubeConfig)

	if _, err := os.Stat(kubeConfigFile); err == nil {
		s.log.Info("Kube config file found: attempting out of cluster configuration.")
		// if the kube config file exists then do an outside of cluster configuration
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigFile)
		if err != nil {
			s.log.Errorf("Could not create out of cluster configuration: %v.", err)
			return nil, err
		}
	} else if os.IsNotExist(err) {
		s.log.Info("Kube config file not found: attempting in cluster configuration.")
		// the kube config file was not found then do in cluster configuration
		config, err = rest.InClusterConfig()
		if err != nil {
			s.log.Errorf("Could not find the K8S client configuration. "+
				"Are you running Sentinel in a container that has not been deployed in Kubernetes?.\n "+
				"The error message was: %v.", err)
			return nil, err
		}
	} else {
		// kube config might be there or not but it failed anyway :(
		if err != nil {
			s.log.Errorf("Could not figure out the Kube client configuration: %v.", err)
			return nil, err
		}
	}
	return config, nil
}
