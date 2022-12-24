package rpc

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RequestService struct{}

func (RequestService) Hello(ctx context.Context, empty *emptypb.Empty) (*proto.HelloResponse, error) {
	panic("implement me")
}
