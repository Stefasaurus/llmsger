BINARY_NAME=llmsgr
.DEFAULT_GOAL := build
.PHONY:fmt vet build
fmt:
	go fmt ../
vet: fmt
	go vet ../
build: vet
	GOARCH=amd64 GOOS=linux go build -ldflags "-X go.szostok.io/version.version=$(shell git describe --tags)" -o linux/amd64/ ../llmsger.go
	GOARCH=386 GOOS=linux go build -ldflags "-X go.szostok.io/version.version=$(shell git describe --tags)" -o linux/x86/ ../llmsger.go
	GOARCH=amd64 GOOS=windows go build -ldflags "-X go.szostok.io/version.version=$(shell git describe --tags)" -o windows/amd64/ ../llmsger.go
	GOARCH=386 GOOS=windows go build -ldflags "-X go.szostok.io/version.version=$(shell git describe --tags)" -o windows/x86/ ../llmsger.go
	
