package config

import (
	"fmt"
	"testing"
)

func TestApplicationConf_Load(t *testing.T) {

	//_ = os.Setenv(ApplicationEnvVar, "/Users/pleuvoir/dev/space/git/leopard-quant/build/application.yml")
	conf := NewApplicationConf()
	if err := conf.Load(); err != nil {
		fmt.Println(err)
	}
}
