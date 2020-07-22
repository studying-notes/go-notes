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
