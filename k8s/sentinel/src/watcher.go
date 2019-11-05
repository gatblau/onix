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
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"strings"
	"time"
)

// a k8s controller that watches for changes to the state of a particular resource
// and triggers the execution of a publisher (e.g. calling a web hook,
// putting a message in a broker, etc.)
type Watcher struct {
	objType   string
	queue     workqueue.RateLimitingInterface
	informer  cache.SharedIndexInformer
	publisher Publisher
	log       *logrus.Entry
}

// creates a new controller to watch for changes in status of a specific resource
func newWatcher(informer cache.SharedIndexInformer, objType string, s Sentinel) *Watcher {
	s.log.Tracef("Creating %s watcher.", strings.ToUpper(objType))
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	var change StatusChange
	var err error

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			meta := getMetaData(obj)
			change.key, err = cache.MetaNamespaceKeyFunc(obj)
			change.Name = meta.Name
			change.Namespace = meta.Namespace
			change.Type = "CREATE"
			change.Kind = objType
			change.Time = time.Now().UTC()
			change.Host = s.config.Platform
			addToQueue(queue, change, err, s.log)
		},
		UpdateFunc: func(obj, new interface{}) {
			meta := getMetaData(obj)
			change.key, err = cache.MetaNamespaceKeyFunc(obj)
			change.Name = meta.Name
			change.Namespace = meta.Namespace
			change.Type = "UPDATE"
			change.Kind = objType
			change.Time = time.Now().UTC()
			change.Host = s.config.Platform
			addToQueue(queue, change, err, s.log)
		},
		DeleteFunc: func(obj interface{}) {
			meta := getMetaData(obj)
			change.key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			change.Name = meta.Name
			change.Namespace = meta.Namespace
			change.Type = "DELETE"
			change.Kind = objType
			change.Time = time.Now().UTC()
			change.Host = s.config.Platform
			addToQueue(queue, change, err, s.log)
		},
	})

	return &Watcher{
		informer:  informer,
		queue:     queue,
		publisher: s.publisher,
		log:       s.log,
		objType:   objType,
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

	w.log.Tracef("%s watcher starting.", strings.ToUpper(w.objType))
	startTime = time.Now().Local()

	// starts and runs the shared informer
	// the informer will be stopped when the stop channel is closed
	go w.informer.Run(stopCh)

	// Wait for all involved caches to be synced, before processing items from the queue is started
	if !cache.WaitForCacheSync(stopCh, w.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync."))
		return
	}

	w.log.Tracef("%s watcher synchronised and ready.", strings.ToUpper(w.objType))

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
		w.log.Tracef("%s queue has shut down.", strings.ToUpper(w.objType))
		return false
	}

	// tells the queue that we are done processing this key
	// this unblocks the key for other workers and allows safe parallel processing because two pods
	// with the same key are never processed in parallel.
	defer w.queue.Done(key)

	// passes the queue item to the registered handler(s)
	err := w.publish(key.(StatusChange))

	// handles the result of the previous operation
	// if something went wrong during the execution of the business logic, triggers a retry mechanism
	w.handleResult(err, key)

	// continues processing
	return true
}

// publish the state change
func (w *Watcher) publish(change StatusChange) error {
	w.log.Tracef("Ready to publish %s changes for %s %s.", change.Type, strings.ToUpper(change.Kind), change.key)
	obj, exists, err := w.informer.GetIndexer().GetByKey(change.key)
	if !exists {
		w.log.Tracef("%s %s does not exist anymore.", strings.ToUpper(change.Kind), change.key)
	}
	if err != nil {
		return fmt.Errorf("failed to retrieve object with key %s: %s", change.key, err)
	} else {
		// get object metadata
		meta := getMetaData(obj)

		// publish events based on its type
		switch change.Type {
		case "CREATE":
			// compare CreationTimestamp and serverStartTime and alert only on latest events
			// Could be Replaced by using Delta or DeltaFIFO
			if meta.CreationTimestamp.Sub(startTime).Seconds() > 0 {
				w.log.Tracef("Calling Publisher.OnCreate(change -> %+v).", change)
				w.publisher.Publish(
					Event{
						Change: change,
						Object: obj,
					})
			} else {
				w.log.Tracef("Change occurred %s before starting Sentinel, so not calling publisher.", meta.CreationTimestamp.Sub(startTime))
			}
		case "UPDATE":
			w.log.Tracef("Calling Publisher.OnUpdate(change -> %+v).", change)
			w.publisher.Publish(Event{
				Change: change,
				Object: obj,
			})
			return nil
		case "DELETE":
			w.log.Tracef("Calling Publisher.OnDelete(change -> %+v).", change)
			w.publisher.Publish(Event{
				Change: change,
				Object: obj,
			})
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
		w.log.Tracef("StatusChange for %s has been processed.", key.(StatusChange).key)
		w.queue.Forget(key)
		return
	} else if w.queue.NumRequeues(key) < maxRetries {
		// this controller retries a specified number of times if something goes wrong
		// after which, stops trying
		w.log.Errorf("Error processing %s (will retry): %s.", key.(StatusChange).key, err)

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
		if err != nil {
			w.log.Errorf("Error processing %s (giving up): %s.", key.(StatusChange).key, err)
		} else {
			w.log.Errorf("Error processing %s: too many retries, giving up!", key.(StatusChange).key)
		}
	}
}

// add a change to the processing queue
func addToQueue(queue workqueue.RateLimitingInterface, change StatusChange, err error, log *logrus.Entry) {
	if err == nil {
		log.Tracef("Queueing %s change for %s %s.", change.Type, strings.ToUpper(change.Kind), change.key)
		queue.Add(change)
	} else {
		log.Errorf("Error adding %s change for %s %s to processing queue.", change.Type, strings.ToUpper(change.Kind), change.key)
	}
}
