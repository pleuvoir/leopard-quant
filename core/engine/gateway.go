package engine

import (
	"leopard-quant/core/engine/model"
	"leopard-quant/core/event"
)

type Gateway struct {
	GatewayName string
	EventEngine *event.Engine
}

func NewGateway(GatewayName string, eventEngine *event.Engine) *Gateway {
	return &Gateway{GatewayName: GatewayName, EventEngine: eventEngine}
}

func (g *Gateway) OnTick(tick model.Tick) {
	g.onEvent(event.Tick, tick)
}

func (g *Gateway) OnOrder(order model.Order) {
	g.onEvent(event.Tick, order)
}

func (g *Gateway) OnTrade(trade model.Trade) {
	g.onEvent(event.Tick, trade)
}

func (g *Gateway) onEvent(eventType event.Type, data any) {
	g.EventEngine.Put(event.NewEvent(eventType, data))
}
