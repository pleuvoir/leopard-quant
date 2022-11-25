package engine

import (
	"leopard-quant/common/model"
	"leopard-quant/core/event"
	"leopard-quant/core/log"
)

const DefaultOrderEngineName = "leopard-order-engine"

type OrderEngine struct {
	name       string
	mainEngine *MainEngine
	tradeMap   map[string]model.Trade
}

// NewOrderEngine 构建订单引擎
func NewOrderEngine(mainEngine *MainEngine) *OrderEngine {
	engine := OrderEngine{mainEngine: mainEngine}
	engine.tradeMap = make(map[string]model.Trade)
	engine.initEngine()
	return &engine
}

func (o *OrderEngine) initEngine() {
	o.name = DefaultOrderEngineName
	o.registerEvent()
}

func (o *OrderEngine) registerEvent() {
	o.addListener(event.Trade, func(event event.Event) {
		trade, _ := event.EventData.(model.Trade)
		o.tradeMap[trade.Id] = trade
		log.Infof("[%d]onEvent. event=%+v", event.EventType, event)
	})
}

func (o *OrderEngine) addListener(t event.Type, f func(e event.Event)) {
	o.mainEngine.RegisterListener(t, f)
}

func (o *OrderEngine) Name() string {
	return o.name
}

func (o *OrderEngine) Start() {
}

func (o *OrderEngine) Stop() {
}
