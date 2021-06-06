package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/xuexiangyou/thor/grpc/interceptor/serverinterceptors"
	pb "github.com/xuexiangyou/thor/grpc/proto/echo"
	"google.golang.org/grpc"
)

var port = flag.Int("port", 50051, "the port to serve on")

type server struct {
	pb.UnimplementedEchoServer
}

func (s *server) UnaryEcho(ctx context.Context, in *pb.EchoRequest) (*pb.EchoResponse, error) {
	fmt.Println("server 执行内容中心")
	// panic("哈哈哈哈")
	return &pb.EchoResponse{Message: in.Message}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		serverinterceptors.UnaryTracingInterceptor(),
		serverinterceptors.UnaryCrashInterceptor(),
	}

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
	}

	s := grpc.NewServer(opts...)
	pb.RegisterEchoServer(s, &server{})
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
