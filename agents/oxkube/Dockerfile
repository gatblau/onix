#
#   Onix Kube - Copyright (c) 2019 by www.gatblau.org
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
#   Unless required by applicable law or agreed to in writing, software distributed under
#   the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
#   either express or implied.
#   See the License for the specific language governing permissions and limitations under the License.
#
#   Contributors to this project, hereby assign copyright in this code to the project,
#   to be licensed under the same terms as the rest of the code.
#
# multi-stage docker build: compile -> package
#
# compile stage: build Sentinel
FROM golang as builder
WORKDIR /app
COPY . .
RUN go get . && CGO_ENABLED=0 GOOS=linux go build -o oxkube .

# package stage: copy the binary into the deployment image
FROM registry.access.redhat.com/ubi8/ubi-minimal
MAINTAINER gatblau <admin@gatblau.org>
LABEL author="gatblau.org"
WORKDIR /app
COPY --from=builder /app/oxkube /app/config.toml ./
USER 20
CMD ["./oxkube"]
