package algorithm

import (
	"leopard-quant/core/engine/model"
	"leopard-quant/core/log"
	"time"
)

type Noop struct {
}

func (n Noop) Name() string {
	return "noop"
}

func (n Noop) OnStart(c context) {
	log.Info(c.getStr("name"))
	log.Info("noop OnStart")
}

func (n Noop) OnStop(c context) {
	log.Info("noop OnStop")

}

func (n Noop) OnBar(c context, bar model.Bar) {
	log.Info("noop OnBar %+v", bar)
}

func (n Noop) OnTick(c context, tick model.Tick) {
	log.Info("noop OnTick %+v", tick)
}

func (n Noop) OnTimer(c context) {
	log.Info("noop OnTimer %s", time.Now())
}
