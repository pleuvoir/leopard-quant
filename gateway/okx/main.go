package okx

import "leopard-quant/gateway"

type OKX struct {
	baseApi   *gateway.BaseApi
	marketApi *MarketApi
}

func New(options *gateway.ApiOptions) *OKX {
	okx := &OKX{}
	baseApi := gateway.NewBaseApiWithOptions(options)
	baseApi.WithUnmarshalerOption(
		gateway.WithGetTickerResponseUnmarshaler(GetTickerResponseUnmarshaler),
		gateway.WithGetKlineResponseUnmarshaler(GetKlineResponseUnmarshaler),
		gateway.WithPongResponseUnmarshaler(PongResponseUnmarshaler),
		gateway.WithSubscribeResponseUnmarshaler(SubscribeResponseUnmarshaler),
	)
	okx.baseApi = baseApi
	return okx
}

func (o *OKX) Subscribe(symbol string, c gateway.ApiCallback) error {
	marketApi := NewMarket(o.baseApi, c)
	if err := marketApi.Start(); err != nil {
		return err
	}
	marketApi.Subscribe(ArgItem{Channel: "tickers", InstId: symbol})
	marketApi.Subscribe(ArgItem{Channel: "candle15m", InstId: symbol})
	o.marketApi = marketApi
	return nil
}

func (o *OKX) CancelSubscribe(symbol string) error {
	o.marketApi.UnSubscribe(ArgItem{Channel: "tickers", InstId: symbol})
	return nil
}

func (o *OKX) Name() string {
	return "okx"
}
