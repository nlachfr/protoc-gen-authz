# protoc-gen-authz

## About

__protoc-gen-authz__ is an authorization plugin for the protocol buffers compiler [protoc](https://github.com/protocolbuffers/protobuf). It relies on the [Common Expression Language](https://github.com/google/cel-spec) specification for writing authorization rules and can use gRPC metadata against protobuf messages.

The only language supported is [Go](https://go.dev/).

## Installation

For installing the plugin, you can simply run the following command :

```shell
go install github.com/Neakxs/protoc-gen-authz/cmd/protoc-gen-go-authz
```

The binary will be placed in your $GOBIN location.

## Usage

1. Create protobuf definition

```protobuf
syntax = "proto3";

package service.v1;
option go_package = "github.com/org/proto/gen/go/service/v1";

import "authorize/authz.proto";

option (neakxs.authz.globals) = {
    functions: [
        {
            key: 'isAdmin'
            value: '"x-admin" in _ctx.metadata'
        }
    ]
};

service OrgService {
    rpc Ping(PingRequest) returns (PingResponse) {};
}

message PingRequest {
    option (neakxs.authz.rule).expr = 'isAdmin() && size(ping) > 0';
    string ping = 1;
}
message PingResponse {
    string pong = 2;
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

```golang
func (*orgServiceServer) Ping(ctx context.Context, in *v1.PingRequest) (*v1.PingResponse, error) {
    if err := in.Authorize(ctx); err != nil {
        return nil, err
    }
    // implement logic
}
```

4. Profit

## Configuration

It is possible to use a configuration file for defining global functions across all your proto definitions.

With the example above, we can create a `config.yml` file :

```yaml
version: v1
globals:
  functions:
    isAdmin: "x-admin" in _ctx.metadata
```

You can then use it with the `--go-authz_opt=config=path/to/config.yml` option.

> When the same function is defined inside a protobuf file and in the configuration, the protobuf one is used.