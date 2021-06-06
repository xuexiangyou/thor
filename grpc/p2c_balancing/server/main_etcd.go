package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/xuexiangyou/thor/core/discov"
	pb "github.com/xuexiangyou/thor/grpc/proto/hello"
	"google.golang.org/grpc"
)

var (
	addr = "127.0.0.1:50071"
	etcdEndpoints = []string{
		"127.0.0.1:2379",
	}
	etcdKey = "add.rpc"
)

type helloServer struct {
	addr string
}

func (s *helloServer) UnaryHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: fmt.Sprintf("%s (from %s)", req.Message, s.addr)}, nil
}

func main() {
	registerEtcd := func() error {
		pubListenOn := addr
		pubClient := discov.NewPublisher(etcdEndpoints, etcdKey, pubListenOn)
		return pubClient.KeepAlive()
	}
	err := registerEtcd()
	if err != nil {
		fmt.Println(err)
		return
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	s := grpc.NewServer()
	pb.RegisterHelloServer(s, &helloServer{addr: addr})
	log.Printf("serving on %s\n", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
