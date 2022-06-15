package authorize

import (
	"context"
	"net/http"
	"testing"

	testdata "github.com/Neakxs/protoc-gen-authz/testdata"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/go-cmp/cmp"
	v1alpha1 "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"google.golang.org/protobuf/testing/protocmp"
)

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
	astBool, _ := env.Compile(`request.ping == "ping" && "hdr" in headers`)
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
			err := NewAuthzInterceptor(tt.Mapping).Authorize(context.Background(), "", http.Header{"hdr": []string{}}, tt.Request)
			if (err != nil && !tt.WantErr) || (err == nil && tt.WantErr) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}
