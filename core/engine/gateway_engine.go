package engine

import (
	"fmt"
	"github.com/gookit/color"
	"leopard-quant/common/model"
	"leopard-quant/core/event"
	"leopard-quant/gateway"
	"sync"
)

type GatewayEngine struct {
	GatewayName string
	EventEngine *event.Engine
	eventPool   sync.Pool //对象池避免创建过多对象
	sub         gateway.Sub
}

func NewGateway(eventEngine *event.Engine, sub gateway.Sub) *GatewayEngine {
	return &GatewayEngine{GatewayName: sub.Name(), EventEngine: eventEngine,
		eventPool: sync.Pool{New: func() any {
			return *new(event.Event)
		}},
		sub: sub}
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

func (g *GatewayEngine) OnTick(tick model.Ticker) {
	g.onEvent(event.Tick, tick)
}

func (g *GatewayEngine) OnBar(kline model.KLine) {
	g.onEvent(event.Bar, kline)
}

func (g *GatewayEngine) OnOrder(order model.Order) {
	g.onEvent(event.Order, order)
}

func (g *GatewayEngine) OnTrade(trade model.Trade) {
	g.onEvent(event.Trade, trade)
}

func (g *GatewayEngine) onEvent(eventType event.Type, data any) {
	e := g.eventPool.Get().(event.Event)
	e.EventType = eventType
	e.EventData = data
	g.EventEngine.Put(e)
	g.eventPool.Put(e)
}

// Subscribe 订阅某币种的所有事件
func (g *GatewayEngine) Subscribe(symbol string) error {
	callback := gateway.ApiCallback{
		TickerCallback: func(tick model.Ticker) {
			g.OnTick(tick)
		},
		KlineCallback: func(k model.KLine) {
			g.OnBar(k)
		}}
	return g.sub.Subscribe(symbol, callback)
}

// CancelSubscribe 取消订阅某币种的所有事件
func (g *GatewayEngine) CancelSubscribe(symbol string) error {
	return g.sub.CancelSubscribe(symbol)
}
