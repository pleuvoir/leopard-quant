package factory

import (
	"fmt"
	"leopard-quant/gateway"
	"leopard-quant/gateway/okx"
)

func NewGateway(name string, options *gateway.ApiOptions) gateway.Sub {

	switch name {
	case "okx":
		return okx.New(options)
	}

	panic(fmt.Sprintf("不支持的交易所 %s", name))
}
