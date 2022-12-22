package bootstrap

import (
	"leopard-quant/bootstrap"
	"leopard-quant/restful"
	"testing"
)

func TestInit(t *testing.T) {

	bootstrap.Init()

	restful.NewServer().ServerStartedListener(nil).Start()

	//
	//forever := make(chan struct{})
	//<-forever
}
