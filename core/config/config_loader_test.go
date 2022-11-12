package config

import (
	"fmt"
	"os"
	"testing"
)

func TestApplicationConf_Load(t *testing.T) {

	//os.Setenv(ApplicationEnvProfileVar, "test")
	os.Setenv(ApplicationEnvVar, "/Users/pleuvoir/dev/space/git/leopard-quant/build/application-prod.yml")
	conf := NewApplicationConf()
	if err := conf.Load(); err != nil {
		fmt.Println(err)
	}
}
