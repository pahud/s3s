GOOSDEV	?= $(shell uname -s)
GOARCH  ?= amd64
CONTAINER ?= $(shell docker run -d myapp false)
GOOS ?= $(shell uname -s |  tr '[:upper:]' '[:lower:]')

build:
ifeq ($(GOOS),darwin)
	@docker run -ti --rm -v $(shell pwd):/go/src/myapp.github.com -w /go/src/myapp.github.com  golang:1.10 /bin/sh -c "make build-darwin"
else
	@docker run -ti --rm -v $(shell pwd):/go/src/myapp.github.com -w /go/src/myapp.github.com  golang:1.10 /bin/sh -c "make build-linux"
endif
	
build-linux:
	@go get -u github.com/golang/dep/cmd/dep
	@[ ! -f ./Gopkg.toml ] && dep init || true
	@dep ensure
	@GOOS=linux GOARCH=amd64 go build -o s3s.linux main.go

build-darwin:
	@go get -u github.com/golang/dep/cmd/dep
	@[ ! -f ./Gopkg.toml ] && dep init || true
	@dep ensure
	@GOOS=darwin GOARCH=amd64 go build -o s3s.darwin main.go
