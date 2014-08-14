SHELL := /bin/bash
REV_VAR := github.com/modcloth-labs/amqp-tee.RevString
VERSION_VAR := github.com/modcloth-labs/amqp-tee.VersionString
REPO_VERSION := $(shell git describe --always --dirty --tags)
REPO_REV := $(shell git rev-parse --sq HEAD)
GOBUILD_VERSION_ARGS := -ldflags "-X $(REV_VAR) $(REPO_REV) -X $(VERSION_VAR) $(REPO_VERSION)"

.PHONY: all
all: build test

.PHONY: clean
clean:
	go clean -x ./...

.PHONY: build
build: deps
	go install $(GOBUILD_VERSION_ARGS) -x ./...

.PHONY: savedeps
savedeps:
	godep save -copy=false ./...

.PHONY: deps
deps:
	godep restore

.PHONY: test
test: deps
	go test ./...

.PHONY: fmtpolice
fmtpolice: fmt lint

.PHONY: fmt
fmt:
	@set -e ; \
	  for f in $(shell git ls-files '*.go'); do \
	  gofmt $$f | diff -u $$f - ; \
	  done

.PHONY: lint
lint:
	@for file in $(shell git ls-files '*.go') ; do \
	  if [[ ! "$$(golint $$file)" =~ ^[[:blank:]]*$$ ]] ; then \
	  $(MAKE) lintv && exit 1 ; fi \
	  done

.PHONY: lintv
lintv:
	@for file in $(shell git ls-files '*.go') ; do $(GOPATH)/bin/golint $$file ; done
