package event

import (
	"github.com/gookit/color"
	"sync"
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

// Engine 事件引擎
type Engine struct {
	Active         bool
	TimerActive    bool
	TimeDuration   time.Duration
	HandlersMap    map[Type][]eventHandler
	CommonHandlers []eventHandler
	EventChan      chan Event
	TimerEventChan chan Event
	stopChan       chan bool
	startMutex     sync.Mutex
	registerMutex  sync.Mutex
}

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

// NewEventEngine 创建引擎
func NewEventEngine() *Engine {
	engine := Engine{
		Active:         true,
		TimerActive:    true,
		TimeDuration:   time.Second,
		HandlersMap:    map[Type][]eventHandler{},
		CommonHandlers: []eventHandler{},
		EventChan:      make(chan Event, 1000),
		TimerEventChan: make(chan Event, 1000),
		stopChan:       make(chan bool),
	}
	return &engine
}

// Process 处理事件
func (e *Engine) Process(event Event) {
	eventHandlers := e.HandlersMap[event.EventType]
	for _, handler := range eventHandlers {
		handler.Process(event)
	}
	for _, handler := range e.CommonHandlers {
		handler.Process(event)
	}
}

// Register 注册事件处理器
func (e *Engine) Register(handler eventHandler) {
	e.registerMutex.Lock()
	defer e.registerMutex.Unlock()
	eventType := handler.GetType()
	HandlersMap := e.HandlersMap
	eventHandlers := HandlersMap[eventType]
	eventHandlers = append(eventHandlers, handler)
	HandlersMap[eventType] = eventHandlers
}

// UnRegister 取消事件处理器
func (e *Engine) UnRegister(handler eventHandler) {
	e.registerMutex.Lock()
	defer e.registerMutex.Unlock()
	eventType := handler.GetType()
	handlersMap := e.HandlersMap
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

// StopAll 停止所有
func (e *Engine) StopAll() {
	e.StopEventConsumer()
	e.StopSchedulerTimer()
}

// StopSchedulerTimer 停止周期引擎
func (e *Engine) StopSchedulerTimer() {
	e.startMutex.Lock()
	defer e.startMutex.Unlock()
	e.stopChan <- true
}

// StopEventConsumer 停止普通事件引擎
func (e *Engine) StopEventConsumer() {
	e.startMutex.Lock()
	defer e.startMutex.Unlock()
	e.Active = false
	close(e.EventChan)
}

// StartConsumer 消费消息  普通消息和定时器消息分开处理
func (e *Engine) StartConsumer() {
	e.startMutex.Lock()
	defer e.startMutex.Unlock()
	go func() {
	over:
		for e.Active {
			select {
			case event, ok := <-e.EventChan:
				if !ok {
					color.Redln("事件引擎已关闭，普通消息消费已终止，不再接收事件。")
					break over
				}
				e.Process(event)
			}
		}
	}()
	go func() {
	over:
		for e.TimerActive {
			select {
			case event, ok := <-e.TimerEventChan:
				if !ok {
					color.Redln("事件引擎已关闭，定时器已终止，不再接收事件。", event)
					break over
				}
				e.Process(event)
			}
		}
	}()
}

// Put 发布事件
func (e *Engine) Put(event Event) {
	if e.Active {
		e.EventChan <- event
	} else {
		color.Redln("事件引擎已关闭，不再接收事件。", event)
	}
}

// StartSchedulerTimer 启动定时器，周期执行
func (e *Engine) StartSchedulerTimer() {
	e.startMutex.Lock()
	defer e.startMutex.Unlock()
	go func() {
		newEvent := NewEvent(Timer, nil)
		ticker := time.NewTicker(e.TimeDuration)
		defer ticker.Stop()
		for e.TimerActive {
			select {
			case <-ticker.C:
				e.TimerEventChan <- newEvent
			case <-e.stopChan:
				e.TimerActive = false
				color.Redln("事件引擎已关闭，定时器已终止，不再接收事件。")
			}
		}
	}()

}
