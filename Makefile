# SPDX-License-Identifier: Apache-2.0
# SPDX-FileCopyrightText: 2021-Present The Zarf Authors

# Provide a default value for the operating system architecture used in tests, e.g. " APPLIANCE_MODE=true|false make test-e2e ARCH=arm64"
ARCH ?= amd64
KEY ?= ""
######################################################################################

# Figure out which Zarf binary we should use based on the operating system we are on
ZARF_UI_BIN := ./build/zarf-ui
ifeq ($(OS),Windows_NT)
	ZARF_UI_BIN := $(addsuffix .exe,$(ZARF_UI_BIN))
else
	UNAME_S := $(shell uname -s)
	UNAME_P := $(shell uname -p)
	ifneq ($(UNAME_S),Linux)
		ifeq ($(UNAME_S),Darwin)
			ZARF_UI_BIN := $(addsuffix -mac,$(ZARF_UI_BIN))
		endif
		ifeq ($(UNAME_P),i386)
			ZARF_UI_BIN := $(addsuffix -intel,$(ZARF_UI_BIN))
		endif
		ifeq ($(UNAME_P),arm)
			ZARF_UI_BIN := $(addsuffix -apple,$(ZARF_UI_BIN))
		endif
	endif
endif

UI_VERSION ?= $(if $(shell git describe --tags),$(shell git describe --tags),"UnknownVersion")
CLI_VERSION ?= $(shell cat go.mod | grep "github.com/defenseunicorns/zarf " | cut -d " " -f 2)
GIT_SHA := $(if $(shell git rev-parse HEAD),$(shell git rev-parse HEAD),"")
BUILD_DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
BUILD_ARGS := -s -w -X 'github.com/defenseunicorns/zarf-ui/src/config.UIVersion=$(UI_VERSION)' \
              -X 'github.com/defenseunicorns/zarf/src/config.CLIVersion=$(CLI_VERSION)' \
			  -X 'github.com/defenseunicorns/zarf/src/config.ActionsCommandZarfPrefix=zarf' \
			  -X 'k8s.io/component-base/version.gitVersion=v0.0.0+zarf$(CLI_VERSION)' \
			  -X 'k8s.io/component-base/version.gitCommit=$(GIT_SHA)' -X 'k8s.io/component-base/version.buildDate=$(BUILD_DATE)'
.DEFAULT_GOAL := help

.PHONY: help
help: ## Display this help information
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	  | sort | awk 'BEGIN {FS = ":.*?## "}; \
	  {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

clean: ## Clean the build directory
	rm -rf build

delete-packages: ## Delete all Zarf package tarballs in the project recursively
	find . -type f -name 'zarf-package-*' -delete

# INTERNAL: used to ensure the ui directory exists
ensure-ui-build-dir:
	mkdir -p build/ui
	touch build/ui/index.html

# INTERNAL: used to build the UI only if necessary
check-ui:
	@ if [ ! -z "$(shell command -v shasum)" ]; then\
		if test "$(shell ./hack/print-ui-diff.sh | shasum)" != "$(shell cat build/ui/git-info.txt | shasum)" ; then\
			$(MAKE) build-ui-frontend;\
			./hack/print-ui-diff.sh > build/ui/git-info.txt;\
		fi;\
	else\
		$(MAKE) build-ui-frontend;\
	fi

build-ui-frontend: ## Build the Zarf UI Frontend
	npm --prefix src/ui ci
	npm --prefix src/ui run build

build-ui-linux-amd: check-ui ## Build the Zarf UI for Linux on AMD64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(BUILD_ARGS)" -o build/zarf-ui .

build-ui-linux-arm: check-ui ## Build the Zarf UI for Linux on ARM
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="$(BUILD_ARGS)" -o build/zarf-ui-arm .

build-ui-mac-intel: check-ui ## Build the Zarf UI for macOS on AMD64
	GOOS=darwin GOARCH=amd64 go build -ldflags="$(BUILD_ARGS)" -o build/zarf-ui-mac-intel .

