package engine

import (
	"fmt"
	"leopard-quant/core/engine/model"
	"leopard-quant/core/log"
)

type StrategyConfig map[string]any

type StrategyTemplate struct {
	StrategyName string
	Config       StrategyConfig
	Active       bool
	RunMode      model.RunMode
}

// NewStrategyTemplate 创建策略模板
func NewStrategyTemplate(StrategyName string, config StrategyConfig) *StrategyTemplate {
	template := StrategyTemplate{StrategyName: StrategyName, Config: config, Active: false}
	return &template
}

func (s *StrategyTemplate) Init() {
	s.LoadRunMode()
}

func (s *StrategyTemplate) Start() {
	s.Active = true
}

func (s *StrategyTemplate) Stop() {
	s.Active = false
}

func (s *StrategyTemplate) UpdateTick(tick model.Tick) {
	if s.Active {
		s.OnTick()
	}
}

func (s *StrategyTemplate) UpdateTimer() {
	if s.Active {
		s.OnTime()
	}
}

func (s *StrategyTemplate) OnTime() {
}

func (s *StrategyTemplate) OnTick() {
}

// LoadRunMode 设置运行模式
func (s *StrategyTemplate) LoadRunMode() {
	mode := s.Config["mode"]
	switch mode {
	case model.Live.Name:
		s.RunMode = model.Live
	case model.DryRun.Name:
		s.RunMode = model.DryRun
	case model.BackTesting.Name:
		s.RunMode = model.BackTesting
	default:
		log.Error(fmt.Sprintf("错误的策略运行模式，mode is %s。默认选中为试运行。", mode))
		s.RunMode = model.DryRun
	}
}
