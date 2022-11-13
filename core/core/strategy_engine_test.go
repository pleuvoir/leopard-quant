package core

import "testing"

func TestStrategyTemplate_OnTick(t *testing.T) {
	engine := NewStrategyEngine()

	engine.initEngine()
	engine.RegisterEvent()

	engine.Start()

	ch := make(chan any)
	<-ch
}
