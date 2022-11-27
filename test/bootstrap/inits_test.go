package bootstrap

import (
	"leopard-quant/bootstrap"
	"testing"
)

func TestInit(t *testing.T) {

	bootstrap.Init()

	forever := make(chan struct{})
	<-forever
}
