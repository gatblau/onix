#!/usr/bin/env bash

# creates an account to run tiller
kubectl create serviceaccount -n kube-system tiller

# makes the tiller account a cluster-admin
kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller

# installs helm on kubernetes
helm init --service-account tiller

# updates the local helm repo info
helm repo update
