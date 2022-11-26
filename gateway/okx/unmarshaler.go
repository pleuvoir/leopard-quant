package okx

import (
	"errors"
	"github.com/tidwall/gjson"
	"leopard-quant/common/model"
)

func GetTickerResponseUnmarshaler(data []byte) (ticker model.Ticker, err error) {
	r := gjson.ParseBytes(data)

	if channelType := r.Get("arg.channel"); channelType.Exists() && channelType.Str == "tickers" {
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
	return ticker, errors.New("不是ticker数据")
}

func GetKlineResponseUnmarshaler(data []byte) (kLine model.KLine, err error) {
	return kLine, err
}
