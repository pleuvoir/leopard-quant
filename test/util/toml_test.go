package util

import (
	"github.com/BurntSushi/toml"
	"github.com/gookit/color"
	"leopard-quant/gateway"
	"testing"
)

func TestToml(t *testing.T) {

	config := make(map[string]*gateway.ApiOptions)

	if _, err := toml.DecodeFile("/Users/pleuvoir/dev/space/git/leopard-quant/build/gateway.toml", &config); err != nil {
		color.Redln("读取网关配置文件toml失败，%s", err)
		panic(err)
	}

}
