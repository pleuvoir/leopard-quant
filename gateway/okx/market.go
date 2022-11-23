package okx

import "C"
import (
	"encoding/json"
	"github.com/gookit/color"
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
	err = m.SubscribeTickers(symbol, c.TickerCallback)
	return err
}

func (m *Market) SubscribeTickers(symbol string, callback gateway.TickerCallback) error {
	//发送订阅消息
	items := []ArgItem{{Channel: "tickers", InstId: symbol}}
	req := SubscribeReq{Op: "subscribe", Args: items}
	err := m.ws.SendJSONTextMessage(req)
	if err != nil {
		return err
	}
	go func() {
		for {
			message, err := m.ws.ReadMessage()
			if err != nil {
				color.Redf("read:", err)
				return
			}
			ticker := convert2Ticker(message)
			log.Infof("接收到okx回调。%+v", ticker)
			callback(ticker)
		}
	}()
	return err
}

func convert2Ticker(data []byte) gateway.Ticker {
	ticker := gateway.Ticker{}
	t := Tickers{}
	_ = json.Unmarshal(data, &t)
	for _, datum := range t.Data {
		float, _ := strconv.ParseFloat(datum.Last, 64)
		ticker.Last = float
	}
	return ticker
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
