package engine

import (
	"fmt"
	"leopard-quant/algorithm"
	. "leopard-quant/core/engine/model"
	"leopard-quant/core/event"
	"leopard-quant/core/log"
)

const DefaultAlgoEngineName = "leopard-algo-engine"

// AlgoEngine 策略引擎
type AlgoEngine struct {
	name                string
	mainEngine          *MainEngine
	symbolTemplateGroup map[string][]string //币种：[]策略模板名称集合
	nameTemplateGroup   map[string]template //名称：策略模板
}

// NewAlgoEngine 构建算法引擎
// 这个类依赖主引擎，因为所有的订单操作都聚合在主引擎中
func NewAlgoEngine(mainEngine *MainEngine) *AlgoEngine {
	engine := AlgoEngine{mainEngine: mainEngine}
	engine.symbolTemplateGroup = make(map[string][]string)
	engine.nameTemplateGroup = make(map[string]template)
	engine.initEngine()
	return &engine
}

func (s *AlgoEngine) initEngine() {
	s.name = DefaultAlgoEngineName
	//加载算法模板
	//elem := reflect.TypeOf(&AlgoTemplate{}).Elem()
	//
	//algoTemplate := reflect.New(elem).Elem().Interface().(AlgoTemplate)

	algoTemplate := algorithm.NewEcho(s)

	s.nameTemplateGroup[algoTemplate.algoName] = algoTemplate

	//注册回调事件
	s.registerEvent()
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

// 注册事件
func (s *AlgoEngine) registerEvent() {
	eventEngine := s.mainEngine.eventEngine
	eventEngine.Register(event.Tick, event.AdaptEventHandlerFunc(func(e event.Event) {
		for _, template := range s.nameTemplateGroup {
			template.onTick(e.EventData.(Tick))
		}
	}))
	eventEngine.Register(event.Bar, event.AdaptEventHandlerFunc(func(e event.Event) {
		for _, template := range s.nameTemplateGroup {
			template.onBar(e.EventData.(Bar))
		}
	}))
	eventEngine.Register(event.Timer, event.AdaptEventHandlerFunc(func(e event.Event) {
		for _, template := range s.nameTemplateGroup {
			template.onTimer()
		}
	}))
}

type AlgoTemplate struct {
	e           *AlgoEngine
	active      bool
	config      map[string]any
	runMode     RunMode
	algoName    string
	ticks       map[string]Tick
	onTickFunc  func(t Tick)
	onBarFunc   func(b Bar)
	onTimerFunc func()
}

func NewAlgoTemplate(algoName string, e *AlgoEngine) *AlgoTemplate {
	template := AlgoTemplate{e: e}
	template.active = false
	template.algoName = algoName
	template.loadConfig()
	return &template
}

func (a *AlgoTemplate) WithOnTicker(f func(t Tick)) {
	a.onTickFunc = f
}

func (a *AlgoTemplate) WithOnBar(f func(b Bar)) {
	a.onBarFunc = f
}

func (a *AlgoTemplate) WithOnTimer(f func()) {
	a.onTimerFunc = f
}

func (a *AlgoTemplate) start() {
	a.active = true
}

func (a *AlgoTemplate) stop() {
	a.active = false
}

func (a *AlgoTemplate) onBar(bar Bar) {
	if a.active {
		a.onBarFunc(bar)
	}
}

func (a *AlgoTemplate) onTick(tick Tick) {
	t, ok := a.ticks[tick.Symbol]
	if ok {
		a.onTickFunc(t)
	} else {
		a.ticks[tick.Symbol] = tick
		a.onTickFunc(tick)
	}
}

func (a *AlgoTemplate) onTimer() {
	if a.active {
		a.onTimerFunc()
	}
}

// loadConfig 加载配置
func (a *AlgoTemplate) loadConfig() {
	a.config = map[string]any{}
	a.onTickFunc = func(t Tick) {}
	a.onBarFunc = func(b Bar) {}
	a.onTimerFunc = func() {}
	a.loadRunMode()
}

// LoadRunMode 设置运行模式
func (a *AlgoTemplate) loadRunMode() {
	mode := a.config["mode"]
	switch mode {
	case Live.Name:
		a.runMode = Live
	case DryRun.Name:
		a.runMode = DryRun
	case BackTesting.Name:
		a.runMode = BackTesting
	default:
		log.Error(fmt.Sprintf("错误的策略运行模式，mode is %s。默认选中为试运行。", mode))
		a.runMode = DryRun
	}
}

type template interface {
	start()
	stop()
	onBar(bar Bar)
	onTick(tick Tick)
	onTimer()
}
