package interceptors

import (
	"context"
	"testing"

	"github.com/Neakxs/protoc-gen-authz/authorize"
	testdata "github.com/Neakxs/protoc-gen-authz/testdata"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"google.golang.org/grpc"
)

func TestNewGRPCUnaryInterceptor(t *testing.T) {
	env, _ := cel.NewEnv(
		cel.Types(&testdata.PingRequest{}),
		cel.Declarations(
			decls.NewVar(
				"headers",
				decls.NewMapType(
					decls.String,
					decls.NewListType(decls.String),
				),
			),
			decls.NewVar(
				"request",
				decls.NewObjectType(string((&testdata.PingRequest{}).ProtoReflect().Descriptor().FullName())),
			),
		),
	)
	astBool, _ := env.Compile(`request.ping == "ping"`)
	pgrBool, _ := env.Program(astBool)
	astString, _ := env.Compile(`request.ping`)
	pgrString, _ := env.Program(astString)
	tests := []struct {
		Name    string
		Mapping map[string]cel.Program
		Request *testdata.PingRequest
		WantErr bool
	}{
		{
			Name: "Permission denied (bool)",
			Mapping: map[string]cel.Program{
				"": pgrBool,
			},
			Request: &testdata.PingRequest{Ping: ""},
			WantErr: true,
		},
		{
			Name: "OK (bool)",
			Mapping: map[string]cel.Program{
				"": pgrBool,
			},
			Request: &testdata.PingRequest{Ping: "ping"},
			WantErr: false,
		},
		{
			Name: "Unknown (str)",
			Mapping: map[string]cel.Program{
				"": pgrString,
			},
			Request: &testdata.PingRequest{Ping: "ping"},
			WantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := NewGRPCUnaryInterceptor(authorize.NewAuthzInterceptor(tt.Mapping))(
				context.Background(),
				tt.Request,
				&grpc.UnaryServerInfo{
					FullMethod: "",
				},
				grpc.UnaryHandler(func(ctx context.Context, req interface{}) (interface{}, error) {
					return nil, nil
				}),
			)
			if (err != nil && !tt.WantErr) || (err == nil && tt.WantErr) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}
