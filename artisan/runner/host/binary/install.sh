#!/bin/bash

WD=/tmp/host-runner-deploy
mkdir -p $WD

# Set standard variables for use below and in template
export RUNNER_HOME=${HOME}/.config/runner
export RUNNER_UID=$(id -u)
export RUNNER_GID=$(id -g)

echo "------------------------------------------------------------"
echo "Stopping the service"

# Stop existing current service
if [ $(systemctl is-active host-runner) == "active" ]
then
    sudo systemctl stop host-runner
fi

echo "------------------------------------------------------------"
echo "Updating config in ${RUNNER_HOME} ..."
if [ ! -d ${RUNNER_HOME} ]
then
    mkdir -p ${RUNNER_HOME}
fi
cp host-runner ${RUNNER_HOME}

# Set up service, replacing with variables
echo "------------------------------------------------------------"
echo "Configuring host-runner service ..."
FILE="$SERVICE_HOME"/.env
echo "env file is ..."$FILE
if [ -f "$FILE" ]; then
  set -a # automatically export all variables
  source "$FILE"
  set +a
  export OX_WAPI_URI="http://localhost:8080"
  export OX_WAPI_INSECURE_SKIP_VERIFY="true"
  export OX_WAPI_USER="$OX_SVC_WAPI_ADMIN_USER"
  export OX_WAPI_PWD="$OX_SVC_WAPI_ADMIN_PWD"
fi
art merge host-runner.service.art
sudo chown root:root host-runner.service && sudo chmod 644 host-runner.service
sudo mv host-runner.service /lib/systemd/system/host-runner.service

echo "------------------------------------------------------------"
echo "Restarting service daemon"
sudo systemctl daemon-reload

# Finish
echo "------------------------------------------------------------"
echo "Completed"