package algorithm

import (
	"leopard-quant/core/engine"
	"leopard-quant/core/engine/model"
	"leopard-quant/core/log"
	"time"
)

var echoTemplate *engine.AlgoTemplate

func NewEcho(e *engine.AlgoEngine) *engine.AlgoTemplate {

	me := engine.NewAlgoTemplate("echo_template", e)
	me.WithOnTicker(onTick)
	me.WithOnBar(onBar)
	me.WithOnTimer(onTimer)

	echoTemplate = me
	return echoTemplate
}

func onBar(t model.Bar) {
	log.Info("接收到 var %+v", t)
}

func onTick(t model.Tick) {
	log.Info("接收到 ticker %+v", t)
}

func onTimer() {
	log.Info("接收到定时 %s", time.Now())

}
