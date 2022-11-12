package config

import (
	"fmt"
	"testing"
)

func TestApplicationConf_Load(t *testing.T) {

	//os.Setenv(ApplicationEnvVar,"")
	conf := NewApplicationConf()
	if err := conf.Load(); err != nil {
		fmt.Println(err)
	}
}
