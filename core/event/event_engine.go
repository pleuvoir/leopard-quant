package event

import (
	"time"
)

// Type 事件类型
type Type int

const (
	Timer Type = iota + 1
	Log
	Tick
	Bar
	Trade
	Order
	Asset
	Position
	Contract
	Account
	Algo
	Error
)

// 事件处理器
type eventHandler interface {
	Process(event Event)
	GetType() Type
}

type Event struct {
	EventType Type
	EventData any
}

// NewEvent 新事件
func NewEvent(eventType Type, data any) Event {
	return Event{eventType, data}
}

// Engine 事件引擎
type Engine struct {
	Active         bool
	TimerActive    bool
	TimeDuration   time.Duration
	HandlersMap    map[Type][]eventHandler
	CommonHandlers []eventHandler
	EventChan      chan Event
}

// NewEventEngine 创建引擎
func NewEventEngine() *Engine {
	engine := Engine{
		Active:         true,
		TimerActive:    true,
		TimeDuration:   time.Second,
		HandlersMap:    map[Type][]eventHandler{},
		CommonHandlers: []eventHandler{},
		EventChan:      make(chan Event, 1000),
	}
	return &engine
}

// Process 处理事件
func (ee *Engine) Process(event Event) {
	eventHandlers := ee.HandlersMap[event.EventType]
	for _, handler := range eventHandlers {
		handler.Process(event)
	}
	for _, handler := range ee.CommonHandlers {
		handler.Process(event)
	}
}

// Register 注册事件处理器
func (ee *Engine) Register(handler eventHandler) {
	eventType := handler.GetType()
	HandlersMap := ee.HandlersMap
	eventHandlers := HandlersMap[eventType]
	eventHandlers = append(eventHandlers, handler)
	HandlersMap[eventType] = eventHandlers
}

// UnRegister 取消事件处理器
func (ee *Engine) UnRegister(handler eventHandler) {
	eventType := handler.GetType()
	handlersMap := ee.HandlersMap
	eventHandlers := handlersMap[eventType]

	var newHandlers []eventHandler
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
func (ee *Engine) StartConsumer() {
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
func (ee *Engine) Put(event Event) {
	ee.EventChan <- event
}

// StartSchedulerTimer 启动定时器，周期执行
func (ee *Engine) StartSchedulerTimer() {
	go func() {
		newEvent := NewEvent(Timer, nil)
		for ee.TimerActive {
			ee.EventChan <- newEvent
			time.Sleep(ee.TimeDuration)
		}
	}()

}
