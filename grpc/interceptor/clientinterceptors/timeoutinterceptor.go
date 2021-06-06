package clientinterceptors

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

func TimeoutInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if timeout <= 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		ctx, cancel := context.WithDeadline(ctx, time.Now().Add(timeout))
		defer cancel()

		done := make(chan error, 1)
		panicChan := make(chan interface{}, 1)
		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			fmt.Println("client timeout拦截器前2")
			done <- invoker(ctx, method, req, reply, cc, opts...)
			fmt.Println("client timeout拦截器后1") // todo client 拦截器监听到server端返回的结果了、后面的拦截器invoker后的逻辑还会执行
		}()

		select {
		case p := <- panicChan:
			panic(p)
		case err := <-done:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
