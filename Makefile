LIBRARIES := github.com/modcloth-labs/amqp-tee
TARGETS := $(LIBRARIES) github.com/modcloth-labs/amqp-tee/amqp-tee
REV_VAR := github.com/modcloth-labs/amqp-tee.RevString
VERSION_VAR := github.com/modcloth-labs/amqp-tee.VersionString
REPO_VERSION := $(shell git describe --always --dirty --tags)
REPO_REV := $(shell git rev-parse --sq HEAD)
GOBUILD_VERSION_ARGS := -ldflags "-X $(REV_VAR) $(REPO_REV) -X $(VERSION_VAR) $(REPO_VERSION)"
JOHNNY_DEPS_VERSION := v0.2.3
DOCKER ?= sudo docker

all: build test

build: deps
	go get -x -n $(TARGETS)
	go install $(GOBUILD_VERSION_ARGS) -x $(TARGETS)

deps: johnny_deps
	./johnny_deps

johnny_deps:
	curl -L -s -o $@ https://raw.github.com/VividCortex/johnny-deps/$(JOHNNY_DEPS_VERSION)/bin/johnny_deps
	chmod +x $@

test:
	go test -i $(LIBRARIES)
	go test -x -v $(LIBRARIES)

clean:
	go clean -x $(LIBS) || true
	if [ -d $${GOPATH%%:*}/pkg ] ; then \
	  find $${GOPATH%%:*}/pkg -name '*amqp-tee*' -exec rm -v {} \; ; \
	fi

distclean: clean
	rm -f ./johnny_deps

container: build
	mkdir -p .build
	cp $${GOPATH%%:*}/bin/amqp-tee .build
	$(DOCKER) build -t quay.io/modcloth/amqp-tee:$(REPO_VERSION) .

.PHONY: all build test container
