#!/usr/bin/env bash
#
#    Sentinel - Copyright (c) 2019 by www.gatblau.org
#
#    Licensed under the Apache License, Version 2.0 (the "License");
#    you may not use this file except in compliance with the License.
#    You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
#    Unless required by applicable law or agreed to in writing, software distributed under
#    the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
#    either express or implied.
#    See the License for the specific language governing permissions and limitations under the License.
#
#    Contributors to this project, hereby assign copyright in this code to the project,
#    to be licensed under the same terms as the rest of the code.
#
#    Description:
#       Use this script to create the Sentinel namespace and configure security ready to deploy
#       the container image.
#       You must be cluster admin to run this script.
#

# creates the cluster role containing the rights to access the k8s API verbs for resources
oc create -f cluster_role.yml

# creates a new namespace to host the Sentinel container
oc new-project sentinel --display-name="Sentinel" --description="Raise notifications for k8s resource changes"

# creates a service account to run the sentinel
oc create serviceaccount sentinel

# binds the sentinel service account to the cluster role created before
oc create clusterrolebinding sentinel-cluster-rule --clusterrole=resource-watcher --serviceaccount=sentinel:sentinel
