package engine

import (
	"leopard-quant/bootstrap"
	"leopard-quant/core/engine"
	"leopard-quant/core/event"
	"testing"
	"time"
)

func TestOrder(t *testing.T) {

	bootstrap.Init()

	mainEngine := engine.NewMainEngine(event.NewEventEngine())
	mainEngine.InitEngines()
	mainEngine.Start()

	//eventEngine := mainEngine.eventEngine
	////eventEngine.StopTimer()
	//
	//for _, listeners := range eventEngine.HandlersMap {
	//	t.Log(listeners)
	//}
	//t.Log(eventEngine.HandlersMap)
	//
	//newEvent := event.NewEvent(event.Trade, model.Trade{Id: "123"})
	//eventEngine.Put(newEvent)

	time.Sleep(time.Second * 5)
}
