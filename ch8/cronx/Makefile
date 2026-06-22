# ------------------------------------------------------------------------------
# Makefile for cronx
#
# Notes:
#   - Run `make help` to display available targets.
#   - Common workflow:
#         make build -> make docker-build -> make docker-push -> make deploy
# ------------------------------------------------------------------------------

SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

# ---------------------------
# Project metadata
# ---------------------------
# Binary name and output directory
APP_NAME := cronx
OUTPUT_DIR := bin
IMAGE := cronx/cronx
TAG ?= latest
DEPLOY_FILE := deployments/cronx.yaml

# ---------------------------
# Build environment
# ---------------------------
# GOOS/GOARCH can be overridden by environment variables:
#   make GOOS=linux GOARCH=amd64 build
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# ---------------------------
# Git and build metadata (injected via -ldflags)
# ---------------------------
# These values will be embedded into cronx/pkg/version
VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse HEAD)
TREESTATE := $(shell test -n "`git status --porcelain`" && echo dirty || echo clean)
BUILDTIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -X cronx/pkg/version.GitVersion=$(VERSION) \
           -X cronx/pkg/version.GitCommit=$(COMMIT) \
           -X cronx/pkg/version.GitTreeState=$(TREESTATE) \
           -X cronx/pkg/version.BuildDate=$(BUILDTIME)

# ---------------------------
# Help (auto-discovered from ## comments)
# ---------------------------
.PHONY: help
help: ## Display this help message.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<TARGET>\033[0m\n"} \
	/^[a-zA-Z0-9_\/\.-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } \
	/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) }' $(MAKEFILE_LIST)

##@ Build

.PHONY: build
build: ## Build the cronx binary (honors GOOS/GOARCH).
	@echo "Building $(APP_NAME) for $(GOOS)/$(GOARCH)..."
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build -ldflags "$(LDFLAGS)" -o $(OUTPUT_DIR)/$(APP_NAME) ./cmd/main.go

##@ Docker

.PHONY: docker-build
docker-build: ## Build the Docker image (supports TAG=xxx). 
	@echo "Building Docker image $(IMAGE):$(TAG) (GOOS=$(GOOS), GOARCH=$(GOARCH))"
	docker build \
		-t $(IMAGE):$(TAG) \
		-f ./build/Dockerfile \
		--build-arg TARGETOS=$(GOOS) \
		--build-arg TARGETARCH=$(GOARCH) \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg TREESTATE=$(TREESTATE) \
		--build-arg BUILDTIME=$(BUILDTIME) \
		.

.PHONY: docker-push
docker-push: ## Push the Docker image to the registry.
	@echo "Pushing $(IMAGE):$(TAG)"
	docker push $(IMAGE):$(TAG)

##@ Docker Buildx (Multi-Arch)

# Default platforms for buildx
PLATFORMS ?= linux/amd64,linux/arm64
BUILDER_NAME ?= cronx-builder

.PHONY: docker-buildx
docker-buildx: ## Build and push multi-arch image using docker buildx.
	@echo "Ensuring docker buildx builder '$(BUILDER_NAME)' exists..."
	@docker buildx inspect $(BUILDER_NAME) >/dev/null 2>&1 || \
		docker buildx create --name $(BUILDER_NAME) --use

	@echo "Building and pushing multi-arch image $(IMAGE):$(TAG)"
	docker buildx build \
		--platform $(PLATFORMS) \
		-t $(IMAGE):$(TAG) \
		-f ./build/Dockerfile \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg TREESTATE=$(TREESTATE) \
		--build-arg BUILDTIME=$(BUILDTIME) \
		--push \
		.

##@ Deploy

.PHONY: deploy
deploy: ## Apply Kubernetes manifests to deploy cronx.
	@echo "Deploying using Kubernetes manifests..."
	sed -e "s|image: .*|image: $(IMAGE):$(TAG)|" $(DEPLOY_FILE) | kubectl apply -f -

##@ Utilities

.PHONY: clean
clean: ## Remove all build artifacts.
	@echo "Cleaning output directory..."
	rm -rf $(OUTPUT_DIR)
