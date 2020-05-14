#
#    Onix DbMan - Copyright (c) 2018-2020 by www.gatblau.org
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
# This dockerfile encapsulates the build process for the Onix Web API
# The builder container is transient and downloads and install maven, package the Java app and extracts the
# Springboot uberjar files to improve startup times
# The release image copy the prepared app files from the builder image

# building stage: compiles dbman
FROM golang as builder
WORKDIR /app
COPY . .
RUN go get . && CGO_ENABLED=0 GOOS=linux go build -o dbman .

# package stage: copy the binary into the deployment image
FROM registry.access.redhat.com/ubi8/ubi-minimal
MAINTAINER gatblau <onix@gatblau.org>
LABEL author="gatblau.org"
WORKDIR /app
COPY --from=builder /app/dbman /app/config.toml ./
RUN rpm -Uvh https://download.postgresql.org/pub/repos/yum/reporpms/EL-8-x86_64/pgdg-redhat-repo-latest.noarch.rpm
RUN microdnf update -y && rm -rf /var/cache/yum && microdnf install postgresql12 && microdnf clean all
USER 20
CMD ["./dbman"]
EXPOSE 8085/tcp
