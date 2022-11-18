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

func (n Noop) OnStart() {
	log.Info("noop OnStart")
}

func (n Noop) OnStop() {
	log.Info("noop OnStop")

}

func (n Noop) OnBar(bar model.Bar) {
	log.Info("noop OnBar %+v", bar)
}

func (n Noop) OnTick(tick model.Tick) {
	log.Info("noop OnTick %+v", tick)
}

func (n Noop) OnTimer() {
	log.Info("noop OnTimer %s", time.Now())
}
