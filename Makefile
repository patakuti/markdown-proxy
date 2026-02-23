BINARY_NAME := markdown-proxy
GO := go

.PHONY: all build clean linux windows test

all: build

build:
	$(GO) build -o $(BINARY_NAME) ./cmd/markdown-proxy

linux:
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BINARY_NAME) ./cmd/markdown-proxy

windows:
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BINARY_NAME).exe ./cmd/markdown-proxy

test:
	$(GO) test ./...

clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe
