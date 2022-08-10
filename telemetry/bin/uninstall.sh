#!/bin/bash

# Set standard variables for use below and in template
export TELEMETRY_HOME=${HOME}/.config/telemetry
export TELEMETRY_UID=$(id -u)
export TELEMETRY_GID=$(id -g)

# Stop any existing service
echo ------------------------------------------------------------
echo Making sure host telemetry is not running
sudo systemctl stop telemetry && sudo systemctl disable telemetry

# Removing service
echo ------------------------------------------------------------
echo Removing host-telemetry service ...
sudo rm /lib/systemd/system/telemetry.service
sudo rm -rf ${TELEMETRY_HOME}
echo Restarting service daemon ...
sudo systemctl daemon-reload

# Finish
echo Completed