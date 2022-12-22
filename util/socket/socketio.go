package socket

import (
	socketIO "github.com/ambelovsky/gosf-socketio"
	"github.com/ambelovsky/gosf-socketio/transport"
	"github.com/gin-gonic/gin"
	"leopard-quant/common/model"
	"sync"
)

var (
	once sync.Once
	sio  *socketIO.Server
)

type requestHandlerFunc func(c *socketIO.Channel, request model.RequestMessage) model.ResultMessage

func InstallSocketIO(g *gin.Engine, requestHandler requestHandlerFunc) {
	once.Do(func() {
		sio = newSocketIO(g, requestHandler)
	})
}

func newSocketIO(g *gin.Engine, requestHandler requestHandlerFunc) *socketIO.Server {
	//使用HTTP代理socketIO
	server := socketIO.NewServer(transport.GetDefaultWebsocketTransport())
	g.Any("/socket.io/*any", gin.WrapH(server))
	_ = server.On(socketIO.OnConnection, func(c *socketIO.Channel) {
	})
	//监听request emit
	if err := server.On("request", requestHandler); err != nil {
		panic(err)
	}
	return server
}

func GetInstance() *socketIO.Server {
	if sio == nil {
		panic("socketIO未初始化，请先调用InstallSocketIO方法")
	}
	return sio
}
