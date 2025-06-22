# Version information
LAST_COMMIT := $(or $(shell git rev-parse --short HEAD 2> /dev/null),"unknown")
VERSION := $(or $(shell git describe --tags --abbrev=0 2> /dev/null),"v0.0.0")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%S%z")

# Build flags
LD_FLAGS := -s -w \
	-X 'main.version=${VERSION}' \
	-X 'main.commit=${LAST_COMMIT}' \
	-X 'main.date=${BUILD_DATE}'

# Tool configurations
GOPATH ?= $(shell go env GOPATH)
STUFFBIN ?= $(GOPATH)/bin/stuffbin
MOCKERY ?= $(GOPATH)/bin/mockery
PNPM ?= pnpm
GO ?= $(shell which go)

# Build configurations
BIN := bin/explorer451

.PHONY: build deps test dev pack-bin \
        release-dry-run release-snapshot release-tag install-goreleaser \
        fmt lint mocks

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

# Start development server
dev:
	CGO_ENABLED=0 $(GO) run -ldflags="${LD_FLAGS}" cmd/*.go

# Install all dependencies
deps: $(STUFFBIN)
	go mod download
	cd frontend && $(PNPM) install

# Run tests
test:
	go test -v ./...

# ==================================================================================== #
# DOCKER
# ==================================================================================== #

# Start compose services
docker-up:
	docker compose up -d
	@echo "Waiting for server to be ready..."
	@sleep 3

# Stop compose services
docker-down:
	docker compose down

# Clean compose volumes
docker-clean:
	docker compose down -v

# Run init scrips
docker-init:
	@echo "==> Initializing Docker Services..."
	./.docker/localstack/init-aws.sh
	@echo "==> LocalStack successfully initialized"

# Reset compose
docker-reset: docker-clean docker-up docker-init

# ==================================================================================== #
# TESTING & MOCKING
# ==================================================================================== #

# Install mockery
install-mockery:
	@echo "==> Installing mockery..."
	go get github.com/vektra/mockery/v2
	go install github.com/vektra/mockery/v2

# Generate mocks
mocks: install-mockery
	@echo "==> Generating mocks..."
	@$(MOCKERY)
	@echo "==> Mocks generated successfully"

# Run tests with coverage
test-coverage:
	@echo "==> Running tests with coverage..."
	@go test -coverprofile=coverage.txt ./...
	@go tool cover -html=coverage.txt

# Clean test cache and generated mocks
clean-test:
	@echo "==> Cleaning test cache and mocks..."
	@go clean -testcache
	@rm -rf internal/mocks coverage.txt

# ==================================================================================== #
# BUILD
# ==================================================================================== #

# Install required tools
$(STUFFBIN):
	go install github.com/knadh/stuffbin/...

# Build the backend
build:
	CGO_ENABLED=0 $(GO) build -o ${BIN} -ldflags="${LD_FLAGS}" cmd/*.go

# Production build with embedded frontend
pack-bin: $(STUFFBIN) build build-frontend
	$(STUFFBIN) -a stuff -in ${BIN} -out ${BIN} ${STATIC}

# ==================================================================================== #
# RELEASE
# ==================================================================================== #

# Install goreleaser
install-goreleaser:
	go install github.com/goreleaser/goreleaser@latest

# Test the release process without publishing
release-dry-run: install-goreleaser
	goreleaser release --snapshot --clean --skip=publish

# Create a snapshot release for testing
release-snapshot: install-goreleaser
	goreleaser release --snapshot --clean

# Create and push a new release tag
release-tag:
	@if [ -z "$(VERSION)" ]; then echo "VERSION is required. Use: make release-tag VERSION=v1.0.0"; exit 1; fi
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

# Format code
fmt:
	@echo "==> Formatting code..."
	go run mvdan.cc/gofumpt@latest -w .
	go run golang.org/x/tools/cmd/goimports@latest -w -local github.com/inbox451/explorer451 .

# Lint code
lint:
	@echo "==> Linting code..."
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run --fix
