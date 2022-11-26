package gateway

import (
	"leopard-quant/gateway"
	"testing"
)

func TestNewMe(t *testing.T) {

	api := gateway.NewBaseApi(gateway.WithConfig("/Users/pleuvoir/dev/space/git/leopard-quant/test/gateway/okx/okx.yml"))
	api.WithUnmarshalerOption(gateway.WithResponseUnmarshaler(func(bytes []byte, a any) error {
		return nil
	}))

}
