package algorithm

import (
	"fmt"
	. "leopard-quant/core/engine/model"
	"leopard-quant/core/log"
	"leopard-quant/util"
	"strconv"
)

// AlgoTemplateBuilder 所有策略统一的回调处理出口类，负责传递对网关的操作和分发回调信息
type AlgoTemplateBuilder struct {
	engine              *AlgoEngine
	config              map[string]any
	runMode             RunMode
	ticks               map[string]Tick
	nameContexts        map[string]context      //名称：算法上下文
	nameConfigs         map[string]configLoader //名称：算法配置
	symbolTemplateGroup map[string][]string     //币种：[]策略模板名称集合
	nameTemplateGroup   map[string]Template     //名称：策略模板
	templates           []Template              //所有模板合集
}

func NewFactory(a *AlgoEngine) *AlgoTemplateBuilder {
	return &AlgoTemplateBuilder{engine: a}
}

func (a *AlgoTemplateBuilder) LoadTemplates() {
	n := append(a.templates, Noop{}) //TODO
	a.templates = n
	//按名称分组
	for _, template := range a.templates {
		name := template.Name()
		a.nameTemplateGroup[name] = template

		loader := DefaultConfigLoader{}
		if err := loader.load(); err != nil {
			log.Error("加载算法配置参数失败，跳过当前算法加载，", err)
			continue
		}
		a.nameConfigs[name] = &loader
		a.nameContexts[name] = context{e: a.engine, configLoader: a.nameConfigs[name]}
	}
}

func (a *AlgoTemplateBuilder) getContext(name string) context {
	return a.nameContexts[name]
}

func (a *AlgoTemplateBuilder) GetTemplates() []Template {
	return a.templates
}

func (a *AlgoTemplateBuilder) Start() {
	for _, template := range a.templates {
		template.OnStart(a.getContext(template.Name()))
	}
}

func (a *AlgoTemplateBuilder) Stop() {
	for _, template := range a.templates {
		template.OnStop(a.getContext(template.Name()))
	}
}

func (a *AlgoTemplateBuilder) OnBar(bar Bar) {
	for _, template := range a.templates {
		template.OnBar(a.getContext(template.Name()), bar)
	}
}

func (a *AlgoTemplateBuilder) OnTick(tick Tick) {
	t, ok := a.ticks[tick.Symbol]
	if ok {
		for _, template := range a.templates {
			template.OnTick(a.getContext(template.Name()), t)
		}
	} else {
		a.ticks[tick.Symbol] = tick
		for _, template := range a.templates {
			template.OnTick(a.getContext(template.Name()), tick)
		}
	}
}

func (a *AlgoTemplateBuilder) OnTimer() {
	for _, template := range a.templates {
		template.OnTimer(a.getContext(template.Name()))
	}
}

// loadConfig 加载配置
func (a *AlgoTemplateBuilder) loadConfig() {
	//每个算法的参数
	a.config = map[string]any{}
	//运行模式
	a.loadRunMode()
	a.symbolTemplateGroup = make(map[string][]string)
	a.nameTemplateGroup = make(map[string]Template)
	a.nameContexts = make(map[string]context)
	a.nameConfigs = make(map[string]configLoader)
}

// LoadRunMode 设置运行模式
func (a *AlgoTemplateBuilder) loadRunMode() {
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

type Template interface {
	Name() string
	OnStart(c context)
	OnStop(c context)
	OnBar(c context, bar Bar)
	OnTick(c context, tick Tick)
	OnTimer(c context)
}

type configLoader interface {
	load() error
	getStr(key string) string
	getStrOrDefault(key string, val string) string
	getInt(key string) (int, error)
	getIntOrDefault(key string, val int) int
	getBool(key string) (bool, error)
	getBoolOrDefault(key string, val bool) bool
}

type DefaultConfigLoader struct {
	raw map[string]string
}

func (d *DefaultConfigLoader) load() error {
	d.raw = make(map[string]string)
	d.raw["name"] = "pleuvoir"
	return nil
}

func (d *DefaultConfigLoader) getStr(key string) string {
	return d.raw[key]
}

func (d *DefaultConfigLoader) getStrOrDefault(key string, val string) string {
	v := d.getStr(key)
	if util.IsBlank(v) {
		return val
	}
	return v
}

func (d *DefaultConfigLoader) getInt(key string) (int, error) {
	return strconv.Atoi(key)
}

func (d *DefaultConfigLoader) getIntOrDefault(key string, val int) int {
	if v, err := d.getInt(key); err == nil {
		return v
	}
	return val
}

func (d *DefaultConfigLoader) getBool(key string) (bool, error) {
	return strconv.ParseBool(key)
}

func (d *DefaultConfigLoader) getBoolOrDefault(key string, val bool) bool {
	if v, err := d.getBool(key); err == nil {
		return v
	}
	return val
}

type context struct {
	configLoader
	e *AlgoEngine
}
