package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/xuexiangyou/thor/core/balancer/p2c"
	"github.com/xuexiangyou/thor/core/resolver"
	ecpb "github.com/xuexiangyou/thor/grpc/proto/hello"
	"google.golang.org/grpc"
)

var etcdAdders1 = []string{"129.211.63.154:2379"}

func BuildDirectTarget1(endpoints []string) string {
	return fmt.Sprintf("%s:///%s", resolver.DirectScheme,
		strings.Join(endpoints, resolver.EndpointSep))
}

func BuildDiscovTarget1(endpoints []string, key string) string {
	return fmt.Sprintf("%s://%s/%s", resolver.DiscovScheme,
		strings.Join(endpoints, resolver.EndpointSep), key)
}


func callUnaryHello(c ecpb.HelloClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.UnaryHello(ctx, &ecpb.HelloRequest{Message: message})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Println(r.Message)
}

func makeHelloRPCs(cc *grpc.ClientConn, n int) {
	hwc := ecpb.NewHelloClient(cc)
	for i := 0; i < n; i++ {
		callUnaryHello(hwc, "this is examples/load_balancing")
	}
}

func main() {
	// var cli client 直链
	// conn, err := grpc.Dial(BuildDirectTarget(addrs),
	// 	grpc.WithBalancerName(p2c.Name),
	// 	grpc.WithInsecure(),
	// 	grpc.WithBlock(),
	// )
	// etcd 链接
	conn, err := grpc.Dial(BuildDiscovTarget1(etcdAdders1, "add.rpc"),
		grpc.WithBalancerName(p2c.Name),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	makeHelloRPCs(conn, 10)
}

func init() {
	resolver.RegisterResolver()
}
