package clientinterceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TracingInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	var pairs  []string
	pairs = append(pairs, "oneKey", "123232")
	ctx = metadata.AppendToOutgoingContext(ctx, pairs...) // todo 自己随便造点数据来演示grpc 拦截器中传输context
	return invoker(ctx, method, req, reply, cc, opts...)
}
