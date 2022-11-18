package algorithm

import (
	"leopard-quant/core/engine"
	. "leopard-quant/core/engine/model"
	"leopard-quant/core/event"
	"leopard-quant/core/log"
)

const DefaultAlgoEngineName = "leopard-algo-engine"

// AlgoEngine 策略引擎
type AlgoEngine struct {
	name       string
	mainEngine *engine.MainEngine
	Active     bool
	factory    *AlgoTemplateBuilder //模板分发工厂
}

// NewAlgoEngine 构建算法引擎
// 这个类依赖主引擎，因为所有的订单操作都聚合在主引擎中
func NewAlgoEngine(mainEngine *engine.MainEngine) *AlgoEngine {
	e := AlgoEngine{name: DefaultAlgoEngineName, mainEngine: mainEngine}
	e.initEngine()
	return &e
}

func (s *AlgoEngine) initEngine() {
	s.factory = NewFactory(s)
	s.factory.loadConfig()
	//加载所有算法
	s.factory.LoadTemplates()
	//注册工厂回调事件
	s.registerFactoryEvent()
}

func (s *AlgoEngine) Name() string {
	return s.name
}

func (s *AlgoEngine) Start() {
	s.Active = true
	s.factory.Start()
	log.Info("算法引擎已启动。")
}

func (s *AlgoEngine) Stop() {
	s.Active = false
	s.factory.Stop()
	log.Info("算法引擎已关闭。")
}

// 注册事件
// 当算法引擎接收到事件回调时，通知工厂
func (s *AlgoEngine) registerFactoryEvent() {
	m := s.mainEngine
	f := s.factory
	m.RegisterListener(event.Tick, func(e event.Event) {
		if s.Active {
			f.OnTick(e.EventData.(Tick))
		} else {
			log.Warnf("算法引擎未开启，丢弃Tick回调事件。")
		}
	})
	m.RegisterListener(event.Bar, func(e event.Event) {
		if s.Active {
			f.OnBar(e.EventData.(Bar))
		} else {
			log.Warnf("算法引擎未开启，丢弃Bar回调事件。")
		}
	})
	m.RegisterListener(event.Timer, func(e event.Event) {
		if s.Active {
			f.OnTimer()
		} else {
			log.Warnf("算法引擎未开启，丢弃Timer回调事件。")
		}
	})
}
