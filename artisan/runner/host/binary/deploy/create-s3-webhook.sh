#!/bin/bash
#==================================================================================================
#environment variables required are MINIO_SVC_MINIO_ROOT_USER,  MINIO_SVC_MINIO_ROOT_PASSWORD,
#HOST_RUNNER_PORT, WEBHOOK_NAME, CMD_KEY, MINIO_BUCKET

# Note:- MINIO_SVC_MINIO_ROOT_USER, MINIO_SVC_MINIO_ROOT_PASSWORD,HOST_RUNNER_PORT will be already 
# available as a part of the host runner setup

#User has to set the value for following environment variable WEBHOOK_NAME, CMD_KEY, MINIO_BUCKET_NAME
#WEBHOOK_NAME and MINIO_BUCKET_NAME can be anything of your choice
#CMD_KEY value must be matching with the ITEM_KEY value in cmdb for your command
#==================================================================================================
FILE="$SERVICE_HOME"/.env
if [ -f "$FILE" ]; then
  set -a # automatically export all variables
  source "$FILE"
  set +a
fi

if [[ -z "${MINIO_SVC_MINIO_ROOT_USER}" ]] || [[ -z "${MINIO_SVC_MINIO_ROOT_PASSWORD}" ]]; then
  echo "value for environment variable MINIO_SVC_MINIO_ROOT_USER and/or MINIO_SVC_MINIO_ROOT_PASSWORD is not set, verify whether these variable are there in $SERVICE_HOME/.env"
  exit 100
fi

if [[ -z "${WEBHOOK_NAME}" ]] || [[ -z "${CMD_KEY}" ]] || [[ -z "${MINIO_BUCKET_NAME}" ]]; then
  echo "value for environment variable WEBHOOK_NAME and/or CMD_KEY, MINIO_BUCKET_NAME is not set by the user before executing artisan function"
  echo "WEBHOOK_NAME and MINIO_BUCKET_NAME can be anything of your choice"
  echo "CMD_KEY value must be matching with the ITEM_KEY value in cmdb for your command"
  exit 100
fi

mc -version
if [ $? -eq 0 ]; then
    echo "minio client is already there"
else
    echo "installing minio client package"
    wget https://dl.minio.io/client/mc/release/linux-amd64/mc
    chmod +x mc
    sudo mv mc /usr/local/bin
    mc alias set local-minio http://localhost:9000 ${MINIO_SVC_MINIO_ROOT_USER} ${MINIO_SVC_MINIO_ROOT_PASSWORD}
fi

LOCAL_IP=$(hostname -I | cut -d' ' -f1)
  mc admin config set local-minio notify_webhook:"${WEBHOOK_NAME}" endpoint="http://${LOCAL_IP}:${HOST_RUNNER_PORT}/host/${CMD_KEY}"
  echo "webhook created with end point http://${LOCAL_IP}:${HOST_RUNNER_PORT}/host/${CMD_KEY}"
  echo "restarting the minio server"
  mc admin service restart local-minio
  sleep 2
  mc event add local-minio/"${MINIO_BUCKET_NAME}" arn:minio:sqs::"${WEBHOOK_NAME}":webhook --event put
