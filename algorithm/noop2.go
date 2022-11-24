package algorithm

import (
	"leopard-quant/core/engine/model"
	"leopard-quant/core/log"
	"time"
)

type Noop2Sub struct {
}

func (n *Noop2Sub) OnStart(c Context) {
	symbol := c.GetStr("symbol")
	err := c.Subscribe(symbol)
	if err != nil {
		log.Errorf("NoopSub2 OnStart，订阅%s失败，%s ", symbol, err)
		panic(err)
	}
	log.Infof("NoopSub2 OnStart，已订阅 %s", symbol)
}

func (n *Noop2Sub) OnTimer(c Context) {
	log.Infof("NoopSub2 OnTimer")
}

func (n *Noop2Sub) OnStop(c Context) {
}

func (n *Noop2Sub) OnTrade(c Context, trade model.Trade) {
}

func (n *Noop2Sub) OnTick(c Context, t model.Tick) {
	log.Infof("NoopSub OnTick，%+v", t)
	time.Sleep(time.Second * 5)

}

func (n *Noop2Sub) OnOrder(c Context, order model.Order) {
}

func (n *Noop2Sub) Name() string {
	return "noop2"
}
