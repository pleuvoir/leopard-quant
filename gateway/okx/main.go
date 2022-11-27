package okx

import "leopard-quant/gateway"

type OKX struct {
	baseApi   *gateway.BaseApi
	marketApi *MarketApi
}

func New(configPath string) *OKX {
	okx := &OKX{}
	baseApi := gateway.NewBaseApi(gateway.WithConfig(configPath))
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
	return nil
}

func (o *OKX) CancelSubscribe(symbol string) error {
	panic("implement me")
}
