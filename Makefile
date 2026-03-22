.PHONY: build test install clean plugin

PLUGIN_NAME = loglint
PLUGIN_PATH = $(shell go env GOPATH)/bin/$(PLUGIN_NAME).so

plugin:
	go build -buildmode=plugin -o $(PLUGIN_PATH) ./cmd/plugin

build:
	go build ./...

test:
	go test -v ./...

install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	make plugin

clean:
	go clean
	rm -f $(PLUGIN_PATH)

lint: build test plugin
	@echo "All done!"