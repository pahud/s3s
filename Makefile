GOOSDEV	?= $(shell uname -s)
GOARCH  ?= amd64
GOOS ?= $(shell uname -s |  tr '[:upper:]' '[:lower:]')
S3TMPBUCKET ?= tmp.pahud.net

build:
ifeq ($(GOOS),darwin)
	@docker run -ti --rm -v $(shell pwd):/go/src/myapp.github.com -w /go/src/myapp.github.com  golang:1.10 /bin/sh -c "make build-darwin"
else
	@docker run -ti --rm -v $(shell pwd):/go/src/myapp.github.com -w /go/src/myapp.github.com  golang:1.10 /bin/sh -c "make build-linux"
endif
	
run:
	@docker run -ti --rm \
	-e BITLY_TOKEN=$$BITLY_TOKEN \
	-v $(HOME)/.aws:/root/.aws \
	-v $(shell pwd):/go/src/myapp.github.com \
	-w /go/src/myapp.github.com  golang:1.10 \
	/bin/sh -c "go run main.go $(S3TMPBUCKET) main.go"
	
build-linux:
	@go get -u github.com/golang/dep/cmd/dep
	@[ ! -f ./Gopkg.toml ] && dep init || true
	@dep ensure
	@GOOS=linux GOARCH=amd64 go build -o s3s main.go

build-darwin:
	@go get -u github.com/golang/dep/cmd/dep
	@[ ! -f ./Gopkg.toml ] && dep init || true
	@dep ensure
	@GOOS=darwin GOARCH=amd64 go build -o s3s main.go
