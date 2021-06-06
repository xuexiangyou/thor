package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/xuexiangyou/thor/grpc/proto/echo"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type Server struct {
	pb.UnimplementedEchoServer
}

func (s *Server) UnaryEcho(ctx context.Context, in *pb.EchoRequest) (*pb.EchoResponse, error) {
	log.Printf("Received: %v", in.GetMessage())
	time.Sleep(2 * time.Second)
	return &pb.EchoResponse{
		Message: "echo " + in.GetMessage(),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterEchoServer(s, &Server{})
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
