package impl

import (
	"leopard-quant/core/engine"
	"leopard-quant/core/engine/model"
	"leopard-quant/core/log"
	"time"
)

type NoopSub struct {
	engine *engine.MainEngine
}

func (n *NoopSub) OnStart(c Context) {
	n.engine = c.MainEngine
	log.Infof("NoopSub OnStart，主引擎设置成功 %s", n.engine.TodayDate)
}

func (n *NoopSub) OnTimer() {
	log.Infof("NoopSub OnTimer，启动时间 start %s", n.engine.TodayDate)
	time.Sleep(time.Second * 5)
	log.Infof("NoopSub OnTimer，启动时间 over %s", n.engine.TodayDate)
}

func (n *NoopSub) OnStop(c Context) {
}

func (n *NoopSub) OnTrade(trade model.Trade) {
}

func (n *NoopSub) OnTick(t model.Tick) {
}

func (n *NoopSub) OnOrder(order model.Order) {
}

func (n *NoopSub) Name() string {
	return "noop"
}
