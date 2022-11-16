package engine

import (
	"fmt"
	"leopard-quant/core/event"
	"leopard-quant/core/log"
	"time"
)

const DefaultAlgoEngineName = "leopard-algo-engine"

// AlgoEngine 策略引擎
type AlgoEngine struct {
	name                  string
	mainEngine            *MainEngine
	symbolAlgoTemplateMap map[string][]string         //币种：[]策略模板名称集合
	algoTemplateMap       map[string]StrategyTemplate //名称：策略模板
}

// NewAlgoEngine 构建算法引擎
// 这个类依赖主引擎，因为所有的订单操作都聚合在主引擎中
func NewAlgoEngine(mainEngine *MainEngine) *AlgoEngine {
	engine := AlgoEngine{mainEngine: mainEngine}
	engine.symbolAlgoTemplateMap = make(map[string][]string)
	engine.algoTemplateMap = make(map[string]StrategyTemplate)
	engine.initEngine()
	return &engine
}

func (s *AlgoEngine) Name() string {
	return s.name
}

func (s *AlgoEngine) Start() {
	log.Info("策略引擎已启动。")
}

func (s *AlgoEngine) Stop() {
	log.Info("策略引擎关闭状态。")
}

func (s *AlgoEngine) initEngine() {
	s.name = DefaultAlgoEngineName
}

func (s *AlgoEngine) RegisterEvent() {
	s.mainEngine.eventEngine.Register(event.Timer, TimerEventInnerHandler{})
}

type TimerEventInnerHandler struct {
	AlgoEngine
}

func (h TimerEventInnerHandler) Process(event event.Event) {
	for _, template := range h.algoTemplateMap {
		template.UpdateTimer()
	}
	fmt.Printf("TimerEventInnerHandler receive event %+v %s \n", event, time.Now())
}
