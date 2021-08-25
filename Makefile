GO_VERSION=$(shell go version)
UNAME_OS=$(shell go env GOOS)
UNAME_ARCH=$(shell go env GOARCH)

build:
	@CGO_ENABLED=0 GOOS=$(UNAME_OS) GOARCH=$(UNAME_ARCH) go build -ldflags="-w -s" -o /go/bin/inmemcache ./cmd/api

run:
	go run main.go

all: build