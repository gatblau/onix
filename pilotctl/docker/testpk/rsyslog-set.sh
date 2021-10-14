#!/bin/bash

ENABLED=$(echo $1 | awk '{print tolower($0)}')
CONF=/etc/rsyslog.d/49-onix.conf


if [ "$ENABLED" == "true" ]; then
  sudo cp rsyslog.conf ${CONF}
else
  if [ -f "$CONF" ]; then
    sudo rm ${CONF}
  fi
fi

sudo systemctl restart rsyslog