#!/usr/bin/env bash
echo 'installing the yum repository for postgres 10...'
rpm -Uvh https://download.postgresql.org/pub/repos/yum/10/redhat/rhel-7-x86_64/pgdg-centos10-10-1.noarch.rpm

echo 'installing the db server...'
yum install -y postgresql10-server postgresql10

echo 'enabling the db server on system start up...'
systemctl enable postgresql-10

echo 'starting the db server...'
sudo systemctl start postgresql-10

echo 'the server status is:'
systemctl status postgresql-10