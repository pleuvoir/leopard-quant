package bootstrap

import (
	socketIO "github.com/ambelovsky/gosf-socketio"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"leopard-quant/algorithm/base"
	"leopard-quant/common/model"
	"leopard-quant/core/config"
	"leopard-quant/core/engine"
	"leopard-quant/core/event"
	"leopard-quant/core/log"
	"leopard-quant/restful/controller"
	"leopard-quant/rpc"
	"leopard-quant/util/socket"
	"time"
)

var Global *globalContent

type globalContent struct {
	ApplicationConf config.ApplicationConfig
	MainEngine      *engine.MainEngine
	RestfulEngine   *gin.Engine
	PushService     *rpc.PushService
}

func Init() {
	Global = &globalContent{}
	//初始化配置文件
	initApplicationConfig()
	//初始化日志
	initLog(&Global.ApplicationConf)
	//初始化主引擎
	initMainEngine(&Global.ApplicationConf)
	//初始化算法
	initAlgoEngine(&Global.ApplicationConf)
	//初始化restful
	initRestfulEngine(&Global.ApplicationConf)
	//初始化socketIO
	initSocketIO()
}

func initAlgoEngine(app *config.ApplicationConfig) {
	e := base.NewAlgoEngine(Global.MainEngine, app.Algo)
	e.Start()
}

func initMainEngine(app *config.ApplicationConfig) {
	mainEngine := engine.NewMainEngine(event.NewEventEngine(), app.Main)
	mainEngine.InitEngines()
	mainEngine.Start()
	Global.MainEngine = mainEngine
}

func initApplicationConfig() {
	conf := config.NewApplicationConf()
	if err := conf.Load(); err != nil {
		panic(err)
	}
	Global.ApplicationConf = conf
}

func initLog(app *config.ApplicationConfig) {
	log.InitLog(app)
}

func initRestfulEngine(app *config.ApplicationConfig) {
	controller.SetMode(app)
	Global.RestfulEngine = gin.New()
	e := Global.RestfulEngine
	e.Use(gin.Recovery())
	e.Use(controller.Logger())
	controller.Validator()
	controller.Router(e)
}

func initSocketIO() {
	e := Global.RestfulEngine
	e.Use(cors.New(cors.Config{
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

	//处理request请求
	requestService := rpc.RequestService{}
	socket.InstallSocketIO(e, func(c *socketIO.Channel, request model.RequestMessage) model.ResultMessage {
		return rpc.RequestHandlerProxy(requestService, request)
	})

	//处理push请求
	pushService := rpc.NewPushService(socket.GetInstance())
	Global.PushService = pushService
}
