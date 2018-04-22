#!/usr/bin/env bash
echo 'installing the yum repository for postgres 10...'
yum install https://download.postgresql.org/pub/repos/yum/10/redhat/rhel-7-x86_64/pgdg-centos10-10-2.noarch.rpm

echo 'installing the client packages...'
yum install -y postgresql10

echo 'installing the server packages...'
yum install -y postgresql10-server

echo 'enabling the server on system start up...'
systemctl enable postgresql-10

echo 'starting the server...'
sudo systemctl start postgresql-10

echo 'the server status is:'
systemctl status postgresql-10