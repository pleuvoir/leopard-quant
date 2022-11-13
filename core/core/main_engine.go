package core

import (
	"fmt"
	"github.com/pkg/errors"
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

// GetEngine 获取引擎
func (m *MainEngine) GetEngine(EngineName string) (Engine, error) {
	e, ok := m.Engines.Load(EngineName)
	if ok {
		return e.(Engine), nil
	}
	return Engine{}, errors.New(fmt.Sprintf("未找到引擎，EngineName[%s]", EngineName))
}

func (m *MainEngine) GetAllEngine() (engines []Engine) {
	r := make(map[string]Engine)
	m.Engines.Range(func(key, value any) bool {
		r[key.(string)] = value.(Engine)
		return true
	})
	for _, engine := range r {
		engines = append(engines, engine)
	}
	return engines
}