build-ui-mac-apple: check-ui ## Build the Zarf UI for macOS on ARM
	GOOS=darwin GOARCH=arm64 go build -ldflags="$(BUILD_ARGS)" -o build/zarf-ui-mac-apple .

build-ui-windows-amd: check-ui ## Build the Zarf UI for Windows on AMD64
	GOOS=windows GOARCH=amd64 go build -ldflags="$(BUILD_ARGS)" -o build/zarf-ui.exe .

build-ui-windows-arm: check-ui ## Build the Zarf UI for Windows on ARM
	GOOS=windows GOARCH=arm64 go build -ldflags="$(BUILD_ARGS)" -o build/zarf-ui-arm.exe .

build-ui-linux: build-ui-linux-amd build-ui-linux-arm ## Build the Zarf UI for Linux on AMD64 and ARM

build-ui: build-ui-linux-amd build-ui-linux-arm build-ui-mac-intel build-ui-mac-apple build-ui-windows-amd build-ui-windows-arm ## Build the UI

dev: ensure-ui-build-dir ## Start a Dev Server for the Zarf UI
	go mod download
	npm --prefix src/ui ci
	npm --prefix src/ui run dev

.PHONY: test-ui
test-ui: ## Run the Zarf UI E2E tests (requires `make build-ui` first) (run with env CI=true to use build/zarf)
	export NODE_PATH=$(CURDIR)/src/ui/node_modules && \
	npm --prefix src/ui run test:pre-init && \
	npm --prefix src/ui run test:init && \
	npm --prefix src/ui run test:post-init && \
	npm --prefix src/ui run test:connect

.PHONY: test-ui-dev-server
# INTERNAL: used to start a dev version of the API server for the Zarf Web UI tests (locally)
test-ui-dev-server:
	API_DEV_PORT=5173 \
		API_PORT=3333 \
		API_TOKEN=insecure \
		go run -ldflags="$(BUILD_ARGS)" main.go -l=trace

.PHONY: test-ui-build-server
# INTERNAL: used to start the built version of the API server for the Zarf Web UI (in CI)
test-ui-build-server:
	API_PORT=3333 API_TOKEN=insecure $(ZARF_UI_BIN)

# INTERNAL: used to test for new CVEs that may have been introduced
test-cves: ensure-ui-build-dir
	go run main.go zarf tools sbom packages . -o json | grype --fail-on low

cve-report: ensure-ui-build-dir ## Create a CVE report for the current project (must `brew install grype` first)
	go run main.go zarf tools sbom packages . -o json | grype -o template -t hack/.templates/grype.tmpl > build/zarf-known-cves.csv

lint-go: ## Run revive to lint the go code (must `brew install revive` first)
	revive -config revive.toml -exclude src/cmd/viper.go -formatter stylish ./src/...

api-schema: ensure-ui-build-dir ## Generate the Zarf UI API Schema
	ZARF_CONFIG=hack/empty-config.toml hack/create-api-schema.sh

# INTERNAL: used to test that a dev has ran `make api-schema` in their PR
test-api-schema:
	$(MAKE) api-schema
	hack/check-api-schema.sh

retrieve-packages: ensure-ui-build-dir ## Retrieve published test packages (can also be built from the Zarf repo)
	@test -s $(ZARF_UI_BIN) || $(MAKE) build-ui
	@test -s ./build/zarf-init-$(ARCH)-$(CLI_VERSION).tar.zst || $(ZARF_UI_BIN) zarf tools download-init -a $(ARCH) -o build
	@test -s ./build/zarf-package-dos-games-$(ARCH)-1.0.0.tar.zst || $(ZARF_UI_BIN) zarf package pull oci://ghcr.io/defenseunicorns/packages/dos-games:1.0.0-$(ARCH) -o build --key=https://zarf.dev/cosign.pub
