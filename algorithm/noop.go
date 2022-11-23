package algorithm

import (
	"leopard-quant/core/engine/model"
	"leopard-quant/core/log"
	"time"
)

type NoopSub struct {
	template *AlgoTemplate
}

func (n *NoopSub) OnStart(c Context) {
	n.template = c.Template
	err := n.template.Subscribe("STARL-USDT")
	if err != nil {
		log.Errorf("NoopSub OnStart，订阅STARL-USDT失败，%s ", err)
		panic(err)
	}
	log.Infof("NoopSub OnStart，已订阅STARL-USDT")
}

func (n *NoopSub) OnTimer() {
	log.Infof("NoopSub OnTimer")
	time.Sleep(time.Second * 5)
}

func (n *NoopSub) OnStop(c Context) {
}

func (n *NoopSub) OnTrade(trade model.Trade) {
}

func (n *NoopSub) OnTick(t model.Tick) {
	log.Infof("NoopSub OnTick，%+v", t)
	time.Sleep(time.Second * 5)

}

func (n *NoopSub) OnOrder(order model.Order) {
}

func (n *NoopSub) Name() string {
	return "noop"
}
