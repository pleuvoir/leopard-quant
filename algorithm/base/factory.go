package base

import (
	"fmt"
	"leopard-quant/algorithm"
	. "leopard-quant/core/engine/model"
	"leopard-quant/core/log"
)

// AlgoTemplateFactory 所有策略统一的回调处理出口类，负责传递对网关的操作和分发回调信息
type AlgoTemplateFactory struct {
	a                   *AlgoEngine
	config              map[string]any
	runMode             RunMode
	ticks               map[string]Tick
	symbolTemplateGroup map[string][]string //币种：[]策略模板名称集合
	nameTemplateGroup   map[string]Template //名称：策略模板
	templates           []Template          //所有模板合集
}

func NewFactory(a *AlgoEngine) *AlgoTemplateFactory {
	f := AlgoTemplateFactory{a: a}
	f.loadConfig()
	return &f
}

func (a *AlgoTemplateFactory) loadTemplates() {
	n := append(a.templates, algorithm.Noop{})
	a.templates = n
	//按名称分组
	for _, template := range a.templates {
		a.nameTemplateGroup[template.Name()] = template
	}
}

func (a *AlgoTemplateFactory) GetTemplates() []Template {
	return a.templates
}

func (a *AlgoTemplateFactory) start() {
	for _, template := range a.templates {
		template.OnStart()
	}
}

func (a *AlgoTemplateFactory) stop() {
	for _, template := range a.templates {
		template.OnStop()
	}
}

func (a *AlgoTemplateFactory) onBar(bar Bar) {
	for _, template := range a.templates {
		template.OnBar(bar)
	}
}

func (a *AlgoTemplateFactory) onTick(tick Tick) {
	t, ok := a.ticks[tick.Symbol]
	if ok {
		for _, template := range a.templates {
			template.OnTick(t)
		}
	} else {
		a.ticks[tick.Symbol] = tick
		for _, template := range a.templates {
			template.OnTick(tick)
		}
	}
}

func (a *AlgoTemplateFactory) onTimer() {
	for _, template := range a.templates {
		template.OnTimer()
	}
}

// loadConfig 加载配置
func (a *AlgoTemplateFactory) loadConfig() {
	a.config = map[string]any{}
	a.loadRunMode()
	a.symbolTemplateGroup = make(map[string][]string)
	a.nameTemplateGroup = make(map[string]Template)
}

// LoadRunMode 设置运行模式
func (a *AlgoTemplateFactory) loadRunMode() {
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
	OnStart()
	OnStop()
	OnBar(bar Bar)
	OnTick(tick Tick)
	OnTimer()
}
