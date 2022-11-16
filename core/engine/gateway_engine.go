package engine

import (
	"fmt"
	"github.com/gookit/color"
	"leopard-quant/core/engine/model"
	"leopard-quant/core/event"
)

type GatewayEngine struct {
	GatewayName string
	EventEngine *event.Engine
}

func NewGateway(GatewayName string, eventEngine *event.Engine) *GatewayEngine {
	return &GatewayEngine{GatewayName: GatewayName, EventEngine: eventEngine}
}

func (g *GatewayEngine) Name() string {
	return g.GatewayName
}

func (g *GatewayEngine) Start() {
	color.Greenln(fmt.Sprintf("[%s]网关引擎已启动", g.Name()))
}

func (g *GatewayEngine) Stop() {
	color.Redln(fmt.Sprintf("[%s]网关引擎关闭状态", g.Name()))
}

func (g *GatewayEngine) OnTick(tick model.Tick) {
	g.onEvent(event.Tick, tick)
}

func (g *GatewayEngine) OnOrder(order model.Order) {
	g.onEvent(event.Tick, order)
}

func (g *GatewayEngine) OnTrade(trade model.Trade) {
	g.onEvent(event.Tick, trade)
}

func (g *GatewayEngine) onEvent(eventType event.Type, data any) {
	g.EventEngine.Put(event.NewEvent(eventType, data))
}
