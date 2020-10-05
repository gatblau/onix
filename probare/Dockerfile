#
#    Onix Probare - Copyright (c) 2018-2020 by www.gatblau.org
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

# building stage: compiles rpbare
FROM golang as builder
WORKDIR /app
COPY . .
RUN export GOOS=linux; export GOARCH=amd64; go build .

# package stage: copy the binary into the deployment image
FROM registry.access.redhat.com/ubi8/ubi-minimal
LABEL author="gatblau"
LABEL maintainer="onix@gatblau.org"
ARG UNAME=probare
ENV UID=1000
ENV GID=1000
RUN microdnf update --disablerepo=* --enablerepo=ubi-8-appstream --enablerepo=ubi-8-baseos -y && \
    microdnf install shadow-utils.x86_64 && \
    groupadd -g $GID -o $UNAME && \
    useradd -M -u $UID -g $GID $UNAME && \
    rm -rf /var/cache/yum && \
    microdnf clean all
WORKDIR /app
# copy binary
COPY --from=builder /app/probare /app/app.toml /app/secrets.toml /app/
COPY --from=builder /app/static/ /app/static/
USER $UNAME
CMD ["/app/probare"]
EXPOSE 3000/tcp
