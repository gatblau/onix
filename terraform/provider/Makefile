#
#    Onix Config Manager - Terraform Provider - Copyright (c) 2018-2020 by www.gatblau.org
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
PROVIDER_VERSION=1.0.0
PROVIDER_NAME=ox

# the path where the ox terraform providers should be installed locally (terraform >= 0.13)
# see knowledge article here: https://www.terraform.io/upgrade-guides/0-13.html#new-filesystem-layout-for-local-copies-of-providers
TF_PATH=terraform.d/plugins/registry.terraform.io/hashicorp/$(PROVIDER_NAME)/$(PROVIDER_VERSION)

# the name of the terraform provider for Onix
BINARY_NAME=terraform-provider-$(PROVIDER_NAME)_v$(PROVIDER_VERSION)

DARWIN_PATH=$(TF_PATH)/darwin_amd64/$(BINARY_NAME)
LINUX_PATH=$(TF_PATH)/linux_amd64/$(BINARY_NAME)
WINDOWS_PATH=$(TF_PATH)/windows_amd64/$(BINARY_NAME)

# the name of the folder where the packaged binaries will be placed after the build
BUILD_FOLDER=build

# run the terraform provider test harness (acceptance tests - require the Onix Web API)
test:
	# starts Onix in containers (requires Docker)
	sh ox_up.sh
	# enables acceptance tests and executes them
	export TF_ACC=true; go test -v
	# removes the Onix containers
	sh ox_down.sh

# compile, install and initialise the terraform provider for onix
# in a local folder and execute the contents of the main.tf file (terraform apply)
# use this function to test changes to the terraform provider (integration test)
install-darwin:
	go fmt && export GOOS=darwin; go build -o $(DARWIN_PATH) -v && chmod 755 $(DARWIN_PATH) && export TF_LOG=DEBUG
	terraform init

install-linux:
	go fmt && export GOOS=linux; go build -o $(LINUX_PATH) -v && chmod 755 $(LINUX_PATH) && export TF_LOG=DEBUG
	terraform init

install-windows:
	go fmt && export GOOS=windows; go build -o $(WINDOWS_PATH) -v && chmod 755 $(WINDOWS_PATH) && export TF_LOG=DEBUG
	#terraform init

# build the terraform provider in the current platform
build:
	go fmt
	go build -o $(BINARY_NAME) -v

# package the terraform provider for all platforms
package:
	go fmt
	$(MAKE) package_linux
	$(MAKE) package_darwin
	$(MAKE) package_windows

# package the terraform provider for linux amd64 platform
package_linux:
	export GOOS=linux; export GOARCH=amd64; go build -o $(BUILD_FOLDER)/$(BINARY_NAME) -v
	zip -mjT $(BUILD_FOLDER)/$(BINARY_NAME)_linux_amd64.zip $(BUILD_FOLDER)/$(BINARY_NAME)

# package the terraform provider for MacOS
package_darwin:
	export GOOS=darwin; export GOARCH=amd64; go build -o $(BUILD_FOLDER)/$(BINARY_NAME) -v
	zip -mjT $(BUILD_FOLDER)/$(BINARY_NAME)_darwin_amd64.zip $(BUILD_FOLDER)/$(BINARY_NAME)

# package the terraform provider for Windows
package_windows:
	export GOOS=windows; export GOARCH=amd64; go build -o $(BUILD_FOLDER)/$(BINARY_NAME) -v
	zip -mjT $(BUILD_FOLDER)/$(BINARY_NAME)_windows_amd64.zip $(BUILD_FOLDER)/$(BINARY_NAME)
