package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/xuexiangyou/thor/grpc/proto/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var port = flag.Int("port", 50051, "the port to serve on")

type Server struct {
	pb.UnimplementedEchoServer
}

func (s *Server) UnaryEcho(ctx context.Context, in *pb.EchoRequest) (*pb.EchoResponse, error) {
	log.Printf("Received: %v", in.GetMessage())
	fmt.Println(metadata.FromIncomingContext(ctx))
	time.Sleep(1 * time.Second)
	return &pb.EchoResponse{
		Message: "echo " + in.GetMessage(),
	}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("server listening at port %v\n", lis.Addr())
	s := grpc.NewServer()
	pb.RegisterEchoServer(s, &Server{})
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
