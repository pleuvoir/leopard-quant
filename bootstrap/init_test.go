package bootstrap

import (
	"leopard-quant/core/log"
	"testing"
)

func TestInit(t *testing.T) {

	conf := Global.ApplicationConf
	t.Logf("%+v", conf)

	log.Info("get")
	log.Error("i am error")
	log.Println("print")

	log.Warnf("i am warn %s", "122")
}
