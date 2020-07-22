# gRPC

## RPC 是什么

在分布式计算，远程过程调用（Remote Procedure Call，缩写为 RPC）是一个计算机通信协议。该协议允许运行于一台计算机的程序调用另一个地址空间（通常为一个开放网络的一台计算机）的子程序，而程序员就像调用本地程序一样，无需额外地为这个交互作用编程（无需关注细节）。

RPC是一种服务器/客户端（Client/Server）模式，经典实现是一个通过 `发送请求-接受回应` 进行信息交互的系统。

## gRPC 是什么

`gRPC` 是一种现代化开源的高性能 RPC 框架，能够运行于任意环境之中。最初由谷歌进行开发。它使用 HTTP/2 作为传输协议。

在 gRPC 里，客户端可以像调用本地方法一样直接调用其他机器上的服务端应用程序的方法，帮助你更容易创建分布式应用程序和服务。与许多 RPC 系统一样，gRPC 是基于定义一个服务，指定一个可以远程调用的带有参数和返回类型的的方法。在服务端程序中实现这个接口并且运行 gRPC 服务处理客户端调用。在客户端，有一个 Stub 提供和服务端相同的方法。

![](grpc.svg)

## 为什么要用 gRPC

使用 gRPC， 我们可以一次性的在一个 `.proto` 文件中定义服务并使用任何支持它的语言去实现客户端和服务端，反过来，它们可以应用在各种场景中，从 Google 的服务器到平板电脑，gRPC 解决了不同语言及环境间通信的复杂性。使用 `protocol buffers` 还能获得其他好处，包括高效的序列号，简单的 IDL 以及容易进行接口更新。总之，使用 gRPC 能让我们更容易编写跨语言的分布式代码。

## 安装 gRPC

```bash
go get -u google.golang.org/grpc
go get -u github.com/golang/protobuf/protoc-gen-go
```

安装用于生成 gRPC 服务代码的协议编译器，见 [初识 Protobuf](../../storage/protobuf/README.md)。

## gRPC 开发分三步

1. 编写 `.proto` 文件，生成指定语言源代码
2. 编写服务端代码
3. 编写客户端代码

## gRPC 入门示例

### 编写 proto 代码

```bash
syntax = "proto3";

package main;

// 定义一个打招呼服务
service Greeter {
    // SayHello 方法
    rpc SayHello (HelloRequest) returns (HelloReply) {}
}

// 包含人名的一个请求消息
message HelloRequest {
    string name = 1;
}

// 包含问候语的响应消息
message HelloReply {
    string message = 1;
}
```

在当前工作区生成文件：

```bash
protoc -I . ./helloworld.proto --go_out=plugins=grpc:.
```

### 编写 Server 端 Go 代码

```go
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
```

### 编写 Client 端 Go 代码

```go
package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":9876", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	cli := NewGreeterClient(conn)
	r, err := cli.SayHello(context.Background(), &HelloRequest{Name: "World"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Greeting: %s\n", r.Message)
}
```

**Server**

```shell
go run server.go helloworld.pb.go  
```

**Client**

```shell
go run client.go helloworld.pb.go 
```
