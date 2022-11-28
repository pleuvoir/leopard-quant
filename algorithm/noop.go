package algorithm

import (
	"fmt"
	"leopard-quant/common/model"
	"leopard-quant/core/log"
	"time"
)

type NoopSub struct {
	cnt int
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
}

func (n *NoopSub) OnStop(c Context) {
}

func (n *NoopSub) OnTrade(c Context, trade model.Trade) {
}

func (n *NoopSub) OnTick(c Context, t model.Ticker) {
	n.cnt = n.cnt + 1
	log.Infoln(fmt.Sprintf("NoopSub OnTick，%+v cnt=%d", t, n.cnt))
	time.Sleep(time.Second * 5)
}

func (n *NoopSub) OnOrder(c Context, order model.Order) {
}

func (n *NoopSub) Name() string {
	return "noop"
}
