package base

import (
	"leopard-quant/algorithm"
	"leopard-quant/algorithm/impl"
	"leopard-quant/core/engine"
	. "leopard-quant/core/engine/model"
	"leopard-quant/core/event"
	"leopard-quant/core/log"
)

const DefaultAlgoEngineName = "leopard-algo-engine"

// AlgoEngine 算法引擎
type AlgoEngine struct {
	name                   string
	mainEngine             *engine.MainEngine
	nameAlgoTemplateMap    map[string]*AlgoTemplate   //模板名称：模板
	orderIdAlgoTemplateMap map[string]*AlgoTemplate   //订单号：模板
	symbolAlgoTemplateMap  map[string][]*AlgoTemplate //币种:[]模板
}

// NewAlgoEngine 构建算法引擎
// 这个类依赖主引擎，因为所有的订单操作都聚合在主引擎中
func NewAlgoEngine(mainEngine *engine.MainEngine) *AlgoEngine {
	e := AlgoEngine{name: DefaultAlgoEngineName, mainEngine: mainEngine}
	e.nameAlgoTemplateMap = map[string]*AlgoTemplate{}
	e.orderIdAlgoTemplateMap = map[string]*AlgoTemplate{}
	e.symbolAlgoTemplateMap = map[string][]*AlgoTemplate{}
	e.initEngine()
	e.registerEvent()
	return &e
}

func (s *AlgoEngine) initEngine() {
	s.loadSetting()
}

func (s *AlgoEngine) loadSetting() {
}

func (s *AlgoEngine) Name() string {
	return s.name
}

func (s *AlgoEngine) Start() {
	//加载配置，初始所有生效的模板
	//TODO

	mainEngine := s.mainEngine
	if _, err := algorithm.MakeInstance("noop"); err == nil {
		//	template := NewAlgoTemplate(mainEngine, sub)
		template := NewAlgoTemplate(mainEngine, &impl.NoopSub{})
		//	sub.template = template
		s.nameAlgoTemplateMap[template.algoName] = template
	}

	for _, template := range s.nameAlgoTemplateMap {
		template.start()
	}
	log.Info("算法引擎已启动。")
}

func (s *AlgoEngine) Stop() {
	for _, template := range s.nameAlgoTemplateMap {
		template.stop()
	}
	log.Info("算法引擎已关闭。")
}

// 注册回调事件，每个模板都会收到通知
func (s *AlgoEngine) registerEvent() {
	mainEngine := s.mainEngine
	mainEngine.RegisterListener(event.Tick, s.tickHandler())
	mainEngine.RegisterListener(event.Timer, s.timerHandler())
	mainEngine.RegisterListener(event.Trade, s.tradeHandler())
	mainEngine.RegisterListener(event.Order, s.orderHandler())
}

// 对应币种模板会收到回调
func (s *AlgoEngine) tickHandler() func(e event.Event) {
	return func(e event.Event) {
		tick := e.EventData.(Tick)
		templates := s.symbolAlgoTemplateMap[tick.Symbol]
		for _, template := range templates {
			s.nameAlgoTemplateMap[template.algoName].updateTick(tick)
		}
	}
}

// 所有模板会收到回调
func (s *AlgoEngine) timerHandler() func(e event.Event) {
	return func(e event.Event) {
		for _, template := range s.nameAlgoTemplateMap {
			template.updateTimer()
		}
	}
}

// 当前交易的模板会收到此回调
func (s *AlgoEngine) tradeHandler() func(e event.Event) {
	return func(e event.Event) {
		trade := e.EventData.(Trade)
		template := s.orderIdAlgoTemplateMap[trade.OrderId]
		if template != nil {
			template.updateTrade(trade)
		}
	}
}

// 当前订单的模板会收到此回调
func (s *AlgoEngine) orderHandler() func(e event.Event) {
	return func(e event.Event) {
		order := e.EventData.(Order)
		template := s.orderIdAlgoTemplateMap[order.Id]
		if template != nil {
			template.updateOrder(order)
		}
	}
}
