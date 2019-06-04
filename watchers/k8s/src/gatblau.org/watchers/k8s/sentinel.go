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
	"time"
)

const maxRetries = 5

var startTime time.Time

/*

 */
type Sentinel struct {
	client kubernetes.Interface
	config *Config
}

func (s *Sentinel) Start() error {
	if s.config == nil {
		return errors.New("Config must me provided.")
	}

	// gets a k8s client
	client, err := getKubeClient()

	if err != nil {
		return err
	}

	if s.config.Observe.Pod {
		// creates a new instance of the listwatcher
		informer := cache.NewSharedIndexInformer(
			// the controller wants to list and watch all pods in all namespaces
			&cache.ListWatch{
				ListFunc: func(options metaV1.ListOptions) (runtime.Object, error) {
					return client.CoreV1().Pods(metaV1.NamespaceAll).List(options)
				},
				WatchFunc: func(options metaV1.ListOptions) (watch.Interface, error) {
					return client.CoreV1().Pods(metaV1.NamespaceAll).Watch(options)
				},
			},
			&apiV1.Pod{},
			0, // skip resync
			cache.Indexers{},
		)

		c := newController(informer, "pod")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.run(stopCh)
	}

	return nil
}
