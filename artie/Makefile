#
#    Onix Config Manager - Artie
#    Copyright (c) 2018-2020 by www.gatblau.org
#    Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
#    Contributors to this project, hereby assign copyright in this code to the project,
#    to be licensed under the same terms as the rest of the code.
#
# the name of the folder where the packaged binaries will be placed after the build
BUILD_FOLDER=output

# the name of the container registry
REPO_NAME=quay.io/gatblau

# the name of the artie binary file
BINARY_NAME=artie

# the name of the go command to use to build the binary
GO_CMD = go

# the version of the application
APP_VER = v0-0-4

DATE = $(shell date -u '+%d%m%y%H%M%S')
HASH = $(shell git rev-parse --short HEAD)
TAG  = $(APP_VER)-$(DATE)-$(HASH)

APP_NAME = $(BINARY_NAME)-linux-$(TAG)

DANGLING_IMGS = $(shell docker images -f dangling=true -q)

build-linux:
	printf "package docs\nconst Version=\"%s\"" "$(TAG)" >> ./docs/version.go
	go fmt && export GOOS=linux; $(GO_CMD) build -o $(BUILD_FOLDER)/$(BINARY_NAME)-linux-$(TAG) -v

build-darwin:
	go fmt && export GOOS=darwin; $(GO_CMD) build -o $(BUILD_FOLDER)/$(BINARY_NAME)-darwin-$(TAG) -v

build-windows:
	go fmt && export GOOS=windows; $(GO_CMD) build -o $(BUILD_FOLDER)/$(BINARY_NAME)-windows-$(TAG) -v

# generate a swagger spec document from the source code
swaggen:
	swag init -d ./server -g server.go

# runs the swagger editor in a docker container
swaggedit:
	docker run -d -p 80:8080 swaggerapi/swagger-editor

snapshot-image:
	$(MAKE) build-linux
	docker pull registry.access.redhat.com/ubi8/ubi-minimal
	docker build --build-arg APP_NAME=$(APP_NAME) -t $(REPO_NAME)/$(BINARY_NAME)-snapshot:$(TAG) .
	docker tag $(REPO_NAME)/$(BINARY_NAME)-snapshot:$(TAG) $(REPO_NAME)/$(BINARY_NAME)-snapshot:latest

snapshot-push:
	docker push $(REPO_NAME)/$(BINARY_NAME)-snapshot:latest

run-snapshot:
	docker run --name artie -d -p 8082:8082 $(REPO_NAME)/$(BINARY_NAME)-snapshot:latest

clean:
	docker ps -a | grep Exit | cut -d ' ' -f 1 | xargs docker rm
	docker rmi $(DANGLING_IMGS) -f

