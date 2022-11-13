package core

import (
	"leopard-quant/core/event"
	"leopard-quant/core/model"
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

func (g *Gateway) onEvent(eventType event.Type, data any) {
	g.EventEngine.Put(event.NewEvent(eventType, data))
}
