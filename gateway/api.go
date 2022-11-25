package gateway

import "leopard-quant/common/model"

type TickerCallback func(ticker model.Ticker)
type KlineCallback func(k model.KLine)

type ComposeCallback struct {
	TickerCallback
	KlineCallback
}
