CURRENT_WORKING_DIR=$(shell pwd)

#------------------------------------------------------------------
# Project build information
#------------------------------------------------------------------
PROJNAME          		:= perimener
VENDOR            		:= swade1987
MAINTAINER        		:= steven@stevenwade.co.uk

GIT_REPO          		:= github.com/$(VENDOR)/$(PROJNAME)
GIT_SHA           		:= $(shell git rev-parse --verify HEAD)
BUILD_DATE        		:= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

DOCKER_USERNAME     	:= "swade1987"
DOCKER_PASSWORD     	?="unknown"

# Construct docker image name.
IMAGE             		:= $(PROJNAME):$(VERSION)

#------------------------------------------------------------------
# Go configuration
#------------------------------------------------------------------
GOCMD             		:= go
GOFMT             		:= gofmt
BIN               		:= bin
VERSION_PKG       		:= $(GIT_REPO)/pkg/runtime/version

#------------------------------------------------------------------
# Build targets
#------------------------------------------------------------------

.PHONY: fmt
fmt: ## Run go fmt against code
	$(GOFMT) -w main.go
	$(GOCMD) fmt ./pkg/...

.PHONY: vet
vet: ## Run go vet against code
	$(GOCMD) vet main.go
	$(GOCMD) vet ./pkg/...

.PHONY: test
test: fmt vet ## Run tests
	$(GOCMD) test ./pkg/... -coverprofile cover.out

.PHONY: perimener
perimener: fmt vet
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOCMD) build -o $(BIN)/linux/$(PROJNAME) $(GIT_REPO)
	env GOOS=darwin GOARCH=amd64 $(GOCMD) build -o $(BIN)/darwin/$(PROJNAME) $(GIT_REPO)

.PHONY: perimener-linux
perimener-linux: fmt vet
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOCMD) build -o $(BIN)/linux/$(PROJNAME) $(GIT_REPO)

#------------------------------------------------------------------
# CI targets
#------------------------------------------------------------------

.PHONY: build
build:
	docker build \
    --build-arg git_repository=`git config --get remote.origin.url` \
    --build-arg git_branch=`git rev-parse --abbrev-ref HEAD` \
    --build-arg git_commit=`git rev-parse HEAD` \
    --build-arg built_on=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
    -t $(IMAGE) .

.PHONY: login
login:
	docker login -u $(DOCKER_USERNAME) -p $(DOCKER_PASSWORD)
