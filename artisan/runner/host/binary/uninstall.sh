#!/bin/bash

# Set standard variables for use below and in template
export RUNNER_HOME=${HOME}/.config/host-runner
export RUNNER_UID=$(id -u)
export RUNNER_GID=$(id -g)

# Stop any existing service
echo ------------------------------------------------------------
echo Making sure host runner is not running
sudo systemctl stop host-runner && sudo systemctl disable host-runner

# Removing service
echo ------------------------------------------------------------
echo Removing host-runner service ...
sudo rm /lib/systemd/system/host-runner.service
echo Restarting service daemon ...
sudo systemctl daemon-reload

# Finish
echo Completed
