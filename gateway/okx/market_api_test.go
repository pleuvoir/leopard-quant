package okx

import (
	"leopard-quant/common/model"
	"leopard-quant/gateway"
	"testing"
	"time"
)

func TestMarketApi(t *testing.T) {
	okx := New("/Users/pleuvoir/dev/space/git/leopard-quant/build/gateway.yml")

	callback := gateway.ApiCallback{
		TickerCallback: func(ticker model.Ticker) {
			t.Log(ticker)
		},
		KlineCallback: func(line model.KLine) {
			t.Log(line)
		},
	}

	if err := okx.Subscribe("BTC-USDT", callback); err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 5)

	<-make(chan struct{})
}
