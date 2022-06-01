#!/bin/bash

#creating home directory for service setup
sudo rm -rf "$SERVICE_HOME"
sudo mkdir -p "$SERVICE_HOME"/registryv2/auth
sudo mkdir -p "$SERVICE_HOME" && sudo chown -R $(id -u):$(id -g) "$SERVICE_HOME"