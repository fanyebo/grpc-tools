package interceptor

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"google.golang.org/grpc"
)

// BeforeInterceptor 调用前拦截
type BeforeInterceptor func(ctx *context.Context, req interface{}, info *grpc.UnaryServerInfo) error

// AfterInterceptor 调用后拦截
type AfterInterceptor func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, resp interface{}, err error) (interface{}, error)

type UnaryServerInterceptorFactory struct {
	BeforeInterceptors []BeforeInterceptor
	AfterInterceptors  []AfterInterceptor
}

func (i UnaryServerInterceptorFactory) Gen() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// recover处理
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 1<<16)
				runtime.Stack(buf, false)
				log.Println(string(buf))
			}
		}()

		var resp interface{}
		var err error
		// 接口调用前拦截器
		for _, bi := range i.BeforeInterceptors {
			err = bi(&ctx, req, info)
			if err != nil {
				log.Panic(err)
				return nil, err
			}
		}
		startTime := time.Now()
		resp, err = handler(ctx, req)
		log.Println(fmt.Sprintf("接口请求时间: %d ms", time.Since(startTime).Milliseconds()))
		// 接口调用后处理
		for _, ai := range i.AfterInterceptors {
			resp, err = ai(ctx, req, info, resp, err)
		}
		return resp, err
	}
}
