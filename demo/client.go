package main

import (
	"context"
	"fmt"
	"github.com/fanyebo/grpc-tools/demo/pb/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

// 拦截器
func unaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...) // invoking RPC method
	return err
}

func main() {
	addr := fmt.Sprintf("%s:%d", "127.0.0.1", 80)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(unaryInterceptor))
	if err != nil {
		log.Panic(fmt.Sprintf("did not connect: %+v", err))
	}

	stub := hello.NewHelloClient(conn)
	resp, err := stub.SayHello(context.Background(), &hello.SayHelloReq{Name: "test"})
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(resp.Reply)
}
