package client

import (
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

// CreateClient 创建grpc客户端
func CreateClient[grpcClientInterface any](createFunc func(cc grpc.ClientConnInterface) grpcClientInterface, addr string, options []grpc.DialOption) (grpcClientInterface, error) {
	var grpcClient grpcClientInterface

	if len(options) == 0 {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	conn, err := grpc.Dial(addr, options...)
	if err != nil {
		log.Panic(err)
		return grpcClient, err
	}

	_client := createFunc(conn)
	client, ok := any(_client).(grpcClientInterface)
	if ok {
		return client, nil
	}
	return grpcClient, errors.New("createFunc return type of client error")
}
