package gateway

import . "leopard-quant/common/model"

// UnmarshalerOptions 解析器
type UnmarshalerOptions struct {
	ResponseUnmarshaler          ResponseUnmarshaler
	GetTickerResponseUnmarshaler GetTickerResponseUnmarshaler
	GetKlineResponseUnmarshaler  GetKlineResponseUnmarshaler
}

func WithResponseUnmarshaler(u ResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.ResponseUnmarshaler = u
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

type ResponseUnmarshaler func([]byte, any) error
type GetTickerResponseUnmarshaler func([]byte) (Ticker, error)
type GetKlineResponseUnmarshaler func([]byte) (KLine, error)
