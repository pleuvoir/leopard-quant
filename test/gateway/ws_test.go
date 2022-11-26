package gateway

import (
	"fmt"
	"github.com/tidwall/gjson"
	"leopard-quant/common/model"
	"leopard-quant/gateway"
	"strconv"
	"testing"
)

func TestWs(t *testing.T) {
	cfg := &gateway.Configuration{
		Addr:          "wss://ws.okx.com:8443/ws/v5/public",
		AutoReconnect: true,
		DebugMode:     false,
	}

	callback := gateway.ComposeCallback{
		TickerCallback: func(ticker model.Ticker) {
			fmt.Println(ticker)
		},
		TickerCallbackConverter: convert2Ticker,
	}

	b := gateway.New(cfg, callback)

	if err := b.Start(); err != nil {
		panic(err)
	}

	b.Subscribe(gateway.ArgItem{Channel: "tickers", InstId: "BTC-USDT"})

	forever := make(chan struct{})
	<-forever
}

func convert2Ticker(data []byte, r gjson.Result) (ticker model.Ticker, err error) {
	v := r.Get("data").Array()[0]
	ticker.Symbol = v.Get("instId").Str
	ticker.Last = parseFloat64(v.Get("last").Str)
	ticker.LastSz = parseFloat64(v.Get("lastSz").Str)
	ticker.AskPx = parseFloat64(v.Get("askPx").Str)
	ticker.AskSz = parseFloat64(v.Get("askSz").Str)
	ticker.BidPx = parseFloat64(v.Get("bidPx").Str)
	ticker.BidSz = parseFloat64(v.Get("bidSz").Str)
	ticker.Open24H = parseFloat64(v.Get("open24h").Str)
	ticker.High24H = parseFloat64(v.Get("high24h").Str)
	ticker.Low24H = parseFloat64(v.Get("low24h").Str)
	ticker.SodUtc0 = parseFloat64(v.Get("sodUtc0").Str)
	ticker.SodUtc8 = parseFloat64(v.Get("sodUtc8").Str)
	ticker.VolCcy24H = parseFloat64(v.Get("volCcy24h").Str)
	ticker.Vol24H = parseFloat64(v.Get("vol24h").Str)
	ticker.Ts = parseUint64(v.Get("ts").Str)
	return ticker, nil
}

func parseFloat64(val string) float64 {
	float, err := strconv.ParseFloat(val, 64)
	if err != nil {
		panic(err)
	}
	return float
}

func parseUint64(val string) uint64 {
	u, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		panic(err)
	}
	return u
}
