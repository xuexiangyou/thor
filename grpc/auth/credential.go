package auth

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

type Credential struct {
	App string
	Token string
}

func (c *Credential) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		appKey:   c.App,
		tokenKey: c.Token,
	}, nil
}

// 表示通讯底层是否要使用安全链接true-是、false-否
func (c *Credential) RequireTransportSecurity() bool {
	return false
}

// ParseCredential parses credential from given ctx.
func ParseCredential(ctx context.Context) Credential {
	var credential Credential

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return credential
	}

	apps, tokens := md[appKey], md[tokenKey]
	if len(apps) == 0 || len(tokens) == 0 {
		return credential
	}

	app, token := apps[0], tokens[0]
	if len(app) == 0 || len(token) == 0 {
		return credential
	}

	credential.App = app
	credential.Token = token

	return credential
}

func EnsureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	fmt.Println("server 拦截器执行")
	md := ParseCredential(ctx)
	if md.App != "testApp" || md.Token != "testToken" {
		return nil, errInvalidToken
	}
	return handler(ctx, req)
}