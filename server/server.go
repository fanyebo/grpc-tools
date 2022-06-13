package server

import (
	"fmt"
	"github.com/fanyebo/grpc-tools/interceptor"
	"google.golang.org/grpc"
	"log"
	"net"
)

type GrpcUnaryService struct {
	Desc            *grpc.ServiceDesc
	ServiceProvider interface{}
}

type GrpcUnaryServer struct {
	svr         *grpc.Server
	Listener    *net.Listener
	ServerOpts  []grpc.ServerOption
	Logger      *log.Logger
	Services    []GrpcUnaryService
	Port        int64
	Interceptor *interceptor.UnaryServerInterceptorFactory
}

// RegisterService 注册服务
func (s *GrpcUnaryServer) RegisterService(sd *grpc.ServiceDesc, ss interface{}) {
	s.Services = append(s.Services, GrpcUnaryService{Desc: sd, ServiceProvider: ss})
}

func (s *GrpcUnaryServer) init() {
	if s.Logger == nil {
		s.Logger = log.Default()
	}
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

	// 拦截器设置
	if s.Interceptor == nil {
		s.Interceptor = new(interceptor.UnaryServerInterceptorFactory)
	}
}

// Start 启动服务
func (s *GrpcUnaryServer) Start() {
	s.init()
	s.ServerOpts = append(s.ServerOpts, grpc.UnaryInterceptor(s.Interceptor.Gen()))
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
