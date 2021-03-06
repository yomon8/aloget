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

package:
	rm -fR ./pkg && mkdir ./pkg ;\
		gox \
		-osarch $(OSARCH) \
		-output "./pkg/{{.OS}}_{{.Arch}}/{{.Dir}}" \
		-ldflags "-X github.com/yomon8/aloget.version=$(VERSION)" \
		./cmd/...;\
	    for d in $$(ls ./pkg);do zip ./pkg/$${d}.zip ./pkg/$${d}/*;done

build:
	go build -o $(BIN) -ldflags "-X main.version=$(VERSION)" ./cmd/...

linuxbuild:
	GOOS=linux GOARCH=amd64 make build

clean:
	go clean
