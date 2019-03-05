NAME := ledgerserver
PROJECT_NAME := ledgerserver-go
PKG := "github.com/bosson/$(PROJECT_NAME)"
VERSION := $(shell git describe --always --abbrev=12)
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/ )
PKG_API_FILES1 := $(shell ls pkg/api/*.go | grep -v _test.go | grep -v data.go)
PKG_API_FILES2 := $(shell ls *.go | grep -v _test.go)

DOCKER_NAME := gcr.io/bosson/ledgerserver-go

.PHONY: all dep build clean test coverage coverhtml lint

all: build

version: ## Get the current version
	@echo ${VERSION}

docker_image_name: ## Get the name of the docker image
	@echo ${DOCKER_NAME}:${VERSION}

get_tools: ## Download and install tools
	@go get -u github.com/golang/lint/golint

lint: get_tools ## Lint the files
	@golint -set_exit_status ${PKG_API_FILES1}
	@golint -set_exit_status ${PKG_API_FILES2}

test: generate ## Run unittests
	@go test -short ${PKG_LIST}

race: generate ## Run data race detector
	@go test -race -short ${PKG_LIST}

msan: generate ## Run memory sanitizer
	@go test -msan -short ${PKG_LIST}

dep: ## Get the dependencies
	@go get -v -d ./...

generate: ## Run generators
	@go generate version.go
	@go generate

build: generate ## Build the binary file
	CGO_ENABLED=1 go build -tags netgo --ldflags '-extldflags "-static"' -o "${NAME}" ./cmd/main.go

build_linux_amd64: generate ## Build for linux amd64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -tags netgo --ldflags '-extldflags "-static"' -o "${NAME}-linux-amd64" ./cmd/main.go

build_all: generate build_linux_amd64 ## Build for all supported platforms
	CC=o64-clang CXX=o64-clang++ GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o ${NAME}-darwin-amd64 ./cmd/main.go

# gcr_login: get_tools ## Login to gcr.io
# 	docker-credential-gcr configure-docker

 docker: ## Build a Docker image
	docker build -t $(DOCKER_NAME):$(VERSION) --build-arg PORT=6502 .

# docker_push: # Push the docker image
# 	docker push $(DOCKER_NAME):$(VERSION)

# docker_tag_latest: ## Tag the docker image as latest
# 	docker tag $(DOCKER_NAME):$(VERSION) $(DOCKER_NAME):latest

run: build # Build and run the code
	./ledgerserver

clean: ## Remove previous build
	@sh -c "sed -i \"s/const Version = \\\".*\\\"/const Version = \\\"\\\"/\" version.go"
	@rm -Rf .cache 2>/dev/null || true
	@rm -Rf .keys 2>/dev/null || true
	@rm idp.db 2>/dev/null || true
	@rm pkg/api/bindata.go 2>/dev/null || true
	@rm -f idp-controller 2>/dev/null || true

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
