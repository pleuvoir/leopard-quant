package rpc

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"leopard-quant/rpc/pb"
)

type RequestService struct{}

func (RequestService) Hello(ctx context.Context, empty *emptypb.Empty) (*pb.HelloResponse, error) {
	panic("implement me")
}
