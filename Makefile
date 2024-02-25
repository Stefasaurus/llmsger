BINARY_NAME=llmsgr
.DEFAULT_GOAL := build
.PHONY:fmt vet build
fmt:
	go fmt
vet: fmt
	go vet
build: vet
	GOARCH=amd64 GOOS=linux go build llmsger.go
	GOARCH=amd64 GOOS=windows go build llmsger.go
	go build