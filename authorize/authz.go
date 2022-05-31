package authorize

import (
	"context"
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	v1alpha1 "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func AuthorizationContextFromContext(ctx context.Context) *AuthorizationContext {
	res := &AuthorizationContext{
		Peer: &AuthorizationContext_Peer{
			Addr:     "",
			AuthInfo: "",
		},
		Metadata: make(map[string]*AuthorizationContext_MetadataValue),
	}
	if p, ok := peer.FromContext(ctx); ok {
		if p.Addr != nil {
			res.Peer.Addr = p.Addr.String()
		}
		if p.AuthInfo != nil {
			res.Peer.AuthInfo = p.AuthInfo.AuthType()
		}
	}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for k, v := range md {
			res.Metadata[k] = &AuthorizationContext_MetadataValue{
				Values: v,
			}
		}
	}
	return res
}

func TypeFromFunctionType(t *FileRule_Overloads_Function_Type) *v1alpha1.Type {
	switch v := t.Type.(type) {
	case *FileRule_Overloads_Function_Type_Primitive_:
		switch v.Primitive {
		case FileRule_Overloads_Function_Type_BOOL:
			return decls.Bool
		case FileRule_Overloads_Function_Type_INT:
			return decls.Int
		case FileRule_Overloads_Function_Type_UINT:
			return decls.Uint
		case FileRule_Overloads_Function_Type_DOUBLE:
			return decls.Double
		case FileRule_Overloads_Function_Type_BYTES:
			return decls.Bytes
		case FileRule_Overloads_Function_Type_STRING:
			return decls.String
		case FileRule_Overloads_Function_Type_DURATION:
			return decls.Duration
		case FileRule_Overloads_Function_Type_TIMESTAMP:
			return decls.Timestamp
		case FileRule_Overloads_Function_Type_ERROR:
			return decls.Error
		case FileRule_Overloads_Function_Type_DYN:
			return decls.Dyn
		case FileRule_Overloads_Function_Type_ANY:
			return decls.Any
		}
	case *FileRule_Overloads_Function_Type_Object:
		return decls.NewObjectType(v.Object)
	case *FileRule_Overloads_Function_Type_Array_:
		return decls.NewListType(TypeFromFunctionType(v.Array.Type))
	case *FileRule_Overloads_Function_Type_Map_:
		return decls.NewMapType(TypeFromFunctionType(v.Map.Key), TypeFromFunctionType(v.Map.Value))
	}
	return decls.Null
}

type AuthzInterceptor interface {
	GetUnaryServerInterceptor() grpc.UnaryServerInterceptor
	GetStreamServerInterceptor() grpc.StreamServerInterceptor
}

func NewAuthzInterceptor(methodProgramMapping map[string]cel.Program) AuthzInterceptor {
	return &authzInterceptor{
		methodProgramMapping: methodProgramMapping,
	}
}

type authzInterceptor struct {
	methodProgramMapping map[string]cel.Program
}

func (i *authzInterceptor) authorize(ctx context.Context, req interface{}, fullMethod string) error {
	if pgr, ok := i.methodProgramMapping[fullMethod]; ok {
		if message, ok := req.(proto.Message); ok {
			val, _, err := pgr.ContextEval(ctx, map[string]interface{}{
				"context": AuthorizationContextFromContext(ctx),
				"request": message,
			})
			if err != nil {
				return err
			}
			if !types.IsBool(val) {
				return status.Error(codes.Unknown, "")
			}
			if !val.Value().(bool) {
				return status.Error(codes.PermissionDenied, fmt.Sprintf(`Permission denied on "%s"`, fullMethod))
			}
		}
	}
	return nil
}

func (i *authzInterceptor) UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if err := i.authorize(ctx, req, info.FullMethod); err != nil {
		return nil, err
	}
	return handler(ctx, req)
}

func (i *authzInterceptor) StreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if err := i.authorize(ss.Context(), nil, info.FullMethod); err != nil {
		return err
	}
	return handler(srv, ss)
}

func (i *authzInterceptor) GetUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return i.UnaryServerInterceptor

}
func (i *authzInterceptor) GetStreamServerInterceptor() grpc.StreamServerInterceptor {
	return i.StreamServerInterceptor
}
