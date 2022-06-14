package main

import (
	"context"
	"fmt"
	"github.com/fanyebo/grpc-tools/client"
	"github.com/fanyebo/grpc-tools/demo/pb/hello"
	"google.golang.org/grpc"
	"log"
)

func main() {
	addr := fmt.Sprintf("%s:%d", "127.0.0.1", 80)
	stub, err := client.CreateClient[hello.HelloClient](hello.NewHelloClient, addr, []grpc.DialOption{})
	resp, err := stub.SayHello(context.Background(), &hello.SayHelloReq{Name: "test"})
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(resp.Reply)
}
