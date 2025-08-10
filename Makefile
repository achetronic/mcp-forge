
# Image URL to use all building/pushing image targets
IMG ?= ghcr.io/achetronic/mcp-forge:placeholder

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Get the current Go OS
GO_OS ?= $(or $(GOOS),$(shell go env GOOS))
# Get the current Go ARCH
GO_ARCH ?= $(or $(GOARCH),$(shell go env GOARCH))

OS=$(shell uname | tr '[:upper:]' '[:lower:]')

# CONTAINER_TOOL defines the container tool to be used for building images.
# Be aware that the target commands are only tested with Docker which is
# scaffolded by default. However, you might want to replace it to use other
# tools. (i.e. podman)
CONTAINER_TOOL ?= docker

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk command is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

GOLANGCI_LINT = $(shell pwd)/bin/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.54.2
golangci-lint:
	@[ -f $(GOLANGCI_LINT) ] || { \
	set -e ;\
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell dirname $(GOLANGCI_LINT)) $(GOLANGCI_LINT_VERSION) ;\
	}

.PHONY: lint
lint: golangci-lint ## Run golangci-lint linter & yamllint
	$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: golangci-lint ## Run golangci-lint linter and perform fixes
	$(GOLANGCI_LINT) run --fix

##@ Build
.PHONY: swagger
swagger: install-swag ## Build Swagger documents.
	$(SWAG) init --dir "./cmd/,."  --outputTypes "go"

.PHONY: build
build: fmt vet ## Build CLI binary.
	go build -o bin/mcp-forge-$(GO_OS)-$(GO_ARCH) cmd/main.go

.PHONY: run
run: fmt vet ## Run a controller from your host.
	go run ./cmd/ --config ./docs/config-http.yaml

# If you wish to build the manager image targeting other platforms you can use the --platform flag.
# (i.e. docker build --platform linux/arm64). However, you must enable docker buildKit for it.
# More info: https://docs.docker.com/develop/develop-images/build_enhancements/
.PHONY: docker-build
docker-build: ## Build docker image with the manager.
	$(CONTAINER_TOOL) build --no-cache -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	$(CONTAINER_TOOL) push ${IMG}

PACKAGE_NAME ?= package.tar.gz
.PHONY: package
package: ## Package binary.
	@printf "\nCreating package at dist/$(PACKAGE_NAME) \n"
	@mkdir -p dist

	@if [ "$(OS)" = "linux" ]; then \
		tar --transform="s/mcp-forge-$(GO_OS)-$(GO_ARCH)/mcp-forge/" -cvzf dist/$(PACKAGE_NAME) -C bin mcp-forge-$(GO_OS)-$(GO_ARCH) -C ../ LICENSE README.md; \
	elif [ "$(OS)" = "darwin" ]; then \
		tar -cvzf dist/$(PACKAGE_NAME) -s '/mcp-forge-$(GO_OS)-$(GO_ARCH)/mcp-forge/' -C bin mcp-forge-$(GO_OS)-$(GO_ARCH) -C ../ LICENSE README.md; \
	else \
		echo "Unsupported OS: $(GO_OS)"; \
		exit 1; \
	fi

.PHONY: package-signature
package-signature: ## Create a signature for the package.
	@printf "\nCreating package signature at dist/$(PACKAGE_NAME).md5 \n"
	md5sum dist/$(PACKAGE_NAME) | awk '{ print $$1 }' > dist/$(PACKAGE_NAME).md5
