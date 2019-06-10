GO=go
GO_BUILD=$(GO) build
GO_CLEAN=$(GO) clean
GO_TEST=$(GO) test
GO_GET=$(GO) get

TARGET_BINARY=build/strict_compiler

all: test build

build:
	$(GO_BUILD) -o $(TARGET_BINARY) -v

test:
	$(GO_TEST) -v ./pkg/...

run:
	$(GO_BUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

deps:
	$(GO_GET) github.com/fatih/color