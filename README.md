# protoc-gen-authz

[![Coverage](https://coveralls.io/repos/Neakxs/protoc-gen-authz/badge.svg?branch=main&service=github)](https://coveralls.io/github/Neakxs/protoc-gen-authz?branch=main) [![GoReportCard](https://goreportcard.com/badge/github.com/Neakxs/protoc-gen-authz)](https://goreportcard.com/badge/github.com/Neakxs/protoc-gen-authz) ![GitHub](https://img.shields.io/github/license/Neakxs/protoc-gen-authz)

## About

__protoc-gen-authz__ is an authorization plugin for the protocol buffers compiler [protoc](https://github.com/protocolbuffers/protobuf). It relies on the [Common Expression Language](https://github.com/google/cel-spec) specification for writing authorization rules and can use gRPC metadata against protobuf messages.

The only language supported is [Go](https://go.dev/).

## Installation

For installing the plugin, you can simply run the following command :

```shell
go install github.com/Neakxs/protoc-gen-authz/cmd/protoc-gen-go-authz
```

The binary will be placed in your $GOBIN location.

## About CEL Environment

When writing rules, the following variables are avaible :
- `headers` (**map[string][]string** type)
- `request` (declared request message type)

The `headers.get(string)` receiver function has been added on the `headers` variable. It allows writing of easier rules by using the go `func (http.Header) Get(string)` function under the hood.

## Usage

1. Create protobuf definition

```protobuf
syntax = "proto3";

package service.v1;
option go_package = "github.com/Neakxs/protoc-gen-authz/example/service/v1";

import "authorize/authz.proto";
import "google/protobuf/empty.proto";

option (authorize.file) = {
    globals: {
        functions: [
            {
                key: 'canPong'
                value: '"X-Pong" in headers'
            }
        ]
    };
	rules: [
        {
            key: "service.v1.OrgService.Pong"
            value: { 
                expr: "canPong() && size(request.pong) > 0"
            }
        }
    ]
};

service OrgService {
    rpc Ping(PingRequest) returns (google.protobuf.Empty) {
        option (authorize.method).expr = '!canPong() && size(request.ping) > 0';
    };
    rpc Pong(PongRequest) returns (google.protobuf.Empty) {};
}

message PingRequest {
    string ping = 1;
}
message PongRequest {
    string pong = 1;
}
```

2. Generate protobuf service

```shell
protoc \
    --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_out=paths=source_relative \
    --go-authz_out=. --go-authz_opt=paths=source_relative \
    github.com/org/proto/gen/go/service/v1/example.proto
```

3. Implement gRPC service

4. Add interceptors to your gRPC server

```golang
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"path"

	v1 "github.com/Neakxs/protoc-gen-authz/example/service/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type orgServer struct {
	v1.UnimplementedOrgServiceServer
}

func (s *orgServer) Ping(context.Context, *v1.PingRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *orgServer) Pong(context.Context, *v1.PongRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func main() {
	authzInterceptor, err := v1.NewOrgServiceAuthzInterceptor()
	if err != nil {
		panic(err)
	}
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(authzInterceptor.GetUnaryServerInterceptor()),
		grpc.StreamInterceptor(authzInterceptor.GetStreamServerInterceptor()),
	)
	v1.RegisterOrgServiceServer(srv, &orgServer{})
	dir, err := ioutil.TempDir("/tmp", "*")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listening on unix://%s/unix.sock...\n", dir)
	lis, err := net.Listen("unix", path.Join(dir, "unix.sock"))
	if err != nil {
		panic(err)
	}
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
```

5. Profit

## Configuration

It is possible to use a configuration file for defining global functions across all your proto definitions.

With the example above, we can create a `config.yml` file :

```yaml
version: v1
globals:
  functions:
    isAdmin: "x-admin" in context.metadata
rules:
  service.v1.OrgService.Pong:
    expr: "canPong() && size(request.pong) > 0"
```

You can then use it with the `--go-authz_opt=config=path/to/config.yml` option.

> When the same function is defined inside a protobuf file and in the configuration, the protobuf one is used.