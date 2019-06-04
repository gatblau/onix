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
	apiV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"time"
)

// a k8s controller that watches for changes to the state of a particular resource
// and triggers the execution of a publisher (e.g. calling a web hook,
// putting a message in a broker, etc.)
type Watcher struct {
	resourceType string
	queue        workqueue.RateLimitingInterface
	informer     cache.SharedIndexInformer
	config       Config
}

// creates a new controller to watch for changes in status of a specific resource
func newWatcher(informer cache.SharedIndexInformer, resourceType string, config Config) *Watcher {
	logrus.Info(fmt.Sprintf("sentinel-%s", resourceType))
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	var event Event
	var err error
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			event.key, err = cache.MetaNamespaceKeyFunc(obj)
			event.eventType = "create"
			event.resourceType = resourceType
			logrus.WithField("pkg", fmt.Sprintf("sentinel-%s", resourceType))
			logrus.Infof("Processing ADD to %v: %s", resourceType, event.key)
			if err == nil {
				queue.Add(event)
			}
		},
		UpdateFunc: func(old, new interface{}) {
			event.key, err = cache.MetaNamespaceKeyFunc(old)
			event.eventType = "update"
			event.resourceType = resourceType
			logrus.WithField("pkg", fmt.Sprintf("sentinel-%s", resourceType))
			logrus.Infof("Processing UPDATE to %v: %s", resourceType, event.key)
			if err == nil {
				queue.Add(event)
			}
		},
		DeleteFunc: func(obj interface{}) {
			event.key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			event.eventType = "delete"
			event.resourceType = resourceType
			event.namespace = getMetaData(obj).Namespace
			logrus.WithField("pkg", fmt.Sprintf("sentinel-%s", resourceType))
			logrus.Infof("Processing DELETE to %v: %s", resourceType, event.key)
			if err == nil {
				queue.Add(event)
			}
		},
	})

	return &Watcher{
		resourceType: resourceType,
		informer:     informer,
		queue:        queue,
		config:       config,
	}
}

// runs the controller
func (w *Watcher) run() {
	// creates a stopCh channel to stop the controller when required
	stopCh := make(chan struct{})
	defer close(stopCh)

	// catches a crash and logs an error
	// TODOs: check if it can be removed as apiserver will handle panics
	defer runtime.HandleCrash()

	// shut downs the queue when it is time
	defer w.queue.ShutDown()

	logrus.Info("Starting a new Sentinel controller")
	startTime = time.Now().Local()

	// starts and runs the shared informer
	// the informer will be stopped when the stop channel is closed
	go w.informer.Run(stopCh)

	// Wait for all involved caches to be synced, before processing items from the queue is started
	if !cache.WaitForCacheSync(stopCh, w.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	logrus.Info("Sentinel controller synchronised and ready")

	// loops until the stop channel is closed, running the worker every second
	wait.Until(w.processQueue, time.Second, stopCh)
}

// process the items in the controller's queue
func (w *Watcher) processQueue() {
	// loops until the worker queue is shut down
	for w.nextItem() {
	}
}

// process the next item in the controller's queue
func (w *Watcher) nextItem() bool {
	// waits until there is a new item in the working queue
	key, shutdown := w.queue.Get()

	// if queue shuts down then quit
	if shutdown {
		return false
	}

	// tells the queue that we are done processing this key
	// this unblocks the key for other workers and allows safe parallel processing because two pods
	// with the same key are never processed in parallel.
	defer w.queue.Done(key)

	// passes the queue item to the registered handler(s)
	err := w.handleEvent(key.(Event))

	// handles the result of the previous operation
	// if something went wrong during the execution of the business logic, triggers a retry mechanism
	w.handleResult(err, key)

	// continues processing
	return true
}

func (w *Watcher) handleEvent(newEvent Event) error {
	obj, exists, err := w.informer.GetIndexer().GetByKey(newEvent.key)
	if err != nil {
		return fmt.Errorf("Error getting object with key '%w' from store: '%v'", newEvent.key, err)
	}
	if !exists {
		logrus.Infof("'%w' with key '%w' does not exist anymore\n", newEvent.resourceType, newEvent.key)
	} else {
		// Note that you also have to check the uid if you have a local controlled ObservedResources, which
		// is dependent on the actual instance, to detect that a Pod was recreated with the same name
		logrus.Infof("Sync/Add/Update for Pod %w\n", obj.(*apiV1.Pod).GetName())

		// get object metadata
		meta := getMetaData(obj)

		// we must have an event Handler configured
		if w.config.Handler == nil {
			return fmt.Errorf("Event Handler not defined.")
		}

		// handleEvent events based on its type
		switch newEvent.eventType {
		case "create":
			// compare CreationTimestamp and serverStartTime and alert only on latest events
			// Could be Replaced by using Delta or DeltaFIFO
			if meta.CreationTimestamp.Sub(startTime).Seconds() > 0 {
				w.config.Handler.OnCreate(newEvent, obj)
				return nil
			}
		case "update":
			w.config.Handler.OnUpdate(newEvent, obj)
			return nil
		case "delete":
			w.config.Handler.OnDelete(newEvent, obj)
			return nil
		}
	}
	return nil
}

// checks if an error has happened triggering retry
// or stops retrying if there is no error
func (w *Watcher) handleResult(err error, key interface{}) {
	if err == nil {
		// indicates that the item is finished being retried.
		// it doesn't matter whether it's for permanent failing or for success,
		// it stops the rate limiter from tracking it.
		w.queue.Forget(key)
		return
	} else if w.queue.NumRequeues(key) < maxRetries {
		// this controller retries a specified number of times if something goes wrong
		// after which, stops trying
		logrus.Errorf("Error processing %w (will retry): %v", key.(Event).key, err)

		// re-queue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		w.queue.AddRateLimited(key)
		return
	} else {
		// err != nil and too many retries, then give up
		w.queue.Forget(key)

		// reports to an external entity that, even after several retries,
		// the key could not be successfully handled
		runtime.HandleError(err)

		// logs the error
		logrus.Errorf("Error processing %w (giving up): %v", key.(Event).key, err)
	}
}
