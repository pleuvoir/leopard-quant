package engine

import (
	"fmt"
	"leopard-quant/bootstrap"
	"leopard-quant/core/event"
	"testing"
	"time"
)

type TimerEventHandler struct{}

func (h TimerEventHandler) WithType() event.Type {
	return event.Timer
}

func (TimerEventHandler) Process(event event.Event) {
	fmt.Printf("TimerEventInnerHandler receive event %+v %s \n", event, time.Now())
}

func TestNewMainEngine(t *testing.T) {

	bootstrap.Init()

	mainEngine := NewMainEngine(event.NewEventEngine())
	mainEngine.InitEngines()

	t.Log(mainEngine.todayDate)

	//ch := make(chan any)

	defer func() {
		mainEngine.Stop()
	}()

}
