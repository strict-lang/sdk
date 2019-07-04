GO=go
GO_GET=$(GO) get
GO_TEST=$(GO) test
GO_BUILD=$(GO) build
GO_CLEAN=$(GO) clean
GO_INSTALL=$(GO) install
BINARY_NAME=$(GOPATH)/bin/strict
all: clean build test

install: build #disabled: install-libstrict
	@-$(GO_INSTALL) ./cmd/strict

build:
	@-$(GO_BUILD) -o $(BINARY_NAME) ./cmd/strict

test:
	@-$(GO_TEST) ./compiler/...

run:
	@-$(GO_BUILD) -o $(BINARY_NAME) -v ./cmd/strict
	@-./$(BINARY_NAME)

clean:
	@-$(GO_CLEAN)
	@-rm -rf $(BINARY_NAME)

build-libstrict:
	@-cd libstrict;         \
		cmake CMakeLists.txt; \
		make > /dev/null

install-libstrict:
	@-cd libstrict;         \
		cmake CMakeLists.txt; \
		make install > /dev/null

deps:
	@-$(GO_GET) github.com/fatih/color
	@-$(GO_GET) github.com/urfave/cli
	@-$(GO_GET) github.com/tonnerre/golang-pretty

.PHONY: all