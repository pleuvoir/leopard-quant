package gateway

import "leopard-quant/common/model"

type TickerCallback func(ticker model.Ticker)
type KlineCallback func(kLine model.KLine)

type ApiCallback struct {
	TickerCallback
	KlineCallback
}
