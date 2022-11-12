package core

import (
	"leopard-quant/core/event"
	"testing"
)

func TestEngine_Start(t *testing.T) {
	type fields struct {
		EngineName  string
		EventEngine *event.Engine
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{name: "测试输出", fields: struct {
			EngineName  string
			EventEngine *event.Engine
		}{EngineName: "引擎名称", EventEngine: event.NewEventEngine()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Engine{
				EngineName:  tt.fields.EngineName,
				EventEngine: tt.fields.EventEngine,
			}
			e.Start()
		})
	}
}
