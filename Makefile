GOPATH := $(PWD)/.cache
PATH := $(PATH):$(GOPATH)/bin
SHELL := env PATH=$(PATH) /bin/bash

PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(GOPATH)/bin/protoc-gen-go-grpc
PROTOC_GEN_GO_AUTHZ := $(GOPATH)/bin/protoc-gen-go-authz

PROTO := $(shell find authorize -name '*.proto')
PROTO_EXAMPLE := $(shell find example -name '*.proto')
GENPROTO_GO := $(PROTO:.proto=.pb.go)
GENPROTO_EXAMPLE_GO := $(PROTO_EXAMPLE:.proto=.pb.go) $(PROTO_EXAMPLE:.proto=_grpc.pb.go) $(PROTO_EXAMPLE:.proto=.pb.authz.go)

.PHONY: all
all: go-genproto

$(PROTOC_GEN_GO):
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28

$(PROTOC_GEN_GO_GRPC):
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

.PHONY: $(PROTOC_GEN_GO_AUTHZ)
$(PROTOC_GEN_GO_AUTHZ):
	go install ./cmd/protoc-gen-go-authz

%.pb.go: %.proto
	protoc --go_out=. --go_opt=paths=source_relative $<

%_grpc.pb.go: %.proto
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative $<

%.pb.authz.go: %.proto
	which protoc-gen-go-authz
	protoc --go-authz_out=. --go-authz_opt=paths=source_relative $<

.PHONY: go-genproto
go-genproto: $(PROTOC_GEN_GO) $(PROTOC_GEN_GO_GRPC) $(PROTOC_GEN_GO_AUTHZ) $(GENPROTO_GO) $(GENPROTO_EXAMPLE_GO)


