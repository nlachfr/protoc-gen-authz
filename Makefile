GOPATH := $(shell go env GOPATH)
PATH := $(PATH):$(GOPATH)/bin
SHELL := env PATH=$(PATH) /bin/bash

PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go

PROTO := $(shell find api -name '*.proto')
GENPROTO_GO := $(PROTO:.proto=.pb.go)

.PHONY: all
all: go-genproto

$(PROTOC_GEN_GO):
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27

%.pb.go: %.proto
	protoc --go_out=. --go_opt=paths=source_relative $<

.PHONY: go-genproto
go-genproto: $(PROTOC_GEN_GO) $(GENPROTO_GO)


