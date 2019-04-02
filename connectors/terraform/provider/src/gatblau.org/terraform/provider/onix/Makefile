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

BINARY_NAME=terraform-provider-ox
BUILD_FOLDER=build

tf_test:
	go fmt
	terraform fmt
	go build -o terraform.d/plugins/darwin_amd64/$(BINARY_NAME) -v
	export TF_LOG=DEBUG
	terraform init
	terraform apply

test:
	go fmt
	go test -v

build:
	go fmt
	go build -o $(BINARY_NAME) -v

package:
	$(MAKE) package_linux
	$(MAKE) package_darwin
	$(MAKE) package_windows

package_linux:
	CGO_ENABLED=0
	GOOS=linux
	GOARCH=amd64
	go fmt
	go build -o $(BUILD_FOLDER)/$(BINARY_NAME) -v
	zip -mjT $(BUILD_FOLDER)/$(BINARY_NAME)_linux_amd64.zip $(BUILD_FOLDER)/$(BINARY_NAME)

package_darwin:
	CGO_ENABLED=0
	GOOS=darwin
	GOARCH=amd64
	go fmt
	go build -o $(BUILD_FOLDER)/$(BINARY_NAME) -v
	zip -mjT $(BUILD_FOLDER)/$(BINARY_NAME)_darwin_amd64.zip $(BUILD_FOLDER)/$(BINARY_NAME)

package_windows:
	$CGO_ENABLED=0
	GOOS=windows
	GOARCH=amd64
	go fmt
	go build -o $(BUILD_FOLDER)/$(BINARY_NAME) -v
	zip -mjT $(BUILD_FOLDER)/$(BINARY_NAME)_windows_amd64.zip $(BUILD_FOLDER)/$(BINARY_NAME)