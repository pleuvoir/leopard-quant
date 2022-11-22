package base

import (
	"leopard-quant/algorithm/impl"
	"leopard-quant/core/config"
	"leopard-quant/core/engine"
	"leopard-quant/core/engine/model"
)

type AlgoTemplate struct {
	sub          impl.TemplateSub
	algoName     string
	engine       *engine.MainEngine
	active       bool
	activeOrders map[string]model.Order
	ticks        map[string]model.Tick
	config       config.Loader
	context      impl.Context
}

// NewAlgoTemplate 创建算法模板
// 可以认为是基类，算法引擎回调的是这个类 sub负责子类逻辑的实现
func NewAlgoTemplate(engine *engine.MainEngine, sub impl.TemplateSub) *AlgoTemplate {
	t := AlgoTemplate{}
	t.sub = sub
	t.algoName = sub.Name()
	t.engine = engine
	t.activeOrders = make(map[string]model.Order)
	t.ticks = make(map[string]model.Tick)
	t.context = impl.Context{MainEngine: engine}
	return &t
}

func (t *AlgoTemplate) updateTick(tick model.Tick) {
	if t.active {
		history, ok := t.ticks[tick.Symbol]
		if ok {
			history.UpdateTick(tick)
			t.sub.OnTick(history)
		} else {
			t.ticks[tick.Symbol] = tick
			t.sub.OnTick(tick)
		}
	}
}

func (t *AlgoTemplate) updateOrder(order model.Order) {
	if t.active {
		if order.IsActive() {
			t.activeOrders[order.Id] = order
		}
		t.sub.OnOrder(order)
	}
}

func (t *AlgoTemplate) updateTrade(trade model.Trade) {
	if t.active {
		t.sub.OnTrade(trade)
	}
}

func (t *AlgoTemplate) init() {

}

func (t *AlgoTemplate) start() {
	context := impl.Context{MainEngine: t.engine}
	t.sub.OnStart(context)
	t.active = true //注意在后面 否则会接受到事件
}

func (t *AlgoTemplate) stop() {
	t.sub.OnStop(t.context)
	defer func() {
		t.active = false
	}()
}

func (t *AlgoTemplate) updateTimer() {
	if t.active {
		t.sub.OnTimer()
	}
}
