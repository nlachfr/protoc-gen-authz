package authorize

import (
	"context"
	"fmt"

	"github.com/google/cel-go/cel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type AuthzInterceptor interface {
	GetUnaryServerInterceptor() grpc.UnaryServerInterceptor
	GetStreamServerInterceptor() grpc.StreamServerInterceptor
}

func NewAuthzInterceptor(methodProgramMapping map[string]cel.Program) AuthzInterceptor {
	return &authzInterceptor{}
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
			if res, ok := val.Value().(bool); ok {
				if !res {
					return status.Error(codes.PermissionDenied, fmt.Sprintf(`Permission denied on "%s"`, fullMethod))
				}
			}
			return status.Error(codes.Unknown, "")
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
