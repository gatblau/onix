#!/bin/bash

# creating htpassword for local docker registry
sudo apt install apache2-utils -y

FILE="$SERVICE_HOME"/.env
if [ -f "$FILE" ]; then
  set -a # automatically export all variables
  source "$FILE"
  set +a
fi

#-B  Force bcrypt encryption of the password (very secure)
#-b  Use the password from the command line
#-c  Create a new file
htpasswd -Bbc "$SERVICE_HOME"/registryv2/auth/registry.password "$REGISTRYV2_SVC_REG_USER" "$REGISTRYV2_SVC_REG_PASSWORD"