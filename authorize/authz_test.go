package authorize

import (
	"context"
	"net"
	"testing"

	"github.com/Neakxs/protoc-gen-authz/testdata"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/go-cmp/cmp"
	v1alpha1 "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestAuthorizationContextFromContext(t *testing.T) {
	tests := []struct {
		Name string
		In   context.Context
		Want *AuthorizationContext
	}{
		{
			Name: "Default",
			In:   context.Background(),
			Want: &AuthorizationContext{
				Peer: &AuthorizationContext_Peer{},
			},
		},
		{
			Name: "Authorization",
			In:   metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"authorization": "Basic user:password"})),
			Want: &AuthorizationContext{
				Peer: &AuthorizationContext_Peer{},
				Metadata: map[string]*AuthorizationContext_MetadataValue{
					"authorization": {
						Values: []string{"Basic user:password"},
					},
				},
			},
		},
		{
			Name: "IP Source",
			In: peer.NewContext(context.Background(), &peer.Peer{
				Addr: &net.IPAddr{IP: net.ParseIP("127.0.0.1")},
			}),
			Want: &AuthorizationContext{
				Peer: &AuthorizationContext_Peer{
					Addr: "127.0.0.1",
				},
			},
		},
		{
			Name: "AuthInfo",
			In: peer.NewContext(context.Background(), &peer.Peer{
				AuthInfo: credentials.TLSInfo{},
			}),
			Want: &AuthorizationContext{
				Peer: &AuthorizationContext_Peer{
					AuthInfo: credentials.TLSInfo{}.AuthType(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			got := AuthorizationContextFromContext(tt.In)
			if !cmp.Equal(got, tt.Want, protocmp.Transform()) {
				t.Errorf("want %v, got %v", tt.Want, got)
			}
		})
	}
}

