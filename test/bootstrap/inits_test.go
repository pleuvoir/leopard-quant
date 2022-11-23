package bootstrap

import (
	"leopard-quant/bootstrap"
	"testing"
	"time"
)

func TestInit(t *testing.T) {

	bootstrap.Init()

	time.Sleep(time.Second * 100)
}
