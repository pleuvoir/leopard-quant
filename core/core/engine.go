package core

import "leopard-quant/core/event"

type Engine struct {
	EngineName  string
	EventEngine *event.Engine
}

func NewEngine(engineName string, eventEngine *event.Engine) Engine {
	return Engine{engineName, eventEngine}
}

func (e Engine) Start() {

}

func (e Engine) Stop() {

}
