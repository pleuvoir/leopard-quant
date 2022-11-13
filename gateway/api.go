package gateway

// IHttpClient 请求接口
type IHttpClient interface {
	Request(method, url string, body string, headers map[string]string) (data []byte, err error)
}

// IMarketRest 行情接口，不需要鉴权
type IMarketRest interface {
	GetExchangeName() string                                                                 //获取交易所名字/域名
	GetDepth(pair CurrencyPair, limit int, opt ...OptionParameter) (*Depth, error)           //获取深度信息
	GetTicker(pair CurrencyPair, opt ...OptionParameter) (*Ticker, error)                    //获取TICKER
	GetKline(pair CurrencyPair, period KlinePeriod, opt ...OptionParameter) ([]Kline, error) //获取蜡烛图
}
