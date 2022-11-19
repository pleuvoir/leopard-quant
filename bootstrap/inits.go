package bootstrap

import (
	"github.com/gin-gonic/gin"
	"leopard-quant/algorithm/base"
	"leopard-quant/core/config"
	"leopard-quant/core/engine"
	"leopard-quant/core/event"
	"leopard-quant/core/log"
	"leopard-quant/restful/controller"
)

var Global *globalContent

type globalContent struct {
	ApplicationConf config.ApplicationConfig
	MainEngine      *engine.MainEngine
	RestfulEngine   *gin.Engine
}

func Init() {
	Global = &globalContent{}
	//初始化配置文件
	initApplicationConfig()
	//初始化日志
	initLog(&Global.ApplicationConf)
	//初始化主引擎
	initMainEngine()
	//初始化算法
	initAlgoEngine()
	//初始化restful
	initRestfulEngine(&Global.ApplicationConf)
}

func initAlgoEngine() {
	e := base.NewAlgoEngine(Global.MainEngine)
	e.Start()
}

func initMainEngine() {
	mainEngine := engine.NewMainEngine(event.NewEventEngine())
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
