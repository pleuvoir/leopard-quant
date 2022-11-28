package model

type Ticker struct {
	Symbol    string  `json:"symbol"`    //BTC-USDT
	Last      float64 `json:"last"`      //最新成交价
	LastSz    float64 `json:"lastSz"`    //最新成交的数量
	AskPx     float64 `json:"askPx"`     //卖一价
	AskSz     float64 `json:"askSz"`     //卖一价对应的量数量
	BidPx     float64 `json:"bidPx"`     //买一价
	BidSz     float64 `json:"bidSz"`     //买一价对应的数量
	Open24H   float64 `json:"open24h"`   //24小时开盘价
	High24H   float64 `json:"high24h"`   //24小时最高价
	Low24H    float64 `json:"low24h"`    //24小时最低价
	SodUtc0   float64 `json:"sodUtc0"`   //UTC 0 时开盘价
	SodUtc8   float64 `json:"sodUtc8"`   //UTC+8 时开盘价
	VolCcy24H float64 `json:"volCcy24h"` //24小时成交量，以币为单位
	Vol24H    float64 `json:"vol24h"`    //24小时成交量，以张为单位
	Ts        uint64  `json:"ts"`        //数据产生时间，Unix时间戳的毫秒数格式，如 1597026383085
}

type KLine struct {
	Ts          uint64  `json:"ts"`          //开始时间，Unix时间戳的毫秒数格式，如 1597026383085
	Open        float64 `json:"open"`        //开盘价格
	Highest     float64 `json:"highest"`     //最高价格
	Lowest      float64 `json:"lowest"`      //最低价格
	Close       float64 `json:"close"`       //收盘价格
	Vol         float64 `json:"vol"`         //交易量，以张为单位
	VolCcy      float64 `json:"volCcy"`      //交易量，交易量，以币为单位
	VolCcyQuote float64 `json:"volCcyQuote"` //交易量，交易量，以计价货币为单位
}

type Trade struct {
	Id      string
	OrderId string
}

type Order struct {
	Id string
}

func (o Order) IsActive() bool {
	return false
}
