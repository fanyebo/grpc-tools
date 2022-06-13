package main

import (
	"context"
	"fmt"
	"github.com/fanyebo/grpc-tools/demo/pb/hello"
	"github.com/fanyebo/grpc-tools/interceptor"
	"github.com/fanyebo/grpc-tools/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ServiceHello struct {
	hello.UnimplementedHelloServer
}

func (s *ServiceHello) SayHello(ctx context.Context, req *hello.SayHelloReq) (*hello.SayHelloResp, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		fmt.Println("metadata: ", md)
	}
	return &hello.SayHelloResp{Reply: fmt.Sprintf("Hello %s", req.Name)}, nil
}

func HelloBeforeInterceptor(ctx *context.Context, req interface{}, info *grpc.UnaryServerInfo) error {
	md, ok := metadata.FromIncomingContext(*ctx)
	if !ok {
		md = metadata.New(map[string]string{})
	}

	fmt.Println(fmt.Sprintf("请求信息ctx=%+v, req=%+v, info:%+v", md, req, info))
	return nil
}

func HellAfterInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, resp interface{}, err error) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(map[string]string{})
	}
	fmt.Println(md.Get(":authority"))
	fmt.Println(fmt.Sprintf("请求结果ctx=%s, req=%+v, info=%+v, resp=%+v, err=%+v", md, req, info, resp, err))
	return resp, err
}

func main() {
	svr := server.GrpcUnaryServer{Port: 80}
	svr.Interceptor = &interceptor.UnaryServerInterceptorFactory{
		BeforeInterceptors: []interceptor.BeforeInterceptor{HelloBeforeInterceptor},
		AfterInterceptors:  []interceptor.AfterInterceptor{HellAfterInterceptor},
	}
	service := new(ServiceHello)
	svr.RegisterService(&hello.Hello_ServiceDesc, service)
	svr.Start()
}
