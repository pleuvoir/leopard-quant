package algorithm

import (
	"leopard-quant/algorithm/base"
	"testing"
)

func TestGet(t *testing.T) {

	if sub, err := base.MakeInstance("noop"); err == nil {

		t.Logf("%+v", sub)

		t.Log(sub.Name())
	}
}

func TestName(t *testing.T) {

}
