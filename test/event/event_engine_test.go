package event

import (
	"fmt"
	"leopard-quant/core/event"
	"testing"
	"time"
)

type TimerEventHandler struct{}

func (TimerEventHandler) Process(event event.Event) {
	fmt.Printf("TimerEventHandler receive event %+v %s \n", event, time.Now())
}

type LogEventHandler struct{}

func (LogEventHandler) Process(event event.Event) {
	fmt.Printf("LogEventHandler receive event %+v  %s \n", event, time.Now())
}

func Process(event event.Event) {
	fmt.Printf("LogEventHandler receive event %+v  %s \n", event, time.Now())
}

func TestHandlerFuncWithReceiver(t *testing.T) {
	engine := event.NewEventEngine()
	engine.Register(event.Log, event.AdaptEventHandlerFunc(Process))

	engine.Put(event.NewEvent(event.Log, "i am log."))

	engine.StartAll()

	engine.Put(event.NewEvent(event.Log, "i am log2."))

	time.Sleep(5 * time.Second)

}

func TestHandlerFunc(t *testing.T) {
	engine := event.NewEventEngine()
	engine.Register(event.Log, event.AdaptEventHandlerFunc(Process))

	engine.Put(event.NewEvent(event.Log, "i am log."))

	engine.StartAll()

	engine.Put(event.NewEvent(event.Log, "i am log2."))

	time.Sleep(5 * time.Second)

	engine.StopAll()

}

func TestNewEventEngine(t *testing.T) {

	engine := event.NewEventEngine()

	timerEventHandler := TimerEventHandler{}
	engine.Register(event.Timer, timerEventHandler)
	engine.Register(event.Log, LogEventHandler{})

	logEvent := event.NewEvent(event.Log, "i am log.")
	engine.Put(logEvent) //事件引擎关闭状态，丢弃事件。 {2 i am log.}

	engine.StartAll()

	time.Sleep(3 * time.Second)

	engine.UnRegister(event.Timer, timerEventHandler)
	time.Sleep(3 * time.Second)

	//	engine.Put(logEvent) //

	//	engine.StopAll()

	//engine.stopEventConsumer()
	//engine.StopTimer()

	engine.Put(logEvent)

	time.Sleep(5 * time.Second)

}