func TestTypeFromOverloadType(t *testing.T) {
	tests := []struct {
		Name string
		In   *FileRule_Overloads_Type
		Out  *v1alpha1.Type
	}{
		{
			Name: "Primitive bool",
			In: &FileRule_Overloads_Type{
				Type: &FileRule_Overloads_Type_Primitive_{
					Primitive: FileRule_Overloads_Type_BOOL,
				},
			},
			Out: decls.Bool,
		},
		{
			Name: "Primitive int",
			In: &FileRule_Overloads_Type{
				Type: &FileRule_Overloads_Type_Primitive_{
					Primitive: FileRule_Overloads_Type_INT,
				},
			},
			Out: decls.Int,
		},
		{
			Name: "Primitive uint",
			In: &FileRule_Overloads_Type{
				Type: &FileRule_Overloads_Type_Primitive_{
					Primitive: FileRule_Overloads_Type_UINT,
				},
			},
			Out: decls.Uint,
		},
		{
			Name: "Primitive double",
			In: &FileRule_Overloads_Type{
				Type: &FileRule_Overloads_Type_Primitive_{
					Primitive: FileRule_Overloads_Type_DOUBLE,
				},
			},
			Out: decls.Double,
		},
		{
			Name: "Primitive bytes",
			In: &FileRule_Overloads_Type{
				Type: &FileRule_Overloads_Type_Primitive_{
					Primitive: FileRule_Overloads_Type_BYTES,
				},
			},
			Out: decls.Bytes,
		},
		{
			Name: "Primitive string",
			In: &FileRule_Overloads_Type{
				Type: &FileRule_Overloads_Type_Primitive_{
					Primitive: FileRule_Overloads_Type_STRING,
				},
			},
			Out: decls.String,
		},
		{
			Name: "Primitive duration",
			In: &FileRule_Overloads_Type{
				Type: &FileRule_Overloads_Type_Primitive_{
					Primitive: FileRule_Overloads_Type_DURATION,
				},
			},
			Out: decls.Duration,
		},
		{
			Name: "Primitive timestamp",
			In: &FileRule_Overloads_Type{
				Type: &FileRule_Overloads_Type_Primitive_{
					Primitive: FileRule_Overloads_Type_TIMESTAMP,
				},
			},
			Out: decls.Timestamp,
		},
		{
			Name: "Primitive error",
			In: &FileRule_Overloads_Type{
				Type: &FileRule_Overloads_Type_Primitive_{
					Primitive: FileRule_Overloads_Type_ERROR,
				},
			},
			Out: decls.Error,
		},
		{
			Name: "Primitive dyn",
			In: &FileRule_Overloads_Type{
				Type: &FileRule_Overloads_Type_Primitive_{
					Primitive: FileRule_Overloads_Type_DYN,
				},
			},
			Out: decls.Dyn,
		},
		{
			Name: "Primitive any",
			In: &FileRule_Overloads_Type{
				Type: &FileRule_Overloads_Type_Primitive_{
					Primitive: FileRule_Overloads_Type_ANY,
				},
			},
			Out: decls.Any,
		},
		{
			Name: "Object",
			In: &FileRule_Overloads_Type{
				Type: &FileRule_Overloads_Type_Object{
					Object: "object",
				},
			},
			Out: decls.NewObjectType("object"),
		},
		{
			Name: "Array",
			In: &FileRule_Overloads_Type{
				Type: &FileRule_Overloads_Type_Array_{
					Array: &FileRule_Overloads_Type_Array{
						Type: &FileRule_Overloads_Type{
							Type: &FileRule_Overloads_Type_Primitive_{
								Primitive: FileRule_Overloads_Type_BOOL,
							},
						},
					},
				},
			},
			Out: decls.NewListType(decls.Bool),
		},
		{
			Name: "Map",
			In: &FileRule_Overloads_Type{
				Type: &FileRule_Overloads_Type_Map_{
					Map: &FileRule_Overloads_Type_Map{
						Key: &FileRule_Overloads_Type{
							Type: &FileRule_Overloads_Type_Primitive_{
								Primitive: FileRule_Overloads_Type_STRING,
							},
						},
						Value: &FileRule_Overloads_Type{
							Type: &FileRule_Overloads_Type_Primitive_{
								Primitive: FileRule_Overloads_Type_STRING,
							},
						},
					},
				},
			},
			Out: decls.NewMapType(decls.String, decls.String),
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			res := TypeFromOverloadType(tt.In)
			if !cmp.Equal(res, tt.Out, protocmp.Transform()) {
				t.Errorf("want %v, got %v", tt.Out, res)
			}
		})
	}
}

func TestAuthzInterceptor(t *testing.T) {
	env, _ := cel.NewEnv(
		cel.Types(&AuthorizationContext{}),
		cel.Types(&testdata.PingRequest{}),
		cel.Declarations(
			decls.NewVar(
				"context",
				decls.NewObjectType(string((&AuthorizationContext{}).ProtoReflect().Descriptor().FullName())),
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
		Want    *status.Status
	}{
		{
			Name: "Permission denied (bool)",
			Mapping: map[string]cel.Program{
				"": pgrBool,
			},
			Request: &testdata.PingRequest{Ping: ""},
			Want:    status.New(codes.PermissionDenied, ""),
		},
		{
			Name: "OK (bool)",
			Mapping: map[string]cel.Program{
				"": pgrBool,
			},
			Request: &testdata.PingRequest{Ping: "ping"},
			Want:    nil,
		},
		{
			Name: "Unknown (str)",
			Mapping: map[string]cel.Program{
				"": pgrString,
			},
			Request: &testdata.PingRequest{Ping: "ping"},
			Want:    status.New(codes.Unknown, ""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := NewAuthzInterceptor(tt.Mapping).(*authzInterceptor).authorize(context.Background(), tt.Request, "")
			if (err == nil && tt.Want != nil) || (tt.Want == nil && err != nil) {
				t.Errorf("want %v, got %v", tt.Want, err)
			} else if err != nil {
				if s, ok := status.FromError(err); ok {
					if s.Code() != tt.Want.Code() {
						t.Errorf("want %v, got %v", tt.Want, err)
					}
				} else {
					t.Errorf("want %v, got %v", tt.Want, err)
				}
			}
		})
	}
}
