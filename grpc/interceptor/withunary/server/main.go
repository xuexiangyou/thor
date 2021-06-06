package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/xuexiangyou/thor/grpc/auth"
	pb "github.com/xuexiangyou/thor/grpc/proto/echo"
	"google.golang.org/grpc"
)

var port = flag.Int("port", 50051, "the port server on")

func main() {
	flag.Parse()
	fmt.Printf("server starting on port %d...\n", *port)

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(auth.EnsureValidToken),
	}
	s := grpc.NewServer(opts...)

	pb.RegisterEchoServer(s, &ecServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *ecServer) UnaryEcho(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	fmt.Println("在拦截器执行后在直接服务逻辑") // todo 在拦截器执行后在直接服务逻辑
	time.Sleep(5 * time.Second)
	return &pb.EchoResponse{Message: req.Message}, nil
}

type ecServer struct {
	pb.UnimplementedEchoServer
}
