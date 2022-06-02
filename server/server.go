package server

import (
	"context"
	"fmt"
	"github.com/fanyebo/grpc-tools/interceptor"
	"log"
	"net"
	"runtime"
	"time"

	"google.golang.org/grpc"
)

type GrpcUnaryService struct {
	Desc            *grpc.ServiceDesc
	ServiceProvider interface{}
}

type GrpcUnaryServer struct {
	svr                *grpc.Server
	Listener           *net.Listener
	BeforeInterceptors []interceptor.BeforeInterceptor
	AfterInterceptors  []interceptor.AfterInterceptor
	ServerOpts         []grpc.ServerOption
	Logger             *log.Logger
	Services           []GrpcUnaryService
	Port               int64
}

// RegisterService 注册服务
func (s *GrpcUnaryServer) RegisterService(sd *grpc.ServiceDesc, ss interface{}) {
	s.Services = append(s.Services, GrpcUnaryService{Desc: sd, ServiceProvider: ss})
}

func (s *GrpcUnaryServer) init() {
	if s.Logger == nil {
		s.Logger = log.Default()
	}
	// 端口未设置默认为80
	if s.Port == 0 {
		s.Logger.Panic("端口未设置")
	}

	// 初始化监听器
	if s.Listener == nil {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
		if err != nil {
			s.Logger.Panic(err)
		}
		s.Listener = &listener
	}
}

// Start 启动服务
func (s *GrpcUnaryServer) Start() {
	s.init()
	// 拦截器
	_interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// recover处理
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 1<<16)
				runtime.Stack(buf, false)
				s.Logger.Println(string(buf))
			}
		}()

		var resp interface{}
		var err error
		// 接口调用前拦截器
		for _, i := range s.BeforeInterceptors {
			err = i.HandleFunc(&ctx, req, info)
			if err != nil {
				s.Logger.Panic(err)
				return nil, err
			}
		}
		startTime := time.Now()
		resp, err = handler(ctx, req)
		s.Logger.Println(fmt.Sprintf("接口请求时间: %d ms", time.Since(startTime).Milliseconds()))
		// 接口调用后处理
		for _, i := range s.AfterInterceptors {
			resp, err = i.HandleFunc(ctx, req, info, resp, err)
		}
		return resp, err
	}

	s.ServerOpts = append(s.ServerOpts, grpc.UnaryInterceptor(_interceptor))
	s.svr = grpc.NewServer(s.ServerOpts...)

	// 注册service到server
	for _, service := range s.Services {
		s.svr.RegisterService(service.Desc, service.ServiceProvider)
	}

	s.Logger.Println(fmt.Sprintf("启动服务[0.0.0.0:%d]...", s.Port))
	if err := s.svr.Serve(*s.Listener); err != nil {
		panic(fmt.Sprintf("服务启动失败: %s", err))
	}
}
