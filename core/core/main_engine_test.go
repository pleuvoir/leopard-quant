package core

import (
	"fmt"
	"leopard-quant/core/event"
	"testing"
	"time"
)

type TimerEventHandler struct{}

func (h TimerEventHandler) GetType() event.Type {
	return event.Timer
}

func (TimerEventHandler) Process(event event.Event) {
	fmt.Printf("TimerEventHandler receive event %+v %s \n", event, time.Now())
}

func TestNewMainEngine(t *testing.T) {

	engine := NewMainEngine(event.NewEventEngine())

	t.Log(engine.TodayDate)

	eventEngine := engine.EventEngine

	timerEventHandler := TimerEventHandler{}
	eventEngine.Register(timerEventHandler)

	eventEngine.Put(event.NewEvent(event.Timer, "我是时间事件"))

	ch := make(chan any)
	<-ch
}
