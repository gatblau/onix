#!/usr/bin/env bash
oc delete dc/onixwapi
oc delete svc/onixwapi
oc delete route/onixwapi
oc delete is/onixwapi
oc delete secret/onix-wapi-admin-user
oc delete secret/onix-wapi-writer-user
oc delete secret/onix-wapi-reader-user

oc delete dc/onixdb
oc delete svc/onixdb
oc delete is/onixdb
oc delete pvc/onixdb
oc delete secret/onix-db-admin
oc delete secret/onix-db-user

oc delete dc/oxkube
oc delete svc/oxkube
oc delete is/oxkube-snapshot
oc delete is/oxkube
oc delete secret/onix-user-secret
oc delete route/oxkube