build-linux:
	go fmt && export GOOS=linux; go build .

build-darwin:
	go fmt && export GOOS=darwin; go build .

build-windows:
	go fmt && export GOOS=windows; go build .

# generate a swagger spec document from the source code
swaggen:
	swag init -d ./server -g server.go

# runs the swagger editor in a docker container
swaggedit:
	docker run -d -p 80:8080 swaggerapi/swagger-editor