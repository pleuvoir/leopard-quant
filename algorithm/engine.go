package algorithm

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"leopard-quant/core/config"
	"leopard-quant/core/engine"
	. "leopard-quant/core/engine/model"
	"leopard-quant/core/event"
	"leopard-quant/core/log"
	"leopard-quant/util"
	"os"
)

const DefaultAlgoEngineName = "leopard-algo-engine"
const DefaultGatewayName = "okx"

// AlgoEngine 算法引擎
type AlgoEngine struct {
	name                   string
	mainEngine             *engine.MainEngine
	nameAlgoTemplateMap    map[string]*AlgoTemplate //模板名称：模板
	nameAlgoConfigMap      map[string]config.Loader //模板名称：配置
	orderIdAlgoTemplateMap map[string]*AlgoTemplate //订单号：模板
	symbolAlgoTemplateMap  map[string][]string      //币种:[]模板名称
	configPath             string                   //配置文件路径
}

// NewAlgoEngine 构建算法引擎
// 这个类依赖主引擎，因为所有的订单操作都聚合在主引擎中
func NewAlgoEngine(mainEngine *engine.MainEngine, algoConfig config.Algo) *AlgoEngine {
	e := AlgoEngine{name: DefaultAlgoEngineName, mainEngine: mainEngine}
	e.nameAlgoTemplateMap = map[string]*AlgoTemplate{}
	e.orderIdAlgoTemplateMap = map[string]*AlgoTemplate{}
	e.symbolAlgoTemplateMap = map[string][]string{}
	e.nameAlgoConfigMap = map[string]config.Loader{}
	e.configPath = algoConfig.ConfigPath
	e.initEngine()
	e.registerEvent()
	return &e
}

func (s *AlgoEngine) initEngine() {
	algoConfig, err := s.loadAlgoConfig()
	if err != nil {
		panic(err)
	}
	for subName, items := range algoConfig {
		s.nameAlgoConfigMap[subName] = config.NewConfigLoader(algoConfigItem{m: items})
	}
}

type algoConfigItem struct {
	m map[string]any
}

func (algoConfigItem) Load() error {
	return nil
}

func (a algoConfigItem) GetStr(key string) string {
	return fmt.Sprint(a.m[key])
}

func (s *AlgoEngine) loadAlgoConfig() (m map[string]map[string]any, err error) {
	c := make(map[string]map[string]any)
	data, err := os.ReadFile(s.configPath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *AlgoEngine) Name() string {
	return s.name
}

func (s *AlgoEngine) Start() {
	for subName, configLoader := range s.nameAlgoConfigMap {
		enable := configLoader.GetBoolOrDefault("enable", false)
		if !enable {
			continue
		}
		//加载配置，初始所有生效的模板
		if sub, err := MakeInstance(subName); err == nil {
			template := NewAlgoTemplate(s, sub, configLoader)
			s.nameAlgoTemplateMap[template.algoName] = template
			template.start()
		}
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
			algoTemplate := s.nameAlgoTemplateMap[template]
			algoTemplate.updateTick(tick)
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

func (s *AlgoEngine) Subscribe(subName string, symbol string) error {
	if !Exist(subName) {
		return errors.Errorf("不存在的模板[%s]", subName)
	}
	if util.IsBlank(symbol) {
		return errors.Errorf("订阅symbol不能为空")
	}
	//TODO 检查支持该币种
	templates := s.symbolAlgoTemplateMap[symbol]
	algoTemplates := append(templates, subName)
	s.symbolAlgoTemplateMap[symbol] = algoTemplates

	return s.mainEngine.Subscribe(DefaultGatewayName, symbol)
}
