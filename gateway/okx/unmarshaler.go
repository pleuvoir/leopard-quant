package okx

import (
	"github.com/tidwall/gjson"
	"leopard-quant/common/model"
	"leopard-quant/util/cast"
)

func GetTickerResponseUnmarshaler(data []byte) (match bool, ticker model.Ticker, err error) {
	r := gjson.ParseBytes(data)
	if channelType := r.Get("arg.channel"); channelType.Exists() && channelType.Str == "tickers" {
		v := r.Get("data").Array()[0]
		ticker.Symbol = v.Get("instId").Str
		ticker.Last = cast.ToFloat64(v.Get("last").Str)
		ticker.LastSz = cast.ToFloat64(v.Get("lastSz").Str)
		ticker.AskPx = cast.ToFloat64(v.Get("askPx").Str)
		ticker.AskSz = cast.ToFloat64(v.Get("askSz").Str)
		ticker.BidPx = cast.ToFloat64(v.Get("bidPx").Str)
		ticker.BidSz = cast.ToFloat64(v.Get("bidSz").Str)
		ticker.Open24H = cast.ToFloat64(v.Get("open24h").Str)
		ticker.High24H = cast.ToFloat64(v.Get("high24h").Str)
		ticker.Low24H = cast.ToFloat64(v.Get("low24h").Str)
		ticker.SodUtc0 = cast.ToFloat64(v.Get("sodUtc0").Str)
		ticker.SodUtc8 = cast.ToFloat64(v.Get("sodUtc8").Str)
		ticker.VolCcy24H = cast.ToFloat64(v.Get("volCcy24h").Str)
		ticker.Vol24H = cast.ToFloat64(v.Get("vol24h").Str)
		ticker.Ts = cast.ToUint64(v.Get("ts").Str)
		return true, ticker, nil
	}
	return false, ticker, nil
}

func GetKlineResponseUnmarshaler(data []byte) (match bool, kLine model.KLine, err error) {
	return false, kLine, nil
}

func PongResponseUnmarshaler(data []byte) (match bool, err error) {
	if s := string(data); s == "pong" {
		return true, nil
	}
	return false, nil
}

func SubscribeResponseUnmarshaler(data []byte) (match bool, err error) {
	ret := gjson.ParseBytes(data)
	if eventValue := ret.Get("event"); eventValue.Exists() {
		return true, nil
	}
	return false, nil
}
