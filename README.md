

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

