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

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"strings"
)

// gets the kube meta-model for Onix
func (c *Client) getModel() Payload {
	return &Data{
		Models: []Model{
			Model{
				Key:         K8SModel,
				Name:        "Kubernetes Resource Model",
				Description: "Defines the item and link types that describe Kubernetes resources in a given Namespace.",
			},
		},
		ItemTypes: []ItemType{
			ItemType{
				Key:         K8SCluster,
				Name:        "Kubernetes Cluster",
				Description: "An open-source system for automating deployment, scaling, and management of containerized applications.",
				Model:       K8SModel,
			},
			ItemType{
				Key:         K8SNamespace,
				Name:        "Namespace",
				Description: "A way to divide cluster resources between multiple users or teams providing virtual areas to deploy project resources.",
				Model:       K8SModel,
			},
			ItemType{
				Key:         K8SResourceQuota,
				Name:        "Resource Quota",
				Description: "A set of constraints that limit aggregate resource consumption per namespace.",
				Model:       K8SModel,
			},
			ItemType{
				Key:         K8SPod,
				Name:        "Pod",
				Description: "Encapsulates an application’s container (or, in some cases, multiple containers), storage resources, a unique network IP, and options that govern how the container(s) should run.",
				Model:       K8SModel,
			},
			ItemType{
				Key:         K8SService,
				Name:        "Service",
				Description: "Exposes an application running on a set of Pods as a network service.",
				Model:       K8SModel,
				Filter: map[string]interface{}{
					"filters": []interface{}{
						map[string]interface{}{
							"selector": []interface{}{
								map[string]interface{}{
									"default": "$.selector",
								},
							},
						},
					},
				},
			},
			ItemType{
				Key:  K8SIngress,
				Name: "Ingress (Route)",
				Description: "Exposes HTTP and HTTPS routes from outside the cluster to services within the cluster.\n" +
					"Traffic routing is controlled by rules defined on the Ingress resource.",
				Model: K8SModel,
			},
			ItemType{
				Key:         K8SReplicationController,
				Name:        "Replication Controller",
				Description: "Ensures that a specified number of pod replicas are running at any one time.",
				Model:       K8SModel,
			},
			ItemType{
				Key:         K8SPersistentVolumeClaim,
				Name:        "Persistent Volume Claim",
				Description: "A claim to a piece of storage in the cluster made by a pod.",
				Model:       K8SModel,
			},
		},
		LinkTypes: []LinkType{
			LinkType{
				Key:         K8SLink,
				Name:        "Kubernetes Resource Link Type",
				Description: "Links Kubernetes resources.",
				Model:       K8SModel,
			},
		},
		LinkRules: []LinkRule{
			LinkRule{
				Key:              fmt.Sprintf("%s->%s", K8SCluster, K8SNamespace),
				Name:             "K8S Cluster to Namespace Rule",
				Description:      "A cluster contains one or more namespaces.",
				LinkTypeKey:      K8SLink,
				StartItemTypeKey: K8SCluster,
				EndItemTypeKey:   K8SNamespace,
			},
			LinkRule{
				Key:              fmt.Sprintf("%s->%s", K8SNamespace, K8SResourceQuota),
				Name:             "K8S Namespace to Resource Quota Rule",
				Description:      "A namespace has a resource quota.",
				LinkTypeKey:      K8SLink,
				StartItemTypeKey: K8SNamespace,
				EndItemTypeKey:   K8SResourceQuota,
			},
			LinkRule{
				Key:              fmt.Sprintf("%s->%s", K8SNamespace, K8SPod),
				Name:             "K8S Namespace to Pod Rule",
				Description:      "A namespace contains one or more pods.",
				LinkTypeKey:      K8SLink,
				StartItemTypeKey: K8SNamespace,
				EndItemTypeKey:   K8SPod,
			},
			LinkRule{
				Key:              fmt.Sprintf("%s->%s", K8SPod, K8SPersistentVolumeClaim),
				Name:             "K8S Pod to Persistent Volume Claim Rule",
				Description:      "A pod makes one or more persistent volume claims.",
				LinkTypeKey:      K8SLink,
				StartItemTypeKey: K8SPod,
				EndItemTypeKey:   K8SPersistentVolumeClaim,
			},
			LinkRule{
				Key:              fmt.Sprintf("%s->%s", K8SPod, K8SReplicationController),
				Name:             "K8S Pod to Replication Controller Rule",
				Description:      "A pod is controlled by a replication controller.",
				LinkTypeKey:      K8SLink,
				StartItemTypeKey: K8SPod,
				EndItemTypeKey:   K8SReplicationController,
			},
			LinkRule{
				Key:              fmt.Sprintf("%s->%s", K8SPod, K8SService),
				Name:             "K8S Pod to Service Rule",
				Description:      "A pod is accessed by a service.",
				LinkTypeKey:      K8SLink,
				StartItemTypeKey: K8SPod,
				EndItemTypeKey:   K8SService,
			},
			LinkRule{
				Key:              fmt.Sprintf("%s->%s", K8SService, K8SIngress),
				Name:             "K8S Service to Ingress Rule",
				Description:      "A service is published via an Ingress route.",
				LinkTypeKey:      K8SLink,
				StartItemTypeKey: K8SService,
				EndItemTypeKey:   K8SIngress,
			},
		},
	}
}

