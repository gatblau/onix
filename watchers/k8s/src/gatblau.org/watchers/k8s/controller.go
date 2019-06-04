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
	"fmt"
	"github.com/Sirupsen/logrus"
	apiV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
	"time"
)

type Controller struct {
	log      *logrus.Entry
	queue    workqueue.RateLimitingInterface
	informer cache.SharedIndexInformer
	config   Config
}

func newController(informer cache.SharedIndexInformer, resourceType string) *Controller {
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	var event Event
	var err error
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			event.key, err = cache.MetaNamespaceKeyFunc(obj)
			event.eventType = "create"
			event.resourceType = resourceType
			logrus.WithField("pkg", fmt.Sprintf("sentinel-%s", resourceType)).Infof("Processing add to %v: %s", resourceType, event.key)
			if err == nil {
				queue.Add(event)
			}
		},
		UpdateFunc: func(old, new interface{}) {
			event.key, err = cache.MetaNamespaceKeyFunc(old)
			event.eventType = "update"
			event.resourceType = resourceType
			logrus.WithField("pkg", fmt.Sprintf("sentinel-%s", resourceType)).Infof("Processing update to %v: %s", resourceType, event.key)
			if err == nil {
				queue.Add(event)
			}
		},
		DeleteFunc: func(obj interface{}) {
			event.key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			event.eventType = "delete"
			event.resourceType = resourceType
			event.namespace = getMetaData(obj).Namespace
			logrus.WithField("pkg", fmt.Sprintf("sentinel-%s", resourceType)).Infof("Processing delete to %v: %s", resourceType, event.key)
			if err == nil {
				queue.Add(event)
			}
		},
	})

	return &Controller{
		log:      logrus.WithField("pkg", fmt.Sprintf("sentinel-%s", resourceType)),
		informer: informer,
		queue:    queue,
	}
}

// runs the controller
func (c *Controller) run(stop <-chan struct{}) {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	c.log.Info("Starting a new Sentinel controller")
	startTime = time.Now().Local()

	go c.informer.Run(stop)

	// Wait for all involved caches to be synced, before processing items from the queue is started
	if !cache.WaitForCacheSync(stop, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	c.log.Info("Sentinel controller synced and ready")

	wait.Until(c.runWorker, time.Second, stop)
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
		// continue looping
	}
}

func (c *Controller) processNextItem() bool {
	// waits until there is a new item in the working queue
	key, shutdown := c.queue.Get()

	// if queue shuts down then quit
	if shutdown {
		return false
	}

	// tells the queue that we are done processing this key
	// this unblocks the key for other workers and allows safe parallel processing because two pods
	// with the same key are never processed in parallel.
	defer c.queue.Done(key)

	err := c.handleEvent(key.(Event))

	// Handle the error if something went wrong during the execution of the business logic
	c.handleErr(err, key)

	//if err == nil {
	//	// No error, reset the ratelimit counters
	//	c.queue.Forget(key)
	//} else if c.queue.NumRequeues(key) < maxRetries {
	//	c.log.Errorf("Error processing %c (will retry): %v", key.(Event).key, err)
	//	c.queue.AddRateLimited(key)
	//} else {
	//	// err != nil and too many retries
	//	c.log.Errorf("Error processing %c (giving up): %v", key.(Event).key, err)
	//	c.queue.Forget(key)
	//	runtime.HandleError(err)
	//}

	return true
}

func (c *Controller) handleEvent(newEvent Event) error {
	obj, exists, err := c.informer.GetIndexer().GetByKey(newEvent.key)
	if err != nil {
		return fmt.Errorf("Error getting object with key %c from store: %v", newEvent.key, err)
	}
	if !exists {
		fmt.Printf("%c with key  %c does not exist anymore\n", newEvent.resourceType, newEvent.key)
	} else {
		// Note that you also have to check the uid if you have a local controlled ObservedResources, which
		// is dependent on the actual instance, to detect that a Pod was recreated with the same name
		fmt.Printf("Sync/Add/Update for Pod %c\n", obj.(*apiV1.Pod).GetName())

		// get object metadata
		meta := getMetaData(obj)

		// we must have an event Handler configured
		if c.config.Handler == nil {
			return errors.New("Event Handler not defined.")
		}

		// handleEvent events based on its type
		switch newEvent.eventType {
		case "create":
			// compare CreationTimestamp and serverStartTime and alert only on latest events
			// Could be Replaced by using Delta or DeltaFIFO
			if meta.CreationTimestamp.Sub(startTime).Seconds() > 0 {
				c.config.Handler.OnCreate(newEvent, obj)
				return nil
			}
		case "update":
			c.config.Handler.OnUpdate(newEvent, obj)
			return nil
		case "delete":
			c.config.Handler.OnDelete(newEvent, obj)
			return nil
		}
	}
	return nil
}

// handleErr checks if an error happened and makes sure we will retry later.
func (c *Controller) handleErr(err error, key interface{}) {
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history.
		c.queue.Forget(key)
		return
	}

	// This controller retries 5 times if something goes wrong. After that, it stops trying.
	if c.queue.NumRequeues(key) < 5 {
		klog.Infof("Error syncing pod %v: %v", key, err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)

	// Report to an external entity that, even after several retries, we could not successfully handleEvent this key
	runtime.HandleError(err)

	klog.Infof("Dropping pod %q out of the queue: %v", key, err)
}
