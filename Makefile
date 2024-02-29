BINARY_NAME=llmsgr
.DEFAULT_GOAL := build
.PHONY:fmt vet build
fmt:
	go fmt
vet: fmt
	go vet
build: vet
	go build -ldflags "-X go.szostok.io/version.version=$(shell git describe --tags)" llmsger.go
