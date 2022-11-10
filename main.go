package main

import (
	"leopard-quant/core/event"
	"time"
)

func main() {

	engine := event.NewEventEngine()
	engine.Register(event.EventLog, event.EventLogHandler{})

	newEvent := event.NewEvent(event.EventLog, "i am log.")
	engine.Put(newEvent)

	engine.StartConsumer()

	time.Sleep(3 * time.Second)
}
