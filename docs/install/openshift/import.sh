#!/usr/bin/env bash
# imports the Onix CMDB templates into OpenShift
oc create -f onix-ephemeral.yml -n openshift
oc create -f onix-persistent.yml -n openshift
oc create -f onix-rds.yml -n openshift
oc create -f oxkube.yml -n openshift
