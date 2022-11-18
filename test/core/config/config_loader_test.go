package config

import (
	"fmt"
	"leopard-quant/core/config"
	"testing"
)

func TestApplicationConf_Load(t *testing.T) {

	//_ = os.Setenv(ApplicationEnvVar, "/Users/pleuvoir/dev/space/git/leopard-quant/build/application.yml")
	conf := config.NewApplicationConf()
	if err := conf.Load(); err != nil {
		fmt.Println(err)
	}
}
