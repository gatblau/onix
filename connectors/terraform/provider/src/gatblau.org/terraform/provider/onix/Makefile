#
#    Onix CMDB - Copyright (c) 2018-2019 by www.gatblau.org
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

# the name of the terraform provider for Onix
BINARY_NAME=terraform-provider-ox

# the name of the folder where the packaged binaries will be placed after the build
BUILD_FOLDER=build

# compile, install and initialise the terraform provider for onix
# in a local folder and execute the contents of the main.tf file (terraform apply)
# use this function to test changes to the terraform provider (integration test)
tf_test:
	go fmt
	terraform fmt
	go build -o terraform.d/plugins/darwin_amd64/$(BINARY_NAME) -v
	export TF_LOG=DEBUG
	terraform init
	terraform apply

# run the terraform provider test harness (unit tests in go only)
test:
	go fmt
	go test -v

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
	export CGO_ENABLED=0; export GOOS=linux; export GOARCH=amd64; go build -o $(BUILD_FOLDER)/$(BINARY_NAME) -v
	zip -mjT $(BUILD_FOLDER)/$(BINARY_NAME)_linux_amd64.zip $(BUILD_FOLDER)/$(BINARY_NAME)

# package the terraform provider for MacOS
package_darwin:
	export CGO_ENABLED=0; export GOOS=darwin; export GOARCH=amd64; go build -o $(BUILD_FOLDER)/$(BINARY_NAME) -v
	zip -mjT $(BUILD_FOLDER)/$(BINARY_NAME)_darwin_amd64.zip $(BUILD_FOLDER)/$(BINARY_NAME)

# package the terraform provider for Windows
package_windows:
	export CGO_ENABLED=0; export GOOS=windows; export GOARCH=amd64; go build -o $(BUILD_FOLDER)/$(BINARY_NAME) -v
	zip -mjT $(BUILD_FOLDER)/$(BINARY_NAME)_windows_amd64.zip $(BUILD_FOLDER)/$(BINARY_NAME)