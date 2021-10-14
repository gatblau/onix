#!bin/bash

sudo apt install nginx -y && \
sudo cp /etc/nginx/nginx.conf /etc/nginx/nginx.conf.backup && \
sudo systemctl enable --now nginx