package algorithm

import (
	"leopard-quant/core/config"
	"leopard-quant/core/engine/model"
)

type TemplateSub interface {
	OnStart(c Context)
	OnStop(c Context)
	OnTimer(c Context)
	OnTrade(c Context, trade model.Trade)
	OnTick(c Context, t model.Tick)
	OnOrder(c Context, order model.Order)
	Name() string
}

type Context struct {
	config.Loader
	Subscribe
}

type Subscribe func(symbol string) error
