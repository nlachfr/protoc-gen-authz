package authorize

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	v1alpha1 "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// func AuthorizationContextFromContext(ctx context.Context) *AuthorizationContext {
// 	res := &AuthorizationContext{
// 		Peer: &AuthorizationContext_Peer{
// 			Addr:     "",
// 			AuthInfo: "",
// 		},
// 		Metadata: make(map[string]*AuthorizationContext_MetadataValue),
// 	}
// 	if p, ok := peer.FromContext(ctx); ok {
// 		if p.Addr != nil {
// 			res.Peer.Addr = p.Addr.String()
// 		}
// 		if p.AuthInfo != nil {
// 			res.Peer.AuthInfo = p.AuthInfo.AuthType()
// 		}
// 	}
// 	if md, ok := metadata.FromIncomingContext(ctx); ok {
// 		for k, v := range md {
// 			res.Metadata[k] = &AuthorizationContext_MetadataValue{
// 				Values: v,
// 			}
// 		}
// 	}
// 	return res
// }

func TypeFromOverloadType(t *FileRule_Overloads_Type) *v1alpha1.Type {
	switch v := t.Type.(type) {
	case *FileRule_Overloads_Type_Primitive_:
		switch v.Primitive {
		case FileRule_Overloads_Type_BOOL:
			return decls.Bool
		case FileRule_Overloads_Type_INT:
			return decls.Int
		case FileRule_Overloads_Type_UINT:
			return decls.Uint
		case FileRule_Overloads_Type_DOUBLE:
			return decls.Double
		case FileRule_Overloads_Type_BYTES:
			return decls.Bytes
		case FileRule_Overloads_Type_STRING:
			return decls.String
		case FileRule_Overloads_Type_DURATION:
			return decls.Duration
		case FileRule_Overloads_Type_TIMESTAMP:
			return decls.Timestamp
		case FileRule_Overloads_Type_ERROR:
			return decls.Error
		case FileRule_Overloads_Type_DYN:
			return decls.Dyn
		case FileRule_Overloads_Type_ANY:
			return decls.Any
		}
	case *FileRule_Overloads_Type_Object:
		return decls.NewObjectType(v.Object)
	case *FileRule_Overloads_Type_Array_:
		return decls.NewListType(TypeFromOverloadType(v.Array.Type))
	case *FileRule_Overloads_Type_Map_:
		return decls.NewMapType(TypeFromOverloadType(v.Map.Key), TypeFromOverloadType(v.Map.Value))
	}
	return decls.Null
}

type AuthzInterceptor interface {
	Authorize(ctx context.Context, method string, headers http.Header, request interface{}) error
}

func NewAuthzInterceptor(methodProgramMapping map[string]cel.Program) AuthzInterceptor {
	return &authzInterceptor{
		methodProgramMapping: methodProgramMapping,
	}
}

type authzInterceptor struct {
	methodProgramMapping map[string]cel.Program
}

func (i *authzInterceptor) Authorize(ctx context.Context, method string, headers http.Header, request interface{}) error {
	if pgr, ok := i.methodProgramMapping[method]; ok {
		if val, _, err := pgr.ContextEval(ctx, map[string]interface{}{
			"headers": headers,
			"request": request,
		}); err != nil {
			return err
		} else if !types.IsBool(val) || !val.Value().(bool) {
			return fmt.Errorf(`permission denied on "%s"`, method)
		}
	}
	return nil
}

// type AuthzInterceptor interface {
// 	GetUnaryServerInterceptor() grpc.UnaryServerInterceptor
// 	GetStreamServerInterceptor() grpc.StreamServerInterceptor
// }

// func NewAuthzInterceptor(methodProgramMapping map[string]cel.Program) AuthzInterceptor {
// 	return &authzInterceptor{
// 		methodProgramMapping: methodProgramMapping,
// 	}
// }

// type authzInterceptor struct {
// 	methodProgramMapping map[string]cel.Program
// }

// func (i *authzInterceptor) authorize(ctx context.Context, req interface{}, fullMethod string) error {
// 	if pgr, ok := i.methodProgramMapping[fullMethod]; ok {
// 		if message, ok := req.(proto.Message); ok {
// 			val, _, err := pgr.ContextEval(ctx, map[string]interface{}{
// 				"context": AuthorizationContextFromContext(ctx),
// 				"request": message,
// 			})
// 			if err != nil {
// 				return err
// 			}
// 			if !types.IsBool(val) {
// 				return status.Error(codes.Unknown, "")
// 			}
// 			if !val.Value().(bool) {
// 				return status.Error(codes.PermissionDenied, fmt.Sprintf(`Permission denied on "%s"`, fullMethod))
// 			}
// 		}
// 	}
// 	return nil
// }

// func (i *authzInterceptor) GetUnaryServerInterceptor() grpc.UnaryServerInterceptor {
// 	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
// 		if err := i.authorize(ctx, req, info.FullMethod); err != nil {
// 			return nil, err
// 		}
// 		return handler(ctx, req)
// 	}
// }

// func (i *authzInterceptor) GetStreamServerInterceptor() grpc.StreamServerInterceptor {
// 	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
// 		if err := i.authorize(ss.Context(), nil, info.FullMethod); err != nil {
// 			return err
// 		}
// 		return handler(srv, ss)
// 	}
// }
