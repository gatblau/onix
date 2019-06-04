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

// configuration info for the Sentinel process
type Config struct {
	Handler Publisher
	Observe ObservedResources
}

// the type of resources that can be observed by the controller
type ObservedResources struct {
	Service          bool
	Pod              bool
	PersistentVolume bool
	Namespace        bool
	//Deployment            bool
	//ReplicationController bool
	//ReplicaSet            bool
	//DaemonSet             bool
	//Job                   bool
	//Secret                bool
	//ConfigMap             bool
	//Ingress               bool
}

// the interface implemented by a controller publisher
type Publisher interface {
	OnCreate(e Event, o interface{})
	OnDelete(e Event, o interface{})
	OnUpdate(e Event, o interface{})
}

type Event struct {
	key          string
	eventType    string
	namespace    string
	resourceType string
}

// Event represent an event got from k8s api server
// Events from different endpoints need to be casted to KubewatchEvent
// before being able to be handled by Handler
type PublishedEvent struct {
	Namespace string
	Kind      string
	Component string
	Host      string
	Reason    string
	Status    string
	Name      string
}
