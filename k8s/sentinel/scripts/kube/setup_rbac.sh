#!/usr/bin/env bash
# creates the cluster role containing the rights to access the k8s API verbs for resources
kubectl create -f cluster_role.yml

# creates a new namespace for Sentinel
kubectl create namespace sentinel

# create service account
kubectl create serviceaccount sentinel -n sentinel

# binds the sentinel service account to the cluster role created before
kubectl create clusterrolebinding sentinel-cluster-rule --clusterrole=resource-watcher --serviceaccount=sentinel:sentinel
