package socket

import (
	"fmt"
	socketIO "github.com/ambelovsky/gosf-socketio"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"leopard-quant/common/model"
	"testing"
	"time"
)

func TestInstallSocketIO(t *testing.T) {

	g := gin.Default()
	g.Use(cors.New(cors.Config{
		AllowAllOrigins:        true,
		AllowMethods:           []string{"POST"},
		AllowHeaders:           []string{"*"},
		AllowCredentials:       false,
		ExposeHeaders:          nil,
		MaxAge:                 12 * time.Hour,
		AllowWildcard:          false,
		AllowBrowserExtensions: false,
		AllowWebSockets:        false,
		AllowFiles:             false,
	}))

	InstallSocketIO(g, func(c *socketIO.Channel, request model.RequestMessage) model.ResultMessage {
		return model.NewSuccess("hello", "world")
	})

	server := GetInstance()

	go func(server *socketIO.Server) {
		var count int64 = 0
		for count < 500 {
			count++
			server.BroadcastToAll("push", model.NewSuccess("methodName", "i am push data"))
			time.Sleep(time.Second)
		}
	}(server)

	t.Log(fmt.Sprintf("完成"))

	if err := g.Run(":8000"); err != nil {
		panic(err)
	}

}
