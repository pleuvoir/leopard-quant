package bootstrap

import "leopard-quant/core/config"

func init() {

	globalConfig := &config.ApplicationConf{}
	err := globalConfig.Load()
	if err != nil {
		panic(err)
	}
}
