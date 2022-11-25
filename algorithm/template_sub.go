package algorithm

import (
	"leopard-quant/common/model"
	"leopard-quant/core/config"
)

type TemplateSub interface {
	OnStart(c Context)
	OnStop(c Context)
	OnTimer(c Context)
	OnTrade(c Context, trade model.Trade)
	OnTick(c Context, ticker model.Ticker)
	OnOrder(c Context, order model.Order)
	Name() string
}

type Context struct {
	config.Loader
	Subscribe
}

type Subscribe func(symbol string) error
