package event

import (
	"testing"
	"time"
)

func TestNewEventEngine(t *testing.T) {

	engine := NewEventEngine()
	engine.Register(EventLog, EventLogHandler{})

	newEvent := NewEvent(EventLog, "i am log.")
	engine.Put(newEvent)

	engine.StartConsumer()

	time.Sleep(3 * time.Second)

}
