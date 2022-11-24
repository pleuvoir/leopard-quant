package algorithm

import (
	"leopard-quant/core/config"
	"leopard-quant/core/engine/model"
)

type AlgoTemplate struct {
	engine       *AlgoEngine
	sub          TemplateSub
	algoName     string
	active       bool
	activeOrders map[string]model.Order
	ticks        map[string]model.Tick
	config       config.Loader
	context      Context
}

// NewAlgoTemplate 创建算法模板
// 可以认为是基类，算法引擎回调的是这个类 sub负责子类逻辑的实现
func NewAlgoTemplate(engine *AlgoEngine, sub TemplateSub, c config.Loader) *AlgoTemplate {
	t := AlgoTemplate{}
	t.engine = engine
	t.sub = sub
	t.config = c
	t.algoName = sub.Name()
	t.activeOrders = make(map[string]model.Order)
	t.ticks = make(map[string]model.Tick)

	//初始化上下文，将一些操作封装在里面，避免子类过度调用模板的方法
	t.context = Context{Loader: t.config,
		Subscribe: func(symbol string) error {
			return t.Subscribe(symbol)
		}}
	return &t
}

func (t *AlgoTemplate) updateTick(tick model.Tick) { //TODO
	if t.active {
		//history, ok := t.ticks[tick.Symbol]
		//if ok {
		//	history.UpdateTick(tick)
		//	t.sub.OnTick(t.context, history)
		//} else {
		//	t.ticks[tick.Symbol] = tick
		t.sub.OnTick(t.context, tick)
		//	}
	}
}

func (t *AlgoTemplate) updateOrder(order model.Order) {
	if t.active {
		if order.IsActive() {
			t.activeOrders[order.Id] = order
		}
		t.sub.OnOrder(t.context, order)
	}
}

func (t *AlgoTemplate) updateTrade(trade model.Trade) {
	if t.active {
		t.sub.OnTrade(t.context, trade)
	}
}

func (t *AlgoTemplate) start() {
	t.sub.OnStart(t.context)
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
		t.sub.OnTimer(t.context)
	}
}

func (t *AlgoTemplate) Subscribe(symbol string) error {
	return t.engine.Subscribe(t.algoName, symbol)
}
