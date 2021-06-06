package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/xuexiangyou/thor/grpc/auth"
	pb "github.com/xuexiangyou/thor/grpc/proto/echo"
	"google.golang.org/grpc"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

var isToken = flag.Bool("isToken", true, "is open token")

// logger is to mock a sophisticated logging system. To simplify the example, we just print out the content.
func logger(format string, a ...interface{}) {
	fmt.Printf("LOG:\t"+format+"\n", a...)
}

func callUnaryEcho(client pb.EchoClient, message string) {
	fmt.Println("这是发送请求执行开始") // todo 这是发送请求执行开始
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.UnaryEcho(ctx, &pb.EchoRequest{Message: message})
	fmt.Println("这是client拦截器执行完后在执行") // todo 这是client拦截器执行完后在执行
	if err != nil {
		log.Fatalf("client.UnaryEcho(_) = _, %v: ", err)
	}
	fmt.Println("UnaryEcho: ", resp.Message)
}

func unaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	var credsConfigured bool
	for _, o := range opts {
		_, ok := o.(grpc.PerRPCCredsCallOption)
		if ok {
			credsConfigured = true
			break
		}
	}

	if *isToken && !credsConfigured {
		opts = append(opts, grpc.PerRPCCredentials(&auth.Credential{
			App: "testApp",
			Token: "testToken",
		}))
	}

	// start := time.Now()
	fmt.Println("这是发送请求前先执行") // todo client request 前先执行
	err := invoker(ctx, method, req, reply, cc, opts...)
	fmt.Printf("这是程序执行server返回结果后会执行到这个: %s", reply.(*pb.EchoResponse).Message) // todo server 返回结果执行
	fmt.Println()
	// end := time.Now()
	// logger("RPC: %s, start time: %s, end time: %s, err: %v", method, start.Format("Basic"), end.Format(time.RFC3339), err)
	return err
	// return errors.New("故意返回错误") // todo 如果返回错误降到导致上层直接回到错误、请求失败
}

func unaryInterceptorTwo(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	fmt.Println("第二个client拦截器执行了")
	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithUnaryInterceptor(unaryInterceptor))
	// conn, err := grpc.Dial(*addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithUnaryInterceptor(unaryInterceptor), grpc.WithUnaryInterceptor(unaryInterceptorTwo)) // todo WithUnaryInterceptor 只能执行一个,请求还是执行后面一个
	if err != nil {
		log.Fatalf("did not connect :%v", err)
	}
	defer conn.Close()

	// Make a echo client and send rpc
	rgc := pb.NewEchoClient(conn)
	callUnaryEcho(rgc, "hello world")
}
