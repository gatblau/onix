#!/bin/bash

sudo cp /etc/nginx/nginx.conf.backup /etc/nginx/nginx.conf && \
sudo systemctl restart nginx
