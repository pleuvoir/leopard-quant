package okx

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gookit/color"
	"github.com/gorilla/websocket"
	"leopard-quant/core/event"
	"leopard-quant/gateway"
	"leopard-quant/util/recws"
	"time"
)

type MarketApi struct {
	unmarshalerOptions *gateway.UnmarshalerOptions
	apiOptions         *gateway.ApiOptions
	conn               *recws.RecConn
	callback           gateway.ApiCallback
	subscribeCmd       []Cmd
	eventEngine        *event.Engine
}

func NewMarket(baseApi *gateway.BaseApi, callback gateway.ApiCallback) *MarketApi {
	m := &MarketApi{}
	m.unmarshalerOptions = baseApi.UnmarshalerOptions
	m.apiOptions = baseApi.ApiOptions
	m.callback = callback
	m.conn = &recws.RecConn{
		KeepAliveTimeout: 60 * time.Second,
		NonVerbose:       true,
	}
	return m
}

func (m *MarketApi) connect() {
	m.conn.Dial(m.apiOptions.Addr, nil)
}

type Cmd struct {
	Op   string    `json:"op"`
	Args []ArgItem `json:"args"`
}

type ArgItem struct {
	Channel string `json:"channel"`
	InstId  string `json:"instId"`
}

func (m *MarketApi) Subscribe(item ArgItem) {
	cmd := Cmd{
		Op:   "subscribe",
		Args: []ArgItem{item},
	}
	m.subscribeCmd = append(m.subscribeCmd, cmd)
	_ = m.SendCmd(cmd)
}

func (m *MarketApi) UnSubscribe(item ArgItem) {
	cmd := Cmd{
		Op:   "unsubscribe",
		Args: []ArgItem{item},
	}
	m.subscribeCmd = append(m.subscribeCmd, cmd)
	_ = m.SendCmd(cmd)
}

func (m *MarketApi) SendCmd(cmd Cmd) error {
	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return m.Send(string(data))
}

func (m *MarketApi) Send(msg string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Ws send error: %v", r))
		}
	}()
	err = m.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	return err
}

func (m *MarketApi) Start() error {
	//连接
	m.connect()
	cancel := make(chan struct{})
	//定时ping
	go func() {
		t := time.NewTicker(time.Second * 5)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				m.ping()
			case <-cancel:
				return
			}
		}
	}()
	//接收消息并处理回调
	go func() {
		defer close(cancel)
		for {
			_, data, err := m.conn.ReadMessage()
			if err != nil {
				color.Redln(fmt.Sprintf("Ws Read error, closing connection: %v", err))
				//	m.conn.Close()
				time.Sleep(time.Second * 3)
				return
			}
			m.processMessagePipeline(data)
		}
	}()
	return nil
}

func (m *MarketApi) processMessagePipeline(data []byte) {

	debugMode := m.apiOptions.DebugMode
	s := string(data)
	if debugMode {
		color.Greenln(fmt.Sprintf("Ws %v", s))
	}

	if match, _ := m.unmarshalerOptions.SubscribeResponseUnmarshaler(data); match {
		if debugMode {
			color.Greenln(fmt.Sprintf("订阅应答 %s", s))
		}
		return
	}

	if match, _ := m.unmarshalerOptions.PongResponseUnmarshaler(data); match {
		if debugMode {
			color.Greenln(fmt.Sprintf("心跳应答 %s", s))
		}
		return
	}

	if match, ticker, _ := m.unmarshalerOptions.GetTickerResponseUnmarshaler(data); match {
		if debugMode {
			color.Greenln(fmt.Sprintf("TickerCallback %s", s))
		}
		m.callback.TickerCallback(ticker)
		return
	}

	if match, kLine, _ := m.unmarshalerOptions.GetKlineResponseUnmarshaler(data); match {
		if debugMode {
			color.Greenln(fmt.Sprintf("KlineCallback %s", s))
		}
		m.callback.KlineCallback(kLine)
		return
	}

}

func (m *MarketApi) ping() {
	defer func() {
		if r := recover(); r != nil {
			color.Redln(fmt.Sprintf("Ws ping error: %v", r))
		}
	}()
	if !m.conn.IsConnected() {
		return
	}
	err := m.conn.WriteMessage(websocket.TextMessage, []byte(`ping`))
	if err != nil {
		color.Redln(fmt.Sprintf("Ws ping error: %v", err))
	}
}
