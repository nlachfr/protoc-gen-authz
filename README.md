# protoc-gen-authz

Authorization plugin for protoc

Using protoc-gen-authz, you can add authorization rules directly inside your protobuf definition using the [CEL specification](https://github.com/google/cel-spec).

## Usage

Here is an example of a protobuf definition :

```protobuf
syntax = "proto3";

package api;

option go_package = "github.com/Neakxs/protoc-gen-authz/example/api";

import "authorize/authz.proto";

service UserService {
    rpc CreateUser(CreateUserRequest) returns (User) {};
    rpc ListUsers(ListUsersRequest) returns (ListUsersResponse) {};
    rpc DeleteUser(DeleteUserRequest) returns (User) {};
}

message User {
    string name = 1;
    string email = 2;
    string display_name = 3;
}

message CreateUserRequest {
    option (neakxs.authz.rule).expr = '"X-Role" in _ctx.metadata && "admin" in _ctx.metadata["X-Role"].values';
    User user = 1;
}
message ListUsersRequest {}
message ListUsersResponse {
    repeated User users = 1;
}
message DeleteUserRequest {
    option (neakxs.authz.rule).expr = '"X-Role" in _ctx.metadata && "admin" in _ctx.metadata["X-Role"].values';
    string name = 1;
}
```

The defined authorization rule can then be checked against a given message by calling the corresponding `Authorize` method.

```go
func (m *CreateUserRequest) Authorize(ctx context.Context) error
func (m *ListUsersRequest) Authorize(ctx context.Context) error
func (m *ListUsersResponse) Authorize(ctx context.Context) error
func (m *DeleteUserRequest) Authorize(ctx context.Context) error
```

For gRPC based request, it is possible to use peer and metadata information in the defined rules.
