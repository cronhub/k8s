PROJ_ROOT := $(shell pwd)

vars:
	echo "PROJ_ROOT=$(PROJ_ROOT)"
	echo "GOROOT=$(GOROOT)"
	echo "GOPATH=$(GOPATH)"
	go version

build:
	go build -o $(PROJ_ROOT)/bin/k8s $(PROJ_ROOT)/cmd/k8s

test:
	go test ./...