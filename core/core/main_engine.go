package core

import (
	"leopard-quant/core/event"
	"sync"
	"time"
)

type MainEngine struct {
	TodayDate   string
	Engines     sync.Map //[string]Engine
	EventEngine *event.Engine
}

func NewMainEngine(eventEngine *event.Engine) *MainEngine {

	todayDate := time.Now().Format("2006-01-02")

	engine := MainEngine{TodayDate: todayDate, Engines: sync.Map{}, EventEngine: eventEngine}
	eventEngine.StartSchedulerTimer()
	eventEngine.StartConsumer()

	engine.InitEngines()

	return &engine
}

func (m *MainEngine) InitEngines() {
}

func (m *MainEngine) AddEngine(engine Engine) {
	m.Engines.Store(engine.EngineName, engine)
}
