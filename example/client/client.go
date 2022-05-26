package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	v1 "github.com/Neakxs/protoc-gen-authz/example/service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	socket = flag.String("target", "", "unix socket of the server")
)

func main() {
	flag.Parse()
	if socket == nil || len(*socket) == 0 {
		os.Exit(1)
	}
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second)
	defer cancelFn()
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("passthrough:///unix://%s", *socket), grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	c := v1.NewOrgServiceClient(conn)
	var md metadata.MD
	var pingReq *v1.PingRequest

	pingReq = &v1.PingRequest{}
	fmt.Printf("Ping with md=%s and req=%s\n", md, pingReq)
	_, err = c.Ping(metadata.NewOutgoingContext(context.Background(), md), pingReq)
	fmt.Println("\t\t\tError:", err)

	pingReq = &v1.PingRequest{Ping: "ok"}
	fmt.Printf("Ping with md=%s and req=%s\n", md, pingReq)
	_, err = c.Ping(metadata.NewOutgoingContext(context.Background(), md), pingReq)
	fmt.Println("\t\t\tError:", err)

	pingReq = &v1.PingRequest{Ping: "ok"}
	md = metadata.MD{"x-pong": []string{"yes"}}
	fmt.Printf("Ping with md=%s and req=%s\n", md, pingReq)
	_, err = c.Ping(metadata.NewOutgoingContext(context.Background(), md), pingReq)
	fmt.Println("\t\t\tError:", err)

	var pongReq *v1.PongRequest

	pongReq = &v1.PongRequest{}
	md = metadata.MD{}
	fmt.Printf("Ping with md=%s and req=%s\n", md, pongReq)
	_, err = c.Pong(metadata.NewOutgoingContext(context.Background(), md), pongReq)
	fmt.Println("\t\t\tError:", err)

	pongReq = &v1.PongRequest{Pong: "ok"}
	fmt.Printf("Ping with md=%s and req=%s\n", md, pongReq)
	_, err = c.Pong(metadata.NewOutgoingContext(context.Background(), md), pongReq)
	fmt.Println("\t\t\tError:", err)

	pongReq = &v1.PongRequest{Pong: "ok"}
	md = metadata.MD{"x-pong": []string{"yes"}}
	fmt.Printf("Ping with md=%s and req=%s\n", md, pongReq)
	_, err = c.Pong(metadata.NewOutgoingContext(context.Background(), md), pongReq)
	fmt.Println("\t\t\tError:", err)
}
