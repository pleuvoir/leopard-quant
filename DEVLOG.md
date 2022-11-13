



# 开发环境

设置代理 GOPROXY=https://goproxy.io


## 创建事件

一个事件类型有多个处理器，需要考虑并发的情况。多个处理器使用切片，但是切片的删除好像没有直接的API。此外，type：List[typeHandler]这种 必须使用线程安全的sync.map
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