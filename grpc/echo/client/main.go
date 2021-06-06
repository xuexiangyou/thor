package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/xuexiangyou/thor/grpc/proto/echo"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
	defaultMessage = "message"
)

func main() {
	// set up a connection to the server
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close() // 关闭链接
	c := pb.NewEchoClient(conn)

	// Contact the server and print out its response.
	message := defaultMessage
	if len(os.Args) > 1 {
		message = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second) // todo 如果server 执行时间在timeout之后返回client 返回报错context deadline exceeded
	defer cancel()
	r, err := c.UnaryEcho(ctx, &pb.EchoRequest{
		Message: message,
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
