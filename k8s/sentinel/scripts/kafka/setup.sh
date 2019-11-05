#!/usr/bin/env bash
#
# Sentinel - Copyright (c) 2019 by www.gatblau.org
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software distributed under
# the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
# either express or implied.
# See the License for the specific language governing permissions and limitations under the License.
#
# Contributors to this project, hereby assign copyright in this code to the project,
# to be licensed under the same terms as the rest of the code.
#
# setup a kafka broker in a vanilla CentOS 7 vagrant box for demo purposes
# NOTE: must be run on the server to avoid dns resolution issues with kafka vs zookeeper
# ALSO: do not call this script directly as it is meant to be called by go.sh

echo "installing prerequisite packages"
sudo yum -y install wget java

# downloads apache kafka
wget https://www-eu.apache.org/dist/kafka/2.2.1/kafka_2.12-2.2.1.tgz

# decompress kafka
tar -xvzf kafka_2.12-2.2.1.tgz

echo "starting zookeeper server"
sh ./kafka_2.12-2.2.1/bin/zookeeper-server-start.sh -daemon ./kafka_2.12-2.2.1/config/zookeeper.properties
sleep 7

echo "starting kafka server"
sh ./kafka_2.12-2.2.1/bin/kafka-server-start.sh -daemon ./kafka_2.12-2.2.1/config/server.properties
sleep 7

echo "creating k8s topic"
sh ./kafka_2.12-2.2.1/bin/kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic k8s