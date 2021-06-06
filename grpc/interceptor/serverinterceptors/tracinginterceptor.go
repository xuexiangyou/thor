package serverinterceptors

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

func UnaryTracingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		fmt.Println("trace 拦截器前奏")
		resp, err = handler(ctx, req)
		fmt.Println("trace 拦截器后奏")
		return resp, err
		// return handler(ctx, req)
	}
}