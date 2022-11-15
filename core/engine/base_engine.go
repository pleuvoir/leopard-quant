package engine

import "leopard-quant/core/event"

type IEngine interface {
	Start()
	Stop()
}

type BaseEngine struct {
	EngineName  string
	EventEngine *event.Engine
}

// NewBaseEngine 构造方法
func NewBaseEngine(engineName string, eventEngine *event.Engine) *BaseEngine {
	return &BaseEngine{engineName, eventEngine}
}

// Start 抽象方法 可以由子类实现
func (b *BaseEngine) Start() {
}

// Stop 抽象方法 可以由子类实现
func (b *BaseEngine) Stop() {
}
