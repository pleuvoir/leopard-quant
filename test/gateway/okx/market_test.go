package okx

import (
	"leopard-quant/gateway"
	"leopard-quant/gateway/okx"
	"testing"
	"time"
)

func TestSubscribe(t *testing.T) {

	ws := okx.NewMarketWS()
	err := ws.Connect()
	if err != nil {
		t.Error(err)
		return
	}

	callback := gateway.ApiCallback{
		TickerCallback: func(tick gateway.Ticker) {
			t.Logf("%+v", tick)
		},

		KlineCallback: func(k gateway.Kline) {

		}}

	err = ws.Subscribe("STARL-USDT", callback)
	if err != nil {
		t.Error(err)
		return
	}

	time.Sleep(time.Second * 30)
}
