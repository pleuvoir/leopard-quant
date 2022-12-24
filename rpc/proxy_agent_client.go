package rpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"leopard-quant/rpc/proto"
)

func Test() {

	conn, err := grpc.Dial("127.0.0.1:8888", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	client := proto.NewHelloClient(conn)

	rsp, err := client.SayHello(context.Background(), &proto.HelloRequest{Name: "pleuvoir"})

	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	fmt.Println(rsp)
}
