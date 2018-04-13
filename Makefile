TARGET ?= darwin
ARCH ?= amd64
DOCKER_REPO=itsdalmo/ami-resource
SRC=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

default: test

build: test
	@echo "== Building check =="
	CGO_ENABLED=0 GOOS=$(TARGET) GOARCH=$(ARCH) go build -o check -v cmd/check/main.go
	@echo "== Building in =="
	CGO_ENABLED=0 GOOS=$(TARGET) GOARCH=$(ARCH) go build -o in -v cmd/in/main.go
	@echo "== Building out =="
	CGO_ENABLED=0 GOOS=$(TARGET) GOARCH=$(ARCH) go build -o out -v cmd/out/main.go

test:
	@echo "== Test =="
	gofmt -s -l -w $(SRC)
	go vet -v ./...
	go test -race -v ./...

clean:
	@echo "== Cleaning =="
	rm -f check
	rm -f in
	rm -f out

lint:
	@echo "== Lint =="
	golint cmd
	golint models
	golint manager

docker:
	@echo "== Docker build =="
	docker build -t $(DOCKER_REPO):dev .

integration: clean build
	@echo "== Integration tests =="
	@# TODO

.PHONY: default build clean test docker integration
