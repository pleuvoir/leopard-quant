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
	*BaseEngine
	mainEngine            *MainEngine
	symbolAlgoTemplateMap map[string][]string         //币种：[]策略模板名称集合
	algoTemplateMap       map[string]StrategyTemplate //名称：策略模板
}

// NewAlgoEngine 构建算法引擎
// 这个类依赖主引擎，因为所有的订单操作都聚合在主引擎中
func NewAlgoEngine(mainEngine *MainEngine) *AlgoEngine {
	//必须这样赋值，直接StrategyEngine.EngineName 会报空，
	//因为本质上父类还没有实例化 AlgoEngine.BaseEngine.EngineName
	engine := AlgoEngine{BaseEngine: NewBaseEngine(DefaultAlgoEngineName, mainEngine.eventEngine)}
	engine.symbolAlgoTemplateMap = make(map[string][]string)
	engine.algoTemplateMap = make(map[string]StrategyTemplate)
	engine.mainEngine = mainEngine
	return &engine
}

func (s *AlgoEngine) Start() {
	log.Info("策略引擎已启动。")
}

func (s *AlgoEngine) Stop() {
	log.Info("策略引擎已关闭。")
}

func (s *AlgoEngine) initEngine() {

}

func (s *AlgoEngine) RegisterEvent() {
	s.EventEngine.Register(TimerEventInnerHandler{})
}

type TimerEventInnerHandler struct {
	AlgoEngine
}

func (h TimerEventInnerHandler) WithType() event.Type {
	return event.Timer
}

func (h TimerEventInnerHandler) Process(event event.Event) {
	for _, template := range h.algoTemplateMap {
		template.UpdateTimer()
	}
	fmt.Printf("TimerEventInnerHandler receive event %+v %s \n", event, time.Now())
}
