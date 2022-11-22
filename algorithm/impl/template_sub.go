package impl

import (
	"leopard-quant/core/engine"
	"leopard-quant/core/engine/model"
)

type TemplateSub interface {
	OnStart(c Context)
	OnStop(c Context)
	OnTimer()
	OnTrade(trade model.Trade)
	OnTick(t model.Tick)
	OnOrder(order model.Order)
	Name() string
}

type Context struct {
	MainEngine *engine.MainEngine
}
