package engine

import (
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"leopard-quant/core/event"
	"leopard-quant/gateway/okx"
	"sync"
	"time"
)

type Engineer interface {
	Name() string
	Start()
	Stop()
}

const DefaultMainEngineName = "leopard-main-engine"

type MainEngine struct {
	TodayDate   string
	eventEngine *event.Engine
	engineMap   sync.Map //[string]Engineer 引擎合集
	gatewayMap  sync.Map //[string]Gateway 网关合集
}

// NewMainEngine 构建主引擎
func NewMainEngine(eventEngine *event.Engine) *MainEngine {
	mainEngine := MainEngine{}
	mainEngine.TodayDate = time.Now().Format("2006-01-02")
	mainEngine.eventEngine = eventEngine
	mainEngine.engineMap = sync.Map{}
	mainEngine.gatewayMap = sync.Map{}
	return &mainEngine
}

func (m *MainEngine) InitEngines() {
	//注册订单引擎
	m.AddEngine(NewOrderEngine(m))
	//注册网关
	m.loadGateway()
}

// 加载网关引擎
func (m *MainEngine) loadGateway() {
	m.AddGateway(NewGateway("okx", m.eventEngine, &okx.Default))
}

// RegisterListener 注册事件
func (m *MainEngine) RegisterListener(t event.Type, f func(e event.Event)) {
	m.eventEngine.Register(t, event.AdaptEventHandlerFunc(f))
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
func (m *MainEngine) AddEngine(engine Engineer) {
	m.engineMap.Store(engine.Name(), engine)
}

// GetEngine 获取引擎
func (m *MainEngine) GetEngine(engineName string) Engineer {
	e, ok := m.engineMap.Load(engineName)
	if ok {
		engine := e.(Engineer)
		return engine
	}
	return nil
}

func (m *MainEngine) GetAllEngine() (engines []Engineer) {
	r := make(map[string]Engineer)
	m.engineMap.Range(func(key, value any) bool {
		r[key.(string)] = value.(Engineer)
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
func (m *MainEngine) GetGateway(engineName string) Engineer {
	e, ok := m.gatewayMap.Load(engineName)
	if ok {
		engine := e.(Engineer)
		return engine
	}
	return nil
}

func (m *MainEngine) GetAllGateway() (engines []Engineer) {
	r := make(map[string]Engineer)
	m.gatewayMap.Range(func(key, value any) bool {
		r[key.(string)] = value.(Engineer)
		return true
	})
	for _, engine := range r {
		engines = append(engines, engine)
	}
	return engines
}

// Subscribe 订阅币种
func (m *MainEngine) Subscribe(gatewayName string, symbol string) error {
	g, ok := m.gatewayMap.Load(gatewayName)
	if !ok {
		return errors.Errorf("未找到该网关，gatewayName=%s", gatewayName)
	}
	return g.(*GatewayEngine).Subscribe(symbol)
}
