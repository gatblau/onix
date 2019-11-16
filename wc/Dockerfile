#!/usr/bin/env bash
#
#    Onix Config Manager - Copyright (c) 2018-2019 by www.gatblau.org
#
#    Licensed under the Apache License, Version 2.0 (the "License");
#    you may not use this file except in compliance with the License.
#    You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
#    Unless required by applicable law or agreed to in writing, software distributed under
#    the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
#    either express or implied.
#    See the License for the specific language governing permissions and limitations under the License.
#
#    Contributors to this project, hereby assign copyright in this code to the project,
#    to be licensed under the same terms as the rest of the code.
#

# Builds the Onix Web Console image6
FROM quay.io/gatblau/node:12-ubi8-min
MAINTAINER Gatblau <onix@gatblau.org>
WORKDIR /usr/src/wc
COPY . .
RUN npm install && npm run build
EXPOSE 3000
USER 20
CMD [ "npm", "start" ]
