package rpc

import (
	"context"
	socketIO "github.com/ambelovsky/gosf-socketio"
	"google.golang.org/protobuf/types/known/emptypb"
	"leopard-quant/common/model"
	"leopard-quant/rpc/pb"
)

type PushService struct {
	sio *socketIO.Server
}

func NewPushService(sio *socketIO.Server) *PushService {
	return &PushService{sio: sio}
}

func (p PushService) push(method string, data any) {
	p.sio.BroadcastToAll("push", model.NewSuccess(method, data))
}

func (p PushService) UpdateCount(ctx context.Context, request *pb.UpdateCountRequest) (*emptypb.Empty, error) {
	p.push("UpdateCount", request)
	return &emptypb.Empty{}, nil
}
