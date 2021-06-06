package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/xuexiangyou/thor/grpc/interceptor/clientinterceptors"
	"google.golang.org/grpc"
	pb "github.com/xuexiangyou/thor/grpc/proto/echo"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

func main() {
	flag.Parse()

	options := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithChainUnaryInterceptor(
			clientinterceptors.TracingInterceptor,
			clientinterceptors.DurationInterceptor,
			clientinterceptors.TimeoutInterceptor(5 * time.Second),
		),
	}
	timeCtx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()
	conn, err := grpc.DialContext(timeCtx, *addr, options...)
	if err != nil {
		fmt.Printf("rpc dial: %s, error: %s, make sure rpc service is alread started",
			*addr, err.Error())
		return
	}
	defer conn.Close()

	rgc := pb.NewEchoClient(conn)
	callUnaryEcho(rgc, "hello world")
}

func callUnaryEcho(client pb.EchoClient, message string) {
	resp, err := client.UnaryEcho(context.Background(), &pb.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("client.UnaryEcho(_) = _, %v: ", err)
	}
	fmt.Println("UnaryEcho: ", resp.Message)
}
