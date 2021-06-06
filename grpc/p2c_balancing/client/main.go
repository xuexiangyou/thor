package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/xuexiangyou/thor/core/balancer/p2c"
	"github.com/xuexiangyou/thor/core/resolver"
	ecpb "github.com/xuexiangyou/thor/grpc/proto/echo"
	"google.golang.org/grpc"
)

var addrs = []string{"localhost:50051", "localhost:50052"}

var etcdAdders = []string{"129.211.63.154:2379"}

type (
	client struct {
		conn *grpc.ClientConn
	}
)

func BuildDirectTarget(endpoints []string) string {
	return fmt.Sprintf("%s:///%s", resolver.DirectScheme,
		strings.Join(endpoints, resolver.EndpointSep))
}

func BuildDiscovTarget(endpoints []string, key string) string {
	return fmt.Sprintf("%s://%s/%s", resolver.DiscovScheme,
		strings.Join(endpoints, resolver.EndpointSep), key)
}

func callUnaryEcho(c ecpb.EchoClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.UnaryEcho(ctx, &ecpb.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Println(r.Message)
}

func makeRPCs(cc *grpc.ClientConn, n int) {
	hwc := ecpb.NewEchoClient(cc)
	for i := 0; i < n; i++ {
		callUnaryEcho(hwc, "this is examples/load_balancing")
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
	conn, err := grpc.Dial(BuildDiscovTarget(etcdAdders, "add.rpc"),
		grpc.WithBalancerName(p2c.Name),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	makeRPCs(conn, 10)
}

func init() {
	resolver.RegisterResolver()
}
