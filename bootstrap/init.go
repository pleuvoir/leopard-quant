package bootstrap

import (
	"leopard-quant/core/config"
	"leopard-quant/core/event"
	"leopard-quant/core/log"
)

var Global *globalContent

type globalContent struct {
	ApplicationConf config.ApplicationConf
	engine          *event.Engine
}

func init() {
	Global = &globalContent{}

	//初始化配置文件
	initApplicationConf()
	//初始化日志
	initLog(&Global.ApplicationConf)
}

func initApplicationConf() {
	conf := config.NewApplicationConf()
	if err := conf.Load(); err != nil {
		panic(err)
	}
	Global.ApplicationConf = conf
}

func initLog(app *config.ApplicationConf) {
	log.InitLog(app)
}
