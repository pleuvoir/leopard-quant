package event

import (
	"fmt"
	"sync"
	"time"
)

// 事件类型
type eventType int

const (
	eventTimer eventType = iota + 1
	eventLog
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
type eventHandler interface {
	process(event event)
}

type event struct {
	EventType eventType
	EventData any
}

// NewEvent 新事件
func NewEvent(eventType eventType, eventData any) event {
	return event{eventType, eventData}
}

// 事件引擎
type eventEngine struct {
	EventType      eventType
	Active         bool
	TimerActive    bool
	TimeDuration   time.Duration
	Handlers       sync.Map //[eventType][]eventHandler
	CommonHandlers []eventHandler
	EventChan      chan event
}

func NewEventEngine() {

}

func (ee *eventEngine) Process(event event) {
	handlers, ok := ee.Handlers.Load(event.EventType)
	if !ok {
		fmt.Printf("未找到处理器，event is %+v", event)
		return
	}
	//转换类型，因为并发安全的MAP里保存的是ANY
	eventHandlers := handlers.([]eventHandler)
	for _, handler := range eventHandlers {
		handler.process(event)
	}
	for _, handler := range ee.CommonHandlers {
		handler.process(event)
	}
}

// Register 注册事件处理器
func (ee *eventEngine) Register(eventType eventType, handler eventHandler) {
	handlers, _ := ee.Handlers.Load(eventType)
	eventHandlers := handlers.([]eventHandler)
	if eventHandlers[eventType] == nil {
		eventHandlers = append(eventHandlers, handler)
		ee.Handlers.Store(eventType, eventHandlers)
	}
}

func (ee *eventEngine) UnRegister(eventType eventType, handler eventHandler) {
	handlers, _ := ee.Handlers.Load(eventType)
	eventHandlers := handlers.([]eventHandler)

	cur := eventHandlers[eventType]
	if cur != nil {
		//移除 TODO
		eventHandlers = append(eventHandlers[:1], eventHandlers[2:]...)
		ee.Handlers.Store(eventType, eventHandlers)
	}

	//没有处理器则将这个类型移除
	if len(eventHandlers) == 0 {
		ee.Handlers.Delete(eventType)
	}

}

// Put 发布事件
func (ee *eventEngine) Put(event event) {
	ee.EventChan <- event
}

// StartSchedulerTimer 启动定时器，周期执行
func (ee *eventEngine) StartSchedulerTimer() {
	duration := ee.TimeDuration
	e := event{EventType: eventTimer, EventData: nil}
	for ee.TimerActive {
		ee.EventChan <- e
		time.Sleep(duration)
	}
}
