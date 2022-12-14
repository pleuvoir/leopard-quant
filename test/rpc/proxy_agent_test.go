package rpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"leopard-quant/rpc/proto"
	"testing"
)

func TestHello(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:8888", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	client := proto.NewAgentClient(conn)

	subscribeRequest := &proto.SubscribeRequest{
		Symbol:   "BTC-USDT",
		Type:     "kline",
		Exchange: "okx",
	}
	rsp, err := client.Subscribe(context.Background(), subscribeRequest)

	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	fmt.Println(rsp)
}
