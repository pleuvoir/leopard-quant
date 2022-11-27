# 开发环境

设置代理 GOPROXY=https://goproxy.io

## 创建事件

一个事件类型有多个处理器，需要考虑并发的情况。多个处理器使用切片，但是切片的删除好像没有直接的API。此外，type：List[typeHandler]
这种 必须使用线程安全的sync.map
但是它又不支持泛型，转换的时候很麻烦。

```go
handlers, ok := ee.Handlers.Load(event.EventType)
if !ok {
fmt.Printf("未找到处理器，event is %+v", event)
return
}
//转换类型，因为并发安全的MAP里保存的是ANY
eventHandlers := handlers.([]eventHandler)
```

在单元测试中不能立马收到消息，但是在main方法中可以。不知道为什么。很神奇，如果在每个里面输出时间，它就能正常打出来。

## 加载配置文件

可以根据当前项目根路径/当前项目可执行路径/环境变量指定/profile 读取全局文件配置

在GOLAND中EDIT CONFIGURATION设置环境变量，这样就可以使用配置文件了

## 引擎

时间格式化，在GOLAND中输入 YY 什么的就会有提示

https://zhuanlan.zhihu.com/p/28441006

## 引入GIN

https://github.com/gin-gonic/gin#using-get-post-put-patch-delete-and-options

## 转换sync.map

TODO
无法得到map中所有元素，需要遍历一次加入到普通的map中，需要研究泛型 提取公共方法以便获取所有的值

## 抽象接口

可以看这篇文章，继承重写都讲了。
https://zhuanlan.zhihu.com/p/88480107

需要注意，GOLAND里不会再提示可以重写这个方法，需要自己手动覆盖，(receiver) 必须是子类自己

    var aa IEngine = engine
    aa.Start()

可以这样转换，先这样。父类符合接收子类做为形参传入。非多态问题
https://zhuanlan.zhihu.com/p/133693915  继承这里应该不能这么用 明天再看

已经通过组合的方式实现，不必纠结。参考 NewConfigLoader 我觉得很完美，如果抽象类也实现接口了就用这个，把子类的行为包装进去。否则要实现的接口太多了。

## 接口适配器

遇到的问题是在订单引擎中注册事件，注册事件必须声明结构体很麻烦，并且事件监听中是需要操作订单引擎中的变量。问题就变成了先要创建结构体，并且结构体持有订单的引用。

```go
// 事件处理器
type eventHandler interface {
Process(event Event)
}

// 接口适配器
// 可以将函数原型为 fun(event Event)的函数直接做为eventHandler接口的实现进行传入
// 可以无需定义结构体
type AdaptEventHandlerFunc func (event Event)

func (funcW AdaptEventHandlerFunc) Process(event Event) {
funcW(event)
}

func Process(event Event) {
fmt.Printf("LogEventHandler receive event %+v \n", event)
}

func TestHandlerFunc(t *testing.T) {
handlerFunc := AdaptEventHandlerFunc(Process)
engine := NewEventEngine()
engine.Register(Log, handlerFunc)
}

```

有了上一步，可以通过匿名函数的形式操作当前作用域下的变量了，不用在传了。

```go
func (o *OrderEngine) registerEvent() {
o.registerListener(event.Trade, func (event event.Event) {
trade, _ := event.EventData.(model.Trade)
o.tradeMap[event.EventType] = trade
})
}

func (o *OrderEngine) registerListener(t event.Type, f func (e event.Event)) {
o.mainEngine.eventEngine.Register(t, event.AdaptEventHandlerFunc(f))
}
```

## 循环依赖 && 反射创建对象

```
package leopard-quant/core/engine
	imports leopard-quant/algorithm
	imports leopard-quant/core/engine: import cycle not allowed
```

可以将接口定义在engine包中，在算法包中实现接口。但是必须所有的接口方法都导出（大写），外部的包才算实现了这个接口。

反射的时候存进去的是指针类型，如果取取出来不是这个类型不会报错，但是会是nil
https://stackoverflow.com/questions/35790935/using-reflection-in-go-to-get-the-name-of-a-struct

## 对象池，降低GC压力

在事件引擎中使用，暂未进行benchmark
https://geektutu.com/post/hpg-sync-pool.html

## 对接OK KLINE

//TODO
重新抽象ws 获取不到kline回调  不能重复readMessage 

gjson一个根据路径获取json值的库，类似于jsonObject



## 	协程池 ants.Pool 限制回调异步任务的创建

在执行处理器时，已经按照事件类型区分了管道，但是一个管道中有阻塞任务 还是会影响别人，因此在处理的时候  启动新协程，完全变成异步。
问题是协程会创建比较多，因此使用协程池限制创建的数量


## toml 格式文件解析

