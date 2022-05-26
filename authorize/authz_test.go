package authorize

import (
	"context"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
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
