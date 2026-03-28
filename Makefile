# Copyright (c) The Kowabunga Project
# Apache License, Version 2.0 (see LICENSE or https://www.apache.org/licenses/LICENSE-2.0.txt)
# SPDX-License-Identifier: Apache-2.0

#export GOOS=linux
#export GOARCH=amd64

# Make sure GOPATH is NOT set to this folder or we'll get an error "$GOPATH/go.mod exists but should not"
#export GOPATH = ""
export GO111MODULE = on
BINDIR = bin
PLUGINS_DIR = plugins
PLUGINS_KAKTUS_PKG_DIR = ./internal/kaktus/plugins

GOLINT = $(BINDIR)/golangci-lint
GOVULNCHECK = $(BINDIR)/govulncheck
GOSEC = $(BINDIR)/gosec

PKGS = $(shell go list ./...)
PKGS_SHORT = $(shell go list ./... | sed 's%github.com/kowabunga-cloud/kaktus/%%')

V = 0
Q = $(if $(filter 1,$V),,@)
PROD = 0
ifeq ($(PROD),1)
DEBUG = -w -s
endif
M = $(shell printf "\033[34;1m▶\033[0m")

ifeq ($(V), 1)
  OUT = ""
else
  OUT = ">/dev/null"
endif

# This is our default target
# it does not build/run the tests
.PHONY: all
all: mod fmt vet lint build ; @ ## Do all
	$Q echo "done"

# This target grabs the necessary go modules
.PHONY: mod
mod: ; $(info $(M) collecting modules…) @
	$Q go mod download
	$Q go mod tidy

# Updates all go modules
update: ; $(info $(M) updating modules…) @
	$Q go get -u ./...
	$Q go mod tidy

# Makes sure bin directory is created
.PHONY: bin
bin: ; $(info $(M) create local bin directory) @
	$Q mkdir -p $(BINDIR)

.PHONY: kaktus
kaktus: ; $(info $(M) building Kaktus agent…) @
	$Q go build \
		-gcflags="internal/...=-e" \
		-ldflags='$(DEBUG)' \
		-o $(BINDIR) ./cmd/kaktus

# Makes sure plugins directory is created
.PHONY: plugins
plugins: ; $(info $(M) create local plugins directory) @
	$Q mkdir -p $(PLUGINS_DIR)

.PHONY: plugin-ceph
plugin-ceph: plugins ; $(info $(M) building Kaktus Ceph plugin…) @
	$Q go build -buildmode=plugin \
		-gcflags="internal/...=-e" \
		-ldflags='$(DEBUG)' \
		-o $(PLUGINS_DIR) $(PLUGINS_KAKTUS_PKG_DIR)/ceph

.PHONY: build
build: kaktus plugin-ceph

.PHONY: tests
tests: ; $(info $(M) test suite…) @
	$Q go test ./... -count=1 -coverprofile=coverage.txt

.PHONY: deb
deb: ; $(info $(M) building Debian package…) @
	$Q ./debian.sh

.PHONY: get-lint
get-lint: ; $(info $(M) downloading go-lint…) @
	$Q test -x $(GOLINT) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s

.PHONY: lint
lint: get-lint ; $(info $(M) running go-lint…) @
	$Q $(GOLINT) run ./... ; exit 0

.PHONY: get-govulncheck
get-govulncheck: ; $(info $(M) downloading govulncheck…) @
	$Q test -x $(GOVULNCHECK) || GOBIN="$(PWD)/$(BINDIR)/" go install golang.org/x/vuln/cmd/govulncheck@latest

.PHONY: vuln
vuln: get-govulncheck ; $(info $(M) running govulncheck…) @ ## Check for known vulnerabilities
	$Q $(GOVULNCHECK) ./... ; exit 0

.PHONY: get-gosec
get-gosec: ; $(info $(M) downloading gosec…) @
	$Q test -x $(GOSEC) || GOBIN="$(PWD)/$(BINDIR)/" go install github.com/securego/gosec/v2/cmd/gosec@latest

.PHONY: sec
sec: get-gosec ; $(info $(M) running gosec…) @ ## AST / SSA code checks
	$Q $(GOSEC) -terse ./... ; exit 0

.PHONY: vet
vet: ; $(info $(M) running go vet…) @
	$Q go vet $(PKGS) ; exit 0

.PHONY: fmt
fmt: ; $(info $(M) running go fmt…) @
	$Q gofmt -w -s $(PKGS_SHORT)

.PHONY: clean
clean: ; $(info $(M) cleaning…)	@ ## Cleanup everything
	$Q rm -rf $(BINDIR)
	$Q rm -rf $(PLUGINS_DIR)

# This target count all the lines of .go files (no matter if empty lines or comments)
.PHONY: lc
lc: ; @ ## Count lines
	@find . -name "*.go" -exec cat {} \; | wc -l | awk '{print $$1}'

# This target count the lines of go code only (ignore empty lines, comments, etc.)
# it requires gosloc
.PHONY: sloc
sloc: ; @ ## Count GO lines
	@find . -name "*.go" -exec cat {} \; | gosloc
