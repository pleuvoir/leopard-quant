package event

import (
	"fmt"
	"time"
)

// 事件类型
type EventType int

const (
	eventTimer EventType = iota + 1
	EventLog
	eventTick
	eventBar
	eventTrade
	eventOrder
	eventAsset
	eventPosition
	eventContract
	eventAccount
	eventAlgo
	eventError
)

// 事件处理器
type EventHandler interface {
	Process(event Event)
}

type Event struct {
	EventType EventType
	EventData any
}

// NewEvent 新事件
func NewEvent(eventType EventType, eventData any) Event {
	return Event{eventType, eventData}
}

// 事件引擎
type eventEngine struct {
	Active         bool
	TimerActive    bool
	TimeDuration   time.Duration
	HandlersMap    map[EventType][]EventHandler
	CommonHandlers []EventHandler
	EventChan      chan Event
}

// NewEventEngine 创建引擎
func NewEventEngine() *eventEngine {
	engine := eventEngine{
		Active:         true,
		TimerActive:    false,
		TimeDuration:   time.Second,
		HandlersMap:    map[EventType][]EventHandler{},
		CommonHandlers: []EventHandler{},
		EventChan:      make(chan Event, 1000),
	}
	return &engine
}

// Process 处理事件
func (ee *eventEngine) Process(event Event) {
	eventHandlers := ee.HandlersMap[event.EventType]
	for _, handler := range eventHandlers {
		handler.Process(event)
	}
	for _, handler := range ee.CommonHandlers {
		handler.Process(event)
	}
}

// Register 注册事件处理器
func (ee *eventEngine) Register(eventType EventType, handler EventHandler) {
	handlers := ee.HandlersMap
	eventHandlers := handlers[eventType]
	eventHandlers = append(eventHandlers, handler)
	handlers[eventType] = eventHandlers
}

// UnRegister 取消事件处理器
func (ee *eventEngine) UnRegister(eventType EventType, handler EventHandler) {
	handlersMap := ee.HandlersMap
	eventHandlers := handlersMap[eventType]

	var newHandlers []EventHandler
	for _, cur := range eventHandlers {
		if cur == handler {
			continue
		}
		newHandlers = append(newHandlers, cur)
	}

	handlersMap[eventType] = newHandlers

	//没有处理器则将这个类型移除
	if len(eventHandlers) == 0 {
		delete(handlersMap, eventType)
	}
}

// StartConsumer 消费消息
func (ee *eventEngine) StartConsumer() {
	go func() {
		for {
			select {
			case e, ok := <-ee.EventChan:
				if !ok {
					return
				}
				ee.Process(e)
			}
		}
	}()
}

// Put 发布事件
func (ee *eventEngine) Put(event Event) {
	ee.EventChan <- event
}

// StartSchedulerTimer 启动定时器，周期执行
func (ee *eventEngine) StartSchedulerTimer() {
	newEvent := NewEvent(eventTimer, nil)
	for ee.TimerActive {
		ee.EventChan <- newEvent
		time.Sleep(ee.TimeDuration)
	}
}

type EventLogHandler struct{}

func (EventLogHandler) Process(event Event) {
	fmt.Printf("received event ===>  %+v", event)
}
