package loader

import (
	"flag"
	"leopard-quant/util"
)

const (
	symbol = "m"
	start  = "s"
	end    = "e"
	period = "p"
)

func Load() {

	param := struct {
		Symbol string
		Start  string
		End    string
		Period string
	}{}

	flag.StringVar(&param.Symbol, symbol, "", "请输入币种代号，如 BTC/USDT")
	flag.StringVar(&param.Symbol, start, "", "请输入开始时间，格式:yyyyMMdd")
	flag.StringVar(&param.Symbol, end, "", "请输入结束时间，格式:yyyyMMdd")
	flag.StringVar(&param.Symbol, period, "", "请输入K线类型，day/min1/min30/min60")

	flag.Parse()

	blank := util.IsAnyBlank(param.Symbol, param.Start, param.End, param.Period)
	if blank {
		flag.Usage()
		return
	}

}
