package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"path"

	"github.com/Neakxs/protoc-gen-authz/authorize"
	v1 "github.com/Neakxs/protoc-gen-authz/example/service/v1"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type orgServer struct {
	v1.UnimplementedOrgServiceServer
}

func (s *orgServer) Ping(context.Context, *v1.PingRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *orgServer) Pong(context.Context, *v1.PongRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func main() {
	authzInterceptor, err := v1.NewOrgServiceAuthzInterceptor(&authorize.Overload{
		Name: "do",
		Function: func(v ...ref.Val) ref.Val {
			fmt.Println("yes")
			return types.Bool(true)
		},
	})
	if err != nil {
		panic(err)
	}
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(authzInterceptor.GetUnaryServerInterceptor()),
		grpc.StreamInterceptor(authzInterceptor.GetStreamServerInterceptor()),
	)
	v1.RegisterOrgServiceServer(srv, &orgServer{})
	dir, err := ioutil.TempDir("/tmp", "*")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listening on unix://%s/unix.sock...\n", dir)
	lis, err := net.Listen("unix", path.Join(dir, "unix.sock"))
	if err != nil {
		panic(err)
	}
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
