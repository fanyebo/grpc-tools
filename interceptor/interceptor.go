package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

// BeforeInterceptor 调用前拦截器
type BeforeInterceptor interface {
	HandleFunc(ctx *context.Context, req interface{}, info *grpc.UnaryServerInfo) error
}

// AfterInterceptor 调用后拦截器
type AfterInterceptor interface {
	HandleFunc(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, resp interface{}, err error) (interface{}, error)
}
