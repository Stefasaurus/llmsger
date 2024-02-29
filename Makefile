BINARY_NAME=llmsgr
.DEFAULT_GOAL := build
.PHONY:fmt vet build
fmt:
	go fmt
vet: fmt
	go vet
build: vet
	GOARCH=amd64 GOOS=linux go build -ldflags "-X go.szostok.io/version.version=$(git describe --tags)" -o build/linux/amd64/ llmsger.go
	GOARCH=386 GOOS=linux go build -ldflags "-X go.szostok.io/version.version=$(git describe --tags)" -o build/linux/x86/ llmsger.go
	GOARCH=amd64 GOOS=windows go build -ldflags "-X go.szostok.io/version.version=$(git describe --tags)" -o build/windows/amd64/ llmsger.go
	GOARCH=386 GOOS=windows go build -ldflags "-X go.szostok.io/version.version=$(git describe --tags)" -o build/windows/x86/ llmsger.go
	go build -ldflags "-X go.szostok.io/version.version=$(git describe --tags)"