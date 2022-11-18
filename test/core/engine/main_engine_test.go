package engine

import (
	"leopard-quant/bootstrap"
	"leopard-quant/core/engine"
	"leopard-quant/core/event"
	"testing"
)

func TestNewMainEngine(t *testing.T) {

	bootstrap.Init()
	mainEngine := engine.NewMainEngine(event.NewEventEngine())
	mainEngine.InitEngines()

	t.Log(mainEngine.TodayDate)

	//ch := make(chan any)

	//mainEngine.Start()
	defer func() {
		mainEngine.Stop()
	}()

}