// generic link payload
func (c *Client) getLink(startItem string, endItem string) Payload {
	return &Link{
		Key:          fmt.Sprintf("%s->%s", startItem, endItem),
		StartItemKey: startItem,
		EndItemKey:   endItem,
		Type:         K8SLink,
	}
}

// cluster item payload
func (c *Client) getClusterItem(event []byte) *Item {
	host := gjson.GetBytes(event, Cluster)
	return &Item{
		Key:         clusterKey(host.String()),
		Name:        fmt.Sprintf("%s Container Platform", strings.ToUpper(host.String())),
		Description: "A Kubernetes Cluster instance.",
		Type:        K8SCluster,
	}
}

func item(event []byte, iType string, oType string) (*Item, error) {
	cluster := gjson.GetBytes(event, Cluster)
	name := gjson.GetBytes(event, Key)
	spec := gjson.GetBytes(event, SpecInfo)
	created := gjson.GetBytes(event, Created)
	namespace := gjson.GetBytes(event, Namespace)
	var key string
	if oType == "ns" {
		// if the object is a namespace then do not repeat it in the key
		key = NS(event)
	} else {
		key = itemKey(event, oType)
	}
	item := &Item{
		Key:       key,
		Name:      name.String(),
		Meta:      MAP{},
		Attribute: MAP{},
		Type:      iType,
	}
	item.Attribute["cluster"] = cluster.String()
	item.Attribute["namespace"] = namespace.String()
	item.Attribute["created"] = created.String()
	addMap(event, item, Labels)
	addMap(event, item, Annotations)
	err := json.Unmarshal([]byte(spec.String()), &item.Meta)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// adds the content of a map in the event to the item attributes
func addMap(event []byte, item *Item, path string) {
	mapObj := gjson.GetBytes(event, path).Map()
	for key, value := range mapObj {
		item.Attribute[key] = value.String()
	}
}

// gets the unique key for a service
func itemKey(event []byte, oType string) string {
	key := gjson.GetBytes(event, Key).String()
	return fmt.Sprintf("%s-%s-%s", NS(event), oType, key)
}

func clusterKey(clusterKey string) string {
	return fmt.Sprintf("k8s-%s", clusterKey)
}

// use to identify the namespace an object is in, in all but Namespace events
func NS(event []byte) string {
	cluster := gjson.GetBytes(event, Cluster).String()
	namespace := gjson.GetBytes(event, Namespace).String()
	if len(namespace) == 0 {
		namespace = gjson.GetBytes(event, Key).String()
	}
	return fmt.Sprintf("%s-ns-%s", clusterKey(cluster), namespace)
}
