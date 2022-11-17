package engine

import (
	"leopard-quant/bootstrap"
	"leopard-quant/core/event"
	"testing"
)

func TestStrategyTemplate_OnTick(t *testing.T) {

	bootstrap.Init()

	mainEngine := NewMainEngine(event.NewEventEngine())
	engine := NewAlgoEngine(mainEngine)

	engine.initEngine()

	engine.Start()

	var aa IEngine = engine
	aa.Start()

	ch := make(chan any)
	<-ch
}
