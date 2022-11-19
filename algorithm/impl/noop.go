package impl

import (
	"fmt"
	"leopard-quant/core/engine"
	"leopard-quant/core/engine/model"
	"leopard-quant/core/log"
)

var count int64 = 0

type NoopSub struct {
	engine *engine.MainEngine
	count  int64
}

func (n NoopSub) OnStart(c Context) {
	n.engine = c.MainEngine
	log.Infof("NoopSub OnStart，主引擎设置成功 %s", n.engine.TodayDate)
	fmt.Printf("%+v", n)

}

func (n NoopSub) OnStop() {
}

func (n NoopSub) OnTimer() {
	count++
	n.count = n.count + 1
	log.Infof("NoopSub OnTimer， %s", n.count)

}

func (n NoopSub) OnTrade(trade model.Trade) {
}

func (n NoopSub) OnTick(t model.Tick) {
}

func (n NoopSub) OnOrder(order model.Order) {
}

func (n NoopSub) Name() string {
	return "noop"
}
