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

mc -version
LOCAL_IP=$(hostname -I | cut -d' ' -f1)
if [ $? -eq 0 ]; then
    echo "minio client is already there"
else
    echo "installing minio client package"
    wget https://dl.minio.io/client/mc/release/linux-amd64/mc
    chmod +x mc
    sudo mv mc /usr/local/bin
fi

mc alias set local-minio http://"${LOCAL_IP}":9000 $MINIO_SVC_MINIO_ROOT_USER $MINIO_SVC_MINIO_ROOT_PASSWORD
echo "alias for this minio set as local-minio"

if [[ -z "${MINIO_BUCKET_NAME}" ]]; then
  echo "value for environment variable MINIO_BUCKET_NAME is not set"
else
  mc mb local-minio/"$MINIO_BUCKET_NAME"
  echo "minio bucket $MINIO_BUCKET_NAME created successfully"
fi