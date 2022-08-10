#!/bin/bash

WD=/tmp/telemetry-deploy
mkdir -p $WD

# Set standard variables for use below and in template
export TELEMETRY_HOME=${HOME}/.config/telemetry
export TELEMETRY_UID=$(id -u)
export TELEMETRY_GID=$(id -g)

echo "------------------------------------------------------------"
echo "Stopping the service"

# Stop existing current service
if [ $(systemctl is-active telemetry) == "active" ]
then
    sudo systemctl stop telemetry
fi

echo "------------------------------------------------------------"
echo "Updating config in ${TELEMETRY_HOME} ..."
if [ ! -d ${TELEMETRY_HOME} ]
then
    mkdir -p ${TELEMETRY_HOME}
fi
cp telemetry ${TELEMETRY_HOME}
cp telem.yaml ${TELEMETRY_HOME}

# Set up service, replacing with variables
echo "------------------------------------------------------------"
echo "Configuring telemetry service ..."
art merge telemetry.service.art
sudo chown root:root telemetry.service && sudo chmod 644 telemetry.service
sudo mv telemetry.service /lib/systemd/system/telemetry.service

echo "------------------------------------------------------------"
echo "Restarting service daemon"
sudo systemctl daemon-reload

# Finish
echo "------------------------------------------------------------"
echo "Completed"