package event

import (
	"fmt"
	"testing"
	"time"
)

type TimerEventHandler struct{}

func (h TimerEventHandler) WithType() Type {
	return Timer
}

func (TimerEventHandler) Process(event Event) {
	fmt.Printf("TimerEventHandler receive event %+v %s \n", event, time.Now())
}

type LogEventHandler struct{}

func (LogEventHandler) Process(event Event) {
	fmt.Printf("LogEventHandler receive event %+v \n", event)
}

func (LogEventHandler) WithType() Type {
	return Log
}

func TestNewEventEngine(t *testing.T) {

	engine := NewEventEngine()

	timerEventHandler := TimerEventHandler{}
	engine.Register(timerEventHandler)
	engine.Register(LogEventHandler{})

	logEvent := NewEvent(Log, "i am log.")
	engine.Put(logEvent)

	engine.StartConsumer()

	engine.StartSchedulerTimer()

	time.Sleep(3 * time.Second)

	//engine.UnRegister(timerEventHandler)

	engine.Put(logEvent)

	//engine.StopAll()

	engine.StopEventConsumer()
	engine.StopSchedulerTimer()

	engine.Put(logEvent)

	time.Sleep(5 * time.Second)

}
