HANDLER ?= main
PACKAGE ?= $(HANDLER)
GOPATH  ?= $(HOME)/go
# GOOS    ?= linux
GOOSDEV	?= $(shell uname -s)
GOARCH  ?= amd64
S3TMPBUCKET	?= pahud-tmp
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

devbuild:
	@echo "Building..."
	docker run -ti -v $(PWD):/go -e GOOS=darwin -e GOARCH=$(GOARCH) -w /go/workdir -e GOPATH:/go/workdir/go  golang:1.11 /bin/bash -c "go mod init s3s; go build" 
	# @GOOS=$(GOOSDEV) GOARCH=$(GOARCH) go build -ldflags='-w -s' -o $(HANDLER)

run:
	docker run -ti -v $(PWD):/go -e GOOS=darwin -e GOARCH=$(GOARCH) -w /go/workdir golang:1.11 /bin/bash -c "go run main.go" 
	# @GOOS=$(GOOSDEV) GOARCH=$(GOARCH) go build -ldflags='-w -s' -o $(HANDLER)