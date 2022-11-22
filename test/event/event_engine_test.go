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
	fmt.Printf("LogEventHandler receive event %+v \n", event)
}

func Process(event event.Event) {
	fmt.Printf("LogEventHandler receive event %+v \n", event)
}

func TestHandlerFuncWithReceiver(t *testing.T) {
	engine := event.NewEventEngine()
	engine.Register(event.Log, event.AdaptEventHandlerFunc(Process))

	engine.Put(event.NewEvent(event.Log, "i am log."))

	engine.StartConsumer()

	engine.Put(event.NewEvent(event.Log, "i am log2."))

	time.Sleep(5 * time.Second)

}

func TestHandlerFunc(t *testing.T) {
	engine := event.NewEventEngine()
	engine.Register(event.Log, event.AdaptEventHandlerFunc(Process))

	engine.Put(event.NewEvent(event.Log, "i am log."))

	engine.StartConsumer()

	engine.Put(event.NewEvent(event.Log, "i am log2."))

	time.Sleep(5 * time.Second)

}

func TestNewEventEngine(t *testing.T) {

	engine := event.NewEventEngine()

	timerEventHandler := TimerEventHandler{}
	engine.Register(event.Timer, timerEventHandler)
	engine.Register(event.Log, LogEventHandler{})

	logEvent := event.NewEvent(event.Log, "i am log.")
	engine.Put(logEvent) //事件引擎关闭状态，丢弃事件。 {2 i am log.}

	engine.StartConsumer()

	engine.StartSchedulerTimer()

	time.Sleep(3 * time.Second)

	//engine.UnRegister(timerEventHandler)

	engine.Put(logEvent) //

	engine.StopAll()

	//engine.StopEventConsumer()
	//engine.StopSchedulerTimer()

	engine.Put(logEvent)

	time.Sleep(5 * time.Second)

}