package bootstrap

import (
	"github.com/gin-gonic/gin"
	"leopard-quant/core/config"
	"leopard-quant/core/event"
	"leopard-quant/core/log"
	"leopard-quant/restful/controller"
)

var Global *globalContent

type globalContent struct {
	ApplicationConf config.ApplicationConfig
	EventEngine     *event.Engine
	RestfulEngine   *gin.Engine
}

func Init() {
	Global = &globalContent{}

	//初始化配置文件
	initApplicationConfig()
	//初始化日志
	initLog(&Global.ApplicationConf)
	//初始化restful
	initRestfulEngine(&Global.ApplicationConf)
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
	engine := Global.RestfulEngine
	engine.Use(gin.Recovery())
	engine.Use(controller.Logger())
	controller.Validator()
	controller.Router(engine)
}
