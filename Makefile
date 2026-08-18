BIN ?= bin/aegis
IMAGE ?= aegis
GO ?= go
GOOS ?=
GOARCH ?=
PLATFORMS ?= linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

MODULE := github.com/tohjustin/aegis
VERSION_PKG := $(MODULE)/internal/version
VERSION ?= $(shell cat VERSION.txt 2>/dev/null || echo dev)
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
GIT_TAG ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo unknown)
GIT_TREE_STATE ?= $(if $(shell git status --porcelain --untracked-files=no 2>/dev/null),dirty,clean)

LDFLAGS := -s -w \
	-X $(VERSION_PKG).Version=$(VERSION) \
	-X $(VERSION_PKG).BuildDate=$(BUILD_DATE) \
	-X $(VERSION_PKG).GitCommit=$(GIT_COMMIT) \
	-X $(VERSION_PKG).GitTag=$(GIT_TAG) \
	-X $(VERSION_PKG).GitTreeState=$(GIT_TREE_STATE)

.PHONY: check-deps build build-all run vet lint tidy test test-e2e test-all docker-build docker-run help

check-deps: ## Verify required tools are on PATH
	@command -v $(GO) >/dev/null || { echo "go is required on PATH"; exit 1; }
	@command -v golangci-lint >/dev/null || { echo "golangci-lint is required on PATH"; exit 1; }

build: ## Build the service binary
	@mkdir -p $(dir $(BIN))
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN) ./cmd/aegis

build-all: ## Cross-compile for linux/darwin amd64/arm64
	@mkdir -p bin
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; arch=$${platform#*/}; \
		echo "building $$os/$$arch"; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch $(GO) build -trimpath -ldflags "$(LDFLAGS)" -o bin/aegis-$$os-$$arch ./cmd/aegis; \
	done

run: build ## Build and run the service
	$(BIN)

vet: ## Run go vet
	$(GO) vet ./...

lint: ## Run golangci-lint
	@command -v golangci-lint >/dev/null || { echo "golangci-lint is required on PATH"; exit 1; }
	golangci-lint run

tidy: ## Run go mod tidy
	$(GO) mod tidy

test: ## Run unit tests
	$(GO) test ./...

test-e2e: ## Run blackbox e2e tests against a real binary
	$(GO) test -tags e2e -count=1 ./test/e2e/...

test-all: test test-e2e ## Run unit and e2e tests

docker-build: ## Build the container image
	@command -v docker >/dev/null || { echo "docker is required on PATH"; exit 1; }
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg GIT_TAG=$(GIT_TAG) \
		--build-arg GIT_TREE_STATE=$(GIT_TREE_STATE) \
		-t $(IMAGE) .

docker-run: docker-build ## Run the container image
	docker run --rm -p 8080:8080 -p 9090:9090 \
		-e GITHUB_ACCESS_TOKEN=$(GITHUB_ACCESS_TOKEN) \
		$(IMAGE)

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-18s\033[0m %s\n", $$1, $$2}'
