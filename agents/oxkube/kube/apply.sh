#!/usr/bin/env bash
oc apply -f namespace.yml
oc project demo-project
oc apply -f quota.yml
oc apply -f replication_ctrl.yml
oc apply -f service.yml
