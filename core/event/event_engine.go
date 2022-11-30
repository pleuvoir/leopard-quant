package event

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/panjf2000/ants"
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
	TimeDuration   time.Duration
	HandlersMap    map[Type][]eventListener
	CommonHandlers []eventListener
	TimerEventChan chan Event
	stopChan       chan struct{}
	registerMutex  sync.Mutex
	queues         map[Type]*eventQueue
}

// 事件处理器
type eventListener interface {
	Process(event Event)
}

// AdaptEventHandlerFunc 接口适配器
// 可以将函数原型为 fun(event Event)的函数直接做为eventHandler接口的实现进行传入
// 可以无需定义结构体
type AdaptEventHandlerFunc func(e Event)

func (funcW AdaptEventHandlerFunc) Process(event Event) {
	funcW(event)
}

type Event struct {
	EventType Type
	EventData any
}

// NewEvent 新事件
func NewEvent(eventType Type, data any) Event {
	return Event{eventType, data}
}

// 每个事件使用单独的队列
type eventQueue struct {
	ch   chan Event
	t    Type
	pool *ants.Pool
}

func (q *eventQueue) shutdown() {
	defer func() {
		close(q.ch)
		_ = q.pool.Release()
	}()
}

func (q *eventQueue) send(e Event) {
	q.ch <- e
}

func newEventQueue(t Type) *eventQueue {
	pool, _ := ants.NewPool(100) //可以控制有多少协程并发处理任务
	return &eventQueue{ch: make(chan Event), t: t, pool: pool}
}

// NewEventEngine 创建引擎
func NewEventEngine() *Engine {
	engine := Engine{
		Active:         false,
		TimeDuration:   time.Second,
		HandlersMap:    map[Type][]eventListener{},
		CommonHandlers: []eventListener{},
		TimerEventChan: make(chan Event, 1000),
		stopChan:       make(chan struct{}),
		queues:         map[Type]*eventQueue{},
	}
	types := []Type{Log, Tick, Bar, Trade, Order, Asset, Position, Contract, Account, Algo, Error}
	for _, t := range types {
		engine.queues[t] = newEventQueue(t)
	}
	return &engine
}

// Process 处理事件
func (e *Engine) Process(event Event) {
	eventHandlers := e.HandlersMap[event.EventType]
	for _, handler := range e.CommonHandlers {
		handler.Process(event)
	}
	for _, handler := range eventHandlers {
		handler.Process(event)
	}
}

// Register 注册事件处理器
func (e *Engine) Register(eventType Type, handler eventListener) {
	e.registerMutex.Lock()
	defer e.registerMutex.Unlock()
	HandlersMap := e.HandlersMap
	eventHandlers := HandlersMap[eventType]
	eventHandlers = append(eventHandlers, handler)
	HandlersMap[eventType] = eventHandlers
}

// UnRegister 取消事件处理器
func (e *Engine) UnRegister(eventType Type, handler eventListener) {
	e.registerMutex.Lock()
	defer e.registerMutex.Unlock()
	handlersMap := e.HandlersMap
	eventHandlers := handlersMap[eventType]
	var newHandlers []eventListener
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
	if !e.Active {
		return
	}
	e.stopConsumer()
	e.stopTimer()
}

// 停止周期引擎
func (e *Engine) stopTimer() {
	e.stopChan <- struct{}{}
}

// 停止普通事件引擎
func (e *Engine) stopConsumer() {
	if e.Active {
		e.Active = false
		for _, e := range e.queues {
			e.shutdown()
		}
	}
}

// StartAll 启动所有
func (e *Engine) StartAll() {
	if e.Active {
		return
	}
	e.Active = true
	e.startTimer()
	e.startConsumer()
}

// 消费消息
func (e *Engine) startConsumer() {
	for _, eq := range e.queues {
		go func(q *eventQueue) {
		over:
			for e.Active {
				select {
				case event, ok := <-q.ch:
					if !ok {
						color.Redln(fmt.Sprintf("[%d]子事件引擎接收到关闭信号，终止事件监听。", q.t))
						break over
					}
					err := q.pool.Submit(func() {
						e.Process(event)
					})
					if err != nil {
						color.Redln(fmt.Sprintf("[%d]子事件引擎在协程池中处理任务失败。err=%s", q.t, err))
					}
				}
			}
		}(eq)
	}
	color.Greenln("事件引擎子事件引擎已全部启动。")
}

func (e *Engine) startTimer() {
	e.startTimerProducer()
	e.startTimerConsumer()
}

// 启动定时器消费者
func (e *Engine) startTimerConsumer() {
	go func() {
	outer:
		for {
			select {
			case event, ok := <-e.TimerEventChan:
				if !ok {
					color.Redln("事件引擎定时器消费者接收到关闭信号，已终止事件监听。")
					break outer
				}
				e.Process(event)
			}
		}
	}()
	color.Greenln("事件引擎定时发布消费者已启动。")
}

// 启动定时器生产者，周期执行
func (e *Engine) startTimerProducer() {
	go func() {
		newEvent := NewEvent(Timer, nil)
		ticker := time.NewTicker(e.TimeDuration)
		defer ticker.Stop()
	outer:
		for {
			select {
			case <-ticker.C:
				e.TimerEventChan <- newEvent
			case <-e.stopChan:
				close(e.TimerEventChan)
				color.Redln("事件引擎定时发布生产者接收到关闭信号，定时器已终止，不再发布时间事件。")
				break outer
			}
		}
	}()
	color.Greenln("事件引擎定时发布生产者已启动。")
}

// Put 发布事件，因为管道自带阻塞特性，为避免满后阻塞，因此没有消费者时不让发布
func (e *Engine) Put(event Event) {
	if e.Active {
		e.queues[event.EventType].send(event)
	} else {
		color.Redln(fmt.Sprintf("事件引擎处于关闭状态，丢弃事件发布。%+v", event))
	}
}
