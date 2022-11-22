package engine

import (
	"fmt"
	"github.com/gookit/color"
	"leopard-quant/core/engine/model"
	"leopard-quant/core/event"
	"sync"
)

type GatewayEngine struct {
	GatewayName string
	EventEngine *event.Engine
	eventPool   sync.Pool //对象池避免创建过多对象
}

func NewGateway(GatewayName string, eventEngine *event.Engine) *GatewayEngine {
	return &GatewayEngine{GatewayName: GatewayName, EventEngine: eventEngine,
		eventPool: sync.Pool{New: func() any {
			return new(event.Event)
		}}}
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
	e := g.eventPool.Get().(event.Event)
	e.EventType = eventType
	e.EventData = data
	g.EventEngine.Put(e)
	g.eventPool.Put(e)
}
