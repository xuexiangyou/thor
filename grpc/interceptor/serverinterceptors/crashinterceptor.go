package serverinterceptors

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryCrashInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer handleCrash(func(r interface{}) {
			err = toPanicError(r) // todo 这个会赋值给handler中返回的err
		})
		fmt.Println("crash 拦截器前奏")
		resp, err = handler(ctx, req)
		fmt.Println("crash 拦截器后奏")
		return resp, err // todo 如果先执行的server 拦截器返回err了 后面的拦截器就走不到后面的handler
	}
}

func handleCrash(handler func(interface{})) {
	if r := recover(); r != nil {
		handler(r)
	}
}

func toPanicError(r interface{}) error {
	// logx.Errorf("%+v %s", r, debug.Stack())
	return status.Errorf(codes.Internal, "panic test: %v", r)
}
