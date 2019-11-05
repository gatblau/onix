/*
   Onix Kube - Copyright (c) 2019 by www.gatblau.org

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
package src

import "time"

// represent the event sent by the publisher
type Event struct {
	// information about the status change
	Change StatusChange
	// the k8s object that changed
	Object interface{}
}

// information about the K8S object status change
type StatusChange struct {
	key       string    `json:"key"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Namespace string    `json:"namespace"`
	Kind      string    `json:"kind"`
	Time      time.Time `json:"time"`
	Host      string    `json:"host"`
}
