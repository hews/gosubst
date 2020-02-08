SRC     := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
TARGET   = gosubst
BUILDS   = ./build
GOFLAGS  = -i -v

COMMIT := $(shell git rev-parse HEAD)
TAG    := $(shell git tag --sort=-v:refname | head -n1)
SPRIG  := $(shell go list -u -m all 2>/dev/null | grep Masterminds/sprig | awk '{ print $$2 }')

LDFLAGS  = -X=main.versionBuild=$(COMMIT) \
           -X=main.versionRelease=$(TAG) \
		   -X=main.versionSprig=$(SPRIG)

PLATFORMS = darwin linux windows

.PHONY: help build build-all
.DEFAULT_GOAL := help

$(BUILDS)/$(TARGET): $(SRC)
	@echo "Building $(TAG)-$(COMMIT) (Sprig: $(SPRIG))..."
	@go build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(BUILDS)/$(TARGET) .

help: ## Print this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: $(BUILDS)/$(TARGET) ## Build the default binary.

build-all: ## Build a series of platform-specific binaries for release.
	@$(foreach GOOS, $(PLATFORMS), \
	   $(shell export GOOS=$(GOOS); go build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(BUILDS)/$(TARGET).$(GOOS) .))
