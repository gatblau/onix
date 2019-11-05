#
#    Sentinel - Copyright (c) 2019 by www.gatblau.org
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
# the name of the container registry repository
REPO_NAME=gatblau

# the name of the Sentinel binary file
BINARY_NAME=sentinel

# the name of the go command to use to build the binary
GO_CMD = go

# the version of the application
APP_VER = v0.0.3

# the name of the folder where the packaged binaries will be placed after the build
BUILD_FOLDER=build

# build the Sentinel binary in the current platform
build:
	$(GO_CMD) fmt
	export GOROOT=/usr/local/go; export GOPATH=$HOME/go; $(GO_CMD) build -o $(BINARY_NAME) -v

# produce a new version tag
set-version:
	sh version.sh $(APP_VER)

# build the Sentinel images as release or snapshot, for both docker.io and quay.io
# make image snapshot=yes fmt=docker|oci reg=docker.io|quay.io
# make image snapshot=no fmt=docker|oci reg=docker.io|quay.io
image:
	sudo bash build.sh $(REPO_NAME) $(BINARY_NAME) $(shell cat ./version) $(snapshot) $(fmt) $(reg)

# push snapshot to specified registry as follows:
#   make snapshot-push registry=docker.io
#   make snapshot-push registry=quay.io
snapshot-push:
	sudo podman login $(registry) -u gatblau
	sudo podman push $(registry)/$(REPO_NAME)/$(BINARY_NAME)-snapshot:$(shell cat ./version)
	sudo podman push $(registry)/$(REPO_NAME)/$(BINARY_NAME)-snapshot:latest

# push release to specified registry as follows:
#   make release-push registry=docker.io
#   make release-push registry=quay.io
release-push:
	sudo podman login $(registry) -u gatblau
	sudo podman push $(registry)/$(REPO_NAME)/$(BINARY_NAME):$(shell cat ./version)
	sudo podman push $(registry)/$(REPO_NAME)/$(BINARY_NAME):latest

# package the Sentinel binary for all platforms
package:
	go fmt
	$(MAKE) package_linux
	$(MAKE) package_darwin
	$(MAKE) package_windows

# package Sentinel for linux amd64 platform
package_linux:
	export GOROOT=/usr/local/go; export GOPATH=$(HOME)/go; export CGO_ENABLED=0; export GOOS=linux; export GOARCH=amd64; $(GO_CMD) build -o $(BUILD_FOLDER)/$(BINARY_NAME) -v
	zip -mjT $(BUILD_FOLDER)/$(BINARY_NAME)_linux_amd64.zip $(BUILD_FOLDER)/$(BINARY_NAME)

# package Sentinel for MacOS
package_darwin:
	export GOROOT=/usr/local/go; export GOPATH=$(HOME)/go; export CGO_ENABLED=0; export GOOS=darwin; export GOARCH=amd64; $(GO_CMD) build -o $(BUILD_FOLDER)/$(BINARY_NAME) -v
	zip -mjT $(BUILD_FOLDER)/$(BINARY_NAME)_darwin_amd64.zip $(BUILD_FOLDER)/$(BINARY_NAME)

# package Sentinel for Windows
package_windows:
	export GOROOT=/usr/local/go; export GOPATH=$(HOME)/go; export CGO_ENABLED=0; export GOOS=windows; export GOARCH=amd64; $(GO_CMD) build -o $(BUILD_FOLDER)/$(BINARY_NAME) -v
	zip -mjT $(BUILD_FOLDER)/$(BINARY_NAME)_windows_amd64.zip $(BUILD_FOLDER)/$(BINARY_NAME)

# creates namespace and roles to run sentinel in openshift
oc-setup:
	cd ./scripts/openshift && sh setup.sh

# imports sentinel template into opneshift
oc-import-template:
	cd ./scripts/openshift && oc create -f sentinel.yml -n openshift

# deletes the sentinel template in opneshift
oc-delete-template:
	oc delete template sentinel -n openshift

# deletes all sentinel resources
oc-cleanup:
	cd ./scripts/openshift && sh cleanup.sh
