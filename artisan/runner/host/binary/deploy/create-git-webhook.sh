#!/bin/bash

#curl -X POST -H "Accept: application/vnd.github.v3+json" https://api.github.com/repos/OWNER/REPO/hooks \
#  -d '{"name":"web","active":true,"events":["push"],"config":{"url":"https://example.com/webhook","content_type":"json","insecure_ssl":"0"}}'

#==================================================================================================
#environment variables required are GIT_SVC_ADMIN_USERNAME,  REPO_NAME,
#GIT_TOKEN, FLOW_KEY, HOST_RUNNER_PORT

# Note:- GIT_SVC_ADMIN_USERNAME, HOST_RUNNER_PORT will be already 
# available as a part of the host runner setup

#User has to set the value for following environment variable GIT_TOKEN, FLOW_KEY, REPO_NAME
#REPO_NAME value must match to the repo already existing in git for user GIT_SVC_ADMIN_USERNAME
#GIT_TOKEN value must match to the git token of user GIT_SVC_ADMIN_USERNAME
#FLOW_KEY value must be matching with the ITEM_KEY value in cmdb for your flow
#==================================================================================================

FILE="$SERVICE_HOME"/.env
if [ -f "$FILE" ]; then
  set -a # automatically export all variables
  source "$FILE"
  set +a
fi

if [[ -z "${GIT_SVC_ADMIN_USERNAME}" ]] || [[ -z "${GIT_SVC_ADMIN_PASSWORD}" ]]; then
  echo "value for environment variable GIT_SVC_ADMIN_USERNAME and or GIT_SVC_ADMIN_PASSWORD is not set, verify whether this variable are there in ${SERVICE_HOME}/.env"
  exit 100
fi

if [[ -z "${REPO_NAME}" ]] || [[ -z "${FLOW_KEY}" ]]; then
  echo "value for environment variable REPO_NAME and/or FLOW_KEY is not set by the user before executing artisan function"
  echo "REPO_NAME must be of the repo which already exists for user $GIT_SVC_ADMIN_USERNAME"
  echo "FLOW_KEY value must be matching with the ITEM_KEY value in cmdb for your flow"
  exit 100
fi

LOCAL_IP=$(ip route get 1.2.3.4 | awk '{print $7}')
GURL=http://"$LOCAL_IP":"${HOST_RUNNER_PORT}"/webhook/"${FLOW_KEY}"/push
curl --location --request POST 'http://localhost:8084/api/v1/repos/'${GIT_SVC_ADMIN_USERNAME}'/'${REPO_NAME}'/hooks' \
--header 'Content-Type: application/json' \
--user "${GIT_SVC_ADMIN_USERNAME}:${GIT_SVC_ADMIN_PASSWORD}" \
--header 'Content-Type: application/json' \
--data-raw '{
  "active": true,
  "branch_filter": "*",
  "config": {
    "content_type": "json",
    "url": "'"$GURL"'",
    "http_method": "post"
  },
  "events": [
    "push"
  ],
  "type": "gitea"
}'