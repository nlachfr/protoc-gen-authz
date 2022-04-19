package authorize

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
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
		res.Peer.Addr = p.Addr.String()
		res.Peer.AuthInfo = p.AuthInfo.AuthType()
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

func BuildProgramVars(ctx context.Context, message proto.Message) interface{} {
	res := map[string]interface{}{
		"_ctx": AuthorizationContextFromContext(ctx),
	}
	fields := message.ProtoReflect().Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		res[field.TextName()] = message.ProtoReflect().Get(field)
	}
	return res
}
