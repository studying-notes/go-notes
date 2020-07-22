package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: "Hello" + in.Name}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":9876")
	if err != nil {
		panic(err)
	}
	// 创建 gRPC 服务器
	srv := grpc.NewServer()
	//注册服务
	RegisterGreeterServer(srv, &server{})
	// 在给定的 gRPC 服务器上注册服务器反射服务
	reflection.Register(srv)

	// Serve 方法在 listener 上接受传入连接，为每个连接
	// 创建一个 ServerTransport 和 server 的 goroutine
	// 该 goroutine 读取 gRPC 请求，然后调用已注册的
	// 处理程序来响应它们
	err = srv.Serve(listener)
	if err != nil {
		panic(err)
	}
}
