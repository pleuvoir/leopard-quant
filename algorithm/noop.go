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
	symbol := c.GetStrOrDefault("symbol", "BTC-USDT")
	err := c.Subscribe(symbol)
	if err != nil {
		log.Errorf("NoopSub OnStart，订阅%s失败，%s ", symbol, err)
		panic(err)
	}
	log.Infof("NoopSub OnStart，已订阅 %s", symbol)
}

func (n *NoopSub) OnTimer(c Context) {
	log.Infof("NoopSub OnTimer")
}

func (n *NoopSub) OnStop(c Context) {
}

func (n *NoopSub) OnTrade(c Context, trade model.Trade) {
}

func (n *NoopSub) OnTick(c Context, t model.Tick) {
	log.Infof("NoopSub OnTick，%+v", t)
	time.Sleep(time.Second * 5)

}

func (n *NoopSub) OnOrder(c Context, order model.Order) {
}

func (n *NoopSub) Name() string {
	return "noop"
}
