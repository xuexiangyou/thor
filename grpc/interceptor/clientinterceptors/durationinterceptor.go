package clientinterceptors

import (
	"context"
	"fmt"
	"path"
	"time"

	"google.golang.org/grpc"
)

var initTime = time.Now().AddDate(-1, -1, -1)

const slowThreshold = time.Millisecond * 500

func DurationInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	serverName := path.Join(cc.Target(), method)
	fmt.Println("哈哈哈" + cc.Target() + "/" + method)
	start := time.Since(initTime)
	fmt.Println("client duration拦截器前1") // todo 多个拦截器的话排在前面的invoker前先执行
	err := invoker(ctx, method, req, reply, cc, opts...)
	fmt.Println("client duration拦截器后2") // todo 多个拦截器的话排在前端的invoker后后执行 invoker 后面的逻辑需要server 返回结果
	if err != nil {
		fmt.Printf("client duration fail - %s - %v - %s", serverName, req, err.Error())
		fmt.Println()
	} else {
		elapsed := time.Since(initTime) - start
		if elapsed > slowThreshold {
			fmt.Printf("[RPC] ok - slowcall - %s - %v - %v", serverName, req, reply)
			fmt.Println()
		}
	}
	return err
}
