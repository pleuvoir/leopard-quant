package bootstrap

import (
	"leopard-quant/bootstrap"
	"testing"
	"time"
)

func TestInit(t *testing.T) {

	bootstrap.Init()
	//conf := bootstrap.Global.ApplicationConf
	//t.Logf("%+v", conf)
	//
	//log.Info("get")
	//log.Error("i am error")
	//log.Println("print")
	//
	//log.Warnf("i am warn %s", "122")

	time.Sleep(time.Second * 50)
}
