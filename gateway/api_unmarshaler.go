package gateway

import (
	. "leopard-quant/common/model"
)

// UnmarshalerOptions 解析器
type UnmarshalerOptions struct {
	SubscribeResponseUnmarshaler SubscribeResponseUnmarshaler
	PongResponseUnmarshaler      PongResponseUnmarshaler
	GetTickerResponseUnmarshaler GetTickerResponseUnmarshaler
	GetKlineResponseUnmarshaler  GetKlineResponseUnmarshaler
}

func WithSubscribeResponseUnmarshaler(u SubscribeResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.SubscribeResponseUnmarshaler = u
	}
}

func WithPongResponseUnmarshaler(u PongResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.PongResponseUnmarshaler = u
	}
}

func WithGetTickerResponseUnmarshaler(u GetTickerResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.GetTickerResponseUnmarshaler = u
	}
}

func WithGetKlineResponseUnmarshaler(u GetKlineResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.GetKlineResponseUnmarshaler = u
	}
}

type UnmarshalerOption func(options *UnmarshalerOptions)

type SubscribeResponseUnmarshaler func([]byte) (match bool, err error)
type PongResponseUnmarshaler func([]byte) (match bool, err error)
type GetTickerResponseUnmarshaler func([]byte) (match bool, ticker Ticker, err error)
type GetKlineResponseUnmarshaler func([]byte) (match bool, kLine KLine, err error)
