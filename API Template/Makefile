ENV_FILE := dev.env
PROJECT_NAME := $(shell basename `git rev-parse --show-toplevel`)
PWD := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

include ${ENV_FILE}

# Container Variables
DOCKER_REGISTRY := registry.rancher.centene.com
DOCKER_PROJECT := capture_provider_rtr
DOCKER_TAG := alpha-1
BUILD_IMAGE := ${DOCKER_REGISTRY}/${DOCKER_PROJECT}/${PROJECT_NAME}-build:0.1
API_IMAGE := ${DOCKER_REGISTRY}/${DOCKER_PROJECT}/${PROJECT_NAME}:${DOCKER_TAG}
SWAGGER_IMAGE=quay.io/goswagger/swagger:0.13.0
REDOC_IMAGE=node:8-alpine

all: help

help:
	@echo
	@echo "Select a command to run for ${PROJECT_NAME}:"
	@echo
	@sed  -n 's/^##//p' ${PWD}/Makefile | column -t -s ':' | sed -e 's/^/ /'
	@echo

## deps: Create the build container and pull down necessary dependencies.
deps:
	sudo docker build -t ${BUILD_IMAGE} -f Dockerfile.build .
	sudo docker run --rm -v "${PWD}":/go/src/bitbucket.centene.com/pdsrtr/${PROJECT_NAME} -w /go/src/bitbucket.centene.com/pdsrtr/${PROJECT_NAME} ${BUILD_IMAGE} sh -c 'bash ./Build.sh deps && chmod -R 777 vendor'

## test: Test the Go code
test: deps
	sudo docker run --rm -v "${PWD}":/go/src/bitbucket.centene.com/pdsrtr/${PROJECT_NAME} -w /go/src/bitbucket.centene.com/pdsrtr/${PROJECT_NAME} ${BUILD_IMAGE} sh -c 'go test -v ./...'

## test-coverage: Test the Go code and generate coverage results
test-coverage: deps
	sudo docker run --rm -v "${PWD}":/go/src/bitbucket.centene.com/pdsrtr/${PROJECT_NAME} -w /go/src/bitbucket.centene.com/pdsrtr/${PROJECT_NAME} ${BUILD_IMAGE} sh -c 'go test -v -cover -coverprofile=${PROJECT_NAME}-coverage.out ./...'

## clean: Remove all build artifacts and generated files.
clean:
	go clean -x
	rm -rf dist/
	rm -rf vendor/
	rm Gopkg.lock

## swagger: Generare the swagger documentation.
swagger: deps
	sudo docker run --rm -it -e GOPATH=${PWD}:/go -v "${PWD}":/go/src/bitbucket.centene.com/pdsrtr/${PROJECT_NAME} -w /go/src/bitbucket.centene.com/pdsrtr/${PROJECT_NAME} ${SWAGGER_IMAGE} generate spec -m -o /go/src/bitbucket.centene.com/pdsrtr/${PROJECT_NAME}/public/docs/swagger.json
	sudo docker run --rm -it -v ${PWD}:/${PROJECT_NAME} -w /${PROJECT_NAME} ${REDOC_IMAGE} npx redoc-cli bundle /${PROJECT_NAME}/public/docs/swagger.json -o /${PROJECT_NAME}/public/api/index.html

## server-docker: Build the application container.
server-docker: swagger
	# Because this command has to be executed as root we need to change the permissions on the vendor directory after running the make.
	sudo docker run --rm -v "${PWD}":/go/src/bitbucket.centene.com/pdsrtr/${PROJECT_NAME} -w /go/src/bitbucket.centene.com/pdsrtr/${PROJECT_NAME} ${BUILD_IMAGE} sh -c 'bash ./Build.sh server && chmod -R 777 vendor && chmod -R 777 dist'
	sudo docker build --build-arg PROJECT_NAME=${PROJECT_NAME} -t ${API_IMAGE} .

## start-docker: Starts the container image
start-docker:
	sudo docker run -d -P --env-file ${ENV_FILE} --name ${PROJECT_NAME} ${API_IMAGE}

## stop-docker: Stops the container, and removes the image from Docker
stop-docker:
	sudo docker stop ${PROJECT_NAME}
	sudo docker rm ${PROJECT_NAME}
	sudo docker images -q ${API_IMAGE} | xargs sudo docker rmi

## deploy-docker: Builds and pushes the API image up to the rancher repository. NOTE: Set DOCKER_TAG to change the image tag.
deploy-docker: server-docker
	@echo Deploying Docker Image: ${API_IMAGE}
	sudo docker login ${DOCKER_REGISTRY}
	sudo docker push ${API_IMAGE}

.PHONY:
