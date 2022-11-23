package engine

import (
	"fmt"
	"github.com/gookit/color"
	"leopard-quant/core/engine/model"
	"leopard-quant/core/event"
	"leopard-quant/core/log"
	"leopard-quant/gateway"
	"sync"
)

type GatewayEngine struct {
	GatewayName string
	EventEngine *event.Engine
	eventPool   sync.Pool //对象池避免创建过多对象
	sub         GatewaySub
}

func NewGateway(GatewayName string, eventEngine *event.Engine, sub GatewaySub) *GatewayEngine {
	return &GatewayEngine{GatewayName: GatewayName, EventEngine: eventEngine,
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

func (g *GatewayEngine) Connect() error {
	err := g.sub.Connect()
	if err != nil {
		color.Greenln("网关连接成功。")
	}
	return err
}

// Subscribe 订阅某币种的所有事件
func (g *GatewayEngine) Subscribe(symbol string) error {
	callback := gateway.ComposeCallback{
		TickerCallback: func(tick gateway.Ticker) {
			log.Infof("网关引擎发布事件。")
			m := model.Tick{Symbol: "STARL-USDT"} //TODO
			g.OnTick(m)
		},
		KlineCallback: func(k gateway.Kline) {

		}}
	err := g.sub.Connect()
	if err != nil {
		color.Errorf("网关连接失败。")
		return err
	}
	return g.sub.Subscribe(symbol, callback)
}

// CancelSubscribe 取消订阅某币种的所有事件
func (g *GatewayEngine) CancelSubscribe(symbol string) error {
	return g.sub.CancelSubscribe(symbol)
}

type GatewaySub interface {
	Connect() error
	Subscribe(symbol string, c gateway.ComposeCallback) error
	CancelSubscribe(symbol string) error
}
