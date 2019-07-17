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
package main

import "fmt"

// get all K8S objects of a specific type in the specified cluster namespace
func (c *Client) getObjectsInNamespace(cluster string, namespace string, objType K8SOBJ) ([]Item, error) {
	filters := map[string]string{
		"type":  objType.String(),
		"attrs": fmt.Sprintf("cluster,%s|namespace,%s", cluster, namespace),
	}
	// get pods in the namespace first
	podsObj, err := c.getResource("item", "", filters)

	if err != nil {
		return nil, err
	}
	// unwraps the response into a list of pod items
	return podsObj.(*ResultList).Values, nil
}
