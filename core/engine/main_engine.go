package engine

import (
	"github.com/gookit/color"
	"leopard-quant/core/event"
	"sync"
	"time"
)

type IEngine interface {
	Name() string
	Start()
	Stop()
}

const DefaultMainEngineName = "leopard-main-engine"

type MainEngine struct {
	todayDate   string
	eventEngine *event.Engine
	engineMap   sync.Map //[string]IEngine 引擎合集
	gatewayMap  sync.Map //[string]Gateway 网关合集
}

// NewMainEngine 构建主引擎
func NewMainEngine(eventEngine *event.Engine) *MainEngine {
	mainEngine := MainEngine{}
	mainEngine.todayDate = time.Now().Format("2006-01-02")
	mainEngine.eventEngine = eventEngine
	mainEngine.engineMap = sync.Map{}
	mainEngine.gatewayMap = sync.Map{}
	return &mainEngine
}

func (m *MainEngine) InitEngines() {
	//注册算法引擎
	m.AddEngine(NewAlgoEngine(m))
	//注册订单引擎
	m.AddEngine(NewOrderEngine(m))
	//注册网关
	m.loadGateway()
}

// 加载网关引擎
func (m *MainEngine) loadGateway() {
	m.AddGateway(NewGateway("mock", m.eventEngine))
}

func (m *MainEngine) Name() string {
	return DefaultMainEngineName
}

// Start 主引擎启动
func (m *MainEngine) Start() {
	//启动事件引擎
	m.eventEngine.StartAll()
	for _, engine := range m.GetAllEngine() {
		engine.Start()
	}
	color.Greenln("主引擎已启动")
}

// Stop 主引擎关闭
func (m *MainEngine) Stop() {
	//关闭事件引擎
	m.eventEngine.StopAll()
	//关闭所有引擎
	for _, engine := range m.GetAllEngine() {
		engine.Stop()
	}
	color.Redln("主引擎已关闭")
}

// AddEngine 增加引擎
func (m *MainEngine) AddEngine(engine IEngine) {
	m.engineMap.Store(engine.Name(), engine)
}

// GetEngine 获取引擎
func (m *MainEngine) GetEngine(engineName string) IEngine {
	e, ok := m.engineMap.Load(engineName)
	if ok {
		engine := e.(IEngine)
		return engine
	}
	return nil
}

func (m *MainEngine) GetAllEngine() (engines []IEngine) {
	r := make(map[string]IEngine)
	m.engineMap.Range(func(key, value any) bool {
		r[key.(string)] = value.(IEngine)
		return true
	})
	for _, engine := range r {
		engines = append(engines, engine)
	}
	return engines
}

// AddGateway 增加网关
func (m *MainEngine) AddGateway(gateway *GatewayEngine) {
	m.gatewayMap.Store(gateway.Name(), gateway)
}

// GetGateway 获取网关
func (m *MainEngine) GetGateway(engineName string) IEngine {
	e, ok := m.gatewayMap.Load(engineName)
	if ok {
		engine := e.(IEngine)
		return engine
	}
	return nil
}

func (m *MainEngine) GetAllGateway() (engines []IEngine) {
	r := make(map[string]IEngine)
	m.gatewayMap.Range(func(key, value any) bool {
		r[key.(string)] = value.(IEngine)
		return true
	})
	for _, engine := range r {
		engines = append(engines, engine)
	}
	return engines
}
