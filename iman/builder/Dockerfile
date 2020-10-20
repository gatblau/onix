#
#    Onix Config Manager - Image Manager - Image Builder
#    Copyright (c) 2018-2020 by www.gatblau.org
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

# package stage: copy the binary into the deployment image
FROM docker.io/centos:7
MAINTAINER gatblau <onix@gatblau.org>
LABEL author="gatblau"
ARG UNAME=image-builder
ENV UID=1000
ENV GID=1000
RUN yum update -y && \
    yum install buildah -y && \
    groupadd -g $GID -o $UNAME && \
    useradd -M -u $UID -g $GID $UNAME && \
    rm -rf /var/cache/yum && \
    yum clean all
USER $UNAME
WORKDIR /build
COPY build.sh ./
ENTRYPOINT ["/build/build.sh"]