BIN      := aloget 
OSARCH   := "darwin/amd64 linux/amd64 windows/amd64"
VERSION  := $(shell git describe --tags)

all: build

test: deps build
	go test ./...

deps:
	go get -d -v -t ./...
	go get github.com/golang/lint/golint
	go get github.com/mitchellh/gox
	go get github.com/aws/aws-sdk-go
	go get github.com/dustin/go-humanize

lint: deps
	go vet ./...
	golint -set_exit_status ./...

crossbuild:
	rm -fR ./pkg && mkdir ./pkg ;\
		gox \
		-osarch $(OSARCH) \
		-output "./pkg/{{.OS}}_{{.Arch}}/{{.Dir}}" \
		-ldflags "-X config.version=$(VERSION)" \
		./cmd/...

build:
	go build -o $(BIN) -ldflags "-X config.version=$(VERSION)" ./cmd/...

linuxbuild:
	GOOS=linux GOARCH=amd64 make build

clean:
	go clean
