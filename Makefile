# Image URL to use all building/pushing image targets
IMG_REGISTRY ?= zdzserver/webhook

# Build binary
BUILD_ROOT = "./bin"

GOOS ?= linux
GOARCH ?= $(shell go env GOARCH)



# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# go option
PKG        := ./...
TAGS       :=
TESTS      := .
TESTFLAGS  :=
LDFLAGS    := -w -s
GOFLAGS    :=

# Rebuild the binary if any of these files change
#SRC := $(shell find . -type f -name '*.go' -print) go.mod go.sum
SRC := go.mod go.sum

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec


# Git information
GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_COMMIT_HASH    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_TREESTATE  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")
BUILDDATE = $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

# Set your version by env or using latest tags from git
VERSION ?= ""
ifeq ($(VERSION), "")
    LATEST_TAG=$(GIT_TAG)
    ifeq ($(LATEST_TAG),)
        # Forked repo may not sync tags from upstream, so give it a default tag to make CI happy.
        VERSION="unknown"
    else
        VERSION=$(LATEST_TAG)
    endif
endif

LDFLAGS +=


.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
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

##@ Build


.PHONY: build
build: fmt vet webhook ## Build binary.

.PHONY: clean
clean: ## Clean binary.
	rm -rf  $(BUILD_ROOT)/*

.PHONY: webhook
webhook: $(SRC)
	CGO_ENABLED=0 GOOS=$(GOOS) go build \
		-ldflags "$(LDFLAGS)" \
		-o $(BUILD_ROOT)/webhook \
		main.go

.PHONY: image
image: build ## build image
	@echo "build image to $(IMG_REGISTRY):$(VERSION)"
	docker build -t $(IMG_REGISTRY):$(VERSION) .


.PHONY: push
push: docker ## Push docker image with the manager.
	@echo "push image to $(IMG_REGISTRY):$(VERSION)"
	docker push $(IMG_REGISTRY):$(VERSION)
