package gateway

import (
	"github.com/tidwall/gjson"
	"leopard-quant/common/model"
)

type TickerCallbackConverter func(data []byte, r gjson.Result) (ticker model.Ticker, err error)

type TickerCallback func(ticker model.Ticker)
type KlineCallback func(k model.KLine)

type ComposeCallback struct {
	TickerCallback
	TickerCallbackConverter
	KlineCallback
}
