package okx

import "leopard-quant/gateway"

var BaseApi *gateway.BaseApi

func init() {
	baseApi := gateway.NewBaseApi(gateway.WithConfig("okx.yml"))
	baseApi.WithUnmarshalerOption(
		gateway.WithGetTickerResponseUnmarshaler(GetTickerResponseUnmarshaler),
		gateway.WithGetKlineResponseUnmarshaler(GetKlineResponseUnmarshaler))
	BaseApi = baseApi
}
