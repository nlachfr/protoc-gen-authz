package authorize

import (
	"context"

	"github.com/Neakxs/protoc-gen-authz/api"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/proto"
)

func AuthorizationContextFromContext(ctx context.Context) *api.AuthorizationContext {
	res := &api.AuthorizationContext{
		Peer: &api.AuthorizationContext_Peer{
			Addr:     "",
			AuthInfo: "",
		},
		Metadata: make(map[string]*api.AuthorizationContext_MetadataValue),
	}
	if p, ok := peer.FromContext(ctx); ok {
		res.Peer.Addr = p.Addr.String()
		res.Peer.AuthInfo = p.AuthInfo.AuthType()
	}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for k, v := range md {
			res.Metadata[k] = &api.AuthorizationContext_MetadataValue{
				Values: v,
			}
		}
	}
	return res
}

func BuildProgramVars(ctx context.Context, message proto.Message) interface{} {
	res := map[string]interface{}{
		".ctx": AuthorizationContextFromContext(ctx),
	}
	fields := message.ProtoReflect().Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		res[field.TextName()] = message.ProtoReflect().Get(field)
	}
	return res
}
