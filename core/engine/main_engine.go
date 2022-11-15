package engine

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"leopard-quant/core/event"
	"sync"
	"time"
)

const DefaultMainEngineName = "leopard-main-engine"

type MainEngine struct {
	todayDate   string
	eventEngine *event.Engine
	engineMap   sync.Map //[string]Engine 引擎合集
	gateWayMap  sync.Map //[string]GateWay 网关合集
}

// NewMainEngine 构建主引擎
func NewMainEngine(eventEngine *event.Engine) *MainEngine {
	mainEngine := MainEngine{}
	mainEngine.todayDate = time.Now().Format("2006-01-02")
	mainEngine.eventEngine = eventEngine
	mainEngine.engineMap = sync.Map{}
	mainEngine.gateWayMap = sync.Map{}
	return &mainEngine
}

func (m *MainEngine) InitEngines() {
	//启动事件引擎
	m.eventEngine.StartAll()
	//注册算法引擎
	//algoEngine := NewAlgoEngine(m)

	//var base BaseEngine = &algoEngine
	//
	//m.AddEngine(algoEngine.BaseEngine) //遗留
}

// Start 主引擎启动
func (m *MainEngine) Start() {
	color.Greenln("主引擎已启动")
}

// Stop 主引擎关闭
func (m *MainEngine) Stop() {
	color.Redln("主引擎已关闭")

}

// AddGateway 增加网关
func (m *MainEngine) AddGateway(gateway *Gateway) {
	m.gateWayMap.Store(gateway.GatewayName, gateway)
}

// GetGateway 获取网关
func (m *MainEngine) GetGateway(gatewayName string) (*Gateway, error) {
	e, ok := m.gateWayMap.Load(gatewayName)
	if ok {
		gateway := e.(Gateway)
		return &gateway, nil
	}
	return &Gateway{}, errors.New(fmt.Sprintf("未找到网关，gatewayName[%s]", gatewayName))
}

// GetAllGateway 获取所有网关
func (m *MainEngine) GetAllGateway() (gateways []*Gateway) {
	r := make(map[string]Gateway)
	m.gateWayMap.Range(func(key, value any) bool {
		r[key.(string)] = value.(Gateway)
		return true
	})
	for _, gateway := range r {
		gateways = append(gateways, &gateway)
	}
	return gateways
}

// AddEngine 增加引擎
func (m *MainEngine) AddEngine(engine *BaseEngine) {
	m.engineMap.Store(engine.EngineName, engine)
}

// GetEngine 获取引擎
func (m *MainEngine) GetEngine(engineName string) (*BaseEngine, error) {
	e, ok := m.engineMap.Load(engineName)
	if ok {
		engine := e.(BaseEngine)
		return &engine, nil
	}
	return &BaseEngine{}, errors.New(fmt.Sprintf("未找到引擎，engineName[%s]", engineName))
}

func (m *MainEngine) GetAllEngine() (engines []*BaseEngine) {
	r := make(map[string]BaseEngine)
	m.engineMap.Range(func(key, value any) bool {
		r[key.(string)] = value.(BaseEngine)
		return true
	})
	for _, engine := range r {
		engines = append(engines, &engine)
	}
	return engines
}
