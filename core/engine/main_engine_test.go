package engine

import (
	"fmt"
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

	mainEngine := NewMainEngine(event.NewEventEngine())

	t.Log(mainEngine.todayDate)

	eventEngine := mainEngine.EventEngine

	timerEventHandler := TimerEventHandler{}
	eventEngine.Register(timerEventHandler)

	eventEngine.Put(event.NewEvent(event.Timer, "我是时间事件"))

	newEngine := NewEngine("default", event.NewEventEngine())
	mainEngine.AddEngine(newEngine)

	defaultEngine, err := mainEngine.GetEngine("default")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("default ==> %+v \n", defaultEngine)
	fmt.Printf("default equal %v \n", defaultEngine == newEngine)

	engines := mainEngine.GetAllEngine()
	for _, engine := range engines {
		fmt.Printf("current engine  ==> %+v \n", engine)

	}

	ch := make(chan any)
	<-ch
}
