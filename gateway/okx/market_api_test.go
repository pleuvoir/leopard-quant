package okx

import (
	"fmt"
	"leopard-quant/common/model"
	"leopard-quant/gateway"
	"testing"
)

func TestMarketApi(t *testing.T) {

	callback := gateway.ApiCallback{
		TickerCallback: func(ticker model.Ticker) {
			fmt.Println(ticker)
		},
	}
	marketApi := NewMarket(nil, callback)

	if err := marketApi.Start(); err != nil {
		panic(err)
	}

	marketApi.Subscribe(ArgItem{Channel: "tickers", InstId: "BTC-USDT"})

	forever := make(chan struct{})
	<-forever
}
