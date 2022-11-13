package core

import (
	"fmt"
	"leopard-quant/core/event"
	"time"
)

const DefaultAlgoEngineName = "leopard-strategy-engine"

// StrategyEngine 策略引擎
type StrategyEngine struct {
	Engine
	SymbolStrategyTemplateMap map[string][]string         //币种：[]策略模板名称集合
	StrategyTemplateMap       map[string]StrategyTemplate //名称：策略模板
}

func NewStrategyEngine() *StrategyEngine {
	return &StrategyEngine{NewEngine(DefaultAlgoEngineName, event.NewEventEngine()),
		map[string][]string{}, map[string]StrategyTemplate{}}
}

func (s *StrategyEngine) initEngine() {

}

func (s *StrategyEngine) Start() {
	s.EventEngine.StartSchedulerTimer()
	s.EventEngine.StartConsumer()
	s.Engine.Start()
}

func (s *StrategyEngine) RegisterEvent() {
	s.EventEngine.Register(TimerEventInnerHandler{})
}

type TimerEventInnerHandler struct {
	StrategyEngine
}

func (h TimerEventInnerHandler) GetType() event.Type {
	return event.Timer
}

func (h TimerEventInnerHandler) Process(event event.Event) {
	for _, template := range h.StrategyTemplateMap {
		template.UpdateTimer()
	}
	fmt.Printf("TimerEventInnerHandler receive event %+v %s \n", event, time.Now())
}
