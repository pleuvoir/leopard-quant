package okx

import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"leopard-quant/common/model"
	"leopard-quant/core/log"
	"leopard-quant/gateway"
	"strconv"
)

type Market struct {
	ws *gateway.Websocket
}

func NewMarketWS() Market {
	m := Market{}
	ws := gateway.NewWebsocket("wss", "ws.okx.com:8443", "/ws/v5/public")
	m.ws = ws
	return m
}

func (m *Market) Connect() error {
	return m.ws.Connect()
}

func (m *Market) Subscribe(symbol string, c gateway.ComposeCallback) (err error) {
	//if err = m.SubscribeTickers(symbol, c.TickerCallback); err != nil {
	//	return err
	//}
	if err = m.SubscribeKLine(symbol, c.KlineCallback); err != nil {
		return err
	}
	return err
}

func (m *Market) SubscribeKLine(symbol string, callback gateway.KlineCallback) error {
	//发送订阅消息
	req := SubscribeReq{Op: "subscribe", Args: []ArgItem{{Channel: "candle3M", InstId: symbol}}}
	err := m.ws.SendJSONTextMessage(req)
	if err != nil {
		return err
	}
	go func() {
		for {
			message, err := m.ws.ReadMessage()
			if err != nil {
				color.Redln(fmt.Sprintf("read:%s", err))
				continue
			}
			log.Infoln(fmt.Sprintf("%s", message))
			break
			//if kLine, err := convert2KLine(message); err == nil {
			//	callback(kLine)
			//}
		}
	}()
	return err
}

func (m *Market) SubscribeTickers(symbol string, callback gateway.TickerCallback) error {
	//发送订阅消息
	req := SubscribeReq{Op: "subscribe", Args: []ArgItem{{Channel: "tickers", InstId: symbol}}}
	err := m.ws.SendJSONTextMessage(req)
	if err != nil {
		return err
	}
	go func() {
		for {
			message, err := m.ws.ReadMessage()
			if err != nil {
				color.Redln(fmt.Sprintf("read:%s", err))
				continue
			}
			//	log.Infoln(fmt.Sprintf("%s", message))
			if ticker, err := convert2Ticker(message); err == nil {
				callback(ticker)
			}
		}
	}()
	return err
}

func convert2KLine(data []byte) (kLine model.KLine, err error) {
	t := Tickers{}
	if err := json.Unmarshal(data, &t); err != nil {
		return kLine, err
	}
	//for _, v := range t.Data {
	//	ticker.Symbol = v.InstId
	//	ticker.Last = parseFloat64(v.Last)
	//	ticker.LastSz = parseFloat64(v.LastSz)
	//	ticker.AskPx = parseFloat64(v.AskPx)
	//	ticker.AskSz = parseFloat64(v.AskSz)
	//	ticker.BidPx = parseFloat64(v.BidPx)
	//	ticker.BidSz = parseFloat64(v.BidSz)
	//	ticker.Open24H = parseFloat64(v.Open24H)
	//	ticker.High24H = parseFloat64(v.High24H)
	//	ticker.Low24H = parseFloat64(v.Low24H)
	//	ticker.SodUtc0 = parseFloat64(v.SodUtc0)
	//	ticker.SodUtc8 = parseFloat64(v.SodUtc8)
	//	ticker.VolCcy24H = parseFloat64(v.VolCcy24H)
	//	ticker.Vol24H = parseFloat64(v.Vol24H)
	//	ticker.Ts = parseUint64(v.Ts)
	//}
	return kLine, nil
}

func convert2Ticker(data []byte) (ticker model.Ticker, err error) {
	t := Tickers{}
	if err := json.Unmarshal(data, &t); err != nil {
		return ticker, err
	}
	for _, v := range t.Data {
		ticker.Symbol = v.InstId
		ticker.Last = parseFloat64(v.Last)
		ticker.LastSz = parseFloat64(v.LastSz)
		ticker.AskPx = parseFloat64(v.AskPx)
		ticker.AskSz = parseFloat64(v.AskSz)
		ticker.BidPx = parseFloat64(v.BidPx)
		ticker.BidSz = parseFloat64(v.BidSz)
		ticker.Open24H = parseFloat64(v.Open24H)
		ticker.High24H = parseFloat64(v.High24H)
		ticker.Low24H = parseFloat64(v.Low24H)
		ticker.SodUtc0 = parseFloat64(v.SodUtc0)
		ticker.SodUtc8 = parseFloat64(v.SodUtc8)
		ticker.VolCcy24H = parseFloat64(v.VolCcy24H)
		ticker.Vol24H = parseFloat64(v.Vol24H)
		ticker.Ts = parseUint64(v.Ts)
	}
	return ticker, nil
}

func parseFloat64(val string) float64 {
	float, err := strconv.ParseFloat(val, 64)
	if err != nil {
		panic(err)
	}
	return float
}

func parseUint64(val string) uint64 {
	u, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		panic(err)
	}
	return u
}

func (m *Market) CancelSubscribe(symbol string) error {
	return nil
}

type ArgItem struct {
	Channel string `json:"channel"`
	InstId  string `json:"instId"`
}

type SubscribeReq struct {
	Op   string    `json:"op"`
	Args []ArgItem `json:"args"`
}

type Tickers struct {
	Arg struct {
		Channel string `json:"channel"`
		InstId  string `json:"instId"`
	} `json:"arg"`
	Data []struct {
		InstType  string `json:"instType"`
		InstId    string `json:"instId"`
		Last      string `json:"last"`
		LastSz    string `json:"lastSz"`
		AskPx     string `json:"askPx"`
		AskSz     string `json:"askSz"`
		BidPx     string `json:"bidPx"`
		BidSz     string `json:"bidSz"`
		Open24H   string `json:"open24h"`
		High24H   string `json:"high24h"`
		Low24H    string `json:"low24h"`
		SodUtc0   string `json:"sodUtc0"`
		SodUtc8   string `json:"sodUtc8"`
		VolCcy24H string `json:"volCcy24h"`
		Vol24H    string `json:"vol24h"`
		Ts        string `json:"ts"`
	} `json:"data"`
}
