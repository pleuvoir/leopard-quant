

### 策略实现

需要实现`TemplateSub`接口。

```go
type TemplateSub interface {
	OnStart(c Context)
	OnStop(c Context)
	OnTimer(c Context)
	OnTrade(c Context, trade model.Trade)
	OnTick(c Context, ticker model.Ticker)
	OnOrder(c Context, order model.Order)
	Name() string
}
```


### 事件引擎

分为普通事件和定时事件。发布事件时不同的种类会被路由到不同的管道，不同类型的事件互不干扰，也因此发布事件是异步的。

由于管道自带阻塞特性，当监听器处理任务执行耗时较长时，发布事件会阻塞。因此实际的任务处理也设计为异步。为了保证不创建过多的协程，选择使用协程池控制。
需要注意的是，同一种类型的事件，如果有多个处理器。则会顺序遍历执行，即他们也会互相影响。

* 普通事件回调会在协程池中异步执行
* 定时事件回调在管道中同步执行

![eventEngine](https://github.com/pleuvoir/leopard-quant/raw/master/docs/eventEngine.png)


### 网关引擎

负责以一套`API`对接多个`Broker`，屏蔽其差异。一般的通信方式为`REST/WebSocket`。举例当接收到`Broker`的`WebSocket Ticker`回调时，会使用事件引擎发布`Tick`事件。
总之，`Broker`状态的变化，网关引擎可以通过发布事件的形式告知其他组件。

### 算法引擎

算法引擎负责与网关交互，监听网关发布的`Tick/Trade`等事件。

#### 算法模板

算法模板则直接作用于算法的实现。为什么需要算法模板这一层，因为不同的算法可能会订阅同一个币种，算法引擎在回调时可以遍历这些模板进行回调。同时将一些操作封装在模板这一层，降低对算法引擎的感知。

注意：如果回调事件中`IO`过长，而回调速度过快，会存在并发的问题，需要用户注意。

### 主引擎

封装了其它引擎中的操作，各组件可调用主引擎访问其他引擎提供的功能。