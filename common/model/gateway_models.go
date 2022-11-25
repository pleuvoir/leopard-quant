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

type KLine struct {
}
