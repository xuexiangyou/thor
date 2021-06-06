package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/xuexiangyou/thor/grpc/auth"
	"google.golang.org/grpc"

	pb "github.com/xuexiangyou/thor/grpc/proto/echo"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

func callUnaryEcho(client pb.EchoClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	resp, err := client.UnaryEcho(ctx, &pb.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("client.UnaryEcho(_) = _, %v", err)
	}
	fmt.Println("UnaryEcho: ", resp.Message)
}

func main() {
	flag.Parse()
	fmt.Printf("server starting on port %s...\n", *addr)

	perRpc := &auth.Credential{
		App: "testApp",
		Token: "testToken",
	}

	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(perRpc),
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial(*addr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	rgc := pb.NewEchoClient(conn)

	callUnaryEcho(rgc, "hello word")
}
